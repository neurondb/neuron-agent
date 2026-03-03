/*-------------------------------------------------------------------------
 *
 * org_middleware.go
 *    Injects org_id from API key into request context for multi-tenant filtering.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/org_middleware.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"net/http"

	"github.com/neurondb/NeuronAgent/internal/auth"
)

/* OrgMiddleware sets org_id in context from the authenticated API key's organization_id.
 * Must run after AuthMiddleware so API key is in context. */
func OrgMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, ok := GetAPIKeyFromContext(r.Context())
		if !ok || apiKey == nil {
			next.ServeHTTP(w, r)
			return
		}
		orgID := auth.ParseOrgIDFromString(apiKey.OrganizationID)
		if orgID != nil {
			ctx := auth.WithOrgID(r.Context(), orgID)
			r = r.WithContext(ctx)
		}
		next.ServeHTTP(w, r)
	})
}
