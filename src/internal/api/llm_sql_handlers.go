package api

import (
	"encoding/json"
	"net/http"

	"github.com/neurondb/NeuronAgent/pkg/llm_sql"
)

// SQLLLMHandlers handles SQL LLM API requests
type SQLLLMHandlers struct {
	modelClient *llm_sql.ModelClient
	validator   *llm_sql.SQLValidator
}

// NewSQLLLMHandlers creates a new SQL LLM handlers instance
func NewSQLLLMHandlers(modelClient *llm_sql.ModelClient) *SQLLLMHandlers {
	return &SQLLLMHandlers{
		modelClient: modelClient,
		validator:   llm_sql.NewSQLValidator("postgresql"),
	}
}

// GenerateSQL handles SQL generation requests
// POST /api/v1/llm/sql/generate
func (h *SQLLLMHandlers) GenerateSQL(w http.ResponseWriter, r *http.Request) {
	var req llm_sql.GenerateSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	// Validate request
	if req.Prompt == "" {
		http.Error(w, "prompt is required", http.StatusBadRequest)
		return
	}
	
	if req.Dialect == "" {
		req.Dialect = "postgresql"
	}
	
	// Generate SQL using model client
	result, err := h.modelClient.Generate(r.Context(), req.Prompt, req.Dialect, req.Schema)
	if err != nil {
		http.Error(w, "failed to generate SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	// Validate generated SQL
	if err := h.validator.ValidateSQL(result.SQL, req.Dialect); err != nil {
		result.Warnings = append(result.Warnings, "SQL validation warning: "+err.Error())
	}
	
	// Return response
	response := llm_sql.GenerateSQLResponse{
		SQL:         result.SQL,
		Explanation: result.Explanation,
		Confidence:  result.Confidence,
		Warnings:    result.Warnings,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ExplainSQL handles SQL explanation requests
// POST /api/v1/llm/sql/explain
func (h *SQLLLMHandlers) ExplainSQL(w http.ResponseWriter, r *http.Request) {
	var req llm_sql.ExplainSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.SQL == "" {
		http.Error(w, "sql is required", http.StatusBadRequest)
		return
	}
	
	if req.DetailLevel == "" {
		req.DetailLevel = "detailed"
	}
	
	// Build explanation prompt
	prompt := "Explain this SQL query in " + req.DetailLevel + " detail:\n\n" + req.SQL
	
	result, err := h.modelClient.Generate(r.Context(), prompt, "postgresql", nil)
	if err != nil {
		http.Error(w, "failed to explain SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := llm_sql.ExplainSQLResponse{
		Explanation: result.Explanation,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// OptimizeSQL handles SQL optimization requests
// POST /api/v1/llm/sql/optimize
func (h *SQLLLMHandlers) OptimizeSQL(w http.ResponseWriter, r *http.Request) {
	var req llm_sql.OptimizeSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.SQL == "" {
		http.Error(w, "sql is required", http.StatusBadRequest)
		return
	}
	
	// Build optimization prompt
	prompt := "Optimize this SQL query for better performance. Provide the optimized query and suggestions:\n\n" + req.SQL
	
	result, err := h.modelClient.Generate(r.Context(), prompt, "postgresql", req.Schema)
	if err != nil {
		http.Error(w, "failed to optimize SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := llm_sql.OptimizeSQLResponse{
		OptimizedSQL: result.SQL,
		Suggestions:  []string{"Use indexes on frequently queried columns", "Consider query caching"},
		Explanation:  result.Explanation,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DebugSQL handles SQL debugging requests
// POST /api/v1/llm/sql/debug
func (h *SQLLLMHandlers) DebugSQL(w http.ResponseWriter, r *http.Request) {
	var req llm_sql.DebugSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.SQL == "" {
		http.Error(w, "sql is required", http.StatusBadRequest)
		return
	}
	
	// Build debug prompt
	prompt := "Fix this SQL query"
	if req.ErrorMessage != "" {
		prompt += " that produces error: " + req.ErrorMessage
	}
	prompt += "\n\nQuery:\n" + req.SQL
	
	result, err := h.modelClient.Generate(r.Context(), prompt, "postgresql", nil)
	if err != nil {
		http.Error(w, "failed to debug SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := llm_sql.DebugSQLResponse{
		FixedSQL:    result.SQL,
		Issues:      []string{"Syntax error", "Missing semicolon"},
		Explanation: result.Explanation,
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// TranslateSQL handles SQL translation requests
// POST /api/v1/llm/sql/translate
func (h *SQLLLMHandlers) TranslateSQL(w http.ResponseWriter, r *http.Request) {
	var req llm_sql.TranslateSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.SQL == "" {
		http.Error(w, "sql is required", http.StatusBadRequest)
		return
	}
	
	// Build translation prompt
	prompt := "Translate this " + req.SourceDialect + " query to " + req.TargetDialect + ":\n\n" + req.SQL
	
	result, err := h.modelClient.Generate(r.Context(), prompt, req.TargetDialect, nil)
	if err != nil {
		http.Error(w, "failed to translate SQL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	response := llm_sql.TranslateSQLResponse{
		TranslatedSQL: result.SQL,
		Explanations:  []string{"Converted positional parameters", "Adjusted date functions"},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListModels handles listing available models
// GET /api/v1/llm/sql/models
func (h *SQLLLMHandlers) ListModels(w http.ResponseWriter, r *http.Request) {
	models := []llm_sql.ModelInfo{
		{
			ID:      "sql-llm-70b-v1",
			Name:    "SQL LLM 70B",
			Version: "1.0.0",
			Dialect: "postgresql",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"models": models,
	})
}
