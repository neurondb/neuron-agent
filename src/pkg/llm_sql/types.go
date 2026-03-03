package llm_sql

// GenerateSQLRequest represents a request to generate SQL
type GenerateSQLRequest struct {
	Prompt  string          `json:"prompt" binding:"required"`
	Dialect string          `json:"dialect" binding:"required,oneof=postgresql mysql"`
	Schema  json.RawMessage `json:"schema,omitempty"`
	Options GenerationOptions `json:"options,omitempty"`
}

// GenerateSQLResponse represents the response from SQL generation
type GenerateSQLResponse struct {
	SQL         string   `json:"sql"`
	Explanation string   `json:"explanation"`
	Confidence  float64  `json:"confidence"`
	Warnings    []string `json:"warnings,omitempty"`
}

// ExplainSQLRequest represents a request to explain SQL
type ExplainSQLRequest struct {
	SQL         string `json:"sql" binding:"required"`
	DetailLevel string `json:"detail_level" binding:"omitempty,oneof=brief detailed expert"`
}

// ExplainSQLResponse represents the response from SQL explanation
type ExplainSQLResponse struct {
	Explanation string `json:"explanation"`
}

// OptimizeSQLRequest represents a request to optimize SQL
type OptimizeSQLRequest struct {
	SQL    string          `json:"sql" binding:"required"`
	Schema json.RawMessage `json:"schema,omitempty"`
}

// OptimizeSQLResponse represents the response from SQL optimization
type OptimizeSQLResponse struct {
	OptimizedSQL string   `json:"optimized_sql"`
	Suggestions  []string `json:"suggestions"`
	Explanation  string   `json:"explanation"`
}

// DebugSQLRequest represents a request to debug SQL
type DebugSQLRequest struct {
	SQL          string `json:"sql" binding:"required"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// DebugSQLResponse represents the response from SQL debugging
type DebugSQLResponse struct {
	FixedSQL    string   `json:"fixed_sql"`
	Issues      []string `json:"issues"`
	Explanation string   `json:"explanation"`
}

// TranslateSQLRequest represents a request to translate SQL
type TranslateSQLRequest struct {
	SQL            string `json:"sql" binding:"required"`
	SourceDialect  string `json:"source_dialect" binding:"required,oneof=postgresql mysql"`
	TargetDialect  string `json:"target_dialect" binding:"required,oneof=postgresql mysql"`
}

// TranslateSQLResponse represents the response from SQL translation
type TranslateSQLResponse struct {
	TranslatedSQL string   `json:"translated_sql"`
	Explanations  []string `json:"explanations"`
}

// GenerationOptions contains optional parameters for generation
type GenerationOptions struct {
	Temperature *float64 `json:"temperature,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
}

// ModelInfo represents information about a model
type ModelInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Dialect string `json:"dialect"`
}
