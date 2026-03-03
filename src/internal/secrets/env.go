/*-------------------------------------------------------------------------
 *
 * env.go
 *    Environment variable secret store implementation
 *
 * Provides environment variable-based secret store (fallback/default).
 *
 * Copyright (c) 2024-2025, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/secrets/env.go
 *
 *-------------------------------------------------------------------------
 */

package secrets

import (
	"context"
	"fmt"
	"os"
	"strings"
)

/* EnvStore implements Store interface using environment variables */
type EnvStore struct {
	prefix string
}

/* NewEnvStore creates a new environment variable secret store */
func NewEnvStore(config Config) (*EnvStore, error) {
	prefix := "NEURONDB_SECRET_"
	if prefixVal, ok := config.Metadata["prefix"].(string); ok {
		prefix = prefixVal
	}

	return &EnvStore{prefix: prefix}, nil
}

/* GetSecret retrieves a secret from environment variables */
func (e *EnvStore) GetSecret(ctx context.Context, key string) (string, error) {
	envKey := e.prefix + strings.ToUpper(strings.ReplaceAll(key, "/", "_"))
	value := os.Getenv(envKey)
	if value == "" {
		return "", fmt.Errorf("secret not found: %s (env key: %s)", key, envKey)
	}
	return value, nil
}

/* PutSecret is not supported for environment variables */
func (e *EnvStore) PutSecret(ctx context.Context, key string, value string) error {
	return fmt.Errorf("put secret not supported for environment variable store")
}

/* DeleteSecret is not supported for environment variables */
func (e *EnvStore) DeleteSecret(ctx context.Context, key string) error {
	return fmt.Errorf("delete secret not supported for environment variable store")
}

/* ListSecrets lists secrets from environment variables */
func (e *EnvStore) ListSecrets(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, e.prefix) {
			keyVal := strings.SplitN(env, "=", 2)
			if len(keyVal) == 2 {
				key := strings.TrimPrefix(keyVal[0], e.prefix)
				key = strings.ToLower(strings.ReplaceAll(key, "_", "/"))
				if prefix == "" || strings.HasPrefix(key, prefix) {
					keys = append(keys, key)
				}
			}
		}
	}
	return keys, nil
}













