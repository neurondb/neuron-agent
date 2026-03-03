/*-------------------------------------------------------------------------
 *
 * resolver.go
 *    Dependency resolution for modules (topological sort).
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/core/resolver.go
 *
 *-------------------------------------------------------------------------
 */

package core

import (
	"fmt"
	"strings"

	"github.com/neurondb/NeuronAgent/pkg/module"
)

/* ResolveOrder returns module names in init/start order (dependencies first). */
func ResolveOrder(modules []module.Module) ([]module.Module, error) {
	nameToMod := make(map[string]module.Module)
	for _, m := range modules {
		nameToMod[m.Name()] = m
	}

	/* Build adjacency list: dep -> dependents */
	depToDependents := make(map[string][]string)
	for _, m := range modules {
		for _, dep := range m.Dependencies() {
			depToDependents[dep] = append(depToDependents[dep], m.Name())
		}
	}

	/* Check all dependencies exist and detect cycles via DFS */
	visited := make(map[string]bool)
	stack := make(map[string]bool)
	var order []string
	var cycle []string

	var visit func(name string) error
	visit = func(name string) error {
		if stack[name] {
			cycle = append(cycle, name)
			return fmt.Errorf("circular dependency involving module: %s", name)
		}
		if visited[name] {
			return nil
		}
		stack[name] = true
		m, ok := nameToMod[name]
		if ok {
			for _, dep := range m.Dependencies() {
				if _, exists := nameToMod[dep]; !exists {
					return fmt.Errorf("module %s depends on unknown module: %s", name, dep)
				}
				if err := visit(dep); err != nil {
					return err
				}
			}
		}
		stack[name] = false
		visited[name] = true
		order = append(order, name)
		return nil
	}

	for _, m := range modules {
		if err := visit(m.Name()); err != nil {
			if len(cycle) > 0 {
				return nil, fmt.Errorf("circular dependency: %s", strings.Join(cycle, " -> "))
			}
			return nil, err
		}
	}

	/* order is dependency-first (leaves first) */
	out := make([]module.Module, 0, len(order))
	for _, name := range order {
		out = append(out, nameToMod[name])
	}
	return out, nil
}
