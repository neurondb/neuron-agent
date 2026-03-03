/*-------------------------------------------------------------------------
 *
 * connection_test.go
 *    Tests for DB connection (skip when no database).
 *
 *-------------------------------------------------------------------------
 */

package db

import (
	"testing"
	"time"
)

func TestNewDB_InvalidConnStr(t *testing.T) {
	_, err := NewDB("invalid", PoolConfig{})
	if err == nil {
		t.Error("NewDB(invalid) expected error")
	}
}

func TestNewDBWithRetry_InvalidConnStr(t *testing.T) {
	_, err := NewDBWithRetry("postgres://invalid:5432/", PoolConfig{}, 1, time.Second)
	if err == nil {
		t.Error("NewDBWithRetry(invalid) expected error")
	}
}

func TestParseConnectionInfo(t *testing.T) {
	info := parseConnectionInfo("host=localhost port=5432 dbname=test user=u password=p")
	if info == nil {
		t.Fatal("parseConnectionInfo returned nil")
	}
	if info.Host != "localhost" {
		t.Errorf("Host = %q", info.Host)
	}
	if info.Port != 5432 {
		t.Errorf("Port = %d", info.Port)
	}
}
