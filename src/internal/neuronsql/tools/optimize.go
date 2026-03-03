/*-------------------------------------------------------------------------
 *
 * optimize.go
 *    tool_optimize_candidates: rewrite options and index suggestions from plan
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/optimize.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
)

/* OptimizeCandidatesResult is the JSON result of tool_optimize_candidates */
type OptimizeCandidatesResult struct {
	RewriteOptions     []RewriteOption   `json:"rewrite_options"`
	IndexSuggestions   []IndexSuggestion `json:"index_suggestions"`
	VerificationQueries []string         `json:"verification_queries"`
	PlanSummary        string           `json:"plan_summary"`
}

type RewriteOption struct {
	SQL          string `json:"sql"`
	Explanation  string `json:"explanation"`
	Risk         string `json:"risk"`
	PlanEvidence string `json:"plan_evidence"`
}

type IndexSuggestion struct {
	Definition       string `json:"definition"`
	EstimatedBenefit string `json:"estimated_benefit"`
	PlanEvidence    string `json:"plan_evidence"`
}

/* RunOptimizeCandidates produces rewrite options and index suggestions from SQL and plan JSON */
func RunOptimizeCandidates(ctx context.Context, conn *SafeConnection, requestID string, sql string, planJSON string) (string, error) {
	var result OptimizeCandidatesResult
	result.PlanSummary = planJSON
	if len(planJSON) > 500 {
		result.PlanSummary = planJSON[:500] + "..."
	}
	result.RewriteOptions = []RewriteOption{
		{SQL: sql, Explanation: "Original query; no rewrite applied.", Risk: "none", PlanEvidence: "input"},
	}
	result.IndexSuggestions = []IndexSuggestion{}
	result.VerificationQueries = []string{"EXPLAIN (FORMAT JSON) " + sql}

	if root, allNodes, err := ParseExplainJSON(planJSON); err == nil {
		_ = root
		if idx := SuggestIndexFromPlan(allNodes); len(idx) > 0 {
			result.IndexSuggestions = idx
		}
		if rew := SuggestRewritesFromPlan(allNodes, sql); len(rew) > 0 {
			result.RewriteOptions = append(result.RewriteOptions, rew...)
		}
	}

	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

/* OptimizeTool implements tools.ToolHandler */
type OptimizeTool struct {
	Factory ConnectionFactory
	Policy  *policy.PolicyEngineImpl
}

func (t *OptimizeTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("optimize_candidates: db_dsn required")
	}
	sql, _ := args["sql"].(string)
	if sql == "" {
		return "", fmt.Errorf("optimize_candidates: sql required")
	}
	planJSON, _ := args["plan_json"].(string)
	requestID, _ := args["request_id"].(string)
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if planJSON == "" {
		planJSON, err = RunExplainJSON(ctx, conn, requestID, sql)
		if err != nil {
			return "", err
		}
	}
	return RunOptimizeCandidates(ctx, conn, requestID, sql, planJSON)
}

func (t *OptimizeTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	if _, ok := args["sql"]; !ok {
		return fmt.Errorf("sql required")
	}
	return nil
}
