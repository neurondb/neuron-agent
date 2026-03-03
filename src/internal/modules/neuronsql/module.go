/*-------------------------------------------------------------------------
 *
 * module.go
 *    NeuronSQL module: first module for the modular NeuronAgent architecture.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/modules/neuronsql/module.go
 *
 *-------------------------------------------------------------------------
 */

package neuronsql

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/neurondb/NeuronAgent/internal/api"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/metrics"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/orchestrator"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/provider"
	neuronsqltools "github.com/neurondb/NeuronAgent/internal/neuronsql/tools"
	"github.com/neurondb/NeuronAgent/internal/tools"
	"github.com/neurondb/NeuronAgent/pkg/module"
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

const moduleName = "neuronsql"
const moduleVersion = "1.0.0"

/* NeuronSQLModule implements module.Module for NeuronSQL. */
type NeuronSQLModule struct {
	orch *orchestrator.Orchestrator
}

/* New returns the NeuronSQL module. */
func New() module.Module {
	return &NeuronSQLModule{}
}

/* Name implements module.Module. */
func (m *NeuronSQLModule) Name() string { return moduleName }

/* Version implements module.Module. */
func (m *NeuronSQLModule) Version() string { return moduleVersion }

/* Dependencies implements module.Module. */
func (m *NeuronSQLModule) Dependencies() []string { return nil }

/* Init implements module.Module. */
func (m *NeuronSQLModule) Init(ctx context.Context, app module.AppContext) error {
	policyEngine := policy.NewPolicyEngineImpl(nil)
	factory := neuronsqltools.NewConnectionFactory(policyEngine, neuronsqltools.DefaultSafeConnectionConfig())
	pglangCfg := provider.DefaultPGLangConfig()
	if cfg := app.ModuleConfig(moduleName); cfg != nil {
		if e, ok := cfg["pglang_endpoint"].(string); ok && e != "" {
			pglangCfg.Endpoint = e
		}
		if k, ok := cfg["pglang_api_key"].(string); ok && k != "" {
			pglangCfg.APIKey = k
		}
	}
	llmProvider := provider.NewPGLangProvider(pglangCfg)
	m.orch = orchestrator.NewOrchestrator(llmProvider, policyEngine, factory, nil)

	app.RegisterRoutes(moduleName, m.routes)
	/* Register tools under neuronsql.* namespace (no generic sql_execute) */
	app.RegisterToolHandler("neuronsql.schema_snapshot", &neuronsqltools.SchemaSnapshotTool{Factory: factory, Policy: policyEngine})
	app.RegisterToolHandler("neuronsql.validate_sql", &neuronsqltools.ValidateSQLTool{Factory: factory, Policy: policyEngine})
	app.RegisterToolHandler("neuronsql.explain_json", &neuronsqltools.ExplainJSONTool{Factory: factory, Policy: policyEngine})
	app.RegisterToolHandler("neuronsql.generate_select", &generateTool{orch: m.orch})
	app.RegisterToolHandler("neuronsql.optimize_select", &optimizeTool{orch: m.orch})
	app.RegisterToolHandler("neuronsql.plpgsql_generate", &plpgsqlTool{orch: m.orch})
	/* Legacy handler names for backwards compatibility */
	app.RegisterToolHandler("neuronsql_generate", &generateTool{orch: m.orch})
	app.RegisterToolHandler("neuronsql_optimize", &optimizeTool{orch: m.orch})
	return nil
}

/* Start implements module.Module. */
func (m *NeuronSQLModule) Start(ctx context.Context) error { return nil }

/* Stop implements module.Module. */
func (m *NeuronSQLModule) Stop(ctx context.Context) error { return nil }

/* Health implements module.Module. */
func (m *NeuronSQLModule) Health(ctx context.Context) module.HealthStatus {
	return module.HealthStatus{Healthy: true}
}

/* routes registers HTTP handlers with the router. */
func (m *NeuronSQLModule) routes(r module.Router) {
	r.POST("/generate", m.handleGenerate)
	r.POST("/optimize", m.handleOptimize)
	r.POST("/validate", m.handleValidate)
	r.POST("/plpgsql", m.handlePLpgSQL)
}

