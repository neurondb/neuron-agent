/*-------------------------------------------------------------------------
 *
 * intent.go
 *    User intent parsing (mode, target tables, constraints)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/prompting/intent.go
 *
 *-------------------------------------------------------------------------
 */

package prompting

import "strings"

/* Intent holds parsed user intent for generation */
type Intent struct {
	Mode    string
	Question string
	Tables  []string
}

/* ParseIntent extracts a simple intent from the user question; mode is set by the API */
func ParseIntent(question string, mode string) Intent {
	return Intent{
		Mode:     mode,
		Question: strings.TrimSpace(question),
		Tables:   nil,
	}
}
