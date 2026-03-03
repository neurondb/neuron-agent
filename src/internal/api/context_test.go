/*-------------------------------------------------------------------------
 *
 * context_test.go
 *    Tests for API context helpers.
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/neurondb/NeuronAgent/internal/auth"
	"github.com/neurondb/NeuronAgent/internal/db"
)

func TestGetAPIKeyFromContext_Empty(t *testing.T) {
	_, ok := GetAPIKeyFromContext(context.Background())
	if ok {
		t.Error("GetAPIKeyFromContext(empty) should not ok")
	}
}

func TestMustGetAPIKeyFromContext_Empty(t *testing.T) {
	_, err := MustGetAPIKeyFromContext(context.Background())
	if err == nil {
		t.Error("MustGetAPIKeyFromContext(empty) expected error")
	}
}

func TestGetPrincipalFromContext_Empty(t *testing.T) {
	principal := auth.GetPrincipal(context.Background())
	if principal != nil {
		t.Error("GetPrincipal(empty) should be nil")
	}
}

func TestGetPrincipalFromContext_WithPrincipal(t *testing.T) {
	p := &db.Principal{ID: uuid.New(), Name: "test", Type: "user"}
	ctx := auth.WithPrincipal(context.Background(), p)
	got := auth.GetPrincipal(ctx)
	if got != p {
		t.Errorf("GetPrincipal(with) = %v", got)
	}
}

func TestGetAuthFromContext_Empty(t *testing.T) {
	apiKey, principal := GetAuthFromContext(context.Background())
	if apiKey != nil || principal != nil {
		t.Errorf("GetAuthFromContext(empty) = %v, %v", apiKey, principal)
	}
}
