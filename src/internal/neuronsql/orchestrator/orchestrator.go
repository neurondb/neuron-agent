/*-------------------------------------------------------------------------
 *
 * orchestrator.go
 *    NeuronSQL orchestrator: generate, validate, optimize, plpgsql pipelines
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/orchestrator/orchestrator.go
 *
 *-------------------------------------------------------------------------
 */

package orchestrator

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/prompting"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/tools"
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* Orchestrator runs the NeuronSQL pipelines */
type Orchestrator struct {
	LLM     neuronsql.LLMProvider
	Policy  neuronsql.PolicyEngine
	Factory tools.ConnectionFactory
	Retrieve neuronsql.Retriever
}

/* NewOrchestrator creates an orchestrator with the given dependencies */
func NewOrchestrator(llm neuronsql.LLMProvider, policy neuronsql.PolicyEngine, factory tools.ConnectionFactory, retriever neuronsql.Retriever) *Orchestrator {
	return &Orchestrator{LLM: llm, Policy: policy, Factory: factory, Retrieve: retriever}
}

/* RunGenerate runs the generate pipeline: intent -> schema snapshot -> retrieve -> LLM -> validate -> response */
func (o *Orchestrator) RunGenerate(ctx context.Context, dsn string, question string, requestID string) (*neuronsql.NeuronSQLResponse, error) {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	resp := &neuronsql.NeuronSQLResponse{RequestID: requestID, Mode: "generate"}
	conn, err := o.Factory(ctx, dsn)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	defer conn.Close()

	schemaJSON, err := tools.RunSchemaSnapshot(ctx, conn, requestID)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	resp.Citations = append(resp.Citations, "tool:schema_snapshot:"+requestID)

	var docContext string
	if o.Retrieve != nil {
		chunks, _ := o.Retrieve.Retrieve(ctx, question, 6)
		for _, c := range chunks {
			docContext += c.Content + "\n\n"
			resp.Citations = append(resp.Citations, "doc:"+c.ID)
		}
	}
	if docContext == "" {
		docContext = "(no docs)"
	}

	intent := prompting.ParseIntent(question, "generate")
	prompt := prompting.BuildGeneratePrompt(schemaJSON, docContext, intent.Question)
	completion, err := o.LLM.Complete(ctx, []neuronsql.Message{{Role: "user", Content: prompt}}, nil, neuronsql.CompletionSettings{MaxTokens: 1024})
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	sql := extractSQL(completion.Text)
	resp.SQL = sql
	resp.Explanation = completion.Text

	const maxRepairIterations = 3
	for iter := 0; iter < maxRepairIterations; iter++ {
		validation, err := tools.RunValidateSQLWithSnapshot(ctx, conn, requestID, sql, schemaJSON)
		if err != nil {
			resp.ValidationReport = &neuronsql.ValidationReport{Valid: false, Errors: []string{err.Error()}}
			return resp, nil
		}
		var vr neuronsql.ValidationReport
		_ = json.Unmarshal([]byte(validation), &vr)
		resp.ValidationReport = &vr
		if vr.Valid {
			break
		}
		if iter == maxRepairIterations-1 {
			resp.Explanation = resp.Explanation + "; Validation: " + strings.Join(vr.Errors, "; ")
			break
		}
		/* Ask LLM to repair */
		repairPrompt := "The following SQL failed validation. Errors: " + strings.Join(vr.Errors, "; ") + "\n\nSQL:\n" + sql + "\n\nReturn only the corrected SQL, no explanation."
		repairCompletion, err := o.LLM.Complete(ctx, []neuronsql.Message{{Role: "user", Content: repairPrompt}}, nil, neuronsql.CompletionSettings{MaxTokens: 1024})
		if err != nil {
			resp.Explanation = resp.Explanation + "; Validation: " + strings.Join(vr.Errors, "; ")
			break
		}
		sql = extractSQL(repairCompletion.Text)
		resp.SQL = sql
	}
	return resp, nil
}

