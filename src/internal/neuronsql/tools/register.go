/*-------------------------------------------------------------------------
 *
 * register.go
 *    Register NeuronSQL tools with the agent tool registry
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/register.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
	agenttools "github.com/neurondb/NeuronAgent/internal/tools"
)

/* RegisterNeuronSQLTools registers all 7 NeuronSQL tools with the agent registry */
func RegisterNeuronSQLTools(registry *agenttools.Registry, factory ConnectionFactory, policyEngine *policy.PolicyEngineImpl, sensitiveTables []string) {
	if factory == nil || policyEngine == nil {
		return
	}
	registry.RegisterHandler("schema_snapshot", &SchemaSnapshotTool{Factory: factory, Policy: policyEngine})
	registry.RegisterHandler("validate_sql", &ValidateSQLTool{Factory: factory, Policy: policyEngine})
	registry.RegisterHandler("explain_json", &ExplainJSONTool{Factory: factory, Policy: policyEngine})
	registry.RegisterHandler("optimize_candidates", &OptimizeTool{Factory: factory, Policy: policyEngine})
	registry.RegisterHandler("table_profile", &TableProfileTool{Factory: factory, Policy: policyEngine})
	registry.RegisterHandler("index_profile", &IndexProfileTool{Factory: factory, Policy: policyEngine})
	registry.RegisterHandler("sample_rows", &SampleRowsTool{Factory: factory, Policy: policyEngine, SensitiveTables: sensitiveTables})
}