func (m *NeuronSQLModule) handleGenerate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := api.GetRequestID(r.Context())
	var req neuronsql.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("generate", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "invalid request body", requestID)
		return
	}
	if req.DBDSN == "" {
		metrics.RecordNeuronSQLRequest("generate", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "db_dsn required", requestID)
		return
	}
	if req.Question == "" {
		metrics.RecordNeuronSQLRequest("generate", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "question required", requestID)
		return
	}
	resp, err := m.orch.RunGenerate(r.Context(), req.DBDSN, req.Question, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("generate", "5xx", time.Since(start))
		writeJSONError(w, http.StatusInternalServerError, "generate failed: "+err.Error(), requestID)
		return
	}
	metrics.RecordNeuronSQLRequest("generate", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (m *NeuronSQLModule) handleOptimize(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := api.GetRequestID(r.Context())
	var req neuronsql.OptimizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("optimize", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "invalid request body", requestID)
		return
	}
	if req.DBDSN == "" {
		metrics.RecordNeuronSQLRequest("optimize", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "db_dsn required", requestID)
		return
	}
	if req.SQL == "" {
		metrics.RecordNeuronSQLRequest("optimize", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "sql required", requestID)
		return
	}
	resp, err := m.orch.RunOptimize(r.Context(), req.DBDSN, req.SQL, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("optimize", "5xx", time.Since(start))
		writeJSONError(w, http.StatusInternalServerError, "optimize failed: "+err.Error(), requestID)
		return
	}
	metrics.RecordNeuronSQLRequest("optimize", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (m *NeuronSQLModule) handleValidate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := api.GetRequestID(r.Context())
	var req neuronsql.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("validate", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "invalid request body", requestID)
		return
	}
	if req.DBDSN == "" {
		metrics.RecordNeuronSQLRequest("validate", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "db_dsn required", requestID)
		return
	}
	if req.SQL == "" {
		metrics.RecordNeuronSQLRequest("validate", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "sql required", requestID)
		return
	}
	resp, err := m.orch.RunValidate(r.Context(), req.DBDSN, req.SQL, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("validate", "5xx", time.Since(start))
		writeJSONError(w, http.StatusInternalServerError, "validate failed: "+err.Error(), requestID)
		return
	}
	metrics.RecordNeuronSQLRequest("validate", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (m *NeuronSQLModule) handlePLpgSQL(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := api.GetRequestID(r.Context())
	var req neuronsql.PLpgSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "invalid request body", requestID)
		return
	}
	if req.DBDSN == "" {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "db_dsn required", requestID)
		return
	}
	if req.Signature == "" {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "signature required", requestID)
		return
	}
	if req.Purpose == "" {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		writeJSONError(w, http.StatusBadRequest, "purpose required", requestID)
		return
	}
	resp, err := m.orch.RunPLpgSQL(r.Context(), req.DBDSN, req.Signature, req.Purpose, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("plpgsql", "5xx", time.Since(start))
		writeJSONError(w, http.StatusInternalServerError, "plpgsql generate failed: "+err.Error(), requestID)
		return
	}
	metrics.RecordNeuronSQLRequest("plpgsql", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func writeJSONError(w http.ResponseWriter, code int, message, requestID string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	body := map[string]string{"error": message}
	if requestID != "" {
		body["request_id"] = requestID
	}
	_ = json.NewEncoder(w).Encode(body)
}

/* generateTool implements tools.ToolHandler for neuronsql_generate. */
type generateTool struct {
	orch *orchestrator.Orchestrator
}

func (t *generateTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	question, _ := args["question"].(string)
	if dsn == "" || question == "" {
		return "", errors.New("neuronsql_generate: db_dsn and question required")
	}
	resp, err := t.orch.RunGenerate(ctx, dsn, question, "")
	if err != nil {
		return "", err
	}
	out, _ := json.Marshal(resp)
	return string(out), nil
}

func (t *generateTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	return tools.ValidateArgs(args, schema)
}

/* optimizeTool implements tools.ToolHandler for neuronsql_optimize. */
type optimizeTool struct {
	orch *orchestrator.Orchestrator
}

func (t *optimizeTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	sqlStr, _ := args["sql"].(string)
	if dsn == "" || sqlStr == "" {
		return "", errors.New("neuronsql_optimize: db_dsn and sql required")
	}
	resp, err := t.orch.RunOptimize(ctx, dsn, sqlStr, "")
	if err != nil {
		return "", err
	}
	out, _ := json.Marshal(resp)
	return string(out), nil
}

func (t *optimizeTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	return tools.ValidateArgs(args, schema)
}

/* plpgsqlTool implements tools.ToolHandler for neuronsql.plpgsql_generate */
type plpgsqlTool struct {
	orch *orchestrator.Orchestrator
}

func (t *plpgsqlTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	signature, _ := args["signature"].(string)
	purpose, _ := args["purpose"].(string)
	if dsn == "" || signature == "" || purpose == "" {
		return "", errors.New("neuronsql.plpgsql_generate: db_dsn, signature and purpose required")
	}
	resp, err := t.orch.RunPLpgSQL(ctx, dsn, signature, purpose, "")
	if err != nil {
		return "", err
	}
	out, _ := json.Marshal(resp)
	return string(out), nil
}

func (t *plpgsqlTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	return tools.ValidateArgs(args, schema)
}
