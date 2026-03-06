/*-------------------------------------------------------------------------
 *
 * statemachine.go
 *    Deterministic agent runtime state machine controller.
 *
 * Replaces the linear Execute() pipeline with an explicit FSM: created ->
 * queued -> planning -> executing -> (awaiting_tool | executing_tool |
 * evaluating_result) -> updating_memory -> completed | failed | canceled.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/neurondb/NeuronAgent/internal/db"
)

/* Run state constants */
const (
	StateCreated          = "created"
	StateQueued           = "queued"
	StatePlanning         = "planning"
	StateExecuting        = "executing"
	StateAwaitingTool     = "awaiting_tool"
	StateExecutingTool    = "executing_tool"
	StateEvaluatingResult = "evaluating_result"
	StateUpdatingMemory   = "updating_memory"
	StateRetrying         = "retrying"
	StateBlocked          = "blocked"
	StateCompleted        = "completed"
	StateFailed           = "failed"
	StateCanceled         = "canceled"
)

/* StateMachine runs the deterministic agent loop using Runtime components. */
type StateMachine struct {
	runtime       *Runtime
	queries       *db.Queries
	traceRecorder TraceRecorder
}

/* NewStateMachine creates a state machine that uses the given runtime and its queries. */
func NewStateMachine(runtime *Runtime) *StateMachine {
	queries := runtime.GetQueries()
	return &StateMachine{
		runtime:       runtime,
		queries:       queries,
		traceRecorder: NewDBTraceRecorder(queries),
	}
}

/* Run executes the state machine loop for the given run until a terminal state. */
func (sm *StateMachine) Run(ctx context.Context, runID uuid.UUID) (*db.AgentRun, error) {
	run, err := sm.queries.GetAgentRun(ctx, runID)
	if err != nil {
		return nil, fmt.Errorf("state machine: get run: %w", err)
	}

	for !isTerminal(run.State) {
		fromState := run.State
		start := time.Now()

		switch run.State {
		case StateCreated:
			if err := sm.handleCreated(ctx, run); err != nil {
				return run, err
			}
		case StateQueued:
			if err := sm.handleQueued(ctx, run); err != nil {
				return run, err
			}
		case StatePlanning:
			if err := sm.handlePlanning(ctx, run); err != nil {
				return run, err
			}
		case StateExecuting:
			if err := sm.handleExecuting(ctx, run); err != nil {
				return run, err
			}
		case StateAwaitingTool:
			if err := sm.handleAwaitingTool(ctx, run); err != nil {
				return run, err
			}
		case StateExecutingTool:
			if err := sm.handleExecutingTool(ctx, run); err != nil {
				return run, err
			}
		case StateEvaluatingResult:
			if err := sm.handleEvaluatingResult(ctx, run); err != nil {
				return run, err
			}
		case StateUpdatingMemory:
			if err := sm.handleUpdatingMemory(ctx, run); err != nil {
				return run, err
			}
		case StateRetrying, StateBlocked:
			/* Phase 1: treat as terminal or no-op; can be extended later */
			if run.State == StateRetrying {
				sm.transition(ctx, run, StateFailed, "max_retries_exceeded")
			}
			return run, nil
		default:
			return run, fmt.Errorf("state machine: unknown state %q", run.State)
		}

		/* Record trace */
		durMs := int(time.Since(start).Milliseconds())
		_ = sm.traceRecorder.RecordTransition(ctx, run.ID, nil, fromState, run.State, "", nil, &durMs)

		/* Reload run after transition */
		run, err = sm.queries.GetAgentRun(ctx, runID)
		if err != nil {
			return run, fmt.Errorf("state machine: reload run: %w", err)
		}
	}

	return run, nil
}

func isTerminal(state string) bool {
	return state == StateCompleted || state == StateFailed || state == StateCanceled
}

/* IsTerminalRunState returns true if the run state is completed, failed, or canceled. */
func IsTerminalRunState(run *db.AgentRun) bool {
	return run != nil && isTerminal(run.State)
}

