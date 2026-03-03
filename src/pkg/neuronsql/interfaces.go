/*-------------------------------------------------------------------------
 *
 * interfaces.go
 *    NeuronSQL public interfaces: LLMProvider, PolicyEngine, Retriever, Evaluator
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/pkg/neuronsql/interfaces.go
 *
 *-------------------------------------------------------------------------
 */

package neuronsql

import "context"

/* Message represents a chat message for LLM completion */
type Message struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"`
}

/* JSONSchema is an optional schema for structured LLM output (e.g. JSON mode) */
type JSONSchema struct {
	Schema map[string]interface{} `json:"schema,omitempty"`
}

/* CompletionSettings holds temperature, max_tokens, etc. */
type CompletionSettings struct {
	Temperature float64 `json:"temperature,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

/* Completion is the result of an LLM completion */
type Completion struct {
	Text    string            `json:"text"`
	Usage   map[string]int    `json:"usage,omitempty"`
	RawJSON string            `json:"raw_json,omitempty"` // best-effort extracted JSON
}

/* LLMProvider is the interface for the SQL/PLpgSQL LLM (e.g. PGLang) */
type LLMProvider interface {
	Complete(ctx context.Context, messages []Message, schema *JSONSchema, settings CompletionSettings) (*Completion, error)
	SupportsToolCalls() bool
	SupportsJSONSchema() bool
}

/* PolicyContext holds context for policy checks (request_id, db alias, etc.) */
type PolicyContext struct {
	RequestID string   `json:"request_id,omitempty"`
	DBAlias   string   `json:"db_alias,omitempty"`
	ToolName  string   `json:"tool_name,omitempty"`
	SensitiveTables []string `json:"sensitive_tables,omitempty"`
}

/* PolicyDecision is the result of a policy check */
type PolicyDecision struct {
	Allowed        bool     `json:"allowed"`
	Reason         string   `json:"reason,omitempty"`
	ReasonCode     string   `json:"reason_code,omitempty"`     /* e.g. blocked_keyword, blocked_function */
	ReasonText     string   `json:"reason_text,omitempty"`     /* human-readable reason */
	BlockedTokens  []string `json:"blocked_tokens,omitempty"` /* list of blocked keywords/functions */
	StatementClass string   `json:"statement_class,omitempty"` /* "select", "explain", "with", "blocked", etc. */
}

/* PolicyEngine checks and sanitizes SQL and tool execution */
type PolicyEngine interface {
	Check(ctx context.Context, sql string, ctxIn PolicyContext) (*PolicyDecision, error)
	Sanitize(input string) string
}

/* Document is a source document for retrieval indexing */
type Document struct {
	ID       string            `json:"id"`
	Source   string            `json:"source"`
	Path     string            `json:"path"`
	Section  string            `json:"section,omitempty"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

/* Chunk is a retrieved chunk with score and citation id */
type Chunk struct {
	ID       string  `json:"id"`
	Content  string  `json:"content"`
	Score    float64 `json:"score"`
	Source   string  `json:"source"`
	Path     string  `json:"path"`
	Section  string  `json:"section,omitempty"`
}

/* Retriever indexes and retrieves document chunks */
type Retriever interface {
	Index(ctx context.Context, docs []Document) error
	Retrieve(ctx context.Context, query string, k int) ([]Chunk, error)
}

/* EvalReport is the result of running an evaluation suite */
type EvalReport struct {
	SuiteName           string             `json:"suite_name"`
	PassRate             float64            `json:"pass_rate"`
	UnsafeRate           float64            `json:"unsafe_rate"`
	SchemaErrorRate      float64            `json:"schema_error_rate"`
	CitationCoverage     float64            `json:"citation_coverage"`
	PlanImprovementRate  float64            `json:"plan_improvement_rate,omitempty"`
	LatencyP50Ms         float64            `json:"latency_p50_ms,omitempty"`
	LatencyP95Ms         float64            `json:"latency_p95_ms,omitempty"`
	TotalTasks           int                `json:"total_tasks"`
	PassedTasks          int                `json:"passed_tasks"`
	FailedTasks          int                `json:"failed_tasks"`
	Details              []EvalTaskResult   `json:"details,omitempty"`
}

/* EvalTaskResult is the result of a single eval task */
type EvalTaskResult struct {
	TaskID   string `json:"task_id"`
	Passed   bool   `json:"passed"`
	Message  string `json:"message,omitempty"`
	LatencyMs float64 `json:"latency_ms,omitempty"`
}

/* Evaluator runs evaluation suites and returns reports */
type Evaluator interface {
	RunSuite(ctx context.Context, suiteName string) (*EvalReport, error)
}
