/*-------------------------------------------------------------------------
 * model_router.go
 *    Model selection by capability registry and routing policy.
 *-------------------------------------------------------------------------*/

package agent

// ModelCapability describes a model's capabilities and costs.
type ModelCapability struct {
	Name               string
	Provider           string
	MaxContextTokens   int
	CostPerMTokenIn    float64
	CostPerMTokenOut   float64
	LatencyP50Ms       int
	LatencyP95Ms       int
	SupportsTools      bool
	SupportsStreaming  bool
	SupportsImages     bool
	ReasoningQuality   string // low, medium, high, very_high
	CodeQuality        string
	TaskTypes          []string
	Score              float64 // set by Router during Route()
}

// RoutingPolicy defines weights for model selection.
type RoutingPolicy struct {
	CostWeight    float64
	QualityWeight float64
	LatencyWeight float64
}

// TaskContext is the context for model selection.
type TaskContext struct {
	EstimatedContextTokens int
	RequiresTools          bool
	RequiresStreaming      bool
	TaskType               string
	CostLimitUSD           float64
}

// ModelSelection is the result of routing.
type ModelSelection struct {
	Primary   *ModelCapability
	Fallbacks []*ModelCapability
	Reason    string
}

// ModelRouter selects the best model for a task.
type ModelRouter interface {
	Route(ctx TaskContext, policy RoutingPolicy, models []ModelCapability) (*ModelSelection, error)
}

// DefaultModelRouter implements ModelRouter with plan-defined scoring.
type DefaultModelRouter struct {
	qualityMap map[string]float64
}

// NewDefaultModelRouter creates a ModelRouter with default quality scores.
func NewDefaultModelRouter() *DefaultModelRouter {
	return &DefaultModelRouter{
		qualityMap: map[string]float64{
			"low": 0.25, "medium": 0.5, "high": 0.75, "very_high": 1.0,
		},
	}
}

// Route returns the best model and fallbacks for the task.
func (r *DefaultModelRouter) Route(ctx TaskContext, policy RoutingPolicy, models []ModelCapability) (*ModelSelection, error) {
	if len(models) == 0 {
		return nil, nil
	}
	candidates := make([]*ModelCapability, 0, len(models))
	for i := range models {
		m := &models[i]
		if m.MaxContextTokens < ctx.EstimatedContextTokens {
			continue
		}
		if ctx.RequiresTools && !m.SupportsTools {
			continue
		}
		if ctx.RequiresStreaming && !m.SupportsStreaming {
			continue
		}
		if ctx.CostLimitUSD > 0 {
			estCost := (float64(ctx.EstimatedContextTokens)/1e6)*m.CostPerMTokenIn + 0.5*m.CostPerMTokenOut
			if estCost > ctx.CostLimitUSD {
				continue
			}
		}
		candidates = append(candidates, m)
	}
	if len(candidates) == 0 {
		return nil, nil
	}
	maxCost := 0.0
	maxLatency := 0
	for _, c := range candidates {
		if c.CostPerMTokenIn > maxCost {
			maxCost = c.CostPerMTokenIn
		}
		if c.LatencyP50Ms > maxLatency {
			maxLatency = c.LatencyP50Ms
		}
	}
	if maxCost == 0 {
		maxCost = 1
	}
	if maxLatency == 0 {
		maxLatency = 1
	}
	for _, c := range candidates {
		costNorm := 1.0 - (c.CostPerMTokenIn / maxCost)
		q := r.qualityMap[c.ReasoningQuality]
		if q == 0 {
			q = 0.5
		}
		if ctx.TaskType == "code" && c.CodeQuality != "" {
			q = r.qualityMap[c.CodeQuality]
			if q == 0 {
				q = 0.5
			}
		}
		latencyNorm := 1.0 - (float64(c.LatencyP50Ms) / float64(maxLatency))
		c.Score = policy.CostWeight*costNorm + policy.QualityWeight*q + policy.LatencyWeight*latencyNorm
	}
	sortModelCandidatesByScore(candidates)
	primary := candidates[0]
	fallbacks := candidates[1:]
	if len(fallbacks) > 3 {
		fallbacks = fallbacks[:3]
	}
	return &ModelSelection{
		Primary:   primary,
		Fallbacks: fallbacks,
		Reason:    "routing_policy",
	}, nil
}

// Score is attached to ModelCapability during Route (internal use).
func sortModelCandidatesByScore(candidates []*ModelCapability) {
	for i := 0; i < len(candidates); i++ {
		for j := i + 1; j < len(candidates); j++ {
			if candidates[j].Score > candidates[i].Score {
				candidates[i], candidates[j] = candidates[j], candidates[i]
			}
		}
	}
}