func (sm *StateMachine) transition(ctx context.Context, run *db.AgentRun, toState, trigger string) error {
	run.State = toState
	if trigger != "" {
		run.ErrorClass = &trigger
	}
	return sm.queries.UpdateAgentRun(ctx, run)
}

func (sm *StateMachine) handleCreated(ctx context.Context, run *db.AgentRun) error {
	agent, err := sm.queries.GetAgentByID(ctx, run.AgentID)
	if err != nil {
		sm.transition(ctx, run, StateFailed, "agent_not_found")
		return err
	}
	if agent.Name == "" {
		sm.transition(ctx, run, StateFailed, "agent_disabled")
		return fmt.Errorf("agent disabled or invalid")
	}
	return sm.transition(ctx, run, StateQueued, "")
}

func (sm *StateMachine) handleQueued(ctx context.Context, run *db.AgentRun) error {
	/* Phase 1: no concurrency limit; move straight to planning */
	return sm.transition(ctx, run, StatePlanning, "")
}

func (sm *StateMachine) handlePlanning(ctx context.Context, run *db.AgentRun) error {
	agent, err := sm.queries.GetAgentByID(ctx, run.AgentID)
	if err != nil {
		sm.transition(ctx, run, StateFailed, "agent_not_found")
		return err
	}

	toolsList := agent.EnabledTools
	if toolsList == nil {
		toolsList = pq.StringArray{}
	}
	planSteps, err := sm.runtime.planner.Plan(ctx, run.TaskInput, toolsList)
	if err != nil {
		/* Fallback to simple single-step plan */
		planSteps = sm.runtime.planner.simplePlan(run.TaskInput)
	}

	/* Store plan: convert []PlanStep to JSONBArray */
	stepsArr := planStepsToJSONBArray(planSteps)
	plan := &db.AgentPlan{
		RunID:     run.ID,
		Version:   1,
		Steps:     stepsArr,
		IsActive:  true,
	}
	if err := sm.queries.CreateAgentPlan(ctx, plan); err != nil {
		sm.transition(ctx, run, StateFailed, "plan_creation_failed")
		return err
	}

	run.PlanID = &plan.ID
	n := len(planSteps)
	run.TotalSteps = &n
	run.CurrentStepIndex = 0
	now := time.Now()
	run.StartedAt = &now
	return sm.transition(ctx, run, StateExecuting, "")
}

func planStepsToJSONBArray(steps []PlanStep) db.JSONBArray {
	out := make(db.JSONBArray, 0, len(steps))
	for _, s := range steps {
		m := map[string]interface{}{
			"action":   s.Action,
			"tool":     s.Tool,
			"payload":  s.Payload,
			"can_parallel": s.CanParallel,
			"retry_strategy": s.RetryStrategy,
		}
		if s.Dependencies != nil {
			m["dependencies"] = s.Dependencies
		}
		out = append(out, m)
	}
	return out
}

func (sm *StateMachine) getPlanSteps(ctx context.Context, run *db.AgentRun) ([]PlanStep, error) {
	if run.PlanID == nil {
		return nil, fmt.Errorf("no plan_id")
	}
	plan, err := sm.queries.GetAgentPlan(ctx, *run.PlanID)
	if err != nil || plan == nil {
		return nil, fmt.Errorf("plan not found")
	}
	return jsonBArrayToPlanSteps(plan.Steps)
}

func jsonBArrayToPlanSteps(arr db.JSONBArray) ([]PlanStep, error) {
	if len(arr) == 0 {
		return nil, nil
	}
	raw, _ := json.Marshal(arr)
	var steps []PlanStep
	if err := json.Unmarshal(raw, &steps); err != nil {
		return nil, err
	}
	return steps, nil
}

