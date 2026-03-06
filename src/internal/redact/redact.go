/*-------------------------------------------------------------------------
 * redact.go
 *    PII-safe logging with configurable redaction rules.
 *-------------------------------------------------------------------------*/

package redact

import (
	"regexp"
	"sync"
)

// Rule is a single redaction rule (pattern and replacement).
type Rule struct {
	Pattern *regexp.Regexp
	Repl    string
}

// RedactConfig holds redaction rules (e.g. email, phone, SSN, credit card).
type RedactConfig struct {
	mu    sync.RWMutex
	rules []Rule
}

// NewRedactConfig creates a config with optional default rules.
func NewRedactConfig(useDefaults bool) *RedactConfig {
	c := &RedactConfig{}
	if useDefaults {
		c.AddRule(regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`), "[EMAIL]")
		c.AddRule(regexp.MustCompile(`\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`), "[PHONE]")
		c.AddRule(regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`), "[SSN]")
		c.AddRule(regexp.MustCompile(`\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`), "[CARD]")
	}
	return c
}

// AddRule adds a redaction rule.
func (c *RedactConfig) AddRule(pattern *regexp.Regexp, replacement string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rules = append(c.rules, Rule{Pattern: pattern, Repl: replacement})
}

// Redact applies all rules to s and returns the redacted string.
func (c *RedactConfig) Redact(s string) string {
	c.mu.RLock()
	rules := make([]Rule, len(c.rules))
	copy(rules, c.rules)
	c.mu.RUnlock()
	for _, r := range rules {
		s = r.Pattern.ReplaceAllString(s, r.Repl)
	}
	return s
}
