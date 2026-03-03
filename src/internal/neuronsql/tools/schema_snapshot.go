/*-------------------------------------------------------------------------
 *
 * schema_snapshot.go
 *    tool_schema_snapshot: tables, columns, indexes, views, extensions, settings
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/schema_snapshot.go
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

/* SchemaSnapshotResult is the JSON result of tool_schema_snapshot */
type SchemaSnapshotResult struct {
	Tables       []TableInfo   `json:"tables"`
	Views        []ViewInfo    `json:"views"`
	Extensions   []string      `json:"extensions"`
	ServerVersion string      `json:"server_version"`
	Settings    []SettingInfo `json:"settings,omitempty"`
}

type TableInfo struct {
	Schema      string       `json:"schema"`
	Name        string       `json:"name"`
	Columns     []ColumnInfo `json:"columns"`
	PrimaryKey  []string     `json:"primary_key,omitempty"`
	ForeignKeys []FKInfo     `json:"foreign_keys,omitempty"`
	Indexes     []IndexInfo  `json:"indexes,omitempty"`
	RowEstimate int64        `json:"row_estimate,omitempty"`
	SizeBytes   int64        `json:"size_bytes,omitempty"`
}

type IndexInfo struct {
	Name    string   `json:"name"`
	Columns []string `json:"columns"`
	Unique  bool     `json:"unique"`
}

type ColumnInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type FKInfo struct {
	Columns    []string `json:"columns"`
	RefTable   string   `json:"ref_table"`
	RefColumns []string `json:"ref_columns"`
}

type ViewInfo struct {
	Schema string `json:"schema"`
	Name   string `json:"name"`
}

type SettingInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

