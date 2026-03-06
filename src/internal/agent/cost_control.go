/*-------------------------------------------------------------------------
 * cost_control.go
 *    Per-agent and per-run cost checks with hard limits.
 *-------------------------------------------------------------------------*/

package agent

// CostLimitFromTaskMetadata returns the cost_limit from task_metadata if set.
func CostLimitFromTaskMetadata(metadata map[string]interface{}) float64 {
	if metadata == nil {
		return 0
	}
	v, ok := metadata["cost_limit"]
	if !ok {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return x
	case int:
		return float64(x)
	default:
		return 0
	}
}

// CheckRunCostLimit returns true if the run is within cost limit (or no limit set).
func CheckRunCostLimit(currentCost *float64, limit float64) bool {
	if limit <= 0 {
		return true
	}
	if currentCost == nil {
		return true
	}
	return *currentCost <= limit
}
