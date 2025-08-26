// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package database

import (
	"bytes"
	"cmp"
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/jackc/pgx/v5"

	"github.com/ISDuBA/ISDuBA/internal/config"
)

//go:embed migrations
var migrations embed.FS

// migration stores the meta information extracted from the
// embedded SQL migration files and their names.
type migration struct {
	version     int64
	description string
	path        string
}

// CheckMigrations checks if the version of the database matches
// migration level of the application.
func CheckMigrations(ctx context.Context, cfg *config.Database) (bool, error) {
	migs, err := listMigrations()
	if err != nil {
		return false, err
	}

	if len(migs) == 0 {
		return false, errors.New("no migrations found")
	}

	checkVersion := func() (int64, error) {
		conn, err := pgx.Connect(ctx, cfg.ConnString())
		if err != nil {
			return -1, err
		}
		defer conn.Close(ctx)

		const selectVersion = `SELECT max(version) from versions`
		version := int64(-1)
		if err := conn.QueryRow(ctx, selectVersion).Scan(&version); err != nil {
			return -1, err
		}

		if current := migs[len(migs)-1].version; version != current {
			return version, fmt.Errorf(
				"db version (%d) mismatches app version (%d)",
				version, current)
		}
		return version, nil
	}
	version, err := checkVersion()
	if err == nil {
		if cfg.Migrate {
			return cfg.TerminateAfterMigration, nil
		}
		return false, nil
	}
	if !cfg.Migrate {
		return true, fmt.Errorf(
			"database version check failed. Maybe starting a migration helps? %w", err)
	}
	slog.Warn("Migration needed", "err", err)

	return doMigrations(ctx, cfg, version, migs)
}

func doMigrations(
	ctx context.Context,
	cfg *config.Database,
	version int64,
	migs []migration,
) (bool, error) {
	if err := func() error {
		conn, err := pgx.Connect(ctx, cfg.AdminConnString())
		if err != nil {
			return err
		}
		defer conn.Close(ctx)
		if err := createUser(ctx, conn, cfg); err != nil {
			return err
		}
		return createDatabase(ctx, conn, cfg)
	}(); err != nil {
		return true, err
	}

	conn, err := pgx.Connect(ctx, cfg.AdminUserConnString())
	if err != nil {
		return true, err
	}
	defer conn.Close(ctx)

	funcs := template.FuncMap{
		"sanitize": func(s string) string {
			return pgx.Identifier{s}.Sanitize()
		},
		"sqlQuote": sqlQuote,
	}

	for i := range migs {
		mig := &migs[i]
		if mig.version <= version {
			continue
		}
		data, err := migrations.ReadFile(mig.path)
		if err != nil {
			return true, fmt.Errorf("loading migration %q failed: %w", mig.path, err)
		}
		tmpl, err := template.New("sql").Funcs(funcs).Parse(string(data))
		if err != nil {
			return true, fmt.Errorf("parsing migration %q failed: %w", mig.path, err)
		}
		var script bytes.Buffer
		if err := tmpl.Execute(&script, cfg); err != nil {
			return true, fmt.Errorf("templating migration %q failed: %w", mig.path, err)
		}
		const insertVersion = `INSERT INTO versions (version, description) VALUES ($1, $2)`
		if err := func() error {
			slog.InfoContext(ctx, "running migration", "name", mig.description)
			ver := mig.version
			if ver == 0 { // Version 0 is special as it is intented to setup directly to lastest.
				ver = migs[len(migs)-1].version
				version = ver
			}
			// Should this script run without a transaction?
			if strings.HasSuffix(mig.description, "_notx") {
				if _, err := conn.Exec(ctx, script.String()); err != nil {
					return fmt.Errorf("executing no tx migration %q failed: %w", mig.path, err)
				}
				if _, err := conn.Exec(ctx, insertVersion, ver, mig.description); err != nil {
					return fmt.Errorf("inserting version of migration %q failed: %w", mig.path, err)
				}
				return nil
			}
			// Run inside a transaction.
			tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
			if err != nil {
				return err
			}
			defer tx.Rollback(ctx)
			if _, err := tx.Exec(ctx, script.String()); err != nil {
				return fmt.Errorf("executing migration %q failed: %w", mig.path, err)
			}
			if _, err := tx.Exec(ctx, insertVersion, ver, mig.description); err != nil {
				return fmt.Errorf("inserting version of migration %q failed: %w", mig.path, err)
			}
			return tx.Commit(ctx)
		}(); err != nil {
			return true, fmt.Errorf("applying migration %q failed: %w", mig.path, err)
		}
	}
	slog.Info("Migrations are done and database is ready")
	return cfg.TerminateAfterMigration, nil
}

func sqlQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}

func createUser(
	ctx context.Context,
	conn *pgx.Conn,
	cfg *config.Database,
) error {
	slog.InfoContext(ctx, "Check if user exists", "user", cfg.User)
	const userExists = `SELECT EXISTS (SELECT FROM pg_roles WHERE rolname = $1)`
	var exists bool
	if err := conn.QueryRow(ctx, userExists, cfg.User).Scan(&exists); err != nil {
		return fmt.Errorf("check if user exists failed: %w err", err)
	}
	if exists {
		slog.InfoContext(ctx, "User exists")
		return nil
	}
	var (
		user     = pgx.Identifier{cfg.User}.Sanitize()
		password = sqlQuote(cfg.Password)
	)
	slog.InfoContext(ctx, "Create user", "user", cfg.User)
	createUser := "CREATE USER " + user + " LOGIN PASSWORD " + password
	if _, err := conn.Exec(ctx, createUser); err != nil {
		return fmt.Errorf("creating user failed: %w err", err)
	}
	return nil
}

func createDatabase(
	ctx context.Context,
	conn *pgx.Conn,
	cfg *config.Database,
) error {
	slog.InfoContext(ctx, "Check if database exists", "database", cfg.Database)
	const dbExists = `SELECT EXISTS (SELECT FROM pg_catalog.pg_database WHERE datname = $1)`
	var exists bool
	if err := conn.QueryRow(ctx, dbExists, cfg.Database).Scan(&exists); err != nil {
		return fmt.Errorf("check if database exists failed: %w", err)
	}
	if exists {
		slog.InfoContext(ctx, "Database exists")
		return nil
	}
	var (
		db   = pgx.Identifier{cfg.Database}.Sanitize()
		user = pgx.Identifier{cfg.User}.Sanitize()
	)
	slog.InfoContext(ctx, "Create database", "database", cfg.Database)
	createDB := "CREATE DATABASE " + db + " OWNER " + user + " ENCODING 'UTF-8'"
	if _, err := conn.Exec(ctx, createDB); err != nil {
		return fmt.Errorf("creating database failed: %w", err)
	}
	return nil
}

func listMigrations() ([]migration, error) {
	entries, err := migrations.ReadDir("migrations")
	if err != nil {
		return nil, err
	}
	migReg, err := regexp.Compile(`^(\d+)-([^.]+)\.sql$`)
	if err != nil {
		return nil, err
	}
	var migs []migration
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		m := migReg.FindStringSubmatch(filepath.Base(entry.Name()))
		if m == nil {
			continue
		}
		version, err := strconv.ParseInt(m[1], 10, 64)
		if err != nil {
			return nil, err
		}
		description := m[2]
		path := "migrations/" + entry.Name()
		migs = append(migs, migration{
			version:     version,
			description: description,
			path:        path,
		})
	}
	slices.SortFunc(migs, func(a, b migration) int {
		return cmp.Compare(a.version, b.version)
	})
	return migs, nil
}
