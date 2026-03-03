package policy

import (
	"testing"
)

func TestSQLClassifier_Classify(t *testing.T) {
	c := NewSQLClassifier()
	tests := []struct {
		sql   string
		class string
	}{
		{"SELECT 1", ClassSelect},
		{"  SELECT * FROM t  ", ClassSelect},
		{"WITH x AS (SELECT 1) SELECT * FROM x", ClassWith},
		{"EXPLAIN SELECT 1", ClassExplain},
		{"EXPLAIN (FORMAT JSON) SELECT 1", ClassExplain},
		{"INSERT INTO t VALUES (1)", ClassDML},
		{"UPDATE t SET x = 1", ClassDML},
		{"DELETE FROM t", ClassDML},
		{"DROP TABLE t", ClassDDL},
		{"TRUNCATE t", ClassDDL},
		{"GRANT SELECT ON t TO u", ClassDDL},
		{"REVOKE SELECT ON t FROM u", ClassDDL},
		{"CREATE FUNCTION f() RETURNS int AS $$ SELECT 1 $$ LANGUAGE sql", ClassDDL},
		{"", ClassBlocked},
		{"SELECT 1; SELECT 2", ClassBlocked},
		{"SELECT 1;", ClassSelect},
		{"  -- comment\nSELECT 1", ClassSelect},
	}
	for _, tt := range tests {
		got := c.Classify(tt.sql)
		if got != tt.class {
			t.Errorf("Classify(%q) = %q, want %q", tt.sql, got, tt.class)
		}
	}
}

func TestSQLClassifier_IsAllowed(t *testing.T) {
	c := NewSQLClassifier()
	allowed := []string{"SELECT 1", "WITH x AS (SELECT 1) SELECT * FROM x", "EXPLAIN SELECT 1"}
	for _, sql := range allowed {
		if !c.IsAllowed(sql) {
			t.Errorf("IsAllowed(%q) = false, want true", sql)
		}
	}
	blocked := []string{"INSERT INTO t VALUES (1)", "DROP TABLE t", "SELECT 1; SELECT 2"}
	for _, sql := range blocked {
		if c.IsAllowed(sql) {
			t.Errorf("IsAllowed(%q) = true, want false", sql)
		}
	}
}

func TestSQLClassifier_CountStatements(t *testing.T) {
	c := NewSQLClassifier()
	if c.CountStatements("SELECT 1") != 1 {
		t.Error("expected 1 statement")
	}
	if c.CountStatements("SELECT 1; SELECT 2") != 2 {
		t.Error("expected 2 statements")
	}
	if c.CountStatements("") != 0 {
		t.Error("expected 0 statements for empty")
	}
}
