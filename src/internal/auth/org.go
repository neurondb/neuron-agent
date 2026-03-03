/*-------------------------------------------------------------------------
 *
 * org.go
 *    Organization (tenant) and principal context for multi-tenancy and RBAC.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/auth/org.go
 *
 *-------------------------------------------------------------------------
 */

package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/neurondb/NeuronAgent/internal/db"
)

/* Principal context (shared so api, tools, workflow can all read it) */
type principalCtxKey struct{}

var principalContextKey = &principalCtxKey{}

/* WithPrincipal returns a context with principal set (used by API middleware). */
func WithPrincipal(ctx context.Context, p *db.Principal) context.Context {
	if p == nil {
		return ctx
	}
	return context.WithValue(ctx, principalContextKey, p)
}

/* GetPrincipal returns the principal from context, if set. */
func GetPrincipal(ctx context.Context) *db.Principal {
	v := ctx.Value(principalContextKey)
	if v == nil {
		return nil
	}
	p, _ := v.(*db.Principal)
	return p
}

type orgContextKey string

const orgIDContextKey orgContextKey = "org_id"

/* WithOrgID returns a context with org_id set (for multi-tenant filtering). */
func WithOrgID(ctx context.Context, orgID *uuid.UUID) context.Context {
	if orgID == nil {
		return ctx
	}
	return context.WithValue(ctx, orgIDContextKey, orgID)
}

/* GetOrgIDFromContext returns the org_id from context, if set. */
func GetOrgIDFromContext(ctx context.Context) (*uuid.UUID, bool) {
	v := ctx.Value(orgIDContextKey)
	if v == nil {
		return nil, false
	}
	id, ok := v.(*uuid.UUID)
	return id, ok
}

/* ParseOrgIDFromString parses organization_id string from API key into UUID for DB use. */
func ParseOrgIDFromString(s *string) *uuid.UUID {
	if s == nil || *s == "" {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}
