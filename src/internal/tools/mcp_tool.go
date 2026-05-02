/*-------------------------------------------------------------------------
 *
 * mcp_tool.go
 *    JSON-RPC tools/call to an HTTP MCP endpoint (synced tools use handler_type mcp).
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/neurondb/NeuronAgent/internal/db"
)

type MCPTool struct {
	client *http.Client
}

func NewMCPTool() *MCPTool {
	return &MCPTool{
		client: &http.Client{Timeout: mcpHTTPTimeout},
	}
}

func (m *MCPTool) Validate(args map[string]interface{}, schema map[string]interface{}) error {
	_ = schema
	_ = args
	return nil
}

func (m *MCPTool) Execute(ctx context.Context, tool *db.Tool, args map[string]interface{}) (string, error) {
	if tool.HandlerConfig == nil {
		return "", fmt.Errorf("mcp tool missing handler_config")
	}
	base, _ := tool.HandlerConfig["mcp_server_url"].(string)
	base = strings.TrimRight(strings.TrimSpace(base), "/")
	if base == "" {
		return "", fmt.Errorf("handler_config.mcp_server_url is required for mcp tools")
	}
	toolName, _ := tool.HandlerConfig["tool_name"].(string)
	if strings.TrimSpace(toolName) == "" {
		return "", fmt.Errorf("handler_config.tool_name is required for mcp tools")
	}

	body := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      1,
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name":      toolName,
			"arguments": args,
		},
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, base, bytes.NewReader(raw))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("mcp HTTP %d: %s", resp.StatusCode, truncateForErr(respBody))
	}

	var envelope struct {
		Result json.RawMessage `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(respBody, &envelope); err != nil {
		return "", fmt.Errorf("mcp response decode: %w", err)
	}
	if envelope.Error != nil {
		return "", fmt.Errorf("mcp tools/call error: %s", envelope.Error.Message)
	}
	return string(envelope.Result), nil
}

func truncateForErr(b []byte) string {
	const max = 512
	s := string(b)
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
