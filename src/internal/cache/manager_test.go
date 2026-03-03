/*-------------------------------------------------------------------------
 *
 * manager_test.go
 *    Tests for cache manager and TTLCache.
 *
 *-------------------------------------------------------------------------
 */

package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewTTLCache_GetSet(t *testing.T) {
	c := NewTTLCache(5*time.Minute, 100)
	defer c.Stop()

	_, ok := c.Get("missing")
	if ok {
		t.Error("Get(missing) should be false")
	}
	if err := c.Set("k", "v", time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, ok := c.Get("k")
	if !ok {
		t.Fatal("Get(k) should be true")
	}
	if val != "v" {
		t.Errorf("Get(k) = %v", val)
	}
	if c.Size() != 1 {
		t.Errorf("Size = %d", c.Size())
	}
}

func TestTTLCache_Delete(t *testing.T) {
	c := NewTTLCache(time.Minute, 10)
	defer c.Stop()
	_ = c.Set("a", 1, time.Minute)
	c.Delete("a")
	_, ok := c.Get("a")
	if ok {
		t.Error("Get(a) after Delete should be false")
	}
}

func TestTTLCache_Clear(t *testing.T) {
	c := NewTTLCache(time.Minute, 10)
	defer c.Stop()
	_ = c.Set("x", 1, time.Minute)
	c.Clear()
	if c.Size() != 0 {
		t.Errorf("Size after Clear = %d", c.Size())
	}
}

func TestNewCacheManager(t *testing.T) {
	cm := NewCacheManager(time.Minute, 10*time.Minute, 5*time.Minute, 50)
	defer cm.Close()
	ctx := context.Background()
	if err := cm.CacheResponse(ctx, "r1", "response1", time.Minute); err != nil {
		t.Fatalf("CacheResponse: %v", err)
	}
	got, ok := cm.GetCachedResponse(ctx, "r1")
	if !ok || got != "response1" {
		t.Errorf("GetCachedResponse = %v, %v", got, ok)
	}
	stats := cm.GetStats()
	if stats["responses_size"].(int) != 1 {
		t.Errorf("stats = %v", stats)
	}
}
