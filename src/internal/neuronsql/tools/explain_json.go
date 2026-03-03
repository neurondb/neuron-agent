/*-------------------------------------------------------------------------
 *
 * explain_json.go
 *    tool_explain_json: EXPLAIN (FORMAT JSON) without ANALYZE
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/explain_json.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* RunExplainJSON runs EXPLAIN (FORMAT JSON) on the SQL and returns the plan JSON */
func RunExplainJSON(ctx context.Context, conn *SafeConnection, requestID string, sql string) (string, error) {
	decision, err := conn.Policy().Check(ctx, sql, neuronsql.PolicyContext{})
	if err != nil {
		return "", err
	}
	if !decision.Allowed {
		return "", fmt.Errorf("explain_json: sql not allowed: %s", decision.Reason)
	}
	var planJSON string
	err = conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		return tx.QueryRowContext(ctx, `EXPLAIN (FORMAT JSON) `+sql).Scan(&planJSON)
	})
	if err != nil {
		return "", err
	}
	return planJSON, nil
}

/* ExplainJSONTool implements tools.ToolHandler */
type ExplainJSONTool struct {
	Factory ConnectionFactory
	Policy  *policy.PolicyEngineImpl
}

func (t *ExplainJSONTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	dsn, _ := args["db_dsn"].(string)
	if dsn == "" {
		return "", fmt.Errorf("explain_json: db_dsn required")
	}
	sql, _ := args["sql"].(string)
	if sql == "" {
		return "", fmt.Errorf("explain_json: sql required")
	}
	requestID, _ := args["request_id"].(string)
	conn, err := t.Factory(ctx, dsn)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return RunExplainJSON(ctx, conn, requestID, sql)
}

func (t *ExplainJSONTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	if _, ok := args["db_dsn"]; !ok {
		return fmt.Errorf("db_dsn required")
	}
	if _, ok := args["sql"]; !ok {
		return fmt.Errorf("sql required")
	}
	return nil
}
