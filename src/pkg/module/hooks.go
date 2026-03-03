/*-------------------------------------------------------------------------
 *
 * hooks.go
 *    Extension point types: Router, Middleware, Worker, HookFunc.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/pkg/module/hooks.go
 *
 *-------------------------------------------------------------------------
 */

package module

import (
	"context"
	"net/http"
)

/* Router allows modules to register HTTP routes under a prefix. */
type Router interface {
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
}

/* Middleware is HTTP middleware (wrap next handler). */
type Middleware func(next http.Handler) http.Handler

/* Worker runs in the background; Start blocks until Stop is called or ctx cancelled. */
type Worker interface {
	Start(ctx context.Context) error
}

/* HookPoint identifies where in the agent runtime a hook runs. */
type HookPoint string

const (
	HookPreLLM      HookPoint = "pre_llm"
	HookPostLLM     HookPoint = "post_llm"
	HookPreTool     HookPoint = "pre_tool"
	HookPostTool    HookPoint = "post_tool"
	HookPreResponse HookPoint = "pre_response"
	HookPostResponse HookPoint = "post_response"
)

/* HookFunc is called at a HookPoint; ctx and data can be inspected/modified. */
type HookFunc func(ctx context.Context, data map[string]interface{}) error

/* Logger provides module-scoped logging. */
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}
