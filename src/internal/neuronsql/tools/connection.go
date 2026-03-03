/*-------------------------------------------------------------------------
 *
 * connection.go
 *    SafeConnection: restricted PG connection with timeouts and read-only
 *
 * Copyright (c) 2024-2026, neurondb, Inc. <support@neurondb.ai>
 *
 * IDENTIFICATION
 *    NeuronAgent/internal/neuronsql/tools/connection.go
 *
 *-------------------------------------------------------------------------
 */

package tools

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/neurondb/NeuronAgent/internal/db"
	"github.com/neurondb/NeuronAgent/internal/neuronsql/policy"
	"github.com/neurondb/NeuronAgent/internal/utils"
)

/* SafeConnectionConfig holds timeouts and limits */
type SafeConnectionConfig struct {
	StatementTimeout time.Duration
	LockTimeout      time.Duration
	IdleInTxTimeout  time.Duration
	MaxRows          int
	MaxResultBytes   int
	MaxSQLLength     int /* max SQL statement length in bytes; 0 = default 64KiB */
}

const defaultMaxSQLLength = 64 * 1024 // 64 KiB

/* DefaultSafeConnectionConfig returns default limits */
func DefaultSafeConnectionConfig() SafeConnectionConfig {
	return SafeConnectionConfig{
		StatementTimeout: 5 * time.Second,
		LockTimeout:      3 * time.Second,
		IdleInTxTimeout:  10 * time.Second,
		MaxRows:          1000,
		MaxResultBytes:   5 * 1024 * 1024, // 5 MiB
		MaxSQLLength:     defaultMaxSQLLength,
	}
}

/* SafeConnection wraps a DB connection with session-level safety settings */
type SafeConnection struct {
	db     *db.DB
	policy *policy.PolicyEngineImpl
	config SafeConnectionConfig
}

/* NewSafeConnection creates a SafeConnection from a DSN */
func NewSafeConnection(dsn string, policyEngine *policy.PolicyEngineImpl, config SafeConnectionConfig) (*SafeConnection, error) {
	database, err := db.NewDB(dsn, db.PoolConfig{
		MaxOpenConns:    5,
		MaxIdleConns:    2,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	})
	if err != nil {
		return nil, err
	}
	return &SafeConnection{
		db:     database,
		policy: policyEngine,
		config: config,
	}, nil
}

/* Close closes the underlying connection pool */
func (c *SafeConnection) Close() error {
	return c.db.Close()
}

/* RunReadOnly runs fn inside a read-only transaction with timeouts; requestID is used for application_name */
func (c *SafeConnection) RunReadOnly(ctx context.Context, requestID string, fn func(tx *sqlx.Tx) error) error {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	appName := "NeuronSQL"
	if requestID != "" {
		appName = "NeuronSQL:" + requestID
	}
	for _, q := range []struct {
		s string
		a []interface{}
	}{
		{"SET LOCAL application_name = $1", []interface{}{appName}},
		{"SET LOCAL statement_timeout = $1", []interface{}{c.config.StatementTimeout.Milliseconds()}},
		{"SET LOCAL lock_timeout = $1", []interface{}{c.config.LockTimeout.Milliseconds()}},
		{"SET LOCAL idle_in_transaction_session_timeout = $1", []interface{}{c.config.IdleInTxTimeout.Milliseconds()}},
		{"SET TRANSACTION READ ONLY", nil},
	} {
		if q.a != nil {
			if _, err := tx.ExecContext(ctx, q.s, q.a...); err != nil {
				return err
			}
		} else {
			if _, err := tx.ExecContext(ctx, q.s); err != nil {
				return err
			}
		}
	}
	return fn(tx)
}

/* DB returns the underlying db for use when not in a transaction (e.g. one-off queries with same session vars) */
func (c *SafeConnection) DB() *db.DB {
	return c.db
}

/* Config returns the safe connection config */
func (c *SafeConnection) Config() SafeConnectionConfig {
	return c.config
}

/* Policy returns the policy engine */
func (c *SafeConnection) Policy() *policy.PolicyEngineImpl {
	return c.policy
}

/* ConnectionFactory creates SafeConnection from DSN (e.g. for orchestrator or registry) */
type ConnectionFactory func(ctx context.Context, dsn string) (*SafeConnection, error)

/* NewConnectionFactory returns a factory that creates new connections each time */
func NewConnectionFactory(policyEngine *policy.PolicyEngineImpl, config SafeConnectionConfig) ConnectionFactory {
	return func(ctx context.Context, dsn string) (*SafeConnection, error) {
		return NewSafeConnection(dsn, policyEngine, config)
	}
}

/* BuildDSN builds a PostgreSQL connection string from components */
func BuildDSN(host string, port int, user, password, database string) string {
	return utils.BuildConnectionString(host, port, user, password, database, "")
}
