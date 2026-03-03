/*-------------------------------------------------------------------------
 *
 * admin_handlers.go
 *    Admin-only endpoints: config dump and diagnostics
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/admin_handlers.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"encoding/json"
	"net/http"

	"github.com/neurondb/NeuronAgent/internal/auth"
	"github.com/neurondb/NeuronAgent/internal/config"
	"github.com/neurondb/NeuronAgent/internal/core"
)

/* AdminHandlers provides admin RBAC-guarded endpoints */
type AdminHandlers struct {
	cfg *config.Config
	app *core.App
}

/* NewAdminHandlers creates admin handlers */
func NewAdminHandlers(cfg *config.Config, app *core.App) *AdminHandlers {
	return &AdminHandlers{cfg: cfg, app: app}
}

/* GetAdminConfig returns redacted config (GET /api/v1/admin/config). Requires admin role. */
func (h *AdminHandlers) GetAdminConfig(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "admin role required", err, requestID, r.URL.Path, r.Method, "admin", "", nil))
		return
	}
	dump := config.ConfigDump(h.cfg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(dump)
}

/* GetAdminDiagnostics returns health of DB, modules, version (GET /api/v1/admin/diagnostics). Requires admin role. */
func (h *AdminHandlers) GetAdminDiagnostics(w http.ResponseWriter, r *http.Request) {
	requestID := GetRequestID(r.Context())
	apiKey, ok := GetAPIKeyFromContext(r.Context())
	if !ok {
		respondError(w, WrapError(ErrUnauthorized, requestID))
		return
	}
	if err := auth.RequireRole(apiKey, auth.RoleAdmin); err != nil {
		respondError(w, NewErrorWithContext(http.StatusForbidden, "admin role required", err, requestID, r.URL.Path, r.Method, "admin", "", nil))
		return
	}
	out := make(map[string]interface{})
	out["version"] = "latest"
	/* DB health */
	if h.app != nil && h.app.DB() != nil {
		if err := h.app.DB().Ping(); err != nil {
			out["database"] = map[string]interface{}{"status": "unhealthy", "error": err.Error()}
		} else {
			out["database"] = map[string]interface{}{"status": "healthy"}
		}
	}
	/* Modules health */
	if h.app != nil {
		modules := make(map[string]interface{})
		for _, m := range h.app.Registry().Ordered() {
			health := m.Health(r.Context())
			modules[m.Name()] = map[string]interface{}{"healthy": health.Healthy, "reason": health.Reason}
		}
		out["modules"] = modules
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}