func (sm *StateMachine) handleExecuting(ctx context.Context, run *db.AgentRun) error {
	planSteps, err := sm.getPlanSteps(ctx, run)
	if err != nil || run.CurrentStepIndex >= len(planSteps) {
		/* No more steps: go to completed */
		run.FinalAnswer = &run.TaskInput
		run.State = StateCompleted
		now := time.Now()
		run.CompletedAt = &now
		return sm.queries.UpdateAgentRun(ctx, run)
	}

	step := planSteps[run.CurrentStepIndex]
	/* Create agent_step record */
	agentStep := &db.AgentStep{
		RunID:      run.ID,
		StepIndex:  run.CurrentStepIndex,
		PlanStepRef: intPtrRun(run.CurrentStepIndex),
		State:      "executing",
		ActionType: "tool_call",
		ActionInput: db.FromMap(step.Payload),
		RetryCount: 0,
	}
	if step.Tool == "" {
		agentStep.ActionType = "model_call"
	}
	if err := sm.queries.CreateAgentStep(ctx, agentStep); err != nil {
		return err
	}

	if step.Tool != "" {
		return sm.transition(ctx, run, StateAwaitingTool, "")
	}
	/* Model-only step: run model and then evaluating_result */
	return sm.handleModelCallStep(ctx, run, agentStep, step)
}

func intPtrRun(i int) *int { return &i }

func (sm *StateMachine) handleAwaitingTool(ctx context.Context, run *db.AgentRun) error {
	/* Validate and move to executing_tool */
	return sm.transition(ctx, run, StateExecutingTool, "")
}

func (sm *StateMachine) handleExecutingTool(ctx context.Context, run *db.AgentRun) error {
	planSteps, err := sm.getPlanSteps(ctx, run)
	if err != nil || run.CurrentStepIndex >= len(planSteps) {
		return sm.transition(ctx, run, StateFailed, "invalid_step_index")
	}
	step := planSteps[run.CurrentStepIndex]
	agent, _ := sm.queries.GetAgentByID(ctx, run.AgentID)
	if agent == nil {
		return sm.transition(ctx, run, StateFailed, "agent_not_found")
	}

	/* Get current step record */
	stepsList, _ := sm.queries.ListAgentStepsByRun(ctx, run.ID)
	var currentStep *db.AgentStep
	for i := range stepsList {
		if stepsList[i].StepIndex == run.CurrentStepIndex {
			currentStep = &stepsList[i]
			break
		}
	}
	if currentStep == nil {
		return sm.transition(ctx, run, StateFailed, "step_not_found")
	}

	toolCall := ToolCall{
		ID:        fmt.Sprintf("call_%s_%d", run.ID.String(), run.CurrentStepIndex),
		Name:      step.Tool,
		Arguments: step.Payload,
	}
	toolCtx := WithSessionID(WithAgentID(ctx, agent.ID), run.SessionID)
	start := time.Now()
	result := sm.runtime.executeSingleTool(toolCtx, agent, toolCall)
	durMs := int(time.Since(start).Milliseconds())

	/* Persist tool invocation */
	status := "success"
	inputValid := true
	outputValid := true
	var errMsg, errCode *string
	if result.Error != nil {
		status = "error"
		outputValid = false
		s := result.Error.Error()
		errMsg = &s
		errCode = strPtr("tool_error")
	}
	inv := &db.RunToolInvocation{
		RunID:        &run.ID,
		StepID:       &currentStep.ID,
		ToolName:     step.Tool,
		InputArgs:    db.FromMap(step.Payload),
		InputValid:   &inputValid,
		OutputResult: db.FromMap(map[string]interface{}{"content": result.Content}),
		OutputValid:  &outputValid,
		Status:       status,
		ErrorCode:    errCode,
		ErrorMessage: errMsg,
		DurationMs:   &durMs,
	}
	_ = sm.traceRecorder.RecordToolInvocation(ctx, inv)

	/* Update step */
	currentStep.State = "completed"
	currentStep.ActionOutput = db.FromMap(map[string]interface{}{"content": result.Content, "error": result.Error != nil})
	currentStep.DurationMs = &durMs
	now := time.Now()
	currentStep.CompletedAt = &now
	_ = sm.queries.UpdateAgentStep(ctx, currentStep)

	return sm.transition(ctx, run, StateEvaluatingResult, "")
}

func strPtr(s string) *string { return &s }

