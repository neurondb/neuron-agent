/*-------------------------------------------------------------------------
 *
 * request_id_test.go
 *    Tests for request ID middleware and context.
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetRequestID_Empty(t *testing.T) {
	if got := GetRequestID(context.Background()); got != "" {
		t.Errorf("GetRequestID(empty) = %q", got)
	}
}

func TestGetRequestID_FromContext(t *testing.T) {
	ctx := context.WithValue(context.Background(), requestIDKey, "req-123")
	if got := GetRequestID(ctx); got != "req-123" {
		t.Errorf("GetRequestID = %q", got)
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		if id == "" {
			t.Error("request ID should be set")
		}
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestIDMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Header().Get("X-Request-ID") == "" {
		t.Error("X-Request-ID header should be set")
	}
	if rec.Header().Get("X-API-Version") != "v1" {
		t.Errorf("X-API-Version = %q", rec.Header().Get("X-API-Version"))
	}
}

func TestRequestIDMiddleware_PreservesHeader(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := GetRequestID(r.Context())
		if id != "custom-id" {
			t.Errorf("request ID = %q", id)
		}
		w.WriteHeader(http.StatusOK)
	})
	handler := RequestIDMiddleware(next)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Request-ID", "custom-id")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Header().Get("X-Request-ID") != "custom-id" {
		t.Errorf("X-Request-ID = %q", rec.Header().Get("X-Request-ID"))
	}
}
