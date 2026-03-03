/*-------------------------------------------------------------------------
 *
 * config.go
 *    Database operations
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/config/config.go
 *
 *-------------------------------------------------------------------------
 */

package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

/* ModuleEntry holds per-module enabled flag and config. */
type ModuleEntry struct {
	Enabled bool                   `yaml:"enabled"`
	Config  map[string]interface{} `yaml:"config"`
}

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Database    DatabaseConfig   `yaml:"database"`
	Auth        AuthConfig       `yaml:"auth"`
	Logging     LoggingConfig    `yaml:"logging"`
	Workflow    WorkflowConfig   `yaml:"workflow"`
	Distributed DistributedConfig `yaml:"distributed"`
	Cache       CacheConfig      `yaml:"cache"`
	Multimodal  MultimodalConfig `yaml:"multimodal"`
	Modules     map[string]ModuleEntry `yaml:"modules"`
	/* Profile selects dev, staging, or prod defaults; overridden by env */
	Profile string `yaml:"profile"`
	/* RejectUnknownFields when true rejects request bodies with unknown fields (schema validation) */
	RejectUnknownFields bool `yaml:"reject_unknown_fields"`
	/* Compliance profile: standard, gdpr, hipaa, sox (memory TTL, audit verbosity, export restrictions) */
	Compliance ComplianceConfig `yaml:"compliance"`
	/* Agent and reliability */
	Agent AgentConfig `yaml:"agent"`
	Tools ToolsConfig `yaml:"tools"`
}

/* ComplianceConfig holds compliance profile name. */
type ComplianceConfig struct {
	Profile string `yaml:"profile"` /* standard | gdpr | hipaa | sox */
}

/* AgentConfig holds agent runtime options. */
type AgentConfig struct {
	DeterministicMode bool `yaml:"deterministic_mode"` /* When true: tool-only, seed-controlled, no freeform LLM */
}

/* ToolsConfig holds tool execution limits. */
type ToolsConfig struct {
	Timeout time.Duration `yaml:"timeout"` /* Per-tool execution timeout; 0 = use default */
}

/* WorkflowConfig already exists with BaseURL; extend with MaxDuration. */

type WorkflowConfig struct {
	BaseURL      string        `yaml:"base_url"`
	MaxDuration  time.Duration `yaml:"max_duration"` /* Max total time for a workflow run; 0 = no limit */
}

type ServerConfig struct {
	Host         string        `yaml:"host"`
	Port         int           `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type DatabaseConfig struct {
	Host            string        `yaml:"host"`
	Port            int           `yaml:"port"`
	Database        string        `yaml:"database"`
	User            string        `yaml:"user"`
	Password        string        `yaml:"password"`
	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time"`
}

