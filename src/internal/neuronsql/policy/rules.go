/*-------------------------------------------------------------------------
 *
 * rules.go
 *    Policy rules: blocklist, denylist, allowlist for NeuronSQL
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/policy/rules.go
 *
 *-------------------------------------------------------------------------
 */

package policy

import (
	"regexp"
	"strings"
)

/* PolicyRules holds configurable blocklist and denylist */
type PolicyRules struct {
	SensitiveTables []string
	BlockedKeywords []string
	BlockedFuncs    []string
}

/* DefaultPolicyRules returns the default NeuronSQL v1 rules */
func DefaultPolicyRules() *PolicyRules {
	return &PolicyRules{
		SensitiveTables: nil,
		BlockedKeywords: []string{
			"DROP", "TRUNCATE", "ALTER SYSTEM", "GRANT", "REVOKE",
			"CREATE FUNCTION", "CREATE OR REPLACE FUNCTION",
			"DO $$", "DO $", "DBLINK_CONNECT", "DBLINK_OPEN",
		},
		BlockedFuncs: []string{
			"pg_read_file", "pg_ls_dir", "lo_export", "lo_import",
			"pg_read_binary_file", "pg_read_file",
		},
	}
}

/* NewPolicyRules creates rules with optional sensitive table list */
func NewPolicyRules(sensitiveTables []string) *PolicyRules {
	r := DefaultPolicyRules()
	r.SensitiveTables = sensitiveTables
	return r
}

/* CheckBlocklist returns a reason string if the SQL is blocked, else "" */
func (r *PolicyRules) CheckBlocklist(sql string) string {
	reason, _, _ := r.CheckBlocklistDetailed(sql)
	return reason
}

/* CheckBlocklistDetailed returns reason string, reason code, and blocked tokens for policy decision */
func (r *PolicyRules) CheckBlocklistDetailed(sql string) (reason string, reasonCode string, blockedTokens []string) {
	upper := strings.ToUpper(sql)
	sqlLower := strings.ToLower(sql)

	for _, kw := range r.BlockedKeywords {
		if strings.Contains(upper, kw) {
			kwLower := strings.ToLower(kw)
			return "blocked_keyword: " + kwLower, "blocked_keyword", []string{kwLower}
		}
	}

	for _, fn := range r.BlockedFuncs {
		if strings.Contains(sqlLower, fn) {
			return "blocked_function: " + fn, "blocked_function", []string{fn}
		}
	}

	/* COPY ... PROGRAM */
	if copyProgramRe.MatchString(upper) {
		return "blocked: COPY ... PROGRAM", "blocked_keyword", []string{"COPY ... PROGRAM"}
	}

	/* DO $$ ... $$ blocks */
	if doBlockRe.MatchString(sqlLower) {
		return "blocked: DO $$ block", "blocked_keyword", []string{"DO $$"}
	}

	return "", "", nil
}

var (
	copyProgramRe = regexp.MustCompile(`COPY\s+.*(FROM|TO)\s+PROGRAM`)
	doBlockRe     = regexp.MustCompile(`do\s+\$`)
)

/* IsSensitiveTable returns true if table (or schema.table) is in the denylist */
func (r *PolicyRules) IsSensitiveTable(tableOrSchemaTable string) bool {
	t := strings.TrimSpace(tableOrSchemaTable)
	for _, s := range r.SensitiveTables {
		if strings.EqualFold(t, s) {
			return true
		}
		/* Also check bare table name */
		if idx := strings.Index(t, "."); idx >= 0 {
			if strings.EqualFold(t[idx+1:], s) {
				return true
			}
		}
	}
	return false
}
