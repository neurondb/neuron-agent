/*-------------------------------------------------------------------------
 * step_evaluator.go
 *    Step result evaluation, failure classification, and repair decisions.
 *-------------------------------------------------------------------------*/

package agent

import (
	"strings"
)

// EvaluationStatus is the result of evaluating a step.
type EvaluationStatus string

const (
	EvalSuccess        EvaluationStatus = "success"
	EvalPartialSuccess EvaluationStatus = "partial_success"
	EvalFailure        EvaluationStatus = "failure"
	EvalNeedsRepair    EvaluationStatus = "needs_repair"
)

// FailureClass is the plan-defined failure taxonomy.
type FailureClass string

const (
	FailureToolValidation  FailureClass = "tool_validation_error"
	FailureToolTimeout     FailureClass = "tool_timeout"
	FailureToolAuth        FailureClass = "tool_auth_failure"
	FailureModelMalformed  FailureClass = "model_malformed_output"
	FailureRetrievalMiss   FailureClass = "retrieval_miss"
	FailureContextOverflow FailureClass = "context_overflow"
	FailureDependency      FailureClass = "dependency_failure"
	FailurePolicyViolation FailureClass = "policy_violation"
	FailureEmptyResult     FailureClass = "empty_result"
	FailureContradictory   FailureClass = "contradictory_evidence"
	FailureMaxRetries      FailureClass = "max_retries_exceeded"
	FailureUnknown         FailureClass = "unknown_failure"
)

// EvaluationResult is the result of StepEvaluator.Evaluate.
type EvaluationResult struct {
	Status        EvaluationStatus
	FailureClass  FailureClass
	Retryable     bool
	ShouldUpdateMemory bool
	Answer        string
	Confidence    float64
	RepairAction  string
}

// StepEvaluator evaluates step outcomes and classifies success/failure.
type StepEvaluator interface {
	Evaluate(stepOutput map[string]interface{}, stepState string) (*EvaluationResult, error)
}

// DefaultStepEvaluator implements StepEvaluator with plan taxonomy.
type DefaultStepEvaluator struct{}

// NewDefaultStepEvaluator returns a StepEvaluator that classifies by output and state.
func NewDefaultStepEvaluator() *DefaultStepEvaluator {
	return &DefaultStepEvaluator{}
}

// Evaluate classifies the step result and sets retryability.
func (e *DefaultStepEvaluator) Evaluate(stepOutput map[string]interface{}, stepState string) (*EvaluationResult, error) {
	res := &EvaluationResult{Confidence: 0.5}
	if stepState == "completed" {
		res.Status = EvalSuccess
		res.ShouldUpdateMemory = true
		if a, ok := stepOutput["content"].(string); ok {
			res.Answer = a
		}
		return res, nil
	}
	if stepState == "failed" {
		res.Status = EvalFailure
		msg, _ := stepOutput["error"].(string)
		res.FailureClass = classifyFailure(msg, stepOutput)
		res.Retryable = isRetryable(res.FailureClass)
		return res, nil
	}
	res.Status = EvalNeedsRepair
	res.FailureClass = FailureUnknown
	res.Retryable = true
	return res, nil
}

func classifyFailure(msg string, out map[string]interface{}) FailureClass {
	if msg == "" {
		if e, ok := out["error_message"].(string); ok {
			msg = e
		}
	}
	msgLower := strings.ToLower(msg)
	switch {
	case containsStr(msgLower, "validation", "schema", "invalid argument"):
		return FailureToolValidation
	case containsStr(msgLower, "timeout", "deadline"):
		return FailureToolTimeout
	case containsStr(msgLower, "401", "403", "unauthorized", "forbidden"):
		return FailureToolAuth
	case containsStr(msgLower, "json", "parse", "malformed"):
		return FailureModelMalformed
	case containsStr(msgLower, "no results", "retrieval", "empty"):
		return FailureRetrievalMiss
	case containsStr(msgLower, "context", "token", "overflow"):
		return FailureContextOverflow
	case containsStr(msgLower, "unavailable", "connection", "dependency"):
		return FailureDependency
	case containsStr(msgLower, "policy", "not permitted"):
		return FailurePolicyViolation
	case containsStr(msgLower, "empty", "null result"):
		return FailureEmptyResult
	default:
		return FailureUnknown
	}
}

func containsStr(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func isRetryable(f FailureClass) bool {
	switch f {
	case FailureToolValidation, FailureToolTimeout, FailureModelMalformed,
		FailureRetrievalMiss, FailureDependency, FailureEmptyResult:
		return true
	default:
		return false
	}
}

// RepairAction is the action to take after a failure.
type RepairAction struct {
	Type         string                 // retry, retry_modified, alternate_tool, fallback_model, skip_step, fail
	Modifications map[string]interface{} // e.g. timeout multiplier, prompt suffix
	DelayMs      int
	AlternateTool string
	Reason       string
}

// RepairHandler decides repair action from evaluation result and step state.
type RepairHandler interface {
	Decide(eval *EvaluationResult, stepRetryCount, maxRetries int) RepairAction
}

// DefaultRepairHandler implements RepairHandler with plan-defined strategies.
type DefaultRepairHandler struct{}

// NewDefaultRepairHandler returns a RepairHandler.
func NewDefaultRepairHandler() *DefaultRepairHandler {
	return &DefaultRepairHandler{}
}

// Decide returns the repair action for the given evaluation.
func (h *DefaultRepairHandler) Decide(eval *EvaluationResult, stepRetryCount, maxRetries int) RepairAction {
	if stepRetryCount >= maxRetries {
		return RepairAction{Type: "fail", Reason: "max_retries_exceeded"}
	}
	if eval.Status != EvalFailure && eval.Status != EvalNeedsRepair {
		return RepairAction{Type: "continue"}
	}
	switch eval.FailureClass {
	case FailureToolValidation:
		return RepairAction{Type: "retry_modified", Modifications: map[string]interface{}{"fix_validation": true}}
	case FailureToolTimeout:
		return RepairAction{Type: "retry_modified", Modifications: map[string]interface{}{"timeout_multiplier": 2}}
	case FailureModelMalformed:
		return RepairAction{Type: "retry_modified", Modifications: map[string]interface{}{"add_prompt_suffix": "Respond ONLY with valid JSON."}}
	case FailureRetrievalMiss:
		return RepairAction{Type: "retry_modified", Modifications: map[string]interface{}{"broaden_query": true, "lower_threshold": 0.1}}
	case FailureDependency:
		return RepairAction{Type: "retry", DelayMs: 1000 * (1 << stepRetryCount)}
	case FailureEmptyResult:
		return RepairAction{Type: "alternate_tool"}
	case FailureToolAuth, FailurePolicyViolation:
		return RepairAction{Type: "fail", Reason: string(eval.FailureClass)}
	default:
		return RepairAction{Type: "retry"}
	}
}
