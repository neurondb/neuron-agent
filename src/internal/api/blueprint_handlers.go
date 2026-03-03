/*-------------------------------------------------------------------------
 *
 * blueprint_handlers.go
 *    API handlers for agent blueprints
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/blueprint_handlers.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"net/http"

	"github.com/neurondb/NeuronAgent/internal/agent"
	"github.com/neurondb/NeuronAgent/internal/auth"
	"github.com/neurondb/NeuronAgent/internal/db"
)

/* ListBlueprints returns all built-in agent blueprints (GET /api/v1/blueprints) */
func (h *Handlers) ListBlueprints(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin, auth.RoleUser); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "insufficient permissions", err, requestID, r.URL.Path, r.Method, "blueprints", "", nil))
		return
	}
	respondJSON(w, http.StatusOK, agent.GetBlueprints())
}

/* CreateAgentFromBlueprintRequest is the body for POST /api/v1/agents/from-blueprint */
type CreateAgentFromBlueprintRequest struct {
	BlueprintID string  `json:"blueprint_id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

/* CreateAgentFromBlueprint creates an agent from a blueprint (POST /api/v1/agents/from-blueprint) */
func (h *Handlers) CreateAgentFromBlueprint(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	endpoint := r.URL.Path
	method := r.Method

	if apiKey, ok := GetAPIKeyFromContext(r.Context()); !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	} else if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin, auth.RoleUser); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "insufficient permissions", err, requestID, endpoint, method, "agent", "", nil))
		return
	}

	var req CreateAgentFromBlueprintRequest
	if err := DecodeJSON(r, &req); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, endpoint, method, "agent", "", nil))
		return
	}
	if req.BlueprintID == "" || req.Name == "" {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "blueprint_id and name are required", nil, requestID, endpoint, method, "agent", "", nil))
		return
	}

	bp := agent.GetBlueprintByID(req.BlueprintID)
	if bp == nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "unknown blueprint_id", nil, requestID, endpoint, method, "agent", "", nil))
		return
	}

	desc := ""
	if req.Description != nil {
		desc = *req.Description
	} else {
		desc = bp.Description
	}

	orgID, _ := auth.GetOrgIDFromContext(r.Context())
	agentRow := &db.Agent{
		OrgID:        orgID,
		Name:         req.Name,
		Description:  &desc,
		SystemPrompt: bp.SystemPrompt,
		ModelName:    bp.ModelName,
		Config:       db.FromMap(bp.Config),
		EnabledTools: bp.EnabledTools,
	}
	if err := h.queries.CreateAgent(r.Context(), agentRow); err != nil {
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "agent creation failed", err, requestID, endpoint, method, "agent", "", nil))
		return
	}
	respondJSON(w, http.StatusCreated, toAgentResponse(agentRow))
}
