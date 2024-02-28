// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
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
	db := &DB{pool: pool}

	return db, nil
}

// Close closes the connection pool.
func (db *DB) Close(context.Context) error {
	if pool := db.pool; pool != nil {
		db.pool = nil
		pool.Close()
	}
	return nil
}

// Run handles a database connection from the connection pool.
func (db *DB) Run(ctx context.Context, fn func(*pgxpool.Conn) error) error {
	return db.pool.AcquireFunc(ctx, fn)
}
