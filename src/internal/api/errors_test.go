/*-------------------------------------------------------------------------
 *
 * errors_test.go
 *    Tests for API error types and helpers.
 *
 *-------------------------------------------------------------------------
 */

package api

import (
	"errors"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	err := NewError(http.StatusBadRequest, "invalid input", nil)
	if err.Code != http.StatusBadRequest {
		t.Errorf("Code = %d", err.Code)
	}
	if err.Message != "invalid input" {
		t.Errorf("Message = %q", err.Message)
	}
	if err.Details == nil {
		t.Error("Details should be non-nil")
	}
}

func TestAPIError_Error(t *testing.T) {
	e := NewError(http.StatusNotFound, "not found", nil)
	s := e.Error()
	if s != "not found" {
		t.Errorf("Error() = %q", s)
	}
	e.Err = errors.New("underlying")
	s = e.Error()
	if s == "" || len(s) < 10 {
		t.Errorf("Error() with Err = %q", s)
	}
}

func TestErrorCodeFromHTTP(t *testing.T) {
	tests := []struct {
		code int
		want string
	}{
		{http.StatusBadRequest, "bad_request"},
		{http.StatusUnauthorized, "unauthorized"},
		{http.StatusForbidden, "forbidden"},
		{http.StatusNotFound, "not_found"},
		{http.StatusConflict, "conflict"},
		{http.StatusTooManyRequests, "rate_limit_exceeded"},
		{http.StatusInternalServerError, "internal_error"},
		{http.StatusServiceUnavailable, "service_unavailable"},
		{http.StatusGatewayTimeout, "gateway_timeout"},
		{499, "client_error"},
		{599, "server_error"},
		{200, "unknown"},
	}
	for _, tt := range tests {
		if got := ErrorCodeFromHTTP(tt.code); got != tt.want {
			t.Errorf("ErrorCodeFromHTTP(%d) = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestRetryableFromHTTP(t *testing.T) {
	if !RetryableFromHTTP(http.StatusTooManyRequests) {
		t.Error("429 should be retryable")
	}
	if !RetryableFromHTTP(http.StatusServiceUnavailable) {
		t.Error("503 should be retryable")
	}
	if !RetryableFromHTTP(http.StatusInternalServerError) {
		t.Error("500 should be retryable")
	}
	if RetryableFromHTTP(http.StatusBadRequest) {
		t.Error("400 should not be retryable")
	}
	if RetryableFromHTTP(http.StatusNotFound) {
		t.Error("404 should not be retryable")
	}
}

func TestWrapError(t *testing.T) {
	e := NewError(http.StatusBadRequest, "bad", nil)
	w := WrapError(e, "req-123")
	if w == nil {
		t.Fatal("WrapError returned nil")
	}
	if w.RequestID != "req-123" {
		t.Errorf("RequestID = %q", w.RequestID)
	}
	if w.Code != e.Code {
		t.Errorf("Code = %d", w.Code)
	}
	if WrapError(nil, "x") != nil {
		t.Error("WrapError(nil) should return nil")
	}
}

func TestNewErrorWithContext(t *testing.T) {
	e := NewErrorWithContext(http.StatusNotFound, "missing", nil, "rid", "/path", "GET", "agent", "id-1", nil)
	if e.RequestID != "rid" || e.Endpoint != "/path" || e.Method != "GET" {
		t.Errorf("context not set: %+v", e)
	}
	if e.Details == nil {
		t.Error("Details should be non-nil")
	}
}