/* RunValidate runs the validate pipeline */
func (o *Orchestrator) RunValidate(ctx context.Context, dsn string, sql string, requestID string) (*neuronsql.NeuronSQLResponse, error) {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	resp := &neuronsql.NeuronSQLResponse{RequestID: requestID, Mode: "validate"}
	decision, err := o.Policy.Check(ctx, sql, neuronsql.PolicyContext{RequestID: requestID})
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	if !decision.Allowed {
		resp.ValidationReport = &neuronsql.ValidationReport{Valid: false, Errors: []string{decision.Reason}, StatementClass: decision.StatementClass}
		resp.Safety = &neuronsql.SafetyInfo{PolicyAllowed: false, BlockedReason: decision.Reason}
		return resp, nil
	}
	conn, err := o.Factory(ctx, dsn)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	defer conn.Close()
	validation, err := tools.RunValidateSQL(ctx, conn, requestID, sql)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	var vr neuronsql.ValidationReport
	_ = json.Unmarshal([]byte(validation), &vr)
	resp.ValidationReport = &vr
	resp.Safety = &neuronsql.SafetyInfo{PolicyAllowed: true}
	return resp, nil
}

/* RunOptimize runs the optimize pipeline */
func (o *Orchestrator) RunOptimize(ctx context.Context, dsn string, sql string, requestID string) (*neuronsql.NeuronSQLResponse, error) {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	resp := &neuronsql.NeuronSQLResponse{RequestID: requestID, Mode: "optimize"}
	conn, err := o.Factory(ctx, dsn)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	defer conn.Close()
	planJSON, err := tools.RunExplainJSON(ctx, conn, requestID, sql)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	resp.PlanSummary = planJSON
	optResult, err := tools.RunOptimizeCandidates(ctx, conn, requestID, sql, planJSON)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	var opt tools.OptimizeCandidatesResult
	_ = json.Unmarshal([]byte(optResult), &opt)
	for _, r := range opt.RewriteOptions {
		resp.RewriteOptions = append(resp.RewriteOptions, neuronsql.RewriteOption{SQL: r.SQL, Explanation: r.Explanation, Risk: r.Risk, PlanEvidence: r.PlanEvidence})
	}
	for _, i := range opt.IndexSuggestions {
		resp.IndexSuggestions = append(resp.IndexSuggestions, neuronsql.IndexSuggestion{Definition: i.Definition, EstimatedBenefit: i.EstimatedBenefit, PlanEvidence: i.PlanEvidence})
	}
	resp.VerificationQueries = opt.VerificationQueries
	return resp, nil
}

/* RunPLpgSQL runs the plpgsql generation pipeline */
func (o *Orchestrator) RunPLpgSQL(ctx context.Context, dsn string, signature string, purpose string, requestID string) (*neuronsql.NeuronSQLResponse, error) {
	if requestID == "" {
		requestID = uuid.New().String()
	}
	resp := &neuronsql.NeuronSQLResponse{RequestID: requestID, Mode: "plpgsql"}
	conn, err := o.Factory(ctx, dsn)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	defer conn.Close()
	schemaJSON, err := tools.RunSchemaSnapshot(ctx, conn, requestID)
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	prompt := prompting.BuildPLpgSQLPrompt(signature, purpose, schemaJSON)
	completion, err := o.LLM.Complete(ctx, []neuronsql.Message{{Role: "user", Content: prompt}}, nil, neuronsql.CompletionSettings{MaxTokens: 2048})
	if err != nil {
		resp.Error = err.Error()
		return resp, err
	}
	resp.PLpgSQLCode = completion.Text
	resp.Explanation = "Generated PL/pgSQL (no CREATE FUNCTION executed in v1)."
	return resp, nil
}

func extractSQL(text string) string {
	text = strings.TrimSpace(text)
	if idx := strings.Index(text, "<sql>"); idx >= 0 {
		if end := strings.Index(text, "</sql>"); end > idx {
			return strings.TrimSpace(text[idx+5 : end])
		}
	}
	if idx := strings.Index(text, "```sql"); idx >= 0 {
		rest := text[idx+6:]
		if end := strings.Index(rest, "```"); end >= 0 {
			return strings.TrimSpace(rest[:end])
		}
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(line), "SELECT") || strings.HasPrefix(strings.ToUpper(line), "WITH") || strings.HasPrefix(strings.ToUpper(line), "EXPLAIN") {
			return line
		}
	}
	return text
}
