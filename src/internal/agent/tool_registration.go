/*-------------------------------------------------------------------------
 *
 * tool_registration.go
 *    Registers optional tools that need runtime-backed implementations.
 *
 *-------------------------------------------------------------------------
 */

package agent

import (
	"github.com/neurondb/NeuronAgent/internal/tools"
)

/* RegisterRetrievalToolOnRegistry wires the retrieval tool when memory is available. */
func RegisterRetrievalToolOnRegistry(registry *tools.Registry, r *Runtime) {
	if registry == nil || r == nil || r.memory == nil {
		return
	}
	ra := NewRetrievalAdapter(r.memory, r.hierMemory, r.relevanceChecker)
	var kr tools.KnowledgeRouterInterface
	if r.knowledgeRouter != nil {
		kr = r.knowledgeRouter
	}
	registry.RegisterHandler("retrieval", tools.NewRetrievalTool(ra, kr, nil, tools.NewHTTPTool()))
}