/* RunSchemaSnapshot runs the schema snapshot tool and returns JSON */
func RunSchemaSnapshot(ctx context.Context, conn *SafeConnection, requestID string) (string, error) {
	var result SchemaSnapshotResult
	err := conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		/* Server version */
		var ver string
		if err := tx.GetContext(ctx, &ver, `SELECT current_setting('server_version')`); err != nil {
			return err
		}
		result.ServerVersion = ver

		/* Extensions */
		var ext []string
		if err := tx.SelectContext(ctx, &ext, `SELECT extname FROM pg_extension ORDER BY extname`); err != nil {
			return err
		}
		result.Extensions = ext

		/* User tables with columns */
		type colRow struct {
			TableSchema string `db:"table_schema"`
			TableName   string `db:"table_name"`
			ColumnName  string `db:"column_name"`
			DataType    string `db:"data_type"`
			IsNullable  string `db:"is_nullable"`
		}
		var cols []colRow
		if err := tx.SelectContext(ctx, &cols, `
			SELECT table_schema, table_name, column_name, data_type, is_nullable
			FROM information_schema.columns
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY table_schema, table_name, ordinal_position
		`); err != nil {
			return err
		}
		tableMap := make(map[string]*TableInfo)
		for _, c := range cols {
			key := c.TableSchema + "." + c.TableName
			if tableMap[key] == nil {
				tableMap[key] = &TableInfo{Schema: c.TableSchema, Name: c.TableName, Columns: nil}
			}
			tableMap[key].Columns = append(tableMap[key].Columns, ColumnInfo{
				Name:     c.ColumnName,
				Type:     c.DataType,
				Nullable: c.IsNullable == "YES",
			})
		}

		/* Primary keys */
		type pkRow struct {
			Schema string `db:"schema_name"`
			Table  string `db:"table_name"`
			Col    string `db:"column_name"`
		}
		var pks []pkRow
		_ = tx.SelectContext(ctx, &pks, `
			SELECT kcu.table_schema AS schema_name, kcu.table_name AS table_name, kcu.column_name
			FROM information_schema.table_constraints tc
			JOIN information_schema.key_column_usage kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema
			WHERE tc.constraint_type = 'PRIMARY KEY' AND tc.table_schema NOT IN ('pg_catalog', 'information_schema')
			ORDER BY kcu.table_schema, kcu.table_name, kcu.ordinal_position
		`)
		for _, p := range pks {
			key := p.Schema + "." + p.Table
			if t := tableMap[key]; t != nil {
				t.PrimaryKey = append(t.PrimaryKey, p.Col)
			}
		}

		/* Row estimates from pg_stat_user_tables */
		type statRow struct {
			Schema string `db:"schemaname"`
			Rel    string `db:"relname"`
			NLive  int64  `db:"n_live_tup"`
		}
		var stats []statRow
		_ = tx.SelectContext(ctx, &stats, `SELECT schemaname, relname, n_live_tup FROM pg_stat_user_tables`)
		for _, s := range stats {
			key := s.Schema + "." + s.Rel
			if t := tableMap[key]; t != nil {
				t.RowEstimate = s.NLive
			}
		}

		/* Indexes per table (name and columns from indexdef) */
		type idxRow struct {
			Schema  string `db:"schemaname"`
			Table   string `db:"tablename"`
			Index   string `db:"indexname"`
			Indexdef string `db:"indexdef"`
		}
		var idxRows []idxRow
		_ = tx.SelectContext(ctx, &idxRows, `
			SELECT schemaname, tablename, indexname, indexdef
			FROM pg_indexes
			WHERE schemaname NOT IN ('pg_catalog', 'information_schema')
		`)
		for _, idx := range idxRows {
			key := idx.Schema + "." + idx.Table
			if t := tableMap[key]; t != nil {
				cols := parseIndexColumns(idx.Indexdef)
				t.Indexes = append(t.Indexes, IndexInfo{Name: idx.Index, Columns: cols, Unique: false})
			}
		}

		/* Table sizes (pg_total_relation_size) */
		for _, t := range tableMap {
			var size int64
			qual := t.Schema + "." + t.Name
			_ = tx.GetContext(ctx, &size, `SELECT COALESCE(pg_total_relation_size($1::regclass), 0)`, qual)
			t.SizeBytes = size
		}

		for _, t := range tableMap {
			result.Tables = append(result.Tables, *t)
		}

		/* Views */
		type viewRow struct {
			Schema string `db:"table_schema"`
			Name   string `db:"table_name"`
		}
		var views []viewRow
		if err := tx.SelectContext(ctx, &views, `
			SELECT table_schema, table_name FROM information_schema.views
			WHERE table_schema NOT IN ('pg_catalog', 'information_schema')
		`); err != nil {
			return err
		}
		for _, v := range views {
			result.Views = append(result.Views, ViewInfo{Schema: v.Schema, Name: v.Name})
		}

		/* Relevant settings */
		var settings []struct {
			Name  string `db:"name"`
			Value string `db:"setting"`
		}
		_ = tx.SelectContext(ctx, &settings, `
			SELECT name, setting FROM pg_settings
			WHERE name IN ('work_mem', 'shared_buffers', 'random_page_cost', 'effective_cache_size', 'max_parallel_workers_per_gather')
		`)
		for _, s := range settings {
			result.Settings = append(result.Settings, SettingInfo{Name: s.Name, Value: s.Value})
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

/* SchemaSnapshotTool implements tools.ToolHandler for schema_snapshot */
type SchemaSnapshotTool struct {
	Factory ConnectionFactory
	Policy  *policy.PolicyEngineImpl
}

/* Execute implements tools.ToolHandler (registry passes *db.Tool and args) */
func (t *SchemaSnapshotTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("schema_snapshot: db_dsn required")
	}
	requestID, _ := args["request_id"].(string)
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return RunSchemaSnapshot(ctx, conn, requestID)
}

/* Validate implements tools.ToolHandler */
func (t *SchemaSnapshotTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	return nil
}

/* parseIndexColumns extracts column names from indexdef (e.g. "CREATE INDEX ... ON t (a, b)" -> ["a", "b"]) */
func parseIndexColumns(indexdef string) []string {
	start := strings.Index(indexdef, "(")
	end := strings.LastIndex(indexdef, ")")
	if start < 0 || end <= start {
		return nil
	}
	inner := strings.TrimSpace(indexdef[start+1 : end])
	if inner == "" {
		return nil
	}
	parts := strings.Split(inner, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		col := strings.TrimSpace(p)
		if idx := strings.Index(col, " "); idx > 0 {
			col = col[:idx]
		}
		if col != "" {
			out = append(out, col)
		}
	}
	return out
}
