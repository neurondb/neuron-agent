/*-------------------------------------------------------------------------
 *
 * claw_handlers.go
 *    Claw gateway: tools/list, tools/run, health (neuronsql.* only)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/claw_handlers.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/neurondb/NeuronAgent/internal/auth"
	"github.com/neurondb/NeuronAgent/internal/tools"
)

/* ClawHandlers provides /claw/v1 endpoints (list tools, run tool, health) */
type ClawHandlers struct {
	registry *tools.Registry
}

/* NewClawHandlers creates Claw handlers */
func NewClawHandlers(registry *tools.Registry) *ClawHandlers {
	return &ClawHandlers{registry: registry}
}

/* ListTools returns only neuronsql.* tools (POST /claw/v1/tools/list) */
func (h *ClawHandlers) ListTools(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin, auth.RoleUser); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "claw: user or admin role required", err, requestID, r.URL.Path, r.Method, "claw", "", nil))
		return
	}
	names := h.registry.ListClawTools()
	out := map[string]interface{}{"tools": names}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

/* RunTool runs a neuronsql.* tool by name (POST /claw/v1/tools/run) */
func (h *ClawHandlers) RunTool(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireAnyRole(apiKey, auth.RoleAdmin, auth.RoleUser); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "claw: user or admin role required", err, requestID, r.URL.Path, r.Method, "claw", "", nil))
		return
	}
	var req struct {
		Tool  string                 `json:"tool"`
		Params map[string]interface{} `json:"params"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, r.URL.Path, r.Method, "claw", "", nil))
		return
	}
	if req.Tool == "" {
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "tool name required", nil, requestID, r.URL.Path, r.Method, "claw", "", nil))
		return
	}
	if !strings.HasPrefix(req.Tool, "neuronsql.") {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "claw: only neuronsql.* tools allowed", nil, requestID, r.URL.Path, r.Method, "claw", "", nil))
		return
	}
	result, err := h.registry.ExecuteByHandlerType(r.Context(), req.Tool, req.Params)
	if err != nil {
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "tool execution failed: "+err.Error(), err, requestID, r.URL.Path, r.Method, "claw", req.Tool, nil))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"result": result})
}

/* Health returns Claw gateway health (GET /claw/v1/health) */
func (h *ClawHandlers) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
