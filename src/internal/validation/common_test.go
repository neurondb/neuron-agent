/*-------------------------------------------------------------------------
 *
 * common_test.go
 *    Tests for common validation helpers.
 *
 *-------------------------------------------------------------------------
 */

package validation

import (
	"net/http"
	"strings"
	"testing"
)

func TestValidateRequired(t *testing.T) {
	if err := ValidateRequired("", "x"); err == nil {
		t.Error("empty should error")
	}
	if err := ValidateRequired("a", "x"); err != nil {
		t.Errorf("non-empty: %v", err)
	}
}

func TestValidateMaxLength(t *testing.T) {
	if err := ValidateMaxLength("abc", "x", 2); err == nil {
		t.Error("over max should error")
	}
	if err := ValidateMaxLength("ab", "x", 2); err != nil {
		t.Errorf("at max: %v", err)
	}
}

func TestValidateLimit(t *testing.T) {
	if err := ValidateLimit(-1); err == nil {
		t.Error("negative should error")
	}
	if err := ValidateLimit(10001); err == nil {
		t.Error("over 10000 should error")
	}
	if err := ValidateLimit(100); err != nil {
		t.Errorf("valid: %v", err)
	}
}

func TestValidateOffset(t *testing.T) {
	if err := ValidateOffset(-1); err == nil {
		t.Error("negative should error")
	}
	if err := ValidateOffset(0); err != nil {
		t.Errorf("zero: %v", err)
	}
}

func TestValidateIntRange(t *testing.T) {
	if err := ValidateIntRange(0, 1, 10, "x"); err == nil {
		t.Error("below min should error")
	}
	if err := ValidateIntRange(11, 1, 10, "x"); err == nil {
		t.Error("above max should error")
	}
	if err := ValidateIntRange(5, 1, 10, "x"); err != nil {
		t.Errorf("in range: %v", err)
	}
}

func TestValidateNonNegative(t *testing.T) {
	if err := ValidateNonNegative(-1, "x"); err == nil {
		t.Error("negative should error")
	}
	if err := ValidateNonNegative(0, "x"); err != nil {
		t.Errorf("zero: %v", err)
	}
}

func TestReadAndValidateBody_NilBody(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", nil)
	req.Body = nil
	_, err := ReadAndValidateBody(req, 100)
	if err == nil {
		t.Error("nil body should error")
	}
}

func TestReadAndValidateBody_WithinLimit(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"a":1}`))
	body, err := ReadAndValidateBody(req, 100)
	if err != nil {
		t.Fatalf("ReadAndValidateBody: %v", err)
	}
	if len(body) == 0 {
		t.Error("body should be read")
	}
}

func TestDecodeJSONBody_EmptyBody(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	var dst map[string]interface{}
	if err := DecodeJSONBody(req, 100, &dst); err == nil {
		t.Error("empty body should error")
	}
}

func TestDecodeJSONBody_ValidJSON(t *testing.T) {
	req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(`{"k":"v"}`))
	req.Header.Set("Content-Type", "application/json")
	var dst map[string]interface{}
	if err := DecodeJSONBody(req, 100, &dst); err != nil {
		t.Fatalf("DecodeJSONBody: %v", err)
	}
	if dst["k"] != "v" {
		t.Errorf("dst = %v", dst)
	}
}
