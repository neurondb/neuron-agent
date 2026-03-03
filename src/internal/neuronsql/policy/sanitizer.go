/*-------------------------------------------------------------------------
 *
 * sanitizer.go
 *    SQL sanitizer: strip comments, normalize whitespace
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/policy/sanitizer.go
 *
 *-------------------------------------------------------------------------
 */

package policy

import (
	"regexp"
	"strings"
	"unicode"
)

/* SQLSanitizer sanitizes user SQL input */
type SQLSanitizer struct{}

/* NewSQLSanitizer creates a new sanitizer */
func NewSQLSanitizer() *SQLSanitizer {
	return &SQLSanitizer{}
}

/* Sanitize strips comments and normalizes whitespace; does not validate encoding */
func (s *SQLSanitizer) Sanitize(input string) string {
	if input == "" {
		return ""
	}
	out := stripCommentsFull(input)
	out = normalizeWhitespace(out)
	return strings.TrimSpace(out)
}

func stripCommentsFull(s string) string {
	/* Block comments */
	blockRe := regexp.MustCompile(`/\*[\s\S]*?\*/`)
	s = blockRe.ReplaceAllString(s, " ")
	/* Line comments (-- to EOL) */
	lines := strings.Split(s, "\n")
	for i, line := range lines {
		inLiteral := false
		var literal rune
		var j int
		for j = 0; j < len(line); j++ {
			r := rune(line[j])
			if inLiteral {
				if r == literal && (literal != '\'' || (j+1 < len(line) && line[j+1] != '\'')) {
					inLiteral = false
				}
				if literal == '\'' && r == '\\' && j+1 < len(line) {
					j++
				}
				continue
			}
			if (r == '\'' || r == '"') && (j == 0 || line[j-1] != '\\') {
				inLiteral = true
				literal = r
				continue
			}
			if r == '-' && j+1 < len(line) && line[j+1] == '-' {
				break
			}
		}
		lines[i] = line[:j]
	}
	return strings.Join(lines, "\n")
}

func normalizeWhitespace(s string) string {
	var b strings.Builder
	lastSpace := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			if !lastSpace {
				b.WriteRune(' ')
				lastSpace = true
			}
			continue
		}
		lastSpace = false
		b.WriteRune(r)
	}
	return b.String()
}
