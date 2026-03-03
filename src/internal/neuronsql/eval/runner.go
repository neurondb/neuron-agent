/*-------------------------------------------------------------------------
 *
 * runner.go
 *    Evaluation runner: run suite, collect metrics
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/eval/runner.go
 *
 *-------------------------------------------------------------------------
 */

package eval

import (
	"context"
	"time"

	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* GenerateRunner is the minimal interface for eval (e.g. *orchestrator.Orchestrator) */
type GenerateRunner interface {
	RunGenerate(ctx context.Context, dsn string, question string, requestID string) (*neuronsql.NeuronSQLResponse, error)
}

/* Runner runs evaluation suites and produces EvalReport */
type Runner struct {
	Orchestrator GenerateRunner
	DSN         string
}

/* RunSuite runs the named suite and returns an EvalReport */
func (r *Runner) RunSuite(ctx context.Context, suiteName string) (*neuronsql.EvalReport, error) {
	suite, ok := GetSuite(suiteName)
	if !ok {
		return &neuronsql.EvalReport{SuiteName: suiteName}, nil
	}
	var results []Result
	var latencies []float64
	passed := 0
	unsafeCount := 0
	schemaErrorCount := 0
	withCitations := 0
	for _, task := range suite.Tasks {
		start := time.Now()
		resp, err := r.Orchestrator.RunGenerate(ctx, r.DSN, task.Question, "eval-"+task.ID)
		elapsed := time.Since(start)
		latencyMs := float64(elapsed.Milliseconds())
		latencies = append(latencies, latencyMs)
		res := Result{TaskID: task.ID, LatencyMs: latencyMs}
		if err != nil {
			res.Passed = false
			res.Message = err.Error()
			schemaErrorCount++
		} else {
			res.Passed = resp.ValidationReport != nil && resp.ValidationReport.Valid
			if resp.Safety != nil && !resp.Safety.PolicyAllowed {
				res.Unsafe = true
				unsafeCount++
			}
			if len(resp.Citations) > 0 {
				res.HasCitation = true
				withCitations++
			}
			if !res.Passed && resp.ValidationReport != nil && len(resp.ValidationReport.Errors) > 0 {
				schemaErrorCount++
			}
		}
		if res.Pass(suite) {
			passed++
		}
		results = append(results, res)
	}
	metrics := ComputeMetrics(len(suite.Tasks), passed, len(suite.Tasks)-passed, unsafeCount, schemaErrorCount, withCitations, latencies)
	report := &neuronsql.EvalReport{
		SuiteName:          suiteName,
		PassRate:           metrics.PassRate,
		UnsafeRate:         metrics.UnsafeRate,
		SchemaErrorRate:    metrics.SchemaErrorRate,
		CitationCoverage:  metrics.CitationCoverage,
		PlanImprovementRate: metrics.PlanImprovementRate,
		LatencyP50Ms:       metrics.LatencyP50Ms,
		LatencyP95Ms:       metrics.LatencyP95Ms,
		TotalTasks:         metrics.TotalTasks,
		PassedTasks:        metrics.PassedTasks,
		FailedTasks:        metrics.FailedTasks,
	}
	for _, res := range results {
		report.Details = append(report.Details, neuronsql.EvalTaskResult{
			TaskID:    res.TaskID,
			Passed:    res.Passed,
			Message:   res.Message,
			LatencyMs: res.LatencyMs,
		})
	}
	return report, nil
}

/* GetSuite returns a suite by name (built-in suites; expand with fixture JSON for 30/15/10 prompts) */
func GetSuite(name string) (*Suite, bool) {
	switch name {
	case "northwind_like":
		return &Suite{Name: name, Tasks: northwindLikeTasks()}, true
	case "saas_basic":
		return &Suite{Name: name, Tasks: saasBasicTasks()}, true
	case "analytics_timeseries":
		return &Suite{Name: name, Tasks: analyticsTimeseriesTasks()}, true
	default:
		return nil, false
	}
}

func northwindLikeTasks() []Task {
	return []Task{
		{ID: "nw1", Question: "List all products", RequireCitation: false},
		{ID: "nw2", Question: "Count orders per customer", RequireCitation: false},
		{ID: "nw3", Question: "Top 10 products by quantity sold", RequireCitation: false},
	}
}

func saasBasicTasks() []Task {
	return []Task{
		{ID: "saas1", Question: "Count users", RequireCitation: false},
		{ID: "saas2", Question: "List tenants", RequireCitation: false},
		{ID: "saas3", Question: "Users per tenant", RequireCitation: false},
	}
}

func analyticsTimeseriesTasks() []Task {
	return []Task{
		{ID: "ts1", Question: "Daily event counts", RequireCitation: false},
		{ID: "ts2", Question: "Events by type", RequireCitation: false},
		{ID: "ts3", Question: "Hourly aggregation", RequireCitation: false},
	}
}
