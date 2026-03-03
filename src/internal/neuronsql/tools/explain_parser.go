/*-------------------------------------------------------------------------
 *
 * explain_parser.go
 *    Parse EXPLAIN (FORMAT JSON) output for plan analysis
 *
 * Extracts node types, costs, row estimates for rewrite and index suggestions.
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/explain_parser.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"encoding/json"
	"fmt"
)

/* PlanNode represents one node in an EXPLAIN plan tree */
type PlanNode struct {
	NodeType    string     `json:"Node Type"`
	TotalCost   float64    `json:"Total Cost"`
	PlanRows    float64    `json:"Plan Rows"`
	RelationName string    `json:"Relation Name,omitempty"`
	IndexName   string     `json:"Index Name,omitempty"`
	Plans       []PlanNode `json:"Plans,omitempty"`
}

/* ExplainPlan is the top-level plan (PostgreSQL returns an array with one element) */
type ExplainPlan struct {
	Plan PlanNode `json:"Plan"`
}

/* ParseExplainJSON parses PostgreSQL EXPLAIN (FORMAT JSON) output.
 * Returns the root plan node and a flat list of all nodes (depth-first).
 * Returns error if JSON is invalid or not in expected shape.
 */
func ParseExplainJSON(planJSON string) (root PlanNode, allNodes []PlanNode, err error) {
	var arr []ExplainPlan
	if err := json.Unmarshal([]byte(planJSON), &arr); err != nil {
		return PlanNode{}, nil, fmt.Errorf("explain_parser: invalid json: %w", err)
	}
	if len(arr) == 0 {
		return PlanNode{}, nil, fmt.Errorf("explain_parser: empty plan array")
	}
	root = arr[0].Plan
	allNodes = flattenPlanNodes(root)
	return root, allNodes, nil
}

func flattenPlanNodes(n PlanNode) []PlanNode {
	out := []PlanNode{n}
	for _, c := range n.Plans {
		out = append(out, flattenPlanNodes(c)...)
	}
	return out
}

/* SuggestIndexFromPlan returns text suggestions for indexes based on Seq Scan nodes with high cost/rows */
func SuggestIndexFromPlan(allNodes []PlanNode) []IndexSuggestion {
	var suggestions []IndexSuggestion
	for _, n := range allNodes {
		if n.NodeType != "Seq Scan" || n.RelationName == "" {
			continue
		}
		if n.PlanRows < 10 && n.TotalCost < 100 {
			continue
		}
		suggestions = append(suggestions, IndexSuggestion{
			Definition:       fmt.Sprintf("CREATE INDEX ON %s (...); /* add columns used in WHERE/JOIN */", n.RelationName),
			EstimatedBenefit: fmt.Sprintf("Seq Scan on %s (rows=%.0f, cost=%.2f); consider index on filter/join columns", n.RelationName, n.PlanRows, n.TotalCost),
			PlanEvidence:     n.NodeType,
		})
	}
	return suggestions
}

/* SuggestRewritesFromPlan returns text suggestions based on plan shape (e.g. suggest LIMIT, avoid SELECT *) */
func SuggestRewritesFromPlan(allNodes []PlanNode, originalSQL string) []RewriteOption {
	var opts []RewriteOption
	for _, n := range allNodes {
		if n.NodeType == "Limit" && n.PlanRows > 0 {
			opts = append(opts, RewriteOption{
				SQL:          originalSQL,
				Explanation:  "Query already uses LIMIT or similar; plan has Limit node.",
				Risk:         "none",
				PlanEvidence: "Limit",
			})
			break
		}
	}
	return opts
}
