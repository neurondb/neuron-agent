/*-------------------------------------------------------------------------
 *
 * config_test.go
 *    Tests for CLI config loading and validation.
 *
 *-------------------------------------------------------------------------
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAgentConfig_MissingFile(t *testing.T) {
	_, err := LoadAgentConfig(filepath.Join(t.TempDir(), "missing.yaml"))
	if err == nil {
		t.Error("LoadAgentConfig(missing) expected error")
	}
}

func TestLoadAgentConfig_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "agent.yaml")
	const yaml = `
name: test-agent
description: test
model:
  name: gpt-4
`
	if err := os.WriteFile(path, []byte(yaml), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	cfg, err := LoadAgentConfig(path)
	if err != nil {
		t.Fatalf("LoadAgentConfig: %v", err)
	}
	if cfg.Name != "test-agent" {
		t.Errorf("Name = %q", cfg.Name)
	}
}

func TestValidateAgentConfig(t *testing.T) {
	c := &AgentConfig{}
	if err := ValidateAgentConfig(c); err == nil {
		t.Error("ValidateAgentConfig(empty name) expected error")
	}
	c.Name = "valid"
	if err := ValidateAgentConfig(c); err != nil {
		t.Errorf("ValidateAgentConfig(valid): %v", err)
	}
}

func TestGetWorkflowTemplates(t *testing.T) {
	templates := GetWorkflowTemplates()
	if templates == nil {
		t.Error("GetWorkflowTemplates should not return nil")
	}
}
