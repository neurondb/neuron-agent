/*-------------------------------------------------------------------------
 *
 * runner.go
 *    Strict SQL runner: single statement, policy check, read-only, timeouts, row/size limits
 *
 * Used by NeuronSQL tools that execute SQL. No multi-statement; blocklist enforced
 * via PolicyEngine before execution. Audit policy_block on reject.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/runner.go
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
	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* ErrMultiStatement is returned when SQL contains multiple statements */
var ErrMultiStatement = fmt.Errorf("runner: only single statement allowed")

/* ErrPolicyBlocked is returned when policy check rejects the SQL */
type ErrPolicyBlocked struct {
	ReasonCode string
	ReasonText string
	Tokens     []string
}

func (e *ErrPolicyBlocked) Error() string {
	return fmt.Sprintf("policy blocked: %s (%s)", e.ReasonCode, e.ReasonText)
}

/* ErrInvalidSQLCharset is returned when SQL contains null bytes or control characters */
var ErrInvalidSQLCharset = fmt.Errorf("runner: sql must not contain null bytes or control characters")

/* ErrSQLTooLong is returned when SQL exceeds max length */
var ErrSQLTooLong = fmt.Errorf("runner: sql exceeds maximum length")

/* ValidateSQLCharsetAndLength rejects SQL with null bytes, control characters, or length > maxLen. maxLen 0 uses no length check. */
func ValidateSQLCharsetAndLength(sql string, maxLen int) error {
	for _, r := range sql {
		if r == 0 || (r < 32 && r != '\t' && r != '\n' && r != '\r') {
			return ErrInvalidSQLCharset
		}
	}
	if maxLen > 0 && len(sql) > maxLen {
		return ErrSQLTooLong
	}
	return nil
}

/* EnforceSingleStatement rejects SQL that contains semicolon as statement separator (outside literals) */
func EnforceSingleStatement(sql string) error {
	s := strings.TrimSpace(sql)
	/* Allow single semicolon at end */
	if strings.HasSuffix(s, ";") {
		s = s[:len(s)-1]
	}
	/* Reject if any semicolon remains (simple check; does not parse string literals) */
	if strings.Contains(s, ";") {
		return ErrMultiStatement
	}
	return nil
}

/* RunSingleReadOnly runs a single SELECT/EXPLAIN statement with policy check and limits.
 * Returns JSON array of rows (map[string]interface{} per row) or error.
 * Policy check is performed first; on block returns ErrPolicyBlocked.
 * maxRows and maxResultBytes cap the result; 0 means use SafeConnection defaults.
 */
func RunSingleReadOnly(ctx context.Context, conn *SafeConnection, requestID string, sql string, maxRows, maxResultBytes int) (string, error) {
	maxSQL := conn.Config().MaxSQLLength
	if maxSQL <= 0 {
		maxSQL = defaultMaxSQLLength
	}
	if err := ValidateSQLCharsetAndLength(sql, maxSQL); err != nil {
		return "", err
	}
	if err := EnforceSingleStatement(sql); err != nil {
		return "", err
	}
	decision, err := conn.Policy().Check(ctx, sql, neuronsql.PolicyContext{RequestID: requestID})
	if err != nil {
		return "", err
	}
	if !decision.Allowed {
		return "", &ErrPolicyBlocked{
			ReasonCode: decision.ReasonCode,
			ReasonText: decision.ReasonText,
			Tokens:     decision.BlockedTokens,
		}
	}
	if maxRows <= 0 {
		maxRows = conn.Config().MaxRows
	}
	if maxResultBytes <= 0 {
		maxResultBytes = conn.Config().MaxResultBytes
	}
	var rows []map[string]interface{}
	err = conn.RunReadOnly(ctx, requestID, func(tx *sqlx.Tx) error {
		rows, err = runQueryWithLimits(ctx, tx, sql, maxRows, maxResultBytes)
		return err
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

/* runQueryWithLimits runs the query and returns rows up to maxRows, truncating if total size exceeds maxResultBytes */
func runQueryWithLimits(ctx context.Context, tx *sqlx.Tx, sql string, maxRows, maxResultBytes int) ([]map[string]interface{}, error) {
	query := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(sql), ";"))
	upper := strings.ToUpper(query)
	if !strings.HasPrefix(upper, "SELECT") && !strings.HasPrefix(upper, "EXPLAIN") {
		return nil, fmt.Errorf("runner: only SELECT or EXPLAIN allowed")
	}
	/* Add LIMIT if not present (simple check) */
	if !strings.Contains(upper, "LIMIT") && strings.HasPrefix(upper, "SELECT") {
		query = query + fmt.Sprintf(" LIMIT %d", maxRows+1)
	}
	rowsi, err := tx.QueryxContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rowsi.Close()
	var out []map[string]interface{}
	var totalBytes int
	for rowsi.Next() && len(out) < maxRows && totalBytes < maxResultBytes {
		m := make(map[string]interface{})
		if err := rowsi.MapScan(m); err != nil {
			return nil, err
		}
		b, _ := json.Marshal(m)
		totalBytes += len(b)
		if totalBytes > maxResultBytes {
			break
		}
		out = append(out, m)
	}
	return out, rowsi.Err()
}
