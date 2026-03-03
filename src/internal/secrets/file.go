/*-------------------------------------------------------------------------
 *
 * file.go
 *    File-based secret store implementation
 *
 * Provides file-based secret store for local development/testing.
 *
 * Copyright (c) 2024-2025, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/secrets/file.go
 *
 *-------------------------------------------------------------------------
 */

package secrets

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/* FileStore implements Store interface using files */
type FileStore struct {
	baseDir string
}

/* NewFileStore creates a new file-based secret store */
func NewFileStore(config Config) (*FileStore, error) {
	baseDir := "/etc/neurondb/secrets"
	if dirVal, ok := config.Metadata["base_dir"].(string); ok {
		baseDir = dirVal
	}

	/* Create directory if it doesn't exist */
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, fmt.Errorf("failed to create secrets directory: %w", err)
	}

	return &FileStore{baseDir: baseDir}, nil
}

/* GetSecret retrieves a secret from file */
func (f *FileStore) GetSecret(ctx context.Context, key string) (string, error) {
	filePath := filepath.Join(f.baseDir, key)
	
	/* Prevent directory traversal */
	if !strings.HasPrefix(filepath.Clean(filePath), f.baseDir) {
		return "", fmt.Errorf("invalid secret key: %s", key)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("secret not found: %s", key)
		}
		return "", fmt.Errorf("failed to read secret file: %w", err)
	}

	return string(data), nil
}

/* PutSecret stores a secret in file */
func (f *FileStore) PutSecret(ctx context.Context, key string, value string) error {
	filePath := filepath.Join(f.baseDir, key)
	
	/* Prevent directory traversal */
	if !strings.HasPrefix(filepath.Clean(filePath), f.baseDir) {
		return fmt.Errorf("invalid secret key: %s", key)
	}

	/* Create parent directory if needed */
	if err := os.MkdirAll(filepath.Dir(filePath), 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	/* Write file with restrictive permissions */
	if err := os.WriteFile(filePath, []byte(value), 0600); err != nil {
		return fmt.Errorf("failed to write secret file: %w", err)
	}

	return nil
}

/* DeleteSecret deletes a secret file */
func (f *FileStore) DeleteSecret(ctx context.Context, key string) error {
	filePath := filepath.Join(f.baseDir, key)
	
	/* Prevent directory traversal */
	if !strings.HasPrefix(filepath.Clean(filePath), f.baseDir) {
		return fmt.Errorf("invalid secret key: %s", key)
	}

	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("secret not found: %s", key)
		}
		return fmt.Errorf("failed to delete secret file: %w", err)
	}

	return nil
}

/* ListSecrets lists secrets from files */
func (f *FileStore) ListSecrets(ctx context.Context, prefix string) ([]string, error) {
	var keys []string
	
	err := filepath.Walk(f.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			relPath, err := filepath.Rel(f.baseDir, path)
			if err != nil {
				return err
			}
			if prefix == "" || strings.HasPrefix(relPath, prefix) {
				keys = append(keys, relPath)
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	return keys, nil
}

