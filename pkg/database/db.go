// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package database implements the handling of the database.
package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/config"
)

// DB implements the handling with the database connection pool.
type DB struct {
	pool *pgxpool.Pool
}

// NewDB creates a new connection pool.
func NewDB(ctx context.Context, cfg *config.Database) (*DB, error) {
	cc, err := pgxpool.ParseConfig(cfg.URL())
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(ctx, cc)
	if err != nil {
		return nil, fmt.Errorf("creating postgresql pool failed: %w", err)
	}
	return &DB{pool: pool}, nil
}

// Close closes the connection pool.
func (db *DB) Close(context.Context) {
	db.pool.Close()
}

// Run a function hands over a database connection from the connection pool.
// If the given timeout is not zero the given context will be cancelled
// after this duration.
func (db *DB) Run(
	ctx context.Context,
	fn func(context.Context, *pgxpool.Conn) error,
	timeout time.Duration,
) error {
	if timeout == 0 {
		return db.pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
			return fn(ctx, conn)
		})
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return db.pool.AcquireFunc(timeoutCtx, func(conn *pgxpool.Conn) error {
		return fn(timeoutCtx, conn)
	})
}
