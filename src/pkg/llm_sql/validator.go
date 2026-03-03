package llm_sql

import (
	"fmt"
	"regexp"
	"strings"
)

// SQLValidator validates SQL queries
type SQLValidator struct {
	dialect string
}

// NewSQLValidator creates a new SQL validator
func NewSQLValidator(dialect string) *SQLValidator {
	return &SQLValidator{
		dialect: dialect,
	}
}

// ValidateSQL validates a SQL query for the configured dialect
func (v *SQLValidator) ValidateSQL(sql, dialect string) error {
	if sql == "" {
		return fmt.Errorf("SQL query is empty")
	}
	
	// Basic SQL injection check
	if err := v.checkSQLInjection(sql); err != nil {
		return err
	}
	
	// Dialect-specific validation
	switch dialect {
	case "postgresql":
		return v.validatePostgreSQL(sql)
	case "mysql":
		return v.validateMySQL(sql)
	default:
		return fmt.Errorf("unsupported dialect: %s", dialect)
	}
}

func (v *SQLValidator) checkSQLInjection(sql string) error {
	// Check for common SQL injection patterns
	dangerousPatterns := []string{
		`(?i);\s*DROP`,
		`(?i);\s*DELETE\s+FROM`,
		`(?i);\s*TRUNCATE`,
		`(?i)UNION.*SELECT`,
		`(?i)--.*password`,
		`(?i)/\*.*\*/`,
	}
	
	for _, pattern := range dangerousPatterns {
		matched, _ := regexp.MatchString(pattern, sql)
		if matched {
			return fmt.Errorf("potential SQL injection detected")
		}
	}
	
	return nil
}

func (v *SQLValidator) validatePostgreSQL(sql string) error {
	sql = strings.TrimSpace(sql)
	
	// Check for basic syntax
	if !strings.HasSuffix(sql, ";") {
		return fmt.Errorf("PostgreSQL query should end with semicolon")
	}
	
	// Check for PostgreSQL-specific syntax
	upperSQL := strings.ToUpper(sql)
	
	// Validate common statement types
	validStarts := []string{"SELECT", "INSERT", "UPDATE", "DELETE", "CREATE", "ALTER", "DROP", "WITH"}
	hasValidStart := false
	for _, start := range validStarts {
		if strings.HasPrefix(upperSQL, start) {
			hasValidStart = true
			break
		}
	}
	
	if !hasValidStart {
		return fmt.Errorf("query must start with a valid SQL statement")
	}
	
	return nil
}

func (v *SQLValidator) validateMySQL(sql string) error {
	sql = strings.TrimSpace(sql)
	
	// Check for basic syntax
	if !strings.HasSuffix(sql, ";") {
		return fmt.Errorf("MySQL query should end with semicolon")
	}
	
	// MySQL-specific validation
	upperSQL := strings.ToUpper(sql)
	
	// Check for MySQL-specific keywords used incorrectly
	if strings.Contains(upperSQL, "RETURNING") {
		return fmt.Errorf("RETURNING clause is not supported in MySQL (use PostgreSQL)")
	}
	
	return nil
}

// SanitizeSQL removes potentially dangerous elements from SQL
func (v *SQLValidator) SanitizeSQL(sql string) string {
	// Remove comments
	sql = regexp.MustCompile(`--.*$`).ReplaceAllString(sql, "")
	sql = regexp.MustCompile(`/\*.*?\*/`).ReplaceAllString(sql, "")
	
	// Remove extra whitespace
	sql = regexp.MustCompile(`\s+`).ReplaceAllString(sql, " ")
	
	return strings.TrimSpace(sql)
}
