/*-------------------------------------------------------------------------
 *
 * trace.go
 *    TraceRecorder interface and DB implementation for execution_traces,
 *    model_calls, and tool_invocations.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

package agent

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurondb/NeuronAgent/internal/db"
)

/* TraceRecorder records runtime events for observability. */
type TraceRecorder interface {
	RecordTransition(ctx context.Context, runID uuid.UUID, stepID *uuid.UUID, fromState, toState, trigger string, metadata map[string]interface{}, durationMs *int) error
	RecordToolInvocation(ctx context.Context, inv *db.RunToolInvocation) error
	RecordModelCall(ctx context.Context, call *db.ModelCall) error
}

/* DBTraceRecorder persists traces to NeuronDB. */
type DBTraceRecorder struct {
	queries *db.Queries
}

/* NewDBTraceRecorder creates a trace recorder that writes to the database. */
func NewDBTraceRecorder(queries *db.Queries) *DBTraceRecorder {
	return &DBTraceRecorder{queries: queries}
}

/* RecordTransition records a state transition. */
func (t *DBTraceRecorder) RecordTransition(ctx context.Context, runID uuid.UUID, stepID *uuid.UUID, fromState, toState, trigger string, metadata map[string]interface{}, durationMs *int) error {
	trace := &db.ExecutionTrace{
		RunID:      runID,
		StepID:     stepID,
		FromState:  &fromState,
		ToState:    toState,
		Trigger:    &trigger,
		Metadata:   db.FromMap(metadata),
		DurationMs: durationMs,
	}
	return t.queries.CreateExecutionTrace(ctx, trace)
}

/* RecordToolInvocation records a tool call. */
func (t *DBTraceRecorder) RecordToolInvocation(ctx context.Context, inv *db.RunToolInvocation) error {
	return t.queries.CreateRunToolInvocation(ctx, inv)
}

/* RecordModelCall records an LLM invocation. */
func (t *DBTraceRecorder) RecordModelCall(ctx context.Context, call *db.ModelCall) error {
	return t.queries.CreateModelCall(ctx, call)
}
