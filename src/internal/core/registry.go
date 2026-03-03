/*-------------------------------------------------------------------------
 *
 * registry.go
 *    Module registry with lifecycle management.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/core/registry.go
 *
 *-------------------------------------------------------------------------
 */

package core

import (
	"context"
	"fmt"
	"sync"

	"github.com/neurondb/NeuronAgent/pkg/module"
)

/* Registry stores modules and runs lifecycle in dependency order. */
type Registry struct {
	mu      sync.RWMutex
	modules []module.Module
	order   []module.Module
	byName  map[string]module.Module
}

/* NewRegistry creates an empty registry. */
func NewRegistry() *Registry {
	return &Registry{
		modules: nil,
		order:   nil,
		byName:  make(map[string]module.Module),
	}
}

/* Register adds a module. Call before InitAll. */
func (r *Registry) Register(m module.Module) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := m.Name()
	if _, exists := r.byName[name]; exists {
		return fmt.Errorf("module already registered: %s", name)
	}
	r.modules = append(r.modules, m)
	r.byName[name] = m
	return nil
}

/* Resolve computes init/start order from dependencies. Call after all Register, before InitAll. */
func (r *Registry) Resolve() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, err := ResolveOrder(r.modules)
	if err != nil {
		return err
	}
	r.order = order
	return nil
}

/* Ordered returns modules in dependency order (after Resolve). */
func (r *Registry) Ordered() []module.Module {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]module.Module, len(r.order))
	copy(out, r.order)
	return out
}

/* Get returns a module by name. */
func (r *Registry) Get(name string) (module.Module, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.byName[name]
	return m, ok
}

/* InitAll calls Init on each module in dependency order. */
func (r *Registry) InitAll(ctx context.Context, app module.AppContext) error {
	ordered := r.Ordered()
	for _, m := range ordered {
		if err := m.Init(ctx, app); err != nil {
			return fmt.Errorf("module %s init: %w", m.Name(), err)
		}
	}
	return nil
}

/* StartAll calls Start on each module in dependency order. */
func (r *Registry) StartAll(ctx context.Context) error {
	ordered := r.Ordered()
	for _, m := range ordered {
		if err := m.Start(ctx); err != nil {
			return fmt.Errorf("module %s start: %w", m.Name(), err)
		}
	}
	return nil
}

/* StopAll calls Stop on each module in reverse order. */
func (r *Registry) StopAll(ctx context.Context) error {
	ordered := r.Ordered()
	for i := len(ordered) - 1; i >= 0; i-- {
		m := ordered[i]
		if err := m.Stop(ctx); err != nil {
			return fmt.Errorf("module %s stop: %w", m.Name(), err)
		}
	}
	return nil
}
