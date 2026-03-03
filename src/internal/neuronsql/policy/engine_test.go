package policy

import (
	"context"
	"testing"

	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

func TestPolicyEngineImpl_Check(t *testing.T) {
	engine := NewPolicyEngineImpl(nil)
	ctx := context.Background()

	allowed, _ := engine.Check(ctx, "SELECT 1", neuronsql.PolicyContext{})
	if !allowed.Allowed {
		t.Errorf("SELECT 1 should be allowed, got reason %q", allowed.Reason)
	}
	if allowed.StatementClass != ClassSelect {
		t.Errorf("statement class = %q, want select", allowed.StatementClass)
	}

	blocked, _ := engine.Check(ctx, "DROP TABLE t", neuronsql.PolicyContext{})
	if blocked.Allowed {
		t.Error("DROP TABLE should be blocked")
	}

	blocked2, _ := engine.Check(ctx, "SELECT 1; DELETE FROM t", neuronsql.PolicyContext{})
	if blocked2.Allowed {
		t.Error("multi-statement should be blocked")
	}

	blocked3, _ := engine.Check(ctx, "SELECT pg_read_file('/etc/passwd')", neuronsql.PolicyContext{})
	if blocked3.Allowed {
		t.Error("pg_read_file should be blocked")
	}
}

func TestPolicyEngineImpl_Sanitize(t *testing.T) {
	engine := NewPolicyEngineImpl(nil)
	got := engine.Sanitize("  SELECT   \n  1  -- comment  ")
	if got != "SELECT 1" {
		t.Errorf("Sanitize = %q, want \"SELECT 1\"", got)
	}
}
