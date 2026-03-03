/*-------------------------------------------------------------------------
 *
 * router.go
 *    Router adapter wrapping gorilla/mux for module route registration.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/core/router.go
 *
 *-------------------------------------------------------------------------
 */

package core

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/neurondb/NeuronAgent/pkg/module"
)

/* MuxRouter adapts a gorilla/mux subrouter to module.Router. */
type MuxRouter struct {
	r *mux.Router
}

/* NewMuxRouter creates a Router that registers onto r. */
func NewMuxRouter(r *mux.Router) *MuxRouter {
	return &MuxRouter{r: r}
}

/* GET registers a GET handler. */
func (m *MuxRouter) GET(path string, handler http.HandlerFunc) {
	m.r.HandleFunc(path, handler).Methods(http.MethodGet)
}

/* POST registers a POST handler. */
func (m *MuxRouter) POST(path string, handler http.HandlerFunc) {
	m.r.HandleFunc(path, handler).Methods(http.MethodPost)
}

/* PUT registers a PUT handler. */
func (m *MuxRouter) PUT(path string, handler http.HandlerFunc) {
	m.r.HandleFunc(path, handler).Methods(http.MethodPut)
}

/* DELETE registers a DELETE handler. */
func (m *MuxRouter) DELETE(path string, handler http.HandlerFunc) {
	m.r.HandleFunc(path, handler).Methods(http.MethodDelete)
}

/* PATCH registers a PATCH handler. */
func (m *MuxRouter) PATCH(path string, handler http.HandlerFunc) {
	m.r.HandleFunc(path, handler).Methods(http.MethodPatch)
}

/* Ensure MuxRouter implements module.Router. */
var _ module.Router = (*MuxRouter)(nil)
