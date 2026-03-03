/*-------------------------------------------------------------------------
 *
 * validate_sql.go
 *    tool_validate_sql: parse/validate SQL against policy and EXPLAIN
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/validate_sql.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* ValidateSQLResult is the JSON result of tool_validate_sql */
type ValidateSQLResult struct {
	Valid             bool     `json:"valid"`
	Errors            []string `json:"errors,omitempty"`
	ReferencedTables  []string `json:"referenced_tables,omitempty"`
	ReferencedColumns []string `json:"referenced_columns,omitempty"`
	StatementClass    string   `json:"statement_class,omitempty"`
	RiskyPatterns     []string `json:"risky_patterns,omitempty"`
	Violations        []string `json:"violations,omitempty"` /* e.g. missing table/column from schema */
}

/* RunValidateSQL runs the validate_sql tool. If schemaSnapshotJSON is non-empty, validates referenced tables exist in snapshot. */
func RunValidateSQL(ctx context.Context, conn *SafeConnection, requestID string, sql string) (string, error) {
	return RunValidateSQLWithSnapshot(ctx, conn, requestID, sql, "")
}

/* RunValidateSQLWithSnapshot runs validate_sql with optional schema snapshot for adherence check */
func RunValidateSQLWithSnapshot(ctx context.Context, conn *SafeConnection, requestID string, sql string, schemaSnapshotJSON string) (string, error) {
	decision, err := conn.Policy().Check(ctx, sql, neuronsql.PolicyContext{RequestID: requestID})
	if err != nil {
		return "", err
	}
	if !decision.Allowed {
		b, _ := json.Marshal(ValidateSQLResult{
			Valid:          false,
			Errors:         []string{decision.Reason},
			StatementClass: decision.StatementClass,
			Violations:     []string{"policy:" + decision.Reason},
		})
		return string(b), nil
	}

	var result ValidateSQLResult
	result.StatementClass = decision.StatementClass
	/* v1: allow only SELECT and EXPLAIN */
	if c := strings.ToUpper(strings.TrimSpace(sql)); !strings.HasPrefix(c, "SELECT") && !strings.HasPrefix(c, "WITH") && !strings.HasPrefix(c, "EXPLAIN") {
		result.Valid = false
		result.Errors = append(result.Errors, "only SELECT and EXPLAIN allowed in v1")
		result.Violations = append(result.Violations, "statement_class: only SELECT/EXPLAIN allowed")
		b, _ := json.Marshal(result)
		return string(b), nil
	}

	if schemaSnapshotJSON != "" {
		tables := extractReferencedTables(sql)
		result.ReferencedTables = tables
		valid, violations := validateTablesAgainstSnapshot(schemaSnapshotJSON, tables)
		if !valid {
			result.Valid = false
			result.Violations = append(result.Violations, violations...)
			result.Errors = append(result.Errors, violations...)
		}
	}

	err = conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		_, err := tx.ExecContext(ctx, `EXPLAIN (FORMAT JSON) `+sql)
		if err != nil {
			result.Valid = false
			result.Errors = append(result.Errors, err.Error())
			if schemaSnapshotJSON != "" {
				result.Violations = append(result.Violations, "explain:"+err.Error())
			}
			return nil
		}
		if len(result.Violations) == 0 {
			result.Valid = true
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

var fromJoinRe = regexp.MustCompile(`(?i)\b(?:FROM|JOIN)\s+([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)?)`)

func extractReferencedTables(sql string) []string {
	seen := make(map[string]bool)
	var out []string
	for _, m := range fromJoinRe.FindAllStringSubmatch(sql, -1) {
		if len(m) > 1 {
			t := strings.TrimSpace(m[1])
			if t != "" && !seen[t] {
				seen[t] = true
				out = append(out, t)
			}
		}
	}
	return out
}

func validateTablesAgainstSnapshot(schemaJSON string, tables []string) (valid bool, violations []string) {
	var snap SchemaSnapshotResult
	if err := json.Unmarshal([]byte(schemaJSON), &snap); err != nil {
		return true, nil
	}
	allowed := make(map[string]bool)
	for _, t := range snap.Tables {
		allowed[t.Schema+"."+t.Name] = true
		allowed[t.Name] = true
	}
	for _, v := range snap.Views {
		allowed[v.Schema+"."+v.Name] = true
		allowed[v.Name] = true
	}
	for _, t := range tables {
		norm := t
		if !strings.Contains(norm, ".") {
			norm = "public." + norm
		}
		if !allowed[t] && !allowed[norm] {
			violations = append(violations, "missing table or view: "+t)
		}
	}
	return len(violations) == 0, violations
}

/* ValidateSQLTool implements tools.ToolHandler */
type ValidateSQLTool struct {
	Factory ConnectionFactory
	Policy  *policy.PolicyEngineImpl
}

func (t *ValidateSQLTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("validate_sql: db_dsn required")
	}
	sql, _ := args["sql"].(string)
	if sql == "" {
		return "", fmt.Errorf("validate_sql: sql required")
	}
	requestID, _ := args["request_id"].(string)
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return RunValidateSQL(ctx, conn, requestID, sql)
}

func (t *ValidateSQLTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	if _, ok := args["sql"]; !ok {
		return fmt.Errorf("sql required")
	}
	return nil
}
