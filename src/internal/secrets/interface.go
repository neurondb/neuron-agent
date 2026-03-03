/*-------------------------------------------------------------------------
 *
 * interface.go
 *    Secret store interface for NeuronAgent
 *
 * Provides interface for secret store implementations (Vault, AWS Secrets Manager, etc.).
 *
 * Copyright (c) 2024-2025, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/secrets/interface.go
 *
 *-------------------------------------------------------------------------
 */

package secrets

import (
	"context"
	"fmt"
)

/* Store defines the interface for secret stores */
type Store interface {
	/* GetSecret retrieves a secret by key */
	GetSecret(ctx context.Context, key string) (string, error)

	/* PutSecret stores a secret */
	PutSecret(ctx context.Context, key string, value string) error

	/* DeleteSecret deletes a secret */
	DeleteSecret(ctx context.Context, key string) error

	/* ListSecrets lists all secret keys (with prefix if supported) */
	ListSecrets(ctx context.Context, prefix string) ([]string, error)
}

/* Config defines configuration for secret stores */
type Config struct {
	Type      string                 `json:"type"`      // "vault", "aws", "env", "file"
	Endpoint  string                 `json:"endpoint"`  // Endpoint URL
	Token     string                 `json:"token"`     // Authentication token
	Region    string                 `json:"region"`    // AWS region
	Namespace string                 `json:"namespace"` // Vault namespace
	Metadata  map[string]interface{} `json:"metadata"`  // Additional config
}

/* NewStore creates a new secret store based on config */
func NewStore(config Config) (Store, error) {
	switch config.Type {
	case "vault":
		return NewVaultStore(config)
	case "aws":
		return NewAWSStore(config)
	case "env":
		return NewEnvStore(config)
	case "file":
		return NewFileStore(config)
	default:
		return nil, fmt.Errorf("unknown secret store type: %s", config.Type)
	}
}

