/*-------------------------------------------------------------------------
 *
 * engine.go
 *    PolicyEngine implementation composing classifier, sanitizer, rules
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/policy/engine.go
 *
 *-------------------------------------------------------------------------
 */

package policy

import (
	"context"

	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* PolicyEngineImpl implements neuronsql.PolicyEngine */
type PolicyEngineImpl struct {
	classifier *SQLClassifier
	sanitizer  *SQLSanitizer
	rules      *PolicyRules
}

/* NewPolicyEngineImpl creates a new policy engine with default rules */
func NewPolicyEngineImpl(rules *PolicyRules) *PolicyEngineImpl {
	if rules == nil {
		rules = DefaultPolicyRules()
	}
	return &PolicyEngineImpl{
		classifier: NewSQLClassifier(),
		sanitizer:  NewSQLSanitizer(),
		rules:      rules,
	}
}

/* Check implements neuronsql.PolicyEngine */
func (e *PolicyEngineImpl) Check(ctx context.Context, sql string, ctxIn neuronsql.PolicyContext) (*neuronsql.PolicyDecision, error) {
	sanitized := e.sanitizer.Sanitize(sql)
	if sanitized == "" {
		return &neuronsql.PolicyDecision{
			Allowed:        false,
			Reason:         "empty_sql",
			ReasonCode:     "empty_sql",
			ReasonText:     "SQL is empty after sanitization",
			StatementClass: ClassBlocked,
		}, nil
	}

	if reason, reasonCode, blockedTokens := e.rules.CheckBlocklistDetailed(sanitized); reason != "" {
		return &neuronsql.PolicyDecision{
			Allowed:        false,
			Reason:         reason,
			ReasonCode:     reasonCode,
			ReasonText:     reason,
			BlockedTokens:  blockedTokens,
			StatementClass: ClassBlocked,
		}, nil
	}

	class := e.classifier.Classify(sanitized)
	allowed := e.classifier.IsAllowed(sanitized)

	reason := ""
	reasonCode := ""
	if !allowed {
		reason = "statement_class_not_allowed: " + class
		reasonCode = "statement_class_not_allowed"
	}

	return &neuronsql.PolicyDecision{
		Allowed:        allowed,
		Reason:         reason,
		ReasonCode:     reasonCode,
		ReasonText:     reason,
		StatementClass: class,
	}, nil
}

/* Sanitize implements neuronsql.PolicyEngine */
func (e *PolicyEngineImpl) Sanitize(input string) string {
	return e.sanitizer.Sanitize(input)
}
