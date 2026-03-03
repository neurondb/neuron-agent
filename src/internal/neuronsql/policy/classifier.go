/*-------------------------------------------------------------------------
 *
 * classifier.go
 *    SQL statement classifier for NeuronSQL policy (Go-native, no external parser)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/policy/classifier.go
 *
 *-------------------------------------------------------------------------
 */

package policy

import (
	"regexp"
	"strings"
	"unicode"
)

/* StatementClass is the classified type of SQL statement */
const (
	ClassSelect   = "select"
	ClassWith     = "with"
	ClassExplain  = "explain"
	ClassBlocked  = "blocked"
	ClassDDL      = "ddl"
	ClassDML      = "dml"
	ClassDCL      = "dcl"
	ClassUnknown  = "unknown"
)

/* SQLClassifier classifies SQL without an external parser */
type SQLClassifier struct{}

/* NewSQLClassifier creates a new classifier */
func NewSQLClassifier() *SQLClassifier {
	return &SQLClassifier{}
}

/* Classify returns the statement class of the first statement in sql */
func (c *SQLClassifier) Classify(sql string) string {
	s := strings.TrimSpace(sql)
	if s == "" {
		return ClassBlocked
	}
	s = stripComments(s)
	s = strings.TrimSpace(s)
	if s == "" {
		return ClassBlocked
	}

	/* Multi-statement: reject if more than one semicolon (allowing trailing) */
	trimmed := strings.TrimRight(s, " \t\n\r;")
	if strings.Contains(trimmed, ";") {
		return ClassBlocked
	}

	/* Normalize for keyword check: first token(s) */
	upper := strings.ToUpper(s)
	firstWord := firstWord(upper)

	switch firstWord {
	case "SELECT":
		return ClassSelect
	case "WITH":
		return ClassWith
	case "EXPLAIN":
		return ClassExplain
	case "INSERT", "UPDATE", "DELETE":
		return ClassDML
	case "DROP", "TRUNCATE", "ALTER", "CREATE", "GRANT", "REVOKE":
		return ClassDDL
	default:
		return ClassUnknown
	}
}

/* IsAllowed returns true only for SELECT, WITH, EXPLAIN */
func (c *SQLClassifier) IsAllowed(sql string) bool {
	class := c.Classify(sql)
	return class == ClassSelect || class == ClassWith || class == ClassExplain
}

/* CountStatements returns the number of statements (semicolon-separated); 0 if empty */
func (c *SQLClassifier) CountStatements(sql string) int {
	s := stripComments(sql)
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	s = strings.TrimRight(s, " \t\n\r;")
	if s == "" {
		return 0
	}
	n := 1
	for _, r := range s {
		if r == ';' {
			n++
		}
	}
	return n
}

func firstWord(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsSpace(r) || r == '(' {
			break
		}
		b.WriteRune(r)
	}
	return b.String()
}

/* stripComments removes single-line (--) and block (/* *\/) comments */
func stripComments(s string) string {
	/* Block comments */
	blockRe := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	s = blockRe.ReplaceAllString(s, " ")
	/* Line comments */
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		if idx := strings.Index(line, "--"); idx >= 0 {
			lines[i] = line[:idx]
		}
	}
	return strings.Join(lines, "\n")
}
