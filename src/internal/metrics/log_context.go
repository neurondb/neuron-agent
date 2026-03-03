/*-------------------------------------------------------------------------
 *
 * log_context.go
 *    Log context helpers for structured logging
 *
 * Provides helpers for consistent structured logging with request_id, agent_id,
 * session_id, tool_id, trace_id fields across all components.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/metrics/log_context.go
 *
 *-------------------------------------------------------------------------
 */

package metrics

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type contextKey string

const (
	requestIDKey    contextKey = "request_id"
	agentIDKey      contextKey = "agent_id"
	sessionIDKey    contextKey = "session_id"
	toolIDKey       contextKey = "tool_id"
	traceIDKey      contextKey = "trace_id"
	orgIDKey        contextKey = "org_id"
	workspaceIDKey  contextKey = "workspace_id"
	toolNameKey     contextKey = "tool_name"
	workflowIDKey   contextKey = "workflow_id"
	latencyMsKey    contextKey = "latency_ms"
	statusKey       contextKey = "status"
	costKey         contextKey = "cost"
)

/* WithLogContext adds logging fields to context (standard fields: request_id, org_id, workspace_id, agent_id, tool_name, workflow_id, latency_ms, status, cost). */
func WithLogContext(ctx context.Context, requestID, agentID, sessionID, toolID, traceID string) context.Context {
	if requestID != "" {
		ctx = context.WithValue(ctx, requestIDKey, requestID)
	}
	if agentID != "" {
		ctx = context.WithValue(ctx, agentIDKey, agentID)
	}
	if sessionID != "" {
		ctx = context.WithValue(ctx, sessionIDKey, sessionID)
	}
	if toolID != "" {
		ctx = context.WithValue(ctx, toolIDKey, toolID)
	}
	if traceID != "" {
		ctx = context.WithValue(ctx, traceIDKey, traceID)
	}
	return ctx
}

/* WithOrgIDLogContext adds org_id to log context */
func WithOrgIDLogContext(ctx context.Context, orgID string) context.Context {
	if orgID != "" {
		return context.WithValue(ctx, orgIDKey, orgID)
	}
	return ctx
}

/* WithWorkspaceIDLogContext adds workspace_id to log context */
func WithWorkspaceIDLogContext(ctx context.Context, workspaceID string) context.Context {
	if workspaceID != "" {
		return context.WithValue(ctx, workspaceIDKey, workspaceID)
	}
	return ctx
}

/* WithToolNameLogContext adds tool_name to log context */
func WithToolNameLogContext(ctx context.Context, toolName string) context.Context {
	if toolName != "" {
		return context.WithValue(ctx, toolNameKey, toolName)
	}
	return ctx
}

/* WithWorkflowIDLogContext adds workflow_id to log context */
func WithWorkflowIDLogContext(ctx context.Context, workflowID string) context.Context {
	if workflowID != "" {
		return context.WithValue(ctx, workflowIDKey, workflowID)
	}
	return ctx
}

/* WithLatencyMsLogContext adds latency_ms to log context */
func WithLatencyMsLogContext(ctx context.Context, latencyMs int64) context.Context {
	return context.WithValue(ctx, latencyMsKey, latencyMs)
}

/* WithStatusLogContext adds status to log context */
func WithStatusLogContext(ctx context.Context, status string) context.Context {
	if status != "" {
		return context.WithValue(ctx, statusKey, status)
	}
	return ctx
}

/* WithCostLogContext adds cost to log context */
func WithCostLogContext(ctx context.Context, cost float64) context.Context {
	return context.WithValue(ctx, costKey, cost)
}

/* WithAgentIDLogContext adds agent ID to log context */
func WithAgentIDLogContext(ctx context.Context, agentID uuid.UUID) context.Context {
	return context.WithValue(ctx, agentIDKey, agentID.String())
}

/* WithSessionIDLogContext adds session ID to log context */
func WithSessionIDLogContext(ctx context.Context, sessionID uuid.UUID) context.Context {
	return context.WithValue(ctx, sessionIDKey, sessionID.String())
}

/* WithToolIDLogContext adds tool ID to log context */
func WithToolIDLogContext(ctx context.Context, toolID string) context.Context {
	return context.WithValue(ctx, toolIDKey, toolID)
}

/* WithTraceIDLogContext adds trace ID to log context */
func WithTraceIDLogContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey, traceID)
}

/* GetRequestIDFromContext gets request ID from context */
func GetRequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(requestIDKey).(string); ok {
		return id
	}
	return ""
}

