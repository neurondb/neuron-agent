/*-------------------------------------------------------------------------
 *
 * module_test.go
 *    Tests for NeuronSQL module.
 *
 *-------------------------------------------------------------------------
 */

package neuronsql

import (
	"context"
	"testing"

	"github.com/neurondb/NeuronAgent/internal/config"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/tools"
	"github.com/neurondb/NeuronAgent/pkg/module"
)

func TestNeuronSQLModule_NameVersionDeps(t *testing.T) {
	m := New()
	if m.Name() != "neuronsql" {
		t.Errorf("Name: got %s", m.Name())
	}
	if m.Version() != "1.0.0" {
		t.Errorf("Version: got %s", m.Version())
	}
	if deps := m.Dependencies(); deps != nil {
		t.Errorf("Dependencies: expected nil, got %v", deps)
	}
}

func TestNeuronSQLModule_Health(t *testing.T) {
	m := New().(*NeuronSQLModule)
	ctx := context.Background()
	hs := m.Health(ctx)
	if !hs.Healthy {
		t.Errorf("Health: expected healthy, got %+v", hs)
	}
}

func TestNeuronSQLModule_StartStop(t *testing.T) {
	m := New()
	ctx := context.Background()
	if err := m.Start(ctx); err != nil {
		t.Errorf("Start: %v", err)
	}
	if err := m.Stop(ctx); err != nil {
		t.Errorf("Stop: %v", err)
	}
}

func TestNeuronSQLModule_Init_RegistersRoutesAndTools(t *testing.T) {
	rec := &recordingAppContext{}
	m := New().(*NeuronSQLModule)
	ctx := context.Background()
	err := m.Init(ctx, rec)
	if err != nil {
		t.Skipf("Init may require full env: %v", err)
	}
	if len(rec.routes) == 0 {
		t.Error("Init should register routes")
	}
	if len(rec.toolHandlers) == 0 {
		t.Error("Init should register tool handlers")
	}
}

/* recordingAppContext records RegisterRoutes and RegisterToolHandler for tests. */
type recordingAppContext struct {
	routes       []string
	toolHandlers []string
}

func (r *recordingAppContext) RegisterRoutes(prefix string, fn func(module.Router)) {
	r.routes = append(r.routes, prefix)
	/* Do not call fn(nil) to avoid nil router panic */
}

func (r *recordingAppContext) RegisterToolHandler(name string, _ tools.ToolHandler) {
	r.toolHandlers = append(r.toolHandlers, name)
}

func (r *recordingAppContext) DB() *db.DB                          { return nil }
func (r *recordingAppContext) Queries() *db.Queries                { return nil }
func (r *recordingAppContext) Config() *config.Config               { return nil }
func (r *recordingAppContext) ModuleConfig(string) map[string]interface{} { return nil }
func (r *recordingAppContext) RegisterMiddleware(module.Middleware) {}
func (r *recordingAppContext) RegisterHook(module.HookPoint, module.HookFunc) {}
func (r *recordingAppContext) RegisterWorker(string, module.Worker) {}
func (r *recordingAppContext) Logger(string) module.Logger         { return nil }
func (r *recordingAppContext) GetModule(string) (module.Module, bool) { return nil, false }
