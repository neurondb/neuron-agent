/*-------------------------------------------------------------------------
 *
 * uuid_test.go
 *    Tests for UUID utilities.
 *
 *-------------------------------------------------------------------------
 */

package utils

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateUUID(t *testing.T) {
	id := GenerateUUID()
	if id == uuid.Nil {
		t.Error("GenerateUUID() returned nil UUID")
	}
	id2 := GenerateUUID()
	if id == id2 {
		t.Error("GenerateUUID() returned duplicate UUIDs")
	}
}

func TestGenerateUUIDString(t *testing.T) {
	s := GenerateUUIDString()
	if s == "" {
		t.Error("GenerateUUIDString() returned empty")
	}
	if _, err := uuid.Parse(s); err != nil {
		t.Errorf("GenerateUUIDString() = %q, not valid UUID: %v", s, err)
	}
}

func TestParseUUID(t *testing.T) {
	valid := "550e8400-e29b-41d4-a716-446655440000"
	parsed, err := ParseUUID(valid)
	if err != nil {
		t.Fatalf("ParseUUID(%q) = %v", valid, err)
	}
	if parsed.String() != valid {
		t.Errorf("ParseUUID(%q) = %s", valid, parsed.String())
	}
	_, err = ParseUUID("not-a-uuid")
	if err == nil {
		t.Error("ParseUUID(invalid) expected error")
	}
}

func TestIsValidUUID(t *testing.T) {
	tests := []struct {
		s    string
		want bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true},
		{"00000000-0000-0000-0000-000000000000", true},
		{"", false},
		{"not-a-uuid", false},
		{"550e8400-e29b-41d4-a716", false},
	}
	for _, tt := range tests {
		if got := IsValidUUID(tt.s); got != tt.want {
			t.Errorf("IsValidUUID(%q) = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestMustParseUUID(t *testing.T) {
	valid := "550e8400-e29b-41d4-a716-446655440000"
	parsed := MustParseUUID(valid)
	if parsed.String() != valid {
		t.Errorf("MustParseUUID(%q) = %s", valid, parsed.String())
	}
	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustParseUUID(invalid) expected panic")
			}
		}()
		MustParseUUID("invalid")
	}()
}
