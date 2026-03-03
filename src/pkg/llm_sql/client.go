package llm_sql

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ModelClient handles communication with the SQL LLM model server
type ModelClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// GenerationConfig contains parameters for SQL generation
type GenerationConfig struct {
	Temperature float64  `json:"temperature,omitempty"`
	TopP        float64  `json:"top_p,omitempty"`
	MaxTokens   int      `json:"max_tokens,omitempty"`
	Stop        []string `json:"stop,omitempty"`
}

// GenerationResult contains the generated SQL and metadata
type GenerationResult struct {
	SQL         string   `json:"sql"`
	Explanation string   `json:"explanation"`
	Confidence  float64  `json:"confidence"`
	Warnings    []string `json:"warnings,omitempty"`
}

// NewModelClient creates a new SQL LLM model client
func NewModelClient(baseURL, apiKey string) *ModelClient {
	return &ModelClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate generates SQL from a natural language prompt
func (c *ModelClient) Generate(ctx context.Context, prompt, dialect string, schema json.RawMessage) (*GenerationResult, error) {
	systemPrompt := buildSystemPrompt(dialect, schema)
	
	requestBody := map[string]interface{}{
		"model": "sql-llm-70b",
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.2,
		"top_p":       0.95,
		"max_tokens":  2048,
		"stop":        []string{"</sql>", "\n\nUser:"},
	}
	
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/v1/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("model server returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	
	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}
	
	content := result.Choices[0].Message.Content
	sql := extractSQL(content)
	explanation := extractExplanation(content)
	
	return &GenerationResult{
		SQL:         sql,
		Explanation: explanation,
		Confidence:  0.95,
	}, nil
}

// GenerateStream generates SQL with streaming response
func (c *ModelClient) GenerateStream(ctx context.Context, prompt, dialect string, schema json.RawMessage, writer io.Writer) error {
	// Implement streaming generation
	// For now, use non-streaming and write once
	result, err := c.Generate(ctx, prompt, dialect, schema)
	if err != nil {
		return err
	}
	
	_, err = writer.Write([]byte(result.SQL))
	return err
}

func buildSystemPrompt(dialect string, schema json.RawMessage) string {
	basePrompt := fmt.Sprintf("You are an expert SQL assistant specializing in %s. Generate accurate, efficient SQL queries based on natural language descriptions.", dialect)
	
	if schema != nil {
		basePrompt += "\n\nDatabase schema:\n" + string(schema)
	}
	
	basePrompt += "\n\nProvide your response in this format:\n<sql>YOUR_SQL_HERE</sql>\n\nExplanation: Brief explanation of the query."
	
	return basePrompt
}

func extractSQL(content string) string {
	// Extract SQL between <sql> tags
	start := bytes.Index([]byte(content), []byte("<sql>"))
	end := bytes.Index([]byte(content), []byte("</sql>"))
	
	if start != -1 && end != -1 && end > start {
		return string(content[start+5 : end])
	}
	
	// Fallback: return content as-is
	return content
}

func extractExplanation(content string) string {
	// Extract explanation after </sql>
	end := bytes.Index([]byte(content), []byte("</sql>"))
	if end != -1 && len(content) > end+6 {
		explanation := content[end+6:]
		// Remove "Explanation: " prefix if present
		explanation = string(bytes.TrimPrefix([]byte(explanation), []byte("Explanation: ")))
		return string(bytes.TrimSpace([]byte(explanation)))
	}
	
	return ""
}
