/*-------------------------------------------------------------------------
 *
 * models.go
 *    NeuronSQL request/response models (NeuronSQLResponse and related)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/pkg/neuronsql/models.go
 *
 *-------------------------------------------------------------------------
 */

package neuronsql

/* Citation format: tool:tool_name:request_id or doc:chunk_id */

/* NeuronSQLResponse is the unified JSON response for all NeuronSQL endpoints */
type NeuronSQLResponse struct {
	RequestID          string                 `json:"request_id"`
	Mode               string                 `json:"mode"` // "generate", "validate", "optimize", "plpgsql"
	SQL                string                 `json:"sql,omitempty"`
	PLpgSQLCode        string                 `json:"plpgsql_code,omitempty"`
	Params             map[string]interface{}  `json:"params,omitempty"`
	Assumptions        []string               `json:"assumptions,omitempty"`
	Explanation        string                 `json:"explanation,omitempty"`
	Citations          []string               `json:"citations,omitempty"` // tool:name:id or doc:chunk_id
	ValidationReport   *ValidationReport      `json:"validation_report,omitempty"`
	PlanSummary        string                 `json:"plan_summary,omitempty"`
	RewriteOptions     []RewriteOption        `json:"rewrite_options,omitempty"`
	IndexSuggestions   []IndexSuggestion      `json:"index_suggestions,omitempty"`
	VerificationQueries []string             `json:"verification_queries,omitempty"`
	Safety             *SafetyInfo            `json:"safety,omitempty"`
	TimingsMs          map[string]float64    `json:"timings_ms,omitempty"`
	Error              string                 `json:"error,omitempty"`
}

/* ValidationReport is the result of tool_validate_sql */
type ValidationReport struct {
	Valid             bool     `json:"valid"`
	Errors            []string `json:"errors,omitempty"`
	ReferencedTables  []string `json:"referenced_tables,omitempty"`
	ReferencedColumns []string `json:"referenced_columns,omitempty"`
	StatementClass    string   `json:"statement_class,omitempty"`
	RiskyPatterns    []string `json:"risky_patterns,omitempty"`
}

/* RewriteOption is a single optimization rewrite with explanation and risk */
type RewriteOption struct {
	SQL       string `json:"sql"`
	Explanation string `json:"explanation"`
	Risk      string `json:"risk,omitempty"`
	PlanEvidence string `json:"plan_evidence,omitempty"`
}

/* IndexSuggestion suggests an index with estimated benefit */
type IndexSuggestion struct {
	Definition      string `json:"definition"`
	EstimatedBenefit string `json:"estimated_benefit"`
	PlanEvidence    string `json:"plan_evidence,omitempty"`
}

/* SafetyInfo summarizes safety checks applied */
type SafetyInfo struct {
	PolicyAllowed bool     `json:"policy_allowed"`
	BlockedReason string   `json:"blocked_reason,omitempty"`
	ChecksApplied []string `json:"checks_applied,omitempty"`
}

/* GenerateRequest is the request body for POST /v1/neuronsql/generate */
type GenerateRequest struct {
	DBDSN         string                 `json:"db_dsn"`
	DBAlias       string                 `json:"db_alias,omitempty"`
	Question      string                 `json:"question"`
	ModeSettings  map[string]interface{} `json:"mode_settings,omitempty"`
}

/* OptimizeRequest is the request body for POST /v1/neuronsql/optimize */
type OptimizeRequest struct {
	DBDSN        string                 `json:"db_dsn"`
	DBAlias      string                 `json:"db_alias,omitempty"`
	SQL          string                 `json:"sql"`
	ModeSettings map[string]interface{} `json:"mode_settings,omitempty"`
}

/* ValidateRequest is the request body for POST /v1/neuronsql/validate */
type ValidateRequest struct {
	DBDSN   string `json:"db_dsn"`
	DBAlias string `json:"db_alias,omitempty"`
	SQL     string `json:"sql"`
}

/* PLpgSQLRequest is the request body for POST /v1/neuronsql/plpgsql */
type PLpgSQLRequest struct {
	DBDSN     string `json:"db_dsn"`
	DBAlias   string `json:"db_alias,omitempty"`
	Signature string `json:"signature"` // e.g. "fn_name(integer, text) returns table(a int, b text)"
	Purpose   string `json:"purpose"`
}
