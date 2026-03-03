/*-------------------------------------------------------------------------
 *
 * client_test.go
 *    Tests for API client.
 *
 *-------------------------------------------------------------------------
 */

package client

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	c := NewClient("http://localhost:8080", "test-key")
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
}
