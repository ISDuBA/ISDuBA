// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package main implements an example bulk importer.
package main

import (
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/ISDuBA/ISDuBA/pkg/version"
)

func processFile(
	ctx context.Context,
	db *database.DB,
	dry bool,
	importer, file string,
) error {
	var actor *string
	if importer != "" {
		actor = &importer
	}
	walk := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.Type().IsRegular() {
			return nil
		}
		lower := strings.ToLower(path)
		if !(strings.HasSuffix(lower, ".json") || strings.HasSuffix(lower, "json.gz")) {
			return nil
		}

		slog.Info("processing document", "file", filepath.Base(path))

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		var r io.Reader
		if strings.HasSuffix(lower, ".gz") {
			if r, err = gzip.NewReader(file); err != nil {
				return err
			}
		} else {
			r = file
		}

		// Store stats in database.
		storeStats := func(ctx context.Context, tx pgx.Tx, docID int64, duplicate bool) error {
			if duplicate {
				return nil
			}
			const insertSQL = `INSERT INTO downloads ` +
				`(documents_id, feeds_id) VALUES ($1, ` +
				`(SELECT id FROM feeds WHERE sources_id = 0 AND label = 'bulk'))`
			_, err := tx.Exec(ctx, insertSQL, docID)
			return err
		}

		var id int64
		if err = db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
			id, err = models.ImportDocument(
				ctx, conn, r, actor,
				nil,
				models.ChainInTx(storeStats, models.StoreFilename(filepath.Base(path))),
				dry)
			return err
		}, 0); err != nil {
			if errors.Is(err, models.ErrAlreadyInDatabase) {
				slog.Warn("advisory already in database", "file", filepath.Base(path))
				err = nil
			}
			return err
		}
		slog.Info("inserted", "id", id)
		return nil
	}
	return filepath.WalkDir(file, walk)
}

func process(creds *config.Database, dry bool, importer string, files []string) error {
	start := time.Now()
	defer func() {
		slog.Info("processing took", "duration", time.Since(start))
	}()
	ctx := context.Background()

	db, err := database.NewDB(ctx, creds)
	if err != nil {
		return err
	}
	defer db.Close(ctx)

	conn, err := pgx.Connect(ctx, creds.URL())
	if err != nil {
		return err
	}
	defer conn.Close(ctx)

	for _, file := range files {
		if err := processFile(ctx, db, dry, importer, file); err != nil {
			return fmt.Errorf("processing %q failed: %w", file, err)
		}
	}
	return nil
}

func check(err error) {
	if err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func userName() string {
	user, err := user.Current()
	if err != nil {
		return "me"
	}
	return user.Username
}

func main() {
	var (
		creds       config.Database
		importer    string
		dry         bool
		showVersion bool
	)
	flag.StringVar(&creds.Database, "database", "isduba", "database name")
	flag.StringVar(&creds.User, "user", "isduba", "database user")
	flag.StringVar(&creds.Password, "password", "isduba", "password")
	flag.StringVar(&creds.Host, "host", "localhost", "database host")
	flag.IntVar(&creds.Port, "port", 5432, "database host")
	flag.BoolVar(&dry, "dry", false, "dont store values")
	flag.BoolVar(&showVersion, "version", false, "show version information")
	flag.StringVar(&importer, "importer", userName(), "importing person")
	flag.Parse()
	if showVersion {
		fmt.Printf("%s version: %s\n", os.Args[0], version.SemVersion)
		os.Exit(0)
	}
	check(process(&creds, dry, importer, flag.Args()))
}
