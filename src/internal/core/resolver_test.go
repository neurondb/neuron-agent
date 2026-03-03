/*-------------------------------------------------------------------------
 *
 * resolver_test.go
 *    Tests for module dependency resolution.
 *
 *-------------------------------------------------------------------------
 */

package core

import (
	"context"
	"strings"
	"testing"

	"github.com/neurondb/NeuronAgent/pkg/module"
)

/* stubModule is a minimal module for testing. */
type stubModule struct {
	name         string
	dependencies []string
}

func (s *stubModule) Name() string                           { return s.name }
func (s *stubModule) Version() string                        { return "0.0.0" }
func (s *stubModule) Dependencies() []string                  { return s.dependencies }
func (s *stubModule) Init(ctx context.Context, app module.AppContext) error { return nil }
func (s *stubModule) Start(ctx context.Context) error         { return nil }
func (s *stubModule) Stop(ctx context.Context) error          { return nil }
func (s *stubModule) Health(ctx context.Context) module.HealthStatus {
	return module.HealthStatus{Healthy: true}
}

func TestResolveOrder_Empty(t *testing.T) {
	order, err := ResolveOrder(nil)
	if err != nil {
		t.Fatalf("ResolveOrder(nil): %v", err)
	}
	if len(order) != 0 {
		t.Errorf("expected empty order, got %d", len(order))
	}
}

func TestResolveOrder_Single(t *testing.T) {
	mod := &stubModule{name: "a"}
	order, err := ResolveOrder([]module.Module{mod})
	if err != nil {
		t.Fatalf("ResolveOrder: %v", err)
	}
	if len(order) != 1 || order[0].Name() != "a" {
		t.Errorf("expected [a], got %v", names(order))
	}
}

func TestResolveOrder_NoDeps(t *testing.T) {
	mods := []module.Module{
		&stubModule{name: "a"},
		&stubModule{name: "b"},
		&stubModule{name: "c"},
	}
	order, err := ResolveOrder(mods)
	if err != nil {
		t.Fatalf("ResolveOrder: %v", err)
	}
	if len(order) != 3 {
		t.Errorf("expected 3 modules, got %d", len(order))
	}
	names := names(order)
	if !sliceContains(names, "a") || !sliceContains(names, "b") || !sliceContains(names, "c") {
		t.Errorf("expected a,b,c in any order, got %v", names)
	}
}

func TestResolveOrder_DepsFirst(t *testing.T) {
	a := &stubModule{name: "a"}
	b := &stubModule{name: "b", dependencies: []string{"a"}}
	c := &stubModule{name: "c", dependencies: []string{"b"}}
	order, err := ResolveOrder([]module.Module{c, a, b})
	if err != nil {
		t.Fatalf("ResolveOrder: %v", err)
	}
	names := names(order)
	idxA := indexOf(names, "a")
	idxB := indexOf(names, "b")
	idxC := indexOf(names, "c")
	if idxA < 0 || idxB < 0 || idxC < 0 {
		t.Fatalf("missing module in %v", names)
	}
	if idxA >= idxB || idxB >= idxC {
		t.Errorf("expected order a,b,c (deps first), got %v", names)
	}
}

func TestResolveOrder_Circular(t *testing.T) {
	a := &stubModule{name: "a", dependencies: []string{"c"}}
	b := &stubModule{name: "b", dependencies: []string{"a"}}
	c := &stubModule{name: "c", dependencies: []string{"b"}}
	_, err := ResolveOrder([]module.Module{a, b, c})
	if err == nil {
		t.Fatal("expected error for circular dependency")
	}
	if !contains(err.Error(), "circular") {
		t.Errorf("error should mention circular dependency: %v", err)
	}
}

func TestResolveOrder_UnknownDep(t *testing.T) {
	a := &stubModule{name: "a", dependencies: []string{"nonexistent"}}
	_, err := ResolveOrder([]module.Module{a})
	if err == nil {
		t.Fatal("expected error for unknown dependency")
	}
	if !contains(err.Error(), "unknown") {
		t.Errorf("error should mention unknown module: %v", err)
	}
}

func names(mods []module.Module) []string {
	out := make([]string, len(mods))
	for i, m := range mods {
		out[i] = m.Name()
	}
	return out
}

func sliceContains(s []string, x string) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func indexOf(s []string, x string) int {
	for i, v := range s {
		if v == x {
			return i
		}
	}
	return -1
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}
