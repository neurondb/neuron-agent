/*-------------------------------------------------------------------------
 *
 * app.go
 *    NeuronAgent kernel: App implements AppContext and owns module registry.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/core/app.go
 *
 *-------------------------------------------------------------------------
 */

package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/gorilla/mux"
	"github.com/neurondb/NeuronAgent/internal/config"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/metrics"
	"github.com/neurondb/NeuronAgent/internal/tools"
	"github.com/neurondb/NeuronAgent/internal/utils"
	"github.com/neurondb/NeuronAgent/pkg/module"
)

/* App is the kernel: config, DB, registry, and extension point collections. */
type App struct {
	cfg      *config.Config
	database *db.DB
	queries  *db.Queries
	registry *Registry

	mu             sync.Mutex
	routeRegs      []routeReg
	toolHandlers   map[string]tools.ToolHandler
	middlewares    []module.Middleware
	hooks          map[module.HookPoint][]module.HookFunc
	workers        map[string]module.Worker
}

type routeReg struct {
	prefix string
	fn     func(module.Router)
}

/* NewApp creates an App and opens the database connection. */
func NewApp(cfg *config.Config) (*App, error) {
	connStr := utils.BuildConnectionString(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
		"",
	)
	database, err := db.NewDB(connStr, db.PoolConfig{
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
		ConnMaxIdleTime: cfg.Database.ConnMaxIdleTime,
	})
	if err != nil {
		return nil, fmt.Errorf("database connection: %w", err)
	}
	queries := db.NewQueries(database.DB)
	queries.SetConnInfoFunc(database.GetConnInfoString)

	return &App{
		cfg:          cfg,
		database:     database,
		queries:      queries,
		registry:     NewRegistry(),
		toolHandlers: make(map[string]tools.ToolHandler),
		hooks:        make(map[module.HookPoint][]module.HookFunc),
		workers:      make(map[string]module.Worker),
	}, nil
}

/* Close closes the database connection. */
func (a *App) Close() error {
	if a.database != nil {
		return a.database.Close()
	}
	return nil
}

/* Registry returns the module registry. */
func (a *App) Registry() *Registry {
	return a.registry
}

/* DB implements module.AppContext. */
func (a *App) DB() *db.DB {
	return a.database
}

/* Queries implements module.AppContext. */
func (a *App) Queries() *db.Queries {
	return a.queries
}

/* Config implements module.AppContext. */
func (a *App) Config() *config.Config {
	return a.cfg
}

/* ModuleConfig implements module.AppContext. */
func (a *App) ModuleConfig(name string) map[string]interface{} {
	if a.cfg.Modules == nil {
		return nil
	}
	ent, ok := a.cfg.Modules[name]
	if !ok || !ent.Enabled {
		return nil
	}
	return ent.Config
}

/* RegisterRoutes implements module.AppContext. */
func (a *App) RegisterRoutes(prefix string, fn func(r module.Router)) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.routeRegs = append(a.routeRegs, routeReg{prefix: prefix, fn: fn})
}

/* RegisterToolHandler implements module.AppContext. */
func (a *App) RegisterToolHandler(name string, handler tools.ToolHandler) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.toolHandlers[name] = handler
}

/* RegisterMiddleware implements module.AppContext. */
func (a *App) RegisterMiddleware(mw module.Middleware) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.middlewares = append(a.middlewares, mw)
}

/* RegisterHook implements module.AppContext. */
func (a *App) RegisterHook(hook module.HookPoint, fn module.HookFunc) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.hooks[hook] = append(a.hooks[hook], fn)
}

/* RegisterWorker implements module.AppContext. */
func (a *App) RegisterWorker(name string, worker module.Worker) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.workers[name] = worker
}

/* Logger implements module.AppContext. */
func (a *App) Logger(moduleName string) module.Logger {
	return &moduleLogger{module: moduleName}
}

/* GetModule implements module.AppContext. */
func (a *App) GetModule(name string) (module.Module, bool) {
	return a.registry.Get(name)
}

/* RegisteredToolHandlers returns tool handlers registered by modules for merging into the main tool registry. */
func (a *App) RegisteredToolHandlers() map[string]tools.ToolHandler {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make(map[string]tools.ToolHandler, len(a.toolHandlers))
	for k, v := range a.toolHandlers {
		out[k] = v
	}
	return out
}

/* ApplyModuleRoutes mounts all module routes onto the given router (e.g. v1 subrouter). */
func (a *App) ApplyModuleRoutes(r *mux.Router) {
	a.mu.Lock()
	regs := make([]routeReg, len(a.routeRegs))
	copy(regs, a.routeRegs)
	a.mu.Unlock()

	for _, reg := range regs {
		sub := r.PathPrefix("/" + reg.prefix).Subrouter()
		reg.fn(NewMuxRouter(sub))
	}
}

/* Middlewares returns registered middlewares in order. */
func (a *App) Middlewares() []module.Middleware {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]module.Middleware(nil), a.middlewares...)
}

/* StartWorkers starts all registered workers in background goroutines. */
func (a *App) StartWorkers(ctx context.Context) {
	a.mu.Lock()
	workers := make(map[string]module.Worker, len(a.workers))
	for k, v := range a.workers {
		workers[k] = v
	}
	a.mu.Unlock()

	for name, w := range workers {
		go func(n string, worker module.Worker) {
			if err := worker.Start(ctx); err != nil && ctx.Err() == nil {
				metrics.ErrorWithContext(ctx, "worker stopped", err, map[string]interface{}{"worker": n})
			}
		}(name, w)
	}
}

/* Ensure App implements module.AppContext. */
var _ module.AppContext = (*App)(nil)

/* moduleLogger adapts zerolog to module.Logger. */
type moduleLogger struct {
	module string
}

func (l *moduleLogger) Debug(msg string, keysAndValues ...interface{}) {
	m := toMap(keysAndValues)
	m["module"] = l.module
	metrics.DebugWithContext(context.Background(), msg, m)
}

func (l *moduleLogger) Info(msg string, keysAndValues ...interface{}) {
	m := toMap(keysAndValues)
	m["module"] = l.module
	metrics.InfoWithContext(context.Background(), msg, m)
}

func (l *moduleLogger) Warn(msg string, keysAndValues ...interface{}) {
	m := toMap(keysAndValues)
	m["module"] = l.module
	metrics.WarnWithContext(context.Background(), msg, m)
}

func (l *moduleLogger) Error(msg string, keysAndValues ...interface{}) {
	m := toMap(keysAndValues)
	m["module"] = l.module
	metrics.ErrorWithContext(context.Background(), msg, fmt.Errorf("%s", msg), m)
}

func toMap(keysAndValues []interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	for i := 0; i+1 < len(keysAndValues); i += 2 {
		if k, ok := keysAndValues[i].(string); ok {
			m[k] = keysAndValues[i+1]
		}
	}
	return m
}
