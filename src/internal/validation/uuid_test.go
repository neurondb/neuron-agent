/*-------------------------------------------------------------------------
 *
 * uuid_test.go
 *    Tests for UUID validation.
 *
 *-------------------------------------------------------------------------
 */

package validation

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateUUID(t *testing.T) {
	if err := ValidateUUID("", "id"); err == nil {
		t.Error("empty string should error")
	}
	if err := ValidateUUID("not-a-uuid", "id"); err == nil {
		t.Error("invalid UUID should error")
	}
	valid := "550e8400-e29b-41d4-a716-446655440000"
	if err := ValidateUUID(valid, "id"); err != nil {
		t.Errorf("valid UUID: %v", err)
	}
	if err := ValidateUUID("  "+valid+"  ", "id"); err != nil {
		t.Errorf("trimmed UUID: %v", err)
	}
}

func TestValidateUUIDRequired(t *testing.T) {
	if err := ValidateUUIDRequired("", "id"); err == nil {
		t.Error("empty should error")
	}
	if err := ValidateUUIDRequired("550e8400-e29b-41d4-a716-446655440000", "id"); err != nil {
		t.Errorf("valid: %v", err)
	}
}

func TestParseUUID(t *testing.T) {
	_, err := ParseUUID("invalid", "id")
	if err == nil {
		t.Error("invalid should error")
	}
	parsed, err := ParseUUID("550e8400-e29b-41d4-a716-446655440000", "id")
	if err != nil {
		t.Fatalf("ParseUUID: %v", err)
	}
	if parsed == uuid.Nil {
		t.Error("parsed should not be Nil")
	}
}
