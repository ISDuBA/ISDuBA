// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"context"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
)

// Downloader fetches advisories from sources.
type Downloader struct {
	cfg  *config.Sources
	db   *database.DB
	fns  chan func(*Downloader)
	done bool
}

// NewDownloader creates a new downloader.
func NewDownloader(cfg *config.Sources, db *database.DB) *Downloader {
	return &Downloader{
		cfg: cfg,
		db:  db,
		fns: make(chan func(*Downloader)),
	}
}

// Boot loads the sources from database.
func (dl *Downloader) Boot() error {
	// TODO: Implement me!:
	return nil
}

// Run runs the downloader. To be used in a Go routine.
func (dl *Downloader) Run(ctx context.Context) {
	for !dl.done {
		select {
		case fn := <-dl.fns:
			fn(dl)
		case <-ctx.Done():
			return
		}
	}
}

// Kill stops the downloader.
func (dl *Downloader) Kill() {
	dl.fns <- func(l *Downloader) { l.done = true }
}
