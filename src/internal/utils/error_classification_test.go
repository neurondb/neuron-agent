/*-------------------------------------------------------------------------
 *
 * error_classification_test.go
 *    Tests for error classification utilities.
 *
 *-------------------------------------------------------------------------
 */

package utils

import (
	"errors"
	"testing"
)

func TestClassifyError(t *testing.T) {
	if got := ClassifyError(nil); got != ErrorTypeNonRetryable {
		t.Errorf("ClassifyError(nil) = %v", got)
	}
	if got := ClassifyError(errors.New("connection refused")); got != ErrorTypeRetryable {
		t.Errorf("ClassifyError(connection) = %v", got)
	}
	if got := ClassifyError(errors.New("timeout")); got != ErrorTypeTimeout {
		t.Errorf("ClassifyError(timeout) = %v", got)
	}
	if got := ClassifyError(errors.New("rate limit exceeded")); got != ErrorTypeRateLimit {
		t.Errorf("ClassifyError(rate limit) = %v", got)
	}
	if got := ClassifyError(errors.New("validation failed")); got != ErrorTypeNonRetryable {
		t.Errorf("ClassifyError(validation) = %v", got)
	}
	if got := ClassifyError(errors.New("unauthorized")); got != ErrorTypeNonRetryable {
		t.Errorf("ClassifyError(unauthorized) = %v", got)
	}
}

func TestIsRetryable(t *testing.T) {
	if IsRetryable(nil) {
		t.Error("IsRetryable(nil) should be false")
	}
	if !IsRetryable(errors.New("connection refused")) {
		t.Error("IsRetryable(connection) should be true")
	}
	if !IsRetryable(errors.New("timeout")) {
		t.Error("IsRetryable(timeout) should be true")
	}
	if IsRetryable(errors.New("invalid input")) {
		t.Error("IsRetryable(invalid) should be false")
	}
}
