/*-------------------------------------------------------------------------
 *
 * index_profile.go
 *    tool_index_profile: index definitions, sizes, usage
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/index_profile.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
)

/* IndexProfileResult is the JSON result of tool_index_profile */
type IndexProfileResult struct {
	Indexes []IndexStatRow `json:"indexes"`
}

type IndexStatRow struct {
	Schema      string `json:"schema" db:"schema"`
	TableName   string `json:"table_name" db:"table_name"`
	IndexName   string `json:"index_name" db:"index_name"`
	IndexDef    string `json:"index_def" db:"index_def"`
	IdxScan     int64  `json:"idx_scan" db:"idx_scan"`
	IdxTupRead  int64  `json:"idx_tup_read" db:"idx_tup_read"`
	IdxTupFetch int64  `json:"idx_tup_fetch" db:"idx_tup_fetch"`
}

/* RunIndexProfile returns index stats for the given table or all */
func RunIndexProfile(ctx context.Context, conn *SafeConnection, requestID string, table string) (string, error) {
	var result IndexProfileResult
	err := conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		query := `
			SELECT s.schemaname AS schema, s.relname AS table_name, s.indexrelname AS index_name,
			       pg_get_indexdef(i.indexrelid) AS index_def,
			       s.idx_scan, s.idx_tup_read, s.idx_tup_fetch
			FROM pg_stat_user_indexes s
			JOIN pg_index i ON i.indexrelid = s.indexrelid
		`
		args := []interface{}{}
		if table != "" {
			query += ` WHERE s.relname = $1`
			args = append(args, table)
		}
		query += ` ORDER BY s.schemaname, s.relname, s.indexrelname`
		var rows []IndexStatRow
		if err := tx.SelectContext(ctx, &rows, query, args...); err != nil {
			return err
		}
		result.Indexes = rows
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

/* IndexProfileTool implements tools.ToolHandler */
type IndexProfileTool struct {
	Factory ConnectionFactory
	Policy  *policy.PolicyEngineImpl
}

func (t *IndexProfileTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("index_profile: db_dsn required")
	}
	table, _ := args["table"].(string)
	requestID, _ := args["request_id"].(string)
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return RunIndexProfile(ctx, conn, requestID, table)
}

func (t *IndexProfileTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	return nil
}
