/*-------------------------------------------------------------------------
 * prompt_templates.go
 *    Prompt template system by task type.
 *-------------------------------------------------------------------------*/

package agent

import (
	"fmt"
	"strings"
)

// TaskType identifies the kind of task for template selection.
type TaskType string

const (
	TaskTypeGeneral    TaskType = "general"
	TaskTypeCode       TaskType = "code"
	TaskTypeAnalysis   TaskType = "analysis"
	TaskTypeSummarization TaskType = "summarization"
	TaskTypeChat       TaskType = "chat"
)

// PromptTemplates holds templates per task type.
type PromptTemplates interface {
	Get(taskType TaskType) string
	Render(taskType TaskType, vars map[string]string) (string, error)
}

// DefaultPromptTemplates implements PromptTemplates with plan-defined sections.
type DefaultPromptTemplates struct {
	templates map[TaskType]string
}

// NewDefaultPromptTemplates returns templates for each task type.
func NewDefaultPromptTemplates() *DefaultPromptTemplates {
	return &DefaultPromptTemplates{
		templates: map[TaskType]string{
			TaskTypeGeneral:        "You are a helpful assistant. Complete the following task.\n\nTask: {{.task}}",
			TaskTypeCode:           "You are a coding assistant. Use tools to read, analyze, or modify code as needed.\n\nTask: {{.task}}",
			TaskTypeAnalysis:       "Analyze the request and use available tools to gather information. Provide a structured analysis.\n\nTask: {{.task}}",
			TaskTypeSummarization:  "Summarize the given content clearly and concisely.\n\nTask: {{.task}}",
			TaskTypeChat:           "Have a helpful conversation. Use tools when needed to answer questions.\n\nTask: {{.task}}",
		},
	}
}

// Get returns the raw template for the task type.
func (t *DefaultPromptTemplates) Get(taskType TaskType) string {
	if s, ok := t.templates[taskType]; ok {
		return s
	}
	return t.templates[TaskTypeGeneral]
}

// Render substitutes vars into the template. Supports {{.varName}}.
func (t *DefaultPromptTemplates) Render(taskType TaskType, vars map[string]string) (string, error) {
	s := t.Get(taskType)
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{{."+k+"}}", v)
	}
	if strings.Contains(s, "{{.") {
		return "", fmt.Errorf("prompt_templates: unresolved template variable")
	}
	return s, nil
}
