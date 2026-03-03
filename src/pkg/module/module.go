/*-------------------------------------------------------------------------
 *
 * module.go
 *    Module and AppContext interfaces for NeuronAgent modular architecture.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/pkg/module/module.go
 *
 *-------------------------------------------------------------------------
 */

package module

import (
	"context"

	"github.com/neurondb/NeuronAgent/internal/config"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/tools"
)

/* Module is the central contract for all NeuronAgent modules. */
type Module interface {
	Name() string
	Version() string
	Dependencies() []string
	Init(ctx context.Context, app AppContext) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Health(ctx context.Context) HealthStatus
}

/* AppContext is the core surface area exposed to modules. */
type AppContext interface {
	DB() *db.DB
	Queries() *db.Queries
	Config() *config.Config
	ModuleConfig(name string) map[string]interface{}
	RegisterRoutes(prefix string, fn func(r Router))
	RegisterToolHandler(name string, handler tools.ToolHandler)
	RegisterMiddleware(mw Middleware)
	RegisterHook(hook HookPoint, fn HookFunc)
	RegisterWorker(name string, worker Worker)
	Logger(module string) Logger
	GetModule(name string) (Module, bool)
}

/* HealthStatus represents module health for readiness/liveness. */
type HealthStatus struct {
	Healthy bool
	Reason  string
}
