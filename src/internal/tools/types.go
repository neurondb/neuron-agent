/*-------------------------------------------------------------------------
 *
 * types.go
 *    Tool implementation for NeuronMCP
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/tools/types.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"

	"github.com/neurondb/NeuronAgent/internal/db"
)

/* ToolHandler is the interface that all tool handlers must implement */
type ToolHandler interface {
	Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error)
	Validate(args map[string]interface{}, schema map[string]interface{}) error
}

/* ExecutionResult represents the result of tool execution */
type ExecutionResult struct {
	Output string
	Error  error
}

/* ToolResultEnvelope is the structured return value for tool execution (status, data, error, metadata). */
type ToolResultEnvelope struct {
	Status   string                 `json:"status"` /* "success" | "error" */
	Data     interface{}            `json:"data,omitempty"`
	Error    *ToolResultError       `json:"error,omitempty"`
	Metadata *ToolResultMetadata   `json:"metadata,omitempty"`
}

/* ToolResultError is the error part of the envelope */
type ToolResultError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Retryable bool   `json:"retryable"`
}

/* ToolResultMetadata holds duration and idempotency key */
type ToolResultMetadata struct {
	DurationMs     int    `json:"duration_ms,omitempty"`
	IdempotencyKey string `json:"idempotency_key,omitempty"`
}
