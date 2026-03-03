/*-------------------------------------------------------------------------
 *
 * hasher_test.go
 *    Tests for API key hashing utilities.
 *
 *-------------------------------------------------------------------------
 */

package auth

import (
	"testing"
)

func TestHashAPIKey(t *testing.T) {
	key := "sk-test-key-12345"
	hash, err := HashAPIKey(key)
	if err != nil {
		t.Fatalf("HashAPIKey: %v", err)
	}
	if hash == "" || hash == key {
		t.Error("HashAPIKey should return a non-empty hash different from key")
	}
}

func TestVerifyAPIKey(t *testing.T) {
	key := "sk-secret-abc"
	hash, err := HashAPIKey(key)
	if err != nil {
		t.Fatalf("HashAPIKey: %v", err)
	}
	if !VerifyAPIKey(key, hash) {
		t.Error("VerifyAPIKey(key, hash) should be true")
	}
	if VerifyAPIKey("wrong-key", hash) {
		t.Error("VerifyAPIKey(wrong, hash) should be false")
	}
	if VerifyAPIKey(key, "not-a-hash") {
		t.Error("VerifyAPIKey(key, invalid hash) should be false")
	}
}

func TestGetKeyPrefix(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"12345678", "12345678"},
		{"123456789", "12345678"},
		{"short", "short"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := GetKeyPrefix(tt.in); got != tt.want {
			t.Errorf("GetKeyPrefix(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}
