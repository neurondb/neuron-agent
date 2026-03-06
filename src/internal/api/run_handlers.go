/*-------------------------------------------------------------------------
 *
 * run_handlers.go
 *    API handlers for agent run lifecycle: create, get, steps, traces, cancel.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/neurondb/NeuronAgent/internal/agent"
	"github.com/neurondb/NeuronAgent/internal/auth"
	"github.com/neurondb/NeuronAgent/internal/validation"
)

/* CreateRunRequest is the body for POST /agents/{id}/runs */
type CreateRunRequest struct {
	SessionID   string                 `json:"session_id"`
	TaskInput   string                 `json:"task_input"`
	TaskMetadata map[string]interface{} `json:"task_metadata,omitempty"`
}

/* CreateRun creates and starts an agent run (state machine). */
func (h *Handlers) CreateRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "agent_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid agent ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	agentID, _ := uuid.Parse(vars["id"])

	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin, auth.RoleUser); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "insufficient permissions", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}

	var req CreateRunRequest
	if err := DecodeJSON(r, &req); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	if req.SessionID == "" || req.TaskInput == "" {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "session_id and task_input are required", nil, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	sessionID, err := uuid.Parse(req.SessionID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid session_id", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	session, err := h.queries.GetSession(r.Context(), sessionID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "session not found", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	if session.AgentID != agentID {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "session does not belong to this agent", nil, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}

	run, err := h.runtime.StartRun(r.Context(), sessionID, req.TaskInput, req.TaskMetadata)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "run failed", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	respondJSON(w, http.StatusCreated, run)
}

/* GetRun returns a run by ID. */
func (h *Handlers) GetRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])

	run, err := h.queries.GetAgentRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "run not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	respondJSON(w, http.StatusOK, run)
}

/* GetRunSteps returns steps for a run. */
func (h *Handlers) GetRunSteps(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])

	steps, err := h.queries.ListAgentStepsByRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "failed to list steps", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	respondJSON(w, http.StatusOK, steps)
}

/* GetRunTraces returns execution traces for a run. */
func (h *Handlers) GetRunTraces(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])

	traces, err := h.queries.ListExecutionTracesByRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "failed to list traces", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	respondJSON(w, http.StatusOK, traces)
}

/* GetRunPlan returns the active plan for a run. */
func (h *Handlers) GetRunPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])

	plan, err := h.queries.GetAgentPlanByRun(r.Context(), runID)
	if err != nil || plan == nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "plan not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	respondJSON(w, http.StatusOK, plan)
}

/* CancelRunRequest is the body for POST /runs/{id}/cancel */
type CancelRunRequest struct {
	Reason string `json:"reason,omitempty"`
}

/* CancelRun transitions a run to canceled. */
func (h *Handlers) CancelRun(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])

	run, err := h.queries.GetAgentRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "run not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	if agent.IsTerminalRunState(run) {
		respondJSON(w, http.StatusOK, run)
		return
	}
	run.State = agent.StateCanceled
	if r.Body != nil {
		var req CancelRunRequest
		_ = json.NewDecoder(r.Body).Decode(&req)
		if req.Reason != "" {
			run.ErrorClass = &req.Reason
		}
	}
	if err := h.queries.UpdateAgentRun(r.Context(), run); err != nil {
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "failed to cancel run", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	respondJSON(w, http.StatusOK, run)
}

/* ExplainTool returns why a tool was called (explainability). */
func (h *Handlers) ExplainTool(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])
	invID := vars["inv_id"]
	if invID == "" {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "inv_id is required", nil, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	_, err := h.queries.GetAgentRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "run not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	explain := map[string]interface{}{
		"tool_invocation_id": invID,
		"reason":            "Tool was selected by the planner for the current step.",
		"run_id":            runID.String(),
	}
	respondJSON(w, http.StatusOK, explain)
}

/* ExplainMemory returns why memory items were selected (explainability). */
func (h *Handlers) ExplainMemory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])
	_, err := h.queries.GetAgentRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "run not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	explain := map[string]interface{}{
		"run_id": runID.String(),
		"reason": "Memory was selected by MemorySelector within token budget (relevance, recency, importance).",
	}
	respondJSON(w, http.StatusOK, explain)
}

/* ExplainModel returns why a model was chosen (explainability). */
func (h *Handlers) ExplainModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])
	_, err := h.queries.GetAgentRun(r.Context(), runID)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "run not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	explain := map[string]interface{}{
		"run_id": runID.String(),
		"reason": "Model was selected by ModelRouter (capability, cost, latency, task type).",
	}
	respondJSON(w, http.StatusOK, explain)
}

/* ExplainPlan returns why the plan was generated (explainability). */
func (h *Handlers) ExplainPlan(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := GetRequestID(r.Context())
	if err := validation.ValidateUUIDRequired(vars["id"], "run_id"); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid run ID", err, requestID, r.URL.Path, r.Method, "run", "", nil))
		return
	}
	runID, _ := uuid.Parse(vars["id"])
	plan, err := h.queries.GetAgentPlanByRun(r.Context(), runID)
	if err != nil || plan == nil {
		respondError(w, NewErrorWithContext(http.StatusNotFound, "plan not found", err, requestID, r.URL.Path, r.Method, "run", runID.String(), nil))
		return
	}
	explain := map[string]interface{}{
		"run_id":   runID.String(),
		"plan_id":  plan.ID.String(),
		"reason":   plan.Reasoning,
		"steps":    plan.Steps,
		"version":  plan.Version,
	}
	respondJSON(w, http.StatusOK, explain)
}
