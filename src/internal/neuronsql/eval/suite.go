/*-------------------------------------------------------------------------
 *
 * suite.go
 *    Test suite definitions and fail rules
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/eval/suite.go
 *
 *-------------------------------------------------------------------------
 */

package eval

/* Task represents a single eval task */
type Task struct {
	ID             string
	Question       string
	ExpectedSQL    string
	ExpectedUnsafe bool
	RequireCitation bool
}

/* Suite defines a named test suite */
type Suite struct {
	Name  string
	Tasks []Task
}

/* Result is the result of running one task */
type Result struct {
	TaskID    string
	Passed    bool
	Message   string
	LatencyMs float64
	Unsafe    bool
	HasCitation bool
}

/* FailRules: any unsafe output fails; any invalid SQL fails; non-trivial claim without citation fails */
func (r *Result) Pass(suite *Suite) bool {
	if r.Unsafe {
		return false
	}
	if !r.Passed && r.Message != "" {
		return false
	}
	for _, t := range suite.Tasks {
		if t.ID == r.TaskID && t.RequireCitation && !r.HasCitation {
			return false
		}
	}
	return r.Passed
}
