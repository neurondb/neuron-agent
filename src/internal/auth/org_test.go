/*-------------------------------------------------------------------------
 *
 * org_test.go
 *    Tests for principal and org context helpers.
 *
 *-------------------------------------------------------------------------
 */

package auth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/neurondb/NeuronAgent/internal/db"
)

func TestWithPrincipal_Nil(t *testing.T) {
	ctx := context.Background()
	out := WithPrincipal(ctx, nil)
	if out != ctx {
		t.Error("WithPrincipal(ctx, nil) should return same context")
	}
}

func TestGetPrincipal_Empty(t *testing.T) {
	if p := GetPrincipal(context.Background()); p != nil {
		t.Error("GetPrincipal(empty) should be nil")
	}
}

func TestWithPrincipal_GetPrincipal(t *testing.T) {
	p := &db.Principal{ID: uuid.New(), Type: "user", Name: "alice"}
	ctx := WithPrincipal(context.Background(), p)
	got := GetPrincipal(ctx)
	if got != p {
		t.Errorf("GetPrincipal = %v", got)
	}
}

func TestWithOrgID_GetOrgIDFromContext(t *testing.T) {
	orgID := uuid.New()
	ctx := WithOrgID(context.Background(), &orgID)
	got, ok := GetOrgIDFromContext(ctx)
	if !ok || got == nil || *got != orgID {
		t.Errorf("GetOrgIDFromContext = %v, %v", got, ok)
	}
	_, ok = GetOrgIDFromContext(context.Background())
	if ok {
		t.Error("GetOrgIDFromContext(empty) should not ok")
	}
}

func TestParseOrgIDFromString(t *testing.T) {
	if ParseOrgIDFromString(nil) != nil {
		t.Error("ParseOrgIDFromString(nil) should be nil")
	}
	s := ""
	if ParseOrgIDFromString(&s) != nil {
		t.Error("ParseOrgIDFromString(empty) should be nil")
	}
	valid := "550e8400-e29b-41d4-a716-446655440000"
	if id := ParseOrgIDFromString(&valid); id == nil || id.String() != valid {
		t.Errorf("ParseOrgIDFromString(valid) = %v", id)
	}
}
