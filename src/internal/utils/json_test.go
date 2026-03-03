/*-------------------------------------------------------------------------
 *
 * json_test.go
 *    Tests for JSON marshaling utilities.
 *
 *-------------------------------------------------------------------------
 */

package utils

import (
	"encoding/json"
	"testing"
)

func TestMarshalUnmarshalJSON(t *testing.T) {
	type T struct {
		A int    `json:"a"`
		B string `json:"b"`
	}
	in := T{A: 42, B: "hello"}
	data, err := MarshalJSON(in)
	if err != nil {
		t.Fatalf("MarshalJSON: %v", err)
	}
	var out T
	if err := UnmarshalJSON(data, &out); err != nil {
		t.Fatalf("UnmarshalJSON: %v", err)
	}
	if out.A != in.A || out.B != in.B {
		t.Errorf("roundtrip: got %+v", out)
	}
}

func TestMarshalJSONString(t *testing.T) {
	m := map[string]int{"x": 1}
	s, err := MarshalJSONString(m)
	if err != nil {
		t.Fatalf("MarshalJSONString: %v", err)
	}
	var out map[string]int
	if err := json.Unmarshal([]byte(s), &out); err != nil {
		t.Errorf("MarshalJSONString output not valid JSON: %q", s)
	}
	if out["x"] != 1 {
		t.Errorf("MarshalJSONString roundtrip: got %v", out)
	}
}

func TestUnmarshalJSONString(t *testing.T) {
	var m map[string]int
	if err := UnmarshalJSONString(`{"a":1}`, &m); err != nil {
		t.Fatalf("UnmarshalJSONString: %v", err)
	}
	if m["a"] != 1 {
		t.Errorf("UnmarshalJSONString: got %v", m)
	}
}

func TestValidateJSON(t *testing.T) {
	if err := ValidateJSON(`{"a":1}`); err != nil {
		t.Errorf("ValidateJSON(valid) = %v", err)
	}
	if err := ValidateJSON("not json"); err == nil {
		t.Error("ValidateJSON(invalid) expected error")
	}
}

func TestMergeJSON(t *testing.T) {
	dst := map[string]interface{}{"a": 1, "b": 2}
	src := map[string]interface{}{"b": 20, "c": 3}
	got := MergeJSON(dst, src)
	if got["a"].(int) != 1 {
		t.Errorf("MergeJSON: a = %v", got["a"])
	}
	if got["b"].(int) != 20 {
		t.Errorf("MergeJSON: b = %v (overwritten)", got["b"])
	}
	if got["c"].(int) != 3 {
		t.Errorf("MergeJSON: c = %v", got["c"])
	}
	/* dst and src unchanged */
	if dst["b"].(int) != 2 {
		t.Error("MergeJSON should not modify dst")
	}
}

func TestGetJSONField(t *testing.T) {
	data := map[string]interface{}{
		"top": map[string]interface{}{
			"nested": "value",
		},
	}
	val, err := GetJSONField(data, "top", "nested")
	if err != nil {
		t.Fatalf("GetJSONField: %v", err)
	}
	if val != "value" {
		t.Errorf("GetJSONField = %v", val)
	}
	_, err = GetJSONField(data, "top", "missing")
	if err == nil {
		t.Error("GetJSONField(missing) expected error")
	}
	_, err = GetJSONField(data, "not_a_map")
	if err == nil {
		t.Error("GetJSONField(invalid path) expected error")
	}
	_, err = GetJSONField(data)
	if err == nil {
		t.Error("GetJSONField(empty path) expected error")
	}
}
