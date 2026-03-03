/*-------------------------------------------------------------------------
 *
 * pglang.go
 *    PGLangProvider: HTTP client for PGLang (memorable) completion API
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/provider/pglang.go
 *
 *-------------------------------------------------------------------------
 */

package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/neurondb/NeuronAgent/pkg/neuronsql"
)

/* PGLangProvider implements neuronsql.LLMProvider via HTTP to PGLang server */
type PGLangProvider struct {
	config PGLangConfig
	client *http.Client
}

/* NewPGLangProvider creates a provider with the given config */
func NewPGLangProvider(config PGLangConfig) *PGLangProvider {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	return &PGLangProvider{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

/* completionRequest is the body for POST /v1/completions */
type completionRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
	Model     string `json:"model"`
}

/* completionResponse is the response from PGLang server */
type completionResponse struct {
	Text  string         `json:"text"`
	Usage map[string]int `json:"usage,omitempty"`
}

/* Complete implements neuronsql.LLMProvider */
func (p *PGLangProvider) Complete(ctx context.Context, messages []neuronsql.Message, schema *neuronsql.JSONSchema, settings neuronsql.CompletionSettings) (*neuronsql.Completion, error) {
	prompt := p.formatMessages(messages)
	maxTokens := settings.MaxTokens
	if maxTokens <= 0 {
		maxTokens = 1024
	}
	reqBody := completionRequest{
		Prompt:    prompt,
		MaxTokens: maxTokens,
		Model:     p.config.ModelName,
	}
	b, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	url := strings.TrimSuffix(p.config.Endpoint, "/") + "/v1/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if p.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.config.APIKey)
	}
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errResponse(resp.StatusCode, resp.Body)
	}
	var out completionResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return &neuronsql.Completion{
		Text:  out.Text,
		Usage: out.Usage,
	}, nil
}

/* formatMessages builds a single prompt string for the model (qlora uses ### Instruction: format) */
func (p *PGLangProvider) formatMessages(messages []neuronsql.Message) string {
	var b strings.Builder
	for _, m := range messages {
		switch m.Role {
		case "system":
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		case "user":
			if p.config.ModelName == "qlora" {
				b.WriteString("### Instruction:\n")
			}
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		case "assistant":
			if p.config.ModelName == "qlora" {
				b.WriteString("### Response:\n")
			}
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		default:
			b.WriteString(m.Content)
			b.WriteString("\n\n")
		}
	}
	if p.config.ModelName == "qlora" {
		b.WriteString("### Response:\n")
	}
	return strings.TrimSpace(b.String())
}

/* SupportsToolCalls implements neuronsql.LLMProvider - PGLang does not support tool calls */
func (p *PGLangProvider) SupportsToolCalls() bool {
	return false
}

/* SupportsJSONSchema implements neuronsql.LLMProvider - best-effort only */
func (p *PGLangProvider) SupportsJSONSchema() bool {
	return false
}

func errResponse(code int, body io.Reader) error {
	return &httpErr{code: code}
}

type httpErr struct {
	code int
}

func (e *httpErr) Error() string {
	return fmt.Sprintf("pglang provider: HTTP %d", e.code)
}
