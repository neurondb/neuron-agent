/*-------------------------------------------------------------------------
 *
 * governance_handlers.go
 *    Governance dashboard APIs: costs, tool risk, policy blocks, memory, performance
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/governance_handlers.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"net/http"
	"time"

	"github.com/neurondb/NeuronAgent/internal/auth"
)

/* GovernanceCostsResponse is the response for GET /api/v1/governance/costs */
type GovernanceCostsResponse struct {
	WorkspaceID string    `json:"workspace_id,omitempty"`
	OrgID      string    `json:"org_id,omitempty"`
	Period     string    `json:"period,omitempty"`
	Cost       float64   `json:"cost"`
	UpdatedAt  time.Time `json:"updated_at"`
}

/* GetGovernanceCosts returns cost analytics by workspace/org/time (GET /api/v1/governance/costs) */
func (h *Handlers) GetGovernanceCosts(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "governance endpoints require admin role", err, requestID, r.URL.Path, r.Method, "governance", "", nil))
		return
	}
	/* Stub: actual implementation would aggregate from cost_per_workspace metric or cost table */
	respondJSON(w, http.StatusOK, []GovernanceCostsResponse{
		{OrgID: "default", Period: "30d", Cost: 0, UpdatedAt: time.Now()},
	})
}

/* ToolRiskResponse is the response for GET /api/v1/governance/tool-risk */
type ToolRiskResponse struct {
	ToolName     string `json:"tool_name"`
	BlockedCalls int64  `json:"blocked_calls,omitempty"`
	ErrorRate    float64 `json:"error_rate,omitempty"`
}

/* GetGovernanceToolRisk returns tool usage risk summary (GET /api/v1/governance/tool-risk) */
func (h *Handlers) GetGovernanceToolRisk(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "governance endpoints require admin role", err, requestID, r.URL.Path, r.Method, "governance", "", nil))
		return
	}
	respondJSON(w, http.StatusOK, []ToolRiskResponse{})
}

/* PolicyBlocksResponse is the response for GET /api/v1/governance/policy-blocks */
type PolicyBlocksResponse struct {
	ReasonCode string `json:"reason_code"`
	Count      int64  `json:"count"`
}

/* GetGovernancePolicyBlocks returns policy block summary by reason_code (GET /api/v1/governance/policy-blocks) */
func (h *Handlers) GetGovernancePolicyBlocks(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "governance endpoints require admin role", err, requestID, r.URL.Path, r.Method, "governance", "", nil))
		return
	}
	respondJSON(w, http.StatusOK, []PolicyBlocksResponse{})
}

/* MemoryGrowthResponse is the response for GET /api/v1/governance/memory-growth */
type MemoryGrowthResponse struct {
	AgentID    string  `json:"agent_id"`
	ChunkCount int     `json:"chunk_count"`
	GrowthRate float64 `json:"growth_rate,omitempty"`
}

/* GetGovernanceMemoryGrowth returns memory growth analytics (GET /api/v1/governance/memory-growth) */
func (h *Handlers) GetGovernanceMemoryGrowth(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "governance endpoints require admin role", err, requestID, r.URL.Path, r.Method, "governance", "", nil))
		return
	}
	respondJSON(w, http.StatusOK, []MemoryGrowthResponse{})
}

/* AgentPerformanceResponse is the response for GET /api/v1/governance/agent-performance */
type AgentPerformanceResponse struct {
	AgentID    string  `json:"agent_id"`
	Score      float64 `json:"score,omitempty"`
	Executions int64   `json:"executions,omitempty"`
	AvgLatency float64 `json:"avg_latency_ms,omitempty"`
}

/* GetGovernanceAgentPerformance returns agent performance scoring (GET /api/v1/governance/agent-performance) */
func (h *Handlers) GetGovernanceAgentPerformance(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "governance endpoints require admin role", err, requestID, r.URL.Path, r.Method, "governance", "", nil))
		return
	}
	respondJSON(w, http.StatusOK, []AgentPerformanceResponse{})
}
