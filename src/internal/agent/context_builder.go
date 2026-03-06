/*-------------------------------------------------------------------------
 * context_builder.go
 *    Section allocation, token budgeting, and prompt hashing for agent context.
 *-------------------------------------------------------------------------*/

package agent

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/neurondb/NeuronAgent/internal/db"
)

const (
	approxCharsPerToken   = 4
	systemContractBudget  = 500
	agentIdentityBudget  = 200
	outputFormatBudget   = 300
	toolContractMax      = 2000
	planSectionMax       = 1500
	taskDefinitionMax    = 1000
	responseReserve      = 1024
)

// ContextSection describes one section of the built context with token usage.
type ContextSection struct {
	Name       string
	Content    string
	TokensUsed int
	ItemsCount int
}

// ContextBuildResult is the result of BuildContext.
type ContextBuildResult struct {
	Prompt    string
	PromptHash string
	Sections  []ContextSection
}

// ContextBuilder assembles the prompt from sections with token budgets.
type ContextBuilder interface {
	Build(opts ContextBuildOptions) (*ContextBuildResult, error)
}

// ContextBuildOptions supplies all inputs for context building.
type ContextBuildOptions struct {
	Agent              *db.Agent
	TaskInput          string
	MemoryItems        []MemoryItem
	ToolSchemasText    string
	PlanSteps          []PlanStep
	CurrentStepIndex   int
	ConversationText   string
	MaxContextTokens   int
	SystemContract     string
	OutputFormatContract string
	SafetyRules        string
}

// DefaultContextBuilder implements ContextBuilder with plan-defined budget allocation.
type DefaultContextBuilder struct{}

// NewDefaultContextBuilder returns a ContextBuilder that allocates sections by plan.
func NewDefaultContextBuilder() *DefaultContextBuilder {
	return &DefaultContextBuilder{}
}

// Build produces the full prompt and hash from opts, respecting token budgets.
func (b *DefaultContextBuilder) Build(opts ContextBuildOptions) (*ContextBuildResult, error) {
	if opts.MaxContextTokens <= 0 {
		opts.MaxContextTokens = 8192
	}
	W := opts.MaxContextTokens - responseReserve
	var sections []ContextSection

	// P0 fixed
	systemContract := opts.SystemContract
	if systemContract == "" {
		systemContract = "You are a helpful assistant. Follow the user's instructions and use tools when needed."
	}
	sections = append(sections, truncateSection("system_contract", systemContract, systemContractBudget))

	identity := ""
	if opts.Agent != nil && opts.Agent.SystemPrompt != "" {
		identity = opts.Agent.SystemPrompt
	}
	sections = append(sections, truncateSection("identity", identity, agentIdentityBudget))

	task := truncateToTokens(opts.TaskInput, taskDefinitionMax)
	sections = append(sections, ContextSection{Name: "task", Content: task, TokensUsed: tokenEstimate(task)})

	outputFormat := opts.OutputFormatContract
	if outputFormat == "" {
		outputFormat = "Respond in plain text. When using tools, follow the tool schema."
	}
	sections = append(sections, truncateSection("output_format", outputFormat, outputFormatBudget))

	safety := opts.SafetyRules
	if safety != "" {
		sections = append(sections, truncateSection("safety", safety, 200))
	}

	// P1 tool contract
	toolBudget := min(int(0.15*float64(W)), toolContractMax)
	toolSection := truncateToTokens(opts.ToolSchemasText, toolBudget)
	sections = append(sections, ContextSection{Name: "tools", Content: toolSection, TokensUsed: tokenEstimate(toolSection)})

	// P1 plan + current step
	planBudget := min(int(0.10*float64(W)), planSectionMax)
	planText := formatPlanForContext(opts.PlanSteps, opts.CurrentStepIndex)
	planSection := truncateToTokens(planText, planBudget)
	sections = append(sections, ContextSection{Name: "plan", Content: planSection, TokensUsed: tokenEstimate(planSection)})

	// P2 memory + retrieved (already budgeted by MemorySelector)
	memoryTokens := 0
	var memoryParts []string
	for _, m := range opts.MemoryItems {
		memoryParts = append(memoryParts, fmt.Sprintf("[%s] %s", m.Tier, m.Content))
		memoryTokens += m.TokenCount
	}
	if len(memoryParts) > 0 {
		memorySection := strings.Join(memoryParts, "\n\n")
		sections = append(sections, ContextSection{
			Name: "memory", Content: memorySection, TokensUsed: memoryTokens, ItemsCount: len(opts.MemoryItems),
		})
	}

	// P2 conversation
	remaining := W
	for _, s := range sections {
		remaining -= s.TokensUsed
	}
	if remaining > 0 && opts.ConversationText != "" {
		convSection := truncateToTokens(opts.ConversationText, remaining)
		sections = append(sections, ContextSection{Name: "conversation", Content: convSection, TokensUsed: tokenEstimate(convSection)})
	}

	// Assemble prompt
	var parts []string
	for _, s := range sections {
		if s.Content == "" {
			continue
		}
		parts = append(parts, "## "+sectionTitle(s.Name)+"\n"+s.Content)
	}
	prompt := strings.Join(parts, "\n\n")
	hash := sha256.Sum256([]byte(prompt))
	return &ContextBuildResult{
		Prompt:     prompt,
		PromptHash: hex.EncodeToString(hash[:]),
		Sections:   sections,
	}, nil
}

func sectionTitle(name string) string {
	switch name {
	case "system_contract":
		return "System"
	case "identity":
		return "Identity"
	case "task":
		return "Task"
	case "output_format":
		return "Output Format"
	case "safety":
		return "Safety"
	case "tools":
		return "Tools"
	case "plan":
		return "Plan"
	case "memory":
		return "Memory"
	case "conversation":
		return "Conversation"
	default:
		return name
	}
}

func formatPlanForContext(steps []PlanStep, current int) string {
	if len(steps) == 0 {
		return "No plan."
	}
	var b strings.Builder
	for i, s := range steps {
		marker := ""
		if i == current {
			marker = " (current)"
		}
		b.WriteString(fmt.Sprintf("Step %d%s: %s", i+1, marker, s.Action))
		if s.Tool != "" {
			b.WriteString(fmt.Sprintf(" [tool: %s]", s.Tool))
		}
		b.WriteString("\n")
	}
	return b.String()
}

func truncateSection(name, content string, maxTokens int) ContextSection {
	truncated := truncateToTokens(content, maxTokens)
	return ContextSection{Name: name, Content: truncated, TokensUsed: tokenEstimate(truncated)}
}

func truncateToTokens(s string, maxTokens int) string {
	targetChars := maxTokens * approxCharsPerToken
	s = strings.TrimSpace(s)
	if len(s) <= targetChars {
		return s
	}
	return s[:targetChars] + "\n[... truncated ...]"
}

func tokenEstimate(s string) int {
	n := len(s) / approxCharsPerToken
	if n < 1 && len(s) > 0 {
		return 1
	}
	return n
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