/* GetAgentIDFromContext gets agent ID from context */
func GetAgentIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(agentIDKey).(string); ok {
		return id
	}
	if id, ok := ctx.Value(agentIDKey).(uuid.UUID); ok {
		return id.String()
	}
	return ""
}

/* GetSessionIDFromContext gets session ID from context */
func GetSessionIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(sessionIDKey).(string); ok {
		return id
	}
	if id, ok := ctx.Value(sessionIDKey).(uuid.UUID); ok {
		return id.String()
	}
	return ""
}

/* GetToolIDFromContext gets tool ID from context */
func GetToolIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(toolIDKey).(string); ok {
		return id
	}
	return ""
}

/* GetTraceIDFromContext gets trace ID from context */
func GetTraceIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(traceIDKey).(string); ok {
		return id
	}
	return ""
}

/* LoggerFromContext creates a zerolog logger with fields from context */
func LoggerFromContext(ctx context.Context) zerolog.Logger {
	logger := *zerolog.Ctx(ctx)
	if logger.GetLevel() == zerolog.Disabled {
		logger = zerolog.Nop()
	}

	/* Add standard context fields (request_id, org_id, workspace_id, agent_id, tool_name, workflow_id, latency_ms, status, cost) */
	requestID := GetRequestIDFromContext(ctx)
	agentID := GetAgentIDFromContext(ctx)
	sessionID := GetSessionIDFromContext(ctx)
	toolID := GetToolIDFromContext(ctx)
	traceID := GetTraceIDFromContext(ctx)
	orgID := getStringFromContext(ctx, orgIDKey)
	workspaceID := getStringFromContext(ctx, workspaceIDKey)
	toolName := getStringFromContext(ctx, toolNameKey)
	workflowID := getStringFromContext(ctx, workflowIDKey)
	status := getStringFromContext(ctx, statusKey)
	latencyMs := getInt64FromContext(ctx, latencyMsKey)
	cost := getFloat64FromContext(ctx, costKey)

	if requestID != "" {
		logger = logger.With().Str("request_id", requestID).Logger()
	}
	if orgID != "" {
		logger = logger.With().Str("org_id", orgID).Logger()
	}
	if workspaceID != "" {
		logger = logger.With().Str("workspace_id", workspaceID).Logger()
	}
	if agentID != "" {
		logger = logger.With().Str("agent_id", agentID).Logger()
	}
	if sessionID != "" {
		logger = logger.With().Str("session_id", sessionID).Logger()
	}
	if toolID != "" {
		logger = logger.With().Str("tool_id", toolID).Logger()
	}
	if toolName != "" {
		logger = logger.With().Str("tool_name", toolName).Logger()
	}
	if workflowID != "" {
		logger = logger.With().Str("workflow_id", workflowID).Logger()
	}
	if traceID != "" {
		logger = logger.With().Str("trace_id", traceID).Logger()
	}
	if status != "" {
		logger = logger.With().Str("status", status).Logger()
	}
	if latencyMs != 0 {
		logger = logger.With().Int64("latency_ms", latencyMs).Logger()
	}
	if cost != 0 {
		logger = logger.With().Float64("cost", cost).Logger()
	}

	return logger
}

func getStringFromContext(ctx context.Context, key contextKey) string {
	if v, ok := ctx.Value(key).(string); ok {
		return v
	}
	return ""
}

func getInt64FromContext(ctx context.Context, key contextKey) int64 {
	if v, ok := ctx.Value(key).(int64); ok {
		return v
	}
	return 0
}

func getFloat64FromContext(ctx context.Context, key contextKey) float64 {
	if v, ok := ctx.Value(key).(float64); ok {
		return v
	}
	return 0
}

/* LogWithContext logs a message with context fields */
func LogWithContext(ctx context.Context, level zerolog.Level, message string, fields map[string]interface{}) {
	logger := LoggerFromContext(ctx)
	event := logger.WithLevel(level)

	for key, value := range fields {
		event = event.Interface(key, value)
	}

	event.Msg(message)
}

/* DebugWithContext logs a debug message with context */
func DebugWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	LogWithContext(ctx, zerolog.DebugLevel, message, fields)
}

/* InfoWithContext logs an info message with context */
func InfoWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	LogWithContext(ctx, zerolog.InfoLevel, message, fields)
}

/* WarnWithContext logs a warning message with context */
func WarnWithContext(ctx context.Context, message string, fields map[string]interface{}) {
	LogWithContext(ctx, zerolog.WarnLevel, message, fields)
}

/* ErrorWithContext logs an error message with context */
func ErrorWithContext(ctx context.Context, message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	if err != nil {
		fields["error"] = err.Error()
	}
	LogWithContext(ctx, zerolog.ErrorLevel, message, fields)
}