type AuthConfig struct {
	APIKeyHeader            string   `yaml:"api_key_header"`
	AllowedOrigins          []string `yaml:"allowed_origins"`           /* Allowed origins for CORS and WebSocket */
	WebSocketAllowedOrigins []string `yaml:"websocket_allowed_origins"` /* Separate WebSocket origins if needed */
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type DistributedConfig struct {
	Enabled     bool          `yaml:"enabled"`
	NodeAddress string        `yaml:"node_address"`
	NodePort    int           `yaml:"node_port"`
	RPCTimeout  time.Duration `yaml:"rpc_timeout"`
	RPCSecret   string        `yaml:"rpc_secret"`  /* Shared secret for node-to-node authentication */
	RPCAPIKey   string        `yaml:"rpc_api_key"` /* Alternative: API key for RPC authentication */
	UseTLS      bool          `yaml:"use_tls"`     /* Use HTTPS for RPC calls */
}

type CacheConfig struct {
	Enabled      bool          `yaml:"enabled"`
	TTL          time.Duration `yaml:"ttl"`
	SyncInterval time.Duration `yaml:"sync_interval"`
}

type MultimodalConfig struct {
	OCRProvider string            `yaml:"ocr_provider"`
	APIKeys     map[string]string `yaml:"api_keys"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	/* Override with environment variables */
	if err := LoadFromEnv(&config); err != nil {
		return nil, fmt.Errorf("failed to load from env: %w", err)
	}
	/* Apply profile defaults if set */
	ApplyProfile(&config)
	if err := ValidateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return &config, nil
}

/* ApplyProfile applies dev/staging/prod profile defaults (only for unset/zero values) */
func ApplyProfile(cfg *Config) {
	switch cfg.Profile {
	case "dev", "development":
		if cfg.Logging.Level == "" {
			cfg.Logging.Level = "debug"
		}
		if cfg.Logging.Format == "" {
			cfg.Logging.Format = "console"
		}
	case "staging":
		if cfg.Logging.Level == "" {
			cfg.Logging.Level = "info"
		}
		if cfg.Logging.Format == "" {
			cfg.Logging.Format = "json"
		}
	case "prod", "production":
		if cfg.Logging.Level == "" {
			cfg.Logging.Level = "warn"
		}
		if cfg.Logging.Format == "" {
			cfg.Logging.Format = "json"
		}
	}
}

/* ValidateConfig checks that config values are valid (e.g. timeouts positive) */
func ValidateConfig(cfg *Config) error {
	if cfg.Server.ReadTimeout < 0 {
		return fmt.Errorf("server read timeout cannot be negative")
	}
	if cfg.Server.WriteTimeout < 0 {
		return fmt.Errorf("server write timeout cannot be negative")
	}
	if cfg.Distributed.Enabled && cfg.Distributed.RPCTimeout <= 0 {
		return fmt.Errorf("RPC timeout must be positive when distributed is enabled")
	}
	return nil
}

/* Redact returns a copy of the config with secret fields replaced by "[REDACTED]" for logging */
func Redact(cfg *Config) map[string]interface{} {
	if cfg == nil {
		return nil
	}
	out := make(map[string]interface{})
	out["server"] = map[string]interface{}{
		"host": cfg.Server.Host, "port": cfg.Server.Port,
		"read_timeout": cfg.Server.ReadTimeout.String(), "write_timeout": cfg.Server.WriteTimeout.String(),
	}
	out["database"] = map[string]interface{}{
		"host": cfg.Database.Host, "port": cfg.Database.Port, "database": cfg.Database.Database,
		"user": cfg.Database.User, "password": "[REDACTED]",
		"max_open_conns": cfg.Database.MaxOpenConns, "max_idle_conns": cfg.Database.MaxIdleConns,
	}
	out["auth"] = map[string]interface{}{
		"api_key_header": cfg.Auth.APIKeyHeader,
		"allowed_origins": cfg.Auth.AllowedOrigins, "websocket_allowed_origins": cfg.Auth.WebSocketAllowedOrigins,
	}
	out["logging"] = map[string]interface{}{"level": cfg.Logging.Level, "format": cfg.Logging.Format}
	out["workflow"] = map[string]interface{}{"base_url": cfg.Workflow.BaseURL}
	out["distributed"] = map[string]interface{}{
		"enabled": cfg.Distributed.Enabled, "node_address": cfg.Distributed.NodeAddress, "node_port": cfg.Distributed.NodePort,
		"rpc_timeout": cfg.Distributed.RPCTimeout.String(), "rpc_secret": "[REDACTED]", "rpc_api_key": "[REDACTED]", "use_tls": cfg.Distributed.UseTLS,
	}
	out["cache"] = map[string]interface{}{
		"enabled": cfg.Cache.Enabled, "ttl": cfg.Cache.TTL.String(), "sync_interval": cfg.Cache.SyncInterval.String(),
	}
	out["profile"] = cfg.Profile
	out["reject_unknown_fields"] = cfg.RejectUnknownFields
	out["multimodal"] = map[string]interface{}{"ocr_provider": cfg.Multimodal.OCRProvider, "api_keys": "[REDACTED]"}
	if cfg.Modules != nil {
		mods := make(map[string]interface{})
		for k, v := range cfg.Modules {
			mods[k] = map[string]interface{}{"enabled": v.Enabled, "config": "[REDACTED]"}
		}
		out["modules"] = mods
	}
	return out
}

/* ConfigDump returns a map suitable for GET /admin/config (non-secret settings only) */
func ConfigDump(cfg *Config) map[string]interface{} {
	return Redact(cfg)
}

/* isSecretField returns true if the struct field name suggests a secret (used for generic redaction) */
func isSecretField(name string) bool {
	lower := strings.ToLower(name)
	return strings.Contains(lower, "password") || strings.Contains(lower, "secret") ||
		strings.Contains(lower, "api_key") || strings.Contains(lower, "token")
}

/* DefaultConfig is now in defaults.go */
