/*-------------------------------------------------------------------------
 *
 * table_profile.go
 *    tool_table_profile: pg_stat_user_tables, bloat signals, vacuum
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/table_profile.go
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

/* TableProfileResult is the JSON result of tool_table_profile */
type TableProfileResult struct {
	Tables []TableStatRow `json:"tables"`
}

type TableStatRow struct {
	Schema      string  `json:"schema" db:"schema"`
	RelName     string  `json:"relname" db:"relname"`
	SeqScan     int64   `json:"seq_scan" db:"seq_scan"`
	IdxScan     int64   `json:"idx_scan" db:"idx_scan"`
	NLiveTup    int64   `json:"n_live_tup" db:"n_live_tup"`
	NDeadTup    int64   `json:"n_dead_tup" db:"n_dead_tup"`
	LastVacuum  *string `json:"last_vacuum,omitempty" db:"last_vacuum"`
	LastAnalyze *string `json:"last_analyze,omitempty" db:"last_analyze"`
}

/* RunTableProfile returns table stats for the given table or all user tables */
func RunTableProfile(ctx context.Context, conn *SafeConnection, requestID string, table string) (string, error) {
	var result TableProfileResult
	err := conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		query := `
			SELECT schemaname AS schema, relname, seq_scan, idx_scan, n_live_tup, n_dead_tup,
			       last_vacuum::text, last_analyze::text
			FROM pg_stat_user_tables
		`
		args := []interface{}{}
		if table != "" {
			query += ` WHERE relname = $1`
			args = append(args, table)
		}
		query += ` ORDER BY schemaname, relname`
		var rows []TableStatRow
		if err := tx.SelectContext(ctx, &rows, query, args...); err != nil {
			return err
		}
		result.Tables = rows
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

/* TableProfileTool implements tools.ToolHandler */
type TableProfileTool struct {
	Factory ConnectionFactory
	Policy  *policy.PolicyEngineImpl
}

func (t *TableProfileTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("table_profile: db_dsn required")
	}
	table, _ := args["table"].(string)
	requestID, _ := args["request_id"].(string)
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return RunTableProfile(ctx, conn, requestID, table)
}

func (t *TableProfileTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	return nil
}
