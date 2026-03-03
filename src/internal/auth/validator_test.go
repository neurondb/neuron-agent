/*-------------------------------------------------------------------------
 *
 * validator_test.go
 *    Tests for in-memory rate limiter.
 *
 *-------------------------------------------------------------------------
 */

package auth

import (
	"testing"
)

func TestNewRateLimiter(t *testing.T) {
	r := NewRateLimiter()
	if r == nil {
		t.Fatal("NewRateLimiter returned nil")
	}
}

func TestRateLimiter_CheckLimit(t *testing.T) {
	r := NewRateLimiter()
	key := "key1"
	limitPerMin := 3
	/* First 3 should be allowed */
	for i := 0; i < limitPerMin; i++ {
		if !r.CheckLimit(key, limitPerMin) {
			t.Errorf("request %d should be allowed", i+1)
		}
	}
	/* 4th should be rejected */
	if r.CheckLimit(key, limitPerMin) {
		t.Error("4th request should be rejected")
	}
	/* Different key should be allowed */
	if !r.CheckLimit("key2", limitPerMin) {
		t.Error("different key should be allowed")
	}
}
