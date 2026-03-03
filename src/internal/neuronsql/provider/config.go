/*-------------------------------------------------------------------------
 *
 * config.go
 *    PGLang provider config (endpoint, key, model, timeout)
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/provider/config.go
 *
 *-------------------------------------------------------------------------
 */

package provider

import (
	"os"
	"strconv"
	"time"
)

/* PGLangConfig holds PGLang (memorable) endpoint and settings */
type PGLangConfig struct {
	Endpoint  string
	APIKey    string
	ModelName string
	Timeout   time.Duration
	Local     bool
	ModelDir  string
}

/* DefaultPGLangConfig returns config from env or defaults */
func DefaultPGLangConfig() PGLangConfig {
	timeout := 30 * time.Second
	if s := os.Getenv("PGLANG_TIMEOUT"); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			timeout = d
		}
	}
	if s := os.Getenv("NEURONSQL_PGLANG_TIMEOUT"); s != "" {
		if d, err := time.ParseDuration(s); err == nil {
			timeout = d
		}
	}
	endpoint := os.Getenv("PGLANG_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("NEURONSQL_PGLANG_ENDPOINT")
	}
	if endpoint == "" {
		endpoint = "http://localhost:9090"
	}
	model := os.Getenv("PGLANG_MODEL_NAME")
	if model == "" {
		model = os.Getenv("NEURONSQL_PGLANG_MODEL_NAME")
	}
	if model == "" {
		model = "scratch"
	}
	local := false
	if s := os.Getenv("PGLANG_LOCAL"); s == "true" || s == "1" {
		local = true
	}
	if s := os.Getenv("NEURONSQL_PGLANG_LOCAL"); s == "true" || s == "1" {
		local = true
	}
	modelDir := os.Getenv("PGLANG_MODEL_DIR")
	if modelDir == "" {
		modelDir = os.Getenv("NEURONSQL_PGLANG_MODEL_DIR")
	}
	apiKey := os.Getenv("PGLANG_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("NEURONSQL_PGLANG_API_KEY")
	}
	return PGLangConfig{
		Endpoint:  endpoint,
		APIKey:    apiKey,
		ModelName: model,
		Timeout:   timeout,
		Local:     local,
		ModelDir:  modelDir,
	}
}

/* LoadFromYAML allows overriding from a config struct (e.g. from main config) */
func (c *PGLangConfig) LoadFromYAML(endpoint, apiKey, modelName string, timeoutSec int) {
	if endpoint != "" {
		c.Endpoint = endpoint
	}
	if apiKey != "" {
		c.APIKey = apiKey
	}
	if modelName != "" {
		c.ModelName = modelName
	}
	if timeoutSec > 0 {
		c.Timeout = time.Duration(timeoutSec) * time.Second
	}
}

func parseInt(s string, defaultVal int) int {
	if s == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return n
}
