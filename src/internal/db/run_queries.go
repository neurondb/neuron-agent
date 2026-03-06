/*-------------------------------------------------------------------------
 *
 * run_queries.go
 *    Queries for agent runtime state machine: agent_runs, agent_plans,
 *    agent_steps, run_tool_invocations, model_calls, execution_traces.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

const (
	createAgentRunQuery = `
		INSERT INTO neurondb_agent.agent_runs
		(agent_id, session_id, task_input, task_metadata, state, org_id)
		VALUES ($1, $2, $3, $4::jsonb, $5, $6)
		RETURNING id, created_at, updated_at`
	getAgentRunQuery    = `SELECT * FROM neurondb_agent.agent_runs WHERE id = $1`
	updateAgentRunQuery = `
		UPDATE neurondb_agent.agent_runs SET
			state = $2, plan_id = $3, current_step_index = $4, total_steps = $5,
			retry_count = $6, final_answer = $7, error_class = $8, error_detail = $9::jsonb,
			tokens_used = $10::jsonb, cost_estimate = $11, started_at = $12, completed_at = $13,
			updated_at = NOW(), checkpoint = $14::jsonb
		WHERE id = $1
		RETURNING *`
	transitionAgentRunStateQuery = `
		UPDATE neurondb_agent.agent_runs SET state = $2, updated_at = NOW() WHERE id = $1 RETURNING id`

	createAgentPlanQuery = `
		INSERT INTO neurondb_agent.agent_plans (run_id, version, steps, reasoning, is_active)
		VALUES ($1, $2, $3::jsonb, $4, $5)
		RETURNING id, created_at`
	getAgentPlanQuery = `SELECT * FROM neurondb_agent.agent_plans WHERE id = $1`
	getAgentPlanByRunQuery = `SELECT * FROM neurondb_agent.agent_plans WHERE run_id = $1 AND is_active = true ORDER BY version DESC LIMIT 1`

	createAgentStepQuery = `
		INSERT INTO neurondb_agent.agent_steps
		(run_id, step_index, plan_step_ref, state, action_type, action_input, retry_count)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7)
		RETURNING id, created_at`
	getAgentStepQuery       = `SELECT * FROM neurondb_agent.agent_steps WHERE id = $1`
	listAgentStepsByRunQuery = `SELECT * FROM neurondb_agent.agent_steps WHERE run_id = $1 ORDER BY step_index ASC`
	updateAgentStepQuery = `
		UPDATE neurondb_agent.agent_steps SET
			state = $2, action_output = $3::jsonb, evaluation = $4::jsonb,
			duration_ms = $5, retry_count = $6, completed_at = $7
		WHERE id = $1
		RETURNING *`

	createRunToolInvocationQuery = `
		INSERT INTO neurondb_agent.run_tool_invocations
		(run_id, step_id, tool_name, tool_version, input_args, input_valid, output_result, output_valid, status, error_code, error_message, retryable, idempotency_key, duration_ms)
		VALUES ($1, $2, $3, $4, $5::jsonb, $6, $7::jsonb, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id, created_at`
	getRunToolInvocationByIdempotencyQuery = `SELECT * FROM neurondb_agent.run_tool_invocations WHERE idempotency_key = $1 AND status = 'success' LIMIT 1`

	createModelCallQuery = `
		INSERT INTO neurondb_agent.model_calls
		(run_id, step_id, model_name, model_provider, prompt_hash, prompt_sections, prompt_tokens, completion_tokens, total_tokens, cost_estimate, latency_ms, finish_reason, routing_reason)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at`

	createExecutionTraceQuery = `
		INSERT INTO neurondb_agent.execution_traces (run_id, step_id, from_state, to_state, trigger, metadata, duration_ms)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb, $7)
		RETURNING id, created_at`
	listExecutionTracesByRunQuery = `SELECT * FROM neurondb_agent.execution_traces WHERE run_id = $1 ORDER BY created_at ASC`
)

func (q *Queries) CreateAgentRun(ctx context.Context, run *AgentRun) error {
	params := []interface{}{run.AgentID, run.SessionID, run.TaskInput, run.TaskMetadata, run.State, run.OrgID}
	err := q.DB.GetContext(ctx, run, createAgentRunQuery, params...)
	if err != nil {
		return q.formatQueryError("INSERT", createAgentRunQuery, len(params), "neurondb_agent.agent_runs", err)
	}
	return nil
}

func (q *Queries) GetAgentRun(ctx context.Context, id uuid.UUID) (*AgentRun, error) {
	var run AgentRun
	err := q.DB.GetContext(ctx, &run, getAgentRunQuery, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent run not found: id=%s: %w", id.String(), err)
	}
	if err != nil {
		return nil, q.formatQueryError("SELECT", getAgentRunQuery, 1, "neurondb_agent.agent_runs", err)
	}
	return &run, nil
}

func (q *Queries) UpdateAgentRun(ctx context.Context, run *AgentRun) error {
	params := []interface{}{
		run.ID, run.State, run.PlanID, run.CurrentStepIndex, run.TotalSteps,
		run.RetryCount, run.FinalAnswer, run.ErrorClass, run.ErrorDetail,
		run.TokensUsed, run.CostEstimate, run.StartedAt, run.CompletedAt, run.Checkpoint,
	}
	err := q.DB.GetContext(ctx, run, updateAgentRunQuery, params...)
	if err != nil {
		return q.formatQueryError("UPDATE", updateAgentRunQuery, len(params), "neurondb_agent.agent_runs", err)
	}
	return nil
}

func (q *Queries) TransitionAgentRunState(ctx context.Context, runID uuid.UUID, state string) error {
	_, err := q.DB.ExecContext(ctx, transitionAgentRunStateQuery, runID, state)
	return err
}

func (q *Queries) CreateAgentPlan(ctx context.Context, plan *AgentPlan) error {
	params := []interface{}{plan.RunID, plan.Version, plan.Steps, plan.Reasoning, plan.IsActive}
	err := q.DB.GetContext(ctx, plan, createAgentPlanQuery, params...)
	if err != nil {
		return q.formatQueryError("INSERT", createAgentPlanQuery, len(params), "neurondb_agent.agent_plans", err)
	}
	return nil
}

func (q *Queries) GetAgentPlan(ctx context.Context, id uuid.UUID) (*AgentPlan, error) {
	var plan AgentPlan
	err := q.DB.GetContext(ctx, &plan, getAgentPlanQuery, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent plan not found: id=%s: %w", id.String(), err)
	}
	if err != nil {
		return nil, q.formatQueryError("SELECT", getAgentPlanQuery, 1, "neurondb_agent.agent_plans", err)
	}
	return &plan, nil
}

func (q *Queries) GetAgentPlanByRun(ctx context.Context, runID uuid.UUID) (*AgentPlan, error) {
	var plan AgentPlan
	err := q.DB.GetContext(ctx, &plan, getAgentPlanByRunQuery, runID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, q.formatQueryError("SELECT", getAgentPlanByRunQuery, 1, "neurondb_agent.agent_plans", err)
	}
	return &plan, nil
}

func (q *Queries) CreateAgentStep(ctx context.Context, step *AgentStep) error {
	params := []interface{}{step.RunID, step.StepIndex, step.PlanStepRef, step.State, step.ActionType, step.ActionInput, step.RetryCount}
	err := q.DB.GetContext(ctx, step, createAgentStepQuery, params...)
	if err != nil {
		return q.formatQueryError("INSERT", createAgentStepQuery, len(params), "neurondb_agent.agent_steps", err)
	}
	return nil
}

func (q *Queries) GetAgentStep(ctx context.Context, id uuid.UUID) (*AgentStep, error) {
	var step AgentStep
	err := q.DB.GetContext(ctx, &step, getAgentStepQuery, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("agent step not found: id=%s: %w", id.String(), err)
	}
	if err != nil {
		return nil, q.formatQueryError("SELECT", getAgentStepQuery, 1, "neurondb_agent.agent_steps", err)
	}
	return &step, nil
}

func (q *Queries) ListAgentStepsByRun(ctx context.Context, runID uuid.UUID) ([]AgentStep, error) {
	var steps []AgentStep
	err := q.DB.SelectContext(ctx, &steps, listAgentStepsByRunQuery, runID)
	if err != nil {
		return nil, q.formatQueryError("SELECT", listAgentStepsByRunQuery, 1, "neurondb_agent.agent_steps", err)
	}
	return steps, nil
}

func (q *Queries) UpdateAgentStep(ctx context.Context, step *AgentStep) error {
	params := []interface{}{step.ID, step.State, step.ActionOutput, step.Evaluation, step.DurationMs, step.RetryCount, step.CompletedAt}
	err := q.DB.GetContext(ctx, step, updateAgentStepQuery, params...)
	if err != nil {
		return q.formatQueryError("UPDATE", updateAgentStepQuery, len(params), "neurondb_agent.agent_steps", err)
	}
	return nil
}

func (q *Queries) CreateRunToolInvocation(ctx context.Context, inv *RunToolInvocation) error {
	params := []interface{}{
		inv.RunID, inv.StepID, inv.ToolName, inv.ToolVersion, inv.InputArgs, inv.InputValid,
		inv.OutputResult, inv.OutputValid, inv.Status, inv.ErrorCode, inv.ErrorMessage,
		inv.Retryable, inv.IdempotencyKey, inv.DurationMs,
	}
	err := q.DB.GetContext(ctx, inv, createRunToolInvocationQuery, params...)
	if err != nil {
		return q.formatQueryError("INSERT", createRunToolInvocationQuery, len(params), "neurondb_agent.run_tool_invocations", err)
	}
	return nil
}

func (q *Queries) GetRunToolInvocationByIdempotency(ctx context.Context, key string) (*RunToolInvocation, error) {
	var inv RunToolInvocation
	err := q.DB.GetContext(ctx, &inv, getRunToolInvocationByIdempotencyQuery, key)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &inv, nil
}

func (q *Queries) CreateModelCall(ctx context.Context, call *ModelCall) error {
	params := []interface{}{
		call.RunID, call.StepID, call.ModelName, call.ModelProvider, call.PromptHash, call.PromptSections,
		call.PromptTokens, call.CompletionTokens, call.TotalTokens, call.CostEstimate, call.LatencyMs,
		call.FinishReason, call.RoutingReason,
	}
	err := q.DB.GetContext(ctx, call, createModelCallQuery, params...)
	if err != nil {
		return q.formatQueryError("INSERT", createModelCallQuery, len(params), "neurondb_agent.model_calls", err)
	}
	return nil
}

func (q *Queries) CreateExecutionTrace(ctx context.Context, trace *ExecutionTrace) error {
	params := []interface{}{trace.RunID, trace.StepID, trace.FromState, trace.ToState, trace.Trigger, trace.Metadata, trace.DurationMs}
	err := q.DB.GetContext(ctx, trace, createExecutionTraceQuery, params...)
	if err != nil {
		return q.formatQueryError("INSERT", createExecutionTraceQuery, len(params), "neurondb_agent.execution_traces", err)
	}
	return nil
}

func (q *Queries) ListExecutionTracesByRun(ctx context.Context, runID uuid.UUID) ([]ExecutionTrace, error) {
	var traces []ExecutionTrace
	err := q.DB.SelectContext(ctx, &traces, listExecutionTracesByRunQuery, runID)
	if err != nil {
		return nil, q.formatQueryError("SELECT", listExecutionTracesByRunQuery, 1, "neurondb_agent.execution_traces", err)
	}
	return traces, nil
}
