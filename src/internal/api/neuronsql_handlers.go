/*-------------------------------------------------------------------------
 *
 * neuronsql_handlers.go
 *    NeuronSQL API: POST /v1/neuronsql/generate, optimize, validate, plpgsql
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/api/neuronsql_handlers.go
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/neurondb/NeuronAgent/internal/metrics"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/orchestrator"
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* NeuronSQLHandlers handles /v1/neuronsql/* requests */
type NeuronSQLHandlers struct {
	Orchestrator *orchestrator.Orchestrator
}

/* NewNeuronSQLHandlers creates handlers with the given orchestrator */
func NewNeuronSQLHandlers(orch *orchestrator.Orchestrator) *NeuronSQLHandlers {
	return &NeuronSQLHandlers{Orchestrator: orch}
}

/* Generate handles POST /v1/neuronsql/generate */
func (h *NeuronSQLHandlers) Generate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := GetRequestID(r.Context())
	var req neuronsql.GenerateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("generate", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	dsn := req.DBDSN
	if dsn == "" {
		metrics.RecordNeuronSQLRequest("generate", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "db_dsn required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	if req.Question == "" {
		metrics.RecordNeuronSQLRequest("generate", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "question required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	resp, err := h.Orchestrator.RunGenerate(r.Context(), dsn, req.Question, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("generate", "5xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "generate failed", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	metrics.RecordNeuronSQLRequest("generate", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* Optimize handles POST /v1/neuronsql/optimize */
func (h *NeuronSQLHandlers) Optimize(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := GetRequestID(r.Context())
	var req neuronsql.OptimizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("optimize", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	dsn := req.DBDSN
	if dsn == "" {
		metrics.RecordNeuronSQLRequest("optimize", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "db_dsn required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	if req.SQL == "" {
		metrics.RecordNeuronSQLRequest("optimize", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "sql required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	resp, err := h.Orchestrator.RunOptimize(r.Context(), dsn, req.SQL, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("optimize", "5xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "optimize failed", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	metrics.RecordNeuronSQLRequest("optimize", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* Validate handles POST /v1/neuronsql/validate */
func (h *NeuronSQLHandlers) Validate(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := GetRequestID(r.Context())
	var req neuronsql.ValidateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("validate", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	dsn := req.DBDSN
	if dsn == "" {
		metrics.RecordNeuronSQLRequest("validate", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "db_dsn required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	if req.SQL == "" {
		metrics.RecordNeuronSQLRequest("validate", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "sql required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	resp, err := h.Orchestrator.RunValidate(r.Context(), dsn, req.SQL, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("validate", "5xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "validate failed", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	metrics.RecordNeuronSQLRequest("validate", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

/* PLpgSQL handles POST /v1/neuronsql/plpgsql */
func (h *NeuronSQLHandlers) PLpgSQL(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := GetRequestID(r.Context())
	var req neuronsql.PLpgSQLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "invalid request body", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	dsn := req.DBDSN
	if dsn == "" {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "db_dsn required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	if req.Signature == "" {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "signature required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	if req.Purpose == "" {
		metrics.RecordNeuronSQLRequest("plpgsql", "4xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusBadRequest, "purpose required", nil, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	resp, err := h.Orchestrator.RunPLpgSQL(r.Context(), dsn, req.Signature, req.Purpose, requestID)
	if err != nil {
		metrics.RecordNeuronSQLRequest("plpgsql", "5xx", time.Since(start))
		respondError(w, NewErrorWithContext(http.StatusInternalServerError, "plpgsql generate failed", err, requestID, r.URL.Path, r.Method, "neuronsql", "", nil))
		return
	}
	metrics.RecordNeuronSQLRequest("plpgsql", "2xx", time.Since(start))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