func (sm *StateMachine) handleModelCallStep(ctx context.Context, run *db.AgentRun, agentStep *db.AgentStep, step PlanStep) error {
	agent, err := sm.queries.GetAgentByID(ctx, run.AgentID)
	if err != nil || agent == nil {
		return sm.transition(ctx, run, StateFailed, "agent_not_found")
	}
	contextLoader := NewContextLoader(sm.queries, sm.runtime.memory, sm.runtime.llm)
	agentContext, err := contextLoader.LoadWithOptions(ctx, run.SessionID, agent.ID, run.TaskInput, 20, 5, false)
	if err != nil {
		agentStep.State = "failed"
		_ = sm.queries.UpdateAgentStep(ctx, agentStep)
		return sm.transition(ctx, run, StateFailed, "context_load_failed")
	}
	prompt, err := sm.runtime.prompt.Build(agent, agentContext, run.TaskInput)
	if err != nil {
		agentStep.State = "failed"
		_ = sm.queries.UpdateAgentStep(ctx, agentStep)
		return sm.transition(ctx, run, StateFailed, "prompt_build_failed")
	}
	start := time.Now()
	llmResponse, err := sm.runtime.llm.Generate(ctx, agent.ModelName, prompt, agent.Config)
	durMs := int(time.Since(start).Milliseconds())
	if err != nil {
		agentStep.State = "failed"
		_ = sm.queries.UpdateAgentStep(ctx, agentStep)
		return sm.transition(ctx, run, StateFailed, "llm_generation_failed")
	}
	/* Record model call */
	pt := EstimateTokens(prompt)
	ct := EstimateTokens(llmResponse.Content)
	total := pt + ct
	_ = sm.traceRecorder.RecordModelCall(ctx, &db.ModelCall{
		RunID:             &run.ID,
		StepID:            &agentStep.ID,
		ModelName:         agent.ModelName,
		PromptTokens:      &pt,
		CompletionTokens:  &ct,
		TotalTokens:       &total,
		LatencyMs:         &durMs,
	})
	agentStep.State = "completed"
	agentStep.ActionOutput = db.FromMap(map[string]interface{}{"content": llmResponse.Content})
	agentStep.DurationMs = &durMs
	now := time.Now()
	agentStep.CompletedAt = &now
	_ = sm.queries.UpdateAgentStep(ctx, agentStep)
	run.FinalAnswer = &llmResponse.Content
	return sm.transition(ctx, run, StateEvaluatingResult, "")
}

func (sm *StateMachine) handleEvaluatingResult(ctx context.Context, run *db.AgentRun) error {
	planSteps, _ := sm.getPlanSteps(ctx, run)
	if run.CurrentStepIndex >= len(planSteps)-1 {
		/* Last step done: complete or update memory */
		if run.FinalAnswer == nil {
			s := ""
			run.FinalAnswer = &s
		}
		run.State = StateCompleted
		now := time.Now()
		run.CompletedAt = &now
		return sm.queries.UpdateAgentRun(ctx, run)
	}
	/* More steps: advance and go to executing */
	run.CurrentStepIndex++
	return sm.transition(ctx, run, StateExecuting, "")
}

func (sm *StateMachine) handleUpdatingMemory(ctx context.Context, run *db.AgentRun) error {
	/* Phase 1: minimal memory write; then complete or next step */
	agent, _ := sm.queries.GetAgentByID(ctx, run.AgentID)
	if agent != nil && sm.runtime.hierMemory != nil && run.FinalAnswer != nil && *run.FinalAnswer != "" {
		_, _ = sm.runtime.hierMemory.StoreSTM(ctx, agent.ID, run.SessionID, *run.FinalAnswer, 0.5)
	}
	planSteps, _ := sm.getPlanSteps(ctx, run)
	if run.CurrentStepIndex >= len(planSteps)-1 {
		run.State = StateCompleted
		now := time.Now()
		run.CompletedAt = &now
		return sm.queries.UpdateAgentRun(ctx, run)
	}
	run.CurrentStepIndex++
	return sm.transition(ctx, run, StateExecuting, "")
}
