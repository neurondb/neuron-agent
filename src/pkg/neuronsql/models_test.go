/*-------------------------------------------------------------------------
 *
 * models_test.go
 *    Tests for NeuronSQL models
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/pkg/neuronsql/models_test.go
 *
 *-------------------------------------------------------------------------
 */

package neuronsql

import (
	"encoding/json"
	"testing"
)

func TestNeuronSQLResponse_JSONRoundTrip(t *testing.T) {
	resp := NeuronSQLResponse{
		RequestID:   "req-123",
		Mode:        "generate",
		SQL:         "SELECT 1",
		Explanation: "test",
		Citations:   []string{"tool:schema_snapshot:req-123", "doc:chunk-1"},
		ValidationReport: &ValidationReport{
			Valid:            true,
			StatementClass:   "select",
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out NeuronSQLResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.RequestID != resp.RequestID || out.Mode != resp.Mode || out.SQL != resp.SQL {
		t.Errorf("round trip mismatch: got %+v", out)
	}
}
