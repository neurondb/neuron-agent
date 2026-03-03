/*-------------------------------------------------------------------------
 *
 * metrics.go
 *    Eval metrics: pass_rate, unsafe_rate, schema_error_rate, citation_coverage, latency
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/eval/metrics.go
 *
 *-------------------------------------------------------------------------
 */

package eval

import "sort"

/* EvalMetrics holds aggregate metrics for a run */
type EvalMetrics struct {
	PassRate            float64
	UnsafeRate          float64
	SchemaErrorRate     float64
	CitationCoverage    float64
	PlanImprovementRate float64
	LatencyP50Ms        float64
	LatencyP95Ms        float64
	TotalTasks          int
	PassedTasks         int
	FailedTasks         int
}

/* ComputeMetrics derives metrics from task results */
func ComputeMetrics(total, passed, failed int, unsafeCount, schemaErrorCount int, withCitations int, latenciesMs []float64) EvalMetrics {
	m := EvalMetrics{TotalTasks: total, PassedTasks: passed, FailedTasks: failed}
	if total > 0 {
		m.PassRate = float64(passed) / float64(total)
		m.UnsafeRate = float64(unsafeCount) / float64(total)
		m.SchemaErrorRate = float64(schemaErrorCount) / float64(total)
		m.CitationCoverage = float64(withCitations) / float64(total)
	}
	if len(latenciesMs) > 0 {
		m.LatencyP50Ms = percentile(latenciesMs, 50)
		m.LatencyP95Ms = percentile(latenciesMs, 95)
	}
	return m
}

func percentile(vals []float64, p int) float64 {
	if len(vals) == 0 {
		return 0
	}
	sorted := make([]float64, len(vals))
	copy(sorted, vals)
	sort.Float64s(sorted)
	idx := (p * len(sorted)) / 100
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}
