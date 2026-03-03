/*-------------------------------------------------------------------------
 *
 * config_test.go
 *    Tests for config defaults, validation, and profile application.
 *
 *-------------------------------------------------------------------------
 */

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg == nil {
		t.Fatal("DefaultConfig() returned nil")
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Server.Port = %d, want 8080", cfg.Server.Port)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("Database.Port = %d, want 5432", cfg.Database.Port)
	}
	if cfg.Logging.Level != "info" {
		t.Errorf("Logging.Level = %q, want info", cfg.Logging.Level)
	}
	if _, ok := cfg.Modules["neuronsql"]; !ok {
		t.Error("DefaultConfig should include neuronsql module")
	}
}

func TestApplyProfile(t *testing.T) {
	tests := []struct {
		profile    string
		wantLevel  string
		wantFormat string
	}{
		{"dev", "debug", "console"},
		{"development", "debug", "console"},
		{"staging", "info", "json"},
		{"prod", "warn", "json"},
		{"production", "warn", "json"},
		{"", "", ""}, /* empty profile: ApplyProfile does not set level/format */
	}
	for _, tt := range tests {
		cfg := DefaultConfig()
		cfg.Profile = tt.profile
		cfg.Logging.Level = ""
		cfg.Logging.Format = ""
		ApplyProfile(cfg)
		if cfg.Logging.Level != tt.wantLevel {
			t.Errorf("Profile %q: Logging.Level = %q, want %q", tt.profile, cfg.Logging.Level, tt.wantLevel)
		}
		if cfg.Logging.Format != tt.wantFormat {
			t.Errorf("Profile %q: Logging.Format = %q, want %q", tt.profile, cfg.Logging.Format, tt.wantFormat)
		}
	}
}

func TestValidateConfig(t *testing.T) {
	valid := DefaultConfig()
	if err := ValidateConfig(valid); err != nil {
		t.Errorf("ValidateConfig(default) = %v", err)
	}

	negRead := DefaultConfig()
	negRead.Server.ReadTimeout = -time.Second
	if err := ValidateConfig(negRead); err == nil {
		t.Error("ValidateConfig(negative read timeout) expected error")
	}

	negWrite := DefaultConfig()
	negWrite.Server.WriteTimeout = -time.Second
	if err := ValidateConfig(negWrite); err == nil {
		t.Error("ValidateConfig(negative write timeout) expected error")
	}

	distEnabledNoTimeout := DefaultConfig()
	distEnabledNoTimeout.Distributed.Enabled = true
	distEnabledNoTimeout.Distributed.RPCTimeout = 0
	if err := ValidateConfig(distEnabledNoTimeout); err == nil {
		t.Error("ValidateConfig(distributed enabled, zero RPC timeout) expected error")
	}
}

func TestRedact(t *testing.T) {
	if got := Redact(nil); got != nil {
		t.Errorf("Redact(nil) = %v, want nil", got)
	}
	cfg := DefaultConfig()
	cfg.Database.Password = "secret"
	out := Redact(cfg)
	if out == nil {
		t.Fatal("Redact(Config) returned nil")
	}
	db, ok := out["database"].(map[string]interface{})
	if !ok {
		t.Fatal("Redact: database section not present")
	}
	if p, ok := db["password"].(string); !ok || p != "[REDACTED]" {
		t.Errorf("Redact: database.password = %v, want [REDACTED]", db["password"])
	}
}

func TestConfigDump(t *testing.T) {
	cfg := DefaultConfig()
	out := ConfigDump(cfg)
	if out == nil {
		t.Fatal("ConfigDump returned nil")
	}
	if _, ok := out["server"]; !ok {
		t.Error("ConfigDump: server section missing")
	}
}

func TestDevelopmentConfig(t *testing.T) {
	cfg := DevelopmentConfig()
	if cfg.Logging.Level != "debug" || cfg.Logging.Format != "console" {
		t.Errorf("DevelopmentConfig: level=%q format=%q", cfg.Logging.Level, cfg.Logging.Format)
	}
	if cfg.Database.MaxOpenConns != 10 {
		t.Errorf("DevelopmentConfig: MaxOpenConns = %d, want 10", cfg.Database.MaxOpenConns)
	}
}

