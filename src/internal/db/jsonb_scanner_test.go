/*-------------------------------------------------------------------------
 *
 * jsonb_scanner_test.go
 *    Tests for JSONB map type.
 *
 *-------------------------------------------------------------------------
 */

package db

import (
	"testing"
)

func TestFromMap(t *testing.T) {
	if got := FromMap(nil); got == nil || len(got) != 0 {
		t.Errorf("FromMap(nil) = %v", got)
	}
	m := map[string]interface{}{"a": 1}
	got := FromMap(m)
	if got["a"] != 1 {
		t.Errorf("FromMap = %v", got)
	}
}

func TestJSONBMap_ToMap(t *testing.T) {
	var j JSONBMap
	if out := j.ToMap(); out == nil {
		t.Error("ToMap() should return non-nil map")
	}
	j = JSONBMap{"x": "y"}
	out := j.ToMap()
	if out["x"] != "y" {
		t.Errorf("ToMap = %v", out)
	}
}

func TestJSONBMap_Value(t *testing.T) {
	var j JSONBMap
	v, err := j.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if s, ok := v.(string); !ok || s != "{}" {
		t.Errorf("Value() = %v", v)
	}
	j = JSONBMap{"a": float64(1)}
	v, err = j.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	if _, ok := v.([]byte); !ok {
		t.Errorf("Value() = %T", v)
	}
}

func TestJSONBMap_Scan(t *testing.T) {
	var j JSONBMap
	if err := j.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if len(j) != 0 {
		t.Errorf("Scan(nil) = %v", j)
	}
	if err := j.Scan([]byte(`{"k":"v"}`)); err != nil {
		t.Fatalf("Scan([]byte): %v", err)
	}
	if j["k"] != "v" {
		t.Errorf("Scan = %v", j)
	}
}
