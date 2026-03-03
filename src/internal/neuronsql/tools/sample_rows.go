/*-------------------------------------------------------------------------
 *
 * sample_rows.go
 *    tool_sample_rows: sample rows with limit and optional denylist
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/sample_rows.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
)

const defaultSampleLimit = 10
const maxSampleLimit = 100

/* RunSampleRows returns sample rows from the table; respects sensitive denylist */
func RunSampleRows(ctx context.Context, conn *SafeConnection, requestID string, table string, limit int, denylist []string) (string, error) {
	if limit <= 0 {
		limit = defaultSampleLimit
	}
	if limit > maxSampleLimit {
		limit = maxSampleLimit
	}
	table = strings.TrimSpace(table)
	if table == "" {
		return "", fmt.Errorf("sample_rows: table required")
	}
	for _, d := range denylist {
		if strings.EqualFold(table, d) {
			return "", fmt.Errorf("sample_rows: table %q is on sensitive denylist", table)
		}
	}
	var rows []map[string]interface{}
	err := conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		query := `SELECT * FROM ` + quoteTable(table) + ` TABLESAMPLE SYSTEM (1) LIMIT $1`
		r, err := tx.QueryContext(ctx, query, limit)
		if err != nil {
			return err
		}
		defer r.Close()
		cols, _ := r.Columns()
		for r.Next() {
			vals := make([]interface{}, len(cols))
			ptrs := make([]interface{}, len(cols))
			for i := range vals {
				ptrs[i] = &vals[i]
			}
			if err := r.Scan(ptrs...); err != nil {
				return err
			}
			row := make(map[string]interface{})
			for i, c := range cols {
				row[c] = vals[i]
			}
			rows = append(rows, row)
		}
		return r.Err()
	})
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(rows)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func quoteTable(name string) string {
	if strings.Index(name, ".") >= 0 {
		parts := strings.SplitN(name, ".", 2)
		return `"` + strings.ReplaceAll(parts[0], `"`, `""`) + `"."` + strings.ReplaceAll(parts[1], `"`, `""`) + `"`
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}

/* SampleRowsTool implements tools.ToolHandler */
type SampleRowsTool struct {
	Factory      ConnectionFactory
	Policy       *policy.PolicyEngineImpl
	SensitiveTables []string
}

func (t *SampleRowsTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("sample_rows: db_dsn required")
	}
	table, _ := args["table"].(string)
	requestID, _ := args["request_id"].(string)
	limit := defaultSampleLimit
	if l, ok := args["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return RunSampleRows(ctx, conn, requestID, table, limit, t.SensitiveTables)
}

func (t *SampleRowsTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	if _, ok := args["table"]; !ok {
		return fmt.Errorf("table required")
	}
	return nil
}
