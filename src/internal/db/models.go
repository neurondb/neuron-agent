/*-------------------------------------------------------------------------
 *
 * models.go
 *    Database models for NeuronAgent
 *
 * Defines data structures for agents, sessions, messages, memory chunks,
 * tools, jobs, and API keys.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/db/models.go
 *
 *-------------------------------------------------------------------------
 */

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Agent struct {
	ID           uuid.UUID      `db:"id"`
	OrgID        *uuid.UUID     `db:"org_id"`
	Name         string         `db:"name"`
	Description  *string        `db:"description"`
	SystemPrompt string         `db:"system_prompt"`
	ModelName    string         `db:"model_name"`
	MemoryTable  *string        `db:"memory_table"`
	EnabledTools pq.StringArray `db:"enabled_tools"`
	Config       JSONBMap       `db:"config"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type Session struct {
	ID             uuid.UUID  `db:"id"`
	OrgID          *uuid.UUID `db:"org_id"`
	AgentID        uuid.UUID  `db:"agent_id"`
	ExternalUserID *string    `db:"external_user_id"`
	Metadata       JSONBMap   `db:"metadata"`
	CreatedAt      time.Time  `db:"created_at"`
	LastActivityAt time.Time  `db:"last_activity_at"`
}

type Message struct {
	ID         int64                  `db:"id"`
	SessionID  uuid.UUID              `db:"session_id"`
	Role       string                 `db:"role"`
	Content    string                 `db:"content"`
	ToolName   *string                `db:"tool_name"`
	ToolCallID *string                `db:"tool_call_id"`
	TokenCount *int                   `db:"token_count"`
	Metadata   map[string]interface{} `db:"metadata"`
	CreatedAt  time.Time              `db:"created_at"`
}

type MemoryChunk struct {
	ID              int64      `db:"id"`
	AgentID         uuid.UUID  `db:"agent_id"`
	SessionID       *uuid.UUID `db:"session_id"`
	MessageID       *int64     `db:"message_id"`
	Content         string     `db:"content"`
	Embedding       []float32  `db:"embedding"`
	ImportanceScore float64    `db:"importance_score"`
	Metadata        JSONBMap   `db:"metadata"`
	CreatedAt       time.Time  `db:"created_at"`
}

/* MemoryChunkWithSimilarity includes similarity score from vector search */
type MemoryChunkWithSimilarity struct {
	MemoryChunk
	Similarity float64 `db:"similarity"`
}

type Tool struct {
	Name          string    `db:"name"`
	Description   string    `db:"description"`
	ArgSchema     JSONBMap  `db:"arg_schema"`
	ResultSchema  JSONBMap  `db:"result_schema"`
	HandlerType   string    `db:"handler_type"`
	HandlerConfig JSONBMap  `db:"handler_config"`
	Enabled       bool      `db:"enabled"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type Job struct {
	ID           int64      `db:"id"`
	AgentID      *uuid.UUID `db:"agent_id"`
	SessionID    *uuid.UUID `db:"session_id"`
	Type         string     `db:"type"`
	Status       string     `db:"status"`
	Priority     int        `db:"priority"`
	Payload      JSONBMap   `db:"payload"`
	Result       JSONBMap   `db:"result"`
	ErrorMessage *string    `db:"error_message"`
	RetryCount   int        `db:"retry_count"`
	MaxRetries   int        `db:"max_retries"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	StartedAt    *time.Time `db:"started_at"`
	CompletedAt  *time.Time `db:"completed_at"`
}

type APIKey struct {
	ID              uuid.UUID      `db:"id"`
	KeyHash         string         `db:"key_hash"`
	KeyPrefix       string         `db:"key_prefix"`
	OrganizationID  *string        `db:"organization_id"`
	UserID          *string        `db:"user_id"`
	PrincipalID     *uuid.UUID     `db:"principal_id"`
	RateLimitPerMin int            `db:"rate_limit_per_minute"`
	Roles           pq.StringArray `db:"roles"`
	Metadata        JSONBMap       `db:"metadata"`
	CreatedAt       time.Time      `db:"created_at"`
	LastUsedAt      *time.Time     `db:"last_used_at"`
	ExpiresAt       *time.Time     `db:"expires_at"`
}

type Principal struct {
	ID        uuid.UUID `db:"id"`
	Type      string    `db:"type"` // 'user', 'org', 'agent', 'tool', 'dataset'
	Name      string    `db:"name"`
	Metadata  JSONBMap  `db:"metadata"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Policy struct {
	ID           uuid.UUID      `db:"id"`
	PrincipalID  uuid.UUID      `db:"principal_id"`
	ResourceType string         `db:"resource_type"`
	ResourceID   *string        `db:"resource_id"`
	Permissions  pq.StringArray `db:"permissions"`
	Conditions   JSONBMap       `db:"conditions"`
	CreatedAt    time.Time      `db:"created_at"`
	UpdatedAt    time.Time      `db:"updated_at"`
}

type ToolPermission struct {
	ID         uuid.UUID `db:"id"`
	AgentID    uuid.UUID `db:"agent_id"`
	ToolName   string    `db:"tool_name"`
	Allowed    bool      `db:"allowed"`
	Conditions JSONBMap  `db:"conditions"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}

type SessionToolPermission struct {
	ID         uuid.UUID `db:"id"`
	SessionID  uuid.UUID `db:"session_id"`
	ToolName   string    `db:"tool_name"`
	Allowed    bool      `db:"allowed"`
	Conditions JSONBMap  `db:"conditions"`
	CreatedAt  time.Time `db:"created_at"`
}

/* PrincipalToolPermission is principal-level tool allow/deny (RBAC). */
type PrincipalToolPermission struct {
	PrincipalID uuid.UUID `db:"principal_id"`
	ToolName    string    `db:"tool_name"`
	Allowed     bool      `db:"allowed"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

/* WorkflowPermission is principal workflow role (RBAC). */
type WorkflowPermission struct {
	ID          uuid.UUID `db:"id"`
	PrincipalID uuid.UUID `db:"principal_id"`
	WorkflowID  uuid.UUID `db:"workflow_id"`
	Role        string    `db:"role"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

/* WorkspacePolicy is workspace-level policy config (RBAC/ABAC). */
type WorkspacePolicy struct {
	ID          uuid.UUID `db:"id"`
	WorkspaceID uuid.UUID `db:"workspace_id"`
	PolicyType  string    `db:"policy_type"`
	Config      JSONBMap  `db:"config"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type DataPermission struct {
	ID          uuid.UUID      `db:"id"`
	PrincipalID uuid.UUID      `db:"principal_id"`
	SchemaName  *string        `db:"schema_name"`
	TableName   *string        `db:"table_name"`
	RowFilter   *string        `db:"row_filter"`
	ColumnMask  JSONBMap       `db:"column_mask"`
	Permissions pq.StringArray `db:"permissions"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type AuditLog struct {
	ID           int64      `db:"id"`
	Timestamp    time.Time  `db:"timestamp"`
	PrincipalID  *uuid.UUID `db:"principal_id"`
	APIKeyID     *uuid.UUID `db:"api_key_id"`
	AgentID      *uuid.UUID `db:"agent_id"`
	SessionID    *uuid.UUID `db:"session_id"`
	Action       string     `db:"action"`
	ResourceType string     `db:"resource_type"`
	ResourceID   *string    `db:"resource_id"`
	InputsHash   *string    `db:"inputs_hash"`
	OutputsHash  *string    `db:"outputs_hash"`
	Inputs       JSONBMap   `db:"inputs"`
	Outputs      JSONBMap   `db:"outputs"`
	Metadata     JSONBMap   `db:"metadata"`
	CreatedAt    time.Time  `db:"created_at"`
}

type ExecutionSnapshot struct {
	ID                uuid.UUID `db:"id"`
	SessionID         uuid.UUID `db:"session_id"`
	AgentID           uuid.UUID `db:"agent_id"`
	UserMessage       string    `db:"user_message"`
	ExecutionState    JSONBMap  `db:"execution_state"`
	DeterministicMode bool      `db:"deterministic_mode"`
	CreatedAt         time.Time `db:"created_at"`
}

type AgentSpecialization struct {
	ID                 uuid.UUID  `db:"id"`
	AgentID            uuid.UUID  `db:"agent_id"`
	SpecializationType string     `db:"specialization_type"`
	Capabilities       pq.StringArray `db:"capabilities"`
	Config             JSONBMap   `db:"config"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

/* AgentRun is the first-class record of a complete agent execution (state machine). */
type AgentRun struct {
	ID               uuid.UUID   `db:"id"`
	AgentID          uuid.UUID   `db:"agent_id"`
	SessionID        uuid.UUID   `db:"session_id"`
	TaskInput        string      `db:"task_input"`
	TaskMetadata     JSONBMap    `db:"task_metadata"`
	State            string      `db:"state"`
	PlanID           *uuid.UUID  `db:"plan_id"`
	CurrentStepIndex int         `db:"current_step_index"`
	TotalSteps       *int        `db:"total_steps"`
	RetryCount       int         `db:"retry_count"`
	FinalAnswer      *string     `db:"final_answer"`
	ErrorClass       *string     `db:"error_class"`
	ErrorDetail      JSONBMap    `db:"error_detail"`
	TokensUsed       JSONBMap    `db:"tokens_used"`
	CostEstimate     *float64    `db:"cost_estimate"`
	StartedAt        *time.Time  `db:"started_at"`
	CompletedAt      *time.Time  `db:"completed_at"`
	CreatedAt        time.Time   `db:"created_at"`
	UpdatedAt        time.Time   `db:"updated_at"`
	OrgID            *uuid.UUID  `db:"org_id"`
	Checkpoint       JSONBMap    `db:"checkpoint"`
}

/* AgentPlan holds a structured plan for a run. */
type AgentPlan struct {
	ID        uuid.UUID  `db:"id"`
	RunID     uuid.UUID  `db:"run_id"`
	Version   int        `db:"version"`
	Steps     JSONBArray `db:"steps"`
	Reasoning *string    `db:"reasoning"`
	IsActive  bool       `db:"is_active"`
	CreatedAt time.Time  `db:"created_at"`
}

/* AgentStep is one step execution within a run. */
type AgentStep struct {
	ID           uuid.UUID  `db:"id"`
	RunID        uuid.UUID  `db:"run_id"`
	StepIndex    int        `db:"step_index"`
	PlanStepRef  *int       `db:"plan_step_ref"`
	State        string     `db:"state"`
	ActionType   string     `db:"action_type"`
	ActionInput  JSONBMap   `db:"action_input"`
	ActionOutput JSONBMap   `db:"action_output"`
	Evaluation   JSONBMap   `db:"evaluation"`
	DurationMs   *int       `db:"duration_ms"`
	RetryCount   int        `db:"retry_count"`
	CreatedAt    time.Time  `db:"created_at"`
	CompletedAt  *time.Time `db:"completed_at"`
}

/* RunToolInvocation records a single tool call within a run/step. */
type RunToolInvocation struct {
	ID             uuid.UUID   `db:"id"`
	RunID          *uuid.UUID  `db:"run_id"`
	StepID         *uuid.UUID  `db:"step_id"`
	ToolName       string      `db:"tool_name"`
	ToolVersion    *int        `db:"tool_version"`
	InputArgs      JSONBMap    `db:"input_args"`
	InputValid     *bool       `db:"input_valid"`
	OutputResult   JSONBMap    `db:"output_result"`
	OutputValid    *bool       `db:"output_valid"`
	Status         string      `db:"status"`
	ErrorCode      *string     `db:"error_code"`
	ErrorMessage   *string     `db:"error_message"`
	Retryable      *bool       `db:"retryable"`
	IdempotencyKey *string     `db:"idempotency_key"`
	DurationMs     *int        `db:"duration_ms"`
	CreatedAt      time.Time   `db:"created_at"`
}

/* ModelCall records one LLM invocation. */
type ModelCall struct {
	ID                uuid.UUID  `db:"id"`
	RunID             *uuid.UUID `db:"run_id"`
	StepID            *uuid.UUID `db:"step_id"`
	ModelName         string     `db:"model_name"`
	ModelProvider     *string    `db:"model_provider"`
	PromptHash        *string    `db:"prompt_hash"`
	PromptSections    JSONBMap   `db:"prompt_sections"`
	PromptTokens      *int       `db:"prompt_tokens"`
	CompletionTokens  *int       `db:"completion_tokens"`
	TotalTokens       *int       `db:"total_tokens"`
	CostEstimate      *float64   `db:"cost_estimate"`
	LatencyMs         *int       `db:"latency_ms"`
	FinishReason      *string    `db:"finish_reason"`
	RoutingReason     *string    `db:"routing_reason"`
	CreatedAt         time.Time  `db:"created_at"`
}

/* ExecutionTrace is one state transition in the runtime FSM. */
type ExecutionTrace struct {
	ID         int64      `db:"id"`
	RunID      uuid.UUID  `db:"run_id"`
	StepID     *uuid.UUID `db:"step_id"`
	FromState  *string    `db:"from_state"`
	ToState    string     `db:"to_state"`
	Trigger    *string    `db:"trigger"`
	Metadata   JSONBMap   `db:"metadata"`
	DurationMs *int       `db:"duration_ms"`
	CreatedAt  time.Time  `db:"created_at"`
}