func TestProductionConfig(t *testing.T) {
	cfg := ProductionConfig()
	if cfg.Logging.Level != "warn" || cfg.Logging.Format != "json" {
		t.Errorf("ProductionConfig: level=%q format=%q", cfg.Logging.Level, cfg.Logging.Format)
	}
	if cfg.Database.MaxOpenConns != 100 {
		t.Errorf("ProductionConfig: MaxOpenConns = %d, want 100", cfg.Database.MaxOpenConns)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Error("LoadConfig(missing file) expected error")
	}
}

func TestLoadConfig_ValidYAML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	const yaml = `
server:
  port: 9000
database:
  host: db.example.com
  port: 5433
logging:
  level: debug
`
	if err := os.WriteFile(path, []byte(yaml), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.Server.Port != 9000 {
		t.Errorf("Server.Port = %d, want 9000", cfg.Server.Port)
	}
	if cfg.Database.Host != "db.example.com" {
		t.Errorf("Database.Host = %q, want db.example.com", cfg.Database.Host)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %q, want debug", cfg.Logging.Level)
	}
}

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", []string{}},
		{"a", []string{"a"}},
		{"a,b,c", []string{"a", "b", "c"}},
		{" a , b , c ", []string{"a", "b", "c"}},
		{"a,,b", []string{"a", "b"}},
	}
	for _, tt := range tests {
		got := parseCommaSeparated(tt.in)
		if len(got) != len(tt.want) {
			t.Errorf("parseCommaSeparated(%q) = %v, want %v", tt.in, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("parseCommaSeparated(%q)[%d] = %q, want %q", tt.in, i, got[i], tt.want[i])
			}
		}
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	key := "NEURON_AGENT_TEST_GETENV_DEFAULT"
	os.Unsetenv(key)
	defer os.Unsetenv(key)
	if got := GetEnvOrDefault(key, "fallback"); got != "fallback" {
		t.Errorf("GetEnvOrDefault(unset) = %q, want fallback", got)
	}
	os.Setenv(key, "set")
	if got := GetEnvOrDefault(key, "fallback"); got != "set" {
		t.Errorf("GetEnvOrDefault(set) = %q, want set", got)
	}
}

func TestGetEnvIntOrDefault(t *testing.T) {
	key := "NEURON_AGENT_TEST_GETENV_INT"
	os.Unsetenv(key)
	defer os.Unsetenv(key)
	if got := GetEnvIntOrDefault(key, 42); got != 42 {
		t.Errorf("GetEnvIntOrDefault(unset) = %d, want 42", got)
	}
	os.Setenv(key, "7")
	if got := GetEnvIntOrDefault(key, 42); got != 7 {
		t.Errorf("GetEnvIntOrDefault(7) = %d, want 7", got)
	}
	os.Setenv(key, "not-a-number")
	if got := GetEnvIntOrDefault(key, 42); got != 42 {
		t.Errorf("GetEnvIntOrDefault(invalid) = %d, want 42", got)
	}
}

func TestGetEnvDurationOrDefault(t *testing.T) {
	key := "NEURON_AGENT_TEST_GETENV_DURATION"
	os.Unsetenv(key)
	defer os.Unsetenv(key)
	want := 5 * time.Second
	if got := GetEnvDurationOrDefault(key, want); got != want {
		t.Errorf("GetEnvDurationOrDefault(unset) = %v, want %v", got, want)
	}
	os.Setenv(key, "10s")
	if got := GetEnvDurationOrDefault(key, want); got != 10*time.Second {
		t.Errorf("GetEnvDurationOrDefault(10s) = %v", got)
	}
}

func TestValidateEnv(t *testing.T) {
	/* ValidateEnv requires DB_* to be set; test that it returns error when unset */
	orig := map[string]string{}
	for _, k := range []string{"DB_HOST", "DB_NAME", "DB_USER", "DB_PASSWORD"} {
		orig[k] = os.Getenv(k)
		os.Unsetenv(k)
	}
	defer func() {
		for k, v := range orig {
			if v != "" {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}
	}()
	if err := ValidateEnv(); err == nil {
		t.Error("ValidateEnv() with unset required vars expected error")
	}
}
