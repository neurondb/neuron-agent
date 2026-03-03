/*-------------------------------------------------------------------------
 *
 * registry_test.go
 *    Tests for module registry lifecycle.
 *
 *-------------------------------------------------------------------------
 */

package core

import (
	"context"
	"testing"

	"github.com/neurondb/NeuronAgent/internal/config"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/tools"
	"github.com/neurondb/NeuronAgent/pkg/module"
)

/* mockAppContext is a minimal AppContext for tests. */
type mockAppContext struct{}

func (mockAppContext) DB() *db.DB                                    { return nil }
func (mockAppContext) Queries() *db.Queries                          { return nil }
func (mockAppContext) Config() *config.Config                        { return nil }
func (mockAppContext) ModuleConfig(name string) map[string]interface{} { return nil }
func (mockAppContext) RegisterRoutes(prefix string, fn func(module.Router)) {}
func (mockAppContext) RegisterToolHandler(name string, handler tools.ToolHandler) {}
func (mockAppContext) RegisterMiddleware(mw module.Middleware)       {}
func (mockAppContext) RegisterHook(hook module.HookPoint, fn module.HookFunc) {}
func (mockAppContext) RegisterWorker(name string, worker module.Worker) {}
func (mockAppContext) Logger(module string) module.Logger            { return nil }
func (mockAppContext) GetModule(name string) (module.Module, bool)    { return nil, false }

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	m := &stubModule{name: "test"}
	if err := r.Register(m); err != nil {
		t.Fatalf("Register: %v", err)
	}
	got, ok := r.Get("test")
	if !ok || got.Name() != "test" {
		t.Errorf("Get: got ok=%v name=%v", ok, got)
	}
	if _, ok := r.Get("missing"); ok {
		t.Error("Get(missing) should be false")
	}
}

func TestRegistry_RegisterDuplicate(t *testing.T) {
	r := NewRegistry()
	m := &stubModule{name: "dup"}
	_ = r.Register(m)
	err := r.Register(m)
	if err == nil {
		t.Error("expected error for duplicate register")
	}
}

func TestRegistry_ResolveAndOrdered(t *testing.T) {
	r := NewRegistry()
	a := &stubModule{name: "a"}
	b := &stubModule{name: "b", dependencies: []string{"a"}}
	_ = r.Register(b)
	_ = r.Register(a)
	if err := r.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	order := r.Ordered()
	if len(order) != 2 {
		t.Fatalf("expected 2 modules, got %d", len(order))
	}
	if order[0].Name() != "a" || order[1].Name() != "b" {
		t.Errorf("expected order [a,b], got [%s,%s]", order[0].Name(), order[1].Name())
	}
}

func TestRegistry_InitAll_StartAll_StopAll(t *testing.T) {
	r := NewRegistry()
	m := &stubModule{name: "life"}
	if err := r.Register(m); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := r.Resolve(); err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	ctx := context.Background()
	var mock mockAppContext
	if err := r.InitAll(ctx, mock); err != nil {
		t.Fatalf("InitAll: %v", err)
	}
	if err := r.StartAll(ctx); err != nil {
		t.Fatalf("StartAll: %v", err)
	}
	if err := r.StopAll(ctx); err != nil {
		t.Fatalf("StopAll: %v", err)
	}
}
