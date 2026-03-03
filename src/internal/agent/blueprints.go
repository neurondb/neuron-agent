/*-------------------------------------------------------------------------
 *
 * blueprints.go
 *    Agent blueprints for enterprise templates
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/agent/blueprints.go
 *
 *-------------------------------------------------------------------------
 */

package agent

/* Blueprint is an enterprise agent template (name, description, system prompt, enabled tools) */
type Blueprint struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	SystemPrompt string                 `json:"system_prompt"`
	ModelName    string                 `json:"model_name,omitempty"`
	Config       map[string]interface{} `json:"config,omitempty"`
	EnabledTools []string               `json:"enabled_tools"`
}

/* GetBlueprints returns built-in enterprise blueprints */
func GetBlueprints() []Blueprint {
	return []Blueprint{
		{
			ID:           "postgres-performance",
			Name:         "Postgres Performance Agent",
			Description:  "Analyzes PostgreSQL performance using NeuronSQL tools and safe read-only SQL",
			SystemPrompt:  "You are a PostgreSQL performance expert. Use neuronsql tools to inspect schema, validate and explain SQL, and suggest optimizations. Never run arbitrary SQL; use only the provided NeuronSQL tools.",
			ModelName:    "gpt-4",
			Config:       map[string]interface{}{"temperature": 0.2, "max_tokens": 2000},
			EnabledTools: []string{"neuronsql.schema_snapshot", "neuronsql.validate_sql", "neuronsql.explain_json", "neuronsql.generate_select", "neuronsql.optimize_select", "sql"},
		},
		{
			ID:           "knowledge-rag",
			Name:         "Knowledge RAG Agent",
			Description:  "RAG agent with memory and web search for knowledge retrieval",
			SystemPrompt:  "You are a knowledge assistant. Use RAG, memory, and web search to answer questions accurately. Cite sources when possible.",
			ModelName:    "gpt-4",
			Config:       map[string]interface{}{"temperature": 0.4, "max_tokens": 1500},
			EnabledTools: []string{"rag", "memory", "web_search"},
		},
		{
			ID:           "compliance-auditor",
			Name:         "Compliance Auditor Agent",
			Description:  "Agent focused on compliance checks, audit trails, and safe SQL",
			SystemPrompt:  "You are a compliance auditor. Use only approved tools (sql, memory) and focus on audit trails, access patterns, and policy-compliant queries. Do not suggest or run destructive operations.",
			ModelName:    "gpt-4",
			Config:       map[string]interface{}{"temperature": 0.1, "max_tokens": 1500},
			EnabledTools: []string{"sql", "memory"},
		},
		{
			ID:           "data-analyst",
			Name:         "Data Analyst Agent",
			Description:  "Data analysis with SQL, visualization, and code",
			SystemPrompt:  "You are a data analyst. Help with SQL queries, visualizations, and analysis code. Prefer read-only queries and cite data sources.",
			ModelName:    "gpt-4",
			Config:       map[string]interface{}{"temperature": 0.3, "max_tokens": 2000},
			EnabledTools: []string{"sql", "visualization", "code"},
		},
	}
}

/* GetBlueprintByID returns a blueprint by id or nil */
func GetBlueprintByID(id string) *Blueprint {
	for _, b := range GetBlueprints() {
		if b.ID == id {
			return &b
		}
	}
	return nil
}
