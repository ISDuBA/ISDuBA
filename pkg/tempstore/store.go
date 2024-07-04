// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package tempstore implements a temporary store for documents.
package tempstore

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"slices"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
)

const cleanupDuration = 5 * time.Minute

// Store implements an in-memory storage for temporary documents.
type Store struct {
	cfg     *config.TempStore
	fns     chan func(*Store)
	done    bool
	total   int
	entries map[string][]entry
}

// Entry represents a file hold in the store.
type Entry struct {
	Inserted time.Time `json:"inserted"`
	Accessed time.Time `json:"accessed"`
	Filename string    `json:"filename"`
	Handle   int64     `json:"handle"`
}

type entry struct {
	Entry
	data []byte
}

// NewStore returns a new store.
func NewStore(cfg *config.TempStore) *Store {
	return &Store{
		cfg:     cfg,
		fns:     make(chan func(*Store)),
		entries: make(map[string][]entry),
	}
}

// Run runs the store. To be used in a Go routine.
func (st *Store) Run(ctx context.Context) {
	ticker := time.NewTicker(cleanupDuration)
	defer ticker.Stop()
	for !st.done {
		select {
		case fn := <-st.fns:
			fn(st)
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			st.cleanup(t)
		}
	}
}

// List lists the entries for a given user.
func (st *Store) List(user string) (entries []Entry) {
	done := make(chan struct{})
	st.fns <- func(st *Store) {
		defer close(done)
		now := time.Now()
		best := now.Add(-st.cfg.StorageDuration)
		userEntries := st.entries[user]
		for i := range userEntries {
			if entry := &userEntries[i]; !entry.Accessed.Before(best) {
				entries = append(entries, entry.Entry)
			}
		}
	}
	<-done
	return
}

// Delete deletes a given file for a given user.
// Returns true is file was really deleted.
func (st *Store) Delete(user string, handle int64) (deleted bool) {
	done := make(chan struct{})
	st.fns <- func(st *Store) {
		defer close(done)
		userEntries := st.entries[user]
		if len(userEntries) == 0 {
			return
		}
		entries := slices.DeleteFunc(userEntries, func(e entry) bool {
			return e.Handle == handle
		})
		switch {
		case len(entries) == 0:
			delete(st.entries, user)
			deleted = true
		case len(entries) != len(userEntries):
			st.entries[user] = entries
			deleted = true
		}
	}
	<-done
	return
}

// Fetch fetches a stored file for a given user and handle.
func (st *Store) Fetch(user string, handle int64) (r io.Reader, err error) {
	done := make(chan struct{})
	st.fns <- func(st *Store) {
		defer close(done)
		userEntries := st.entries[user]
		if len(userEntries) == 0 {
			err = errors.New("not found")
			return
		}
		now := time.Now()
		best := now.Add(-st.cfg.StorageDuration)
		for i := range userEntries {
			if entry := &userEntries[i]; entry.Handle == handle {
				if entry.Accessed.Before(best) {
					err = errors.New("expired")
				} else {
					entry.Accessed = now
					r, err = gzip.NewReader(bytes.NewReader(entry.data))
				}
				return
			}
		}
	}
	<-done
	return
}

// Store stores a file with a filename for a given user.
// Returns a unique handle to fetch it afterwards.
func (st *Store) Store(user, filename string, r io.Reader) (handle int64, err error) {
	var buf bytes.Buffer
	var w *gzip.Writer
	if w, err = gzip.NewWriterLevel(&buf, gzip.BestSpeed); err != nil {
		err = fmt.Errorf("init compression failed: %w", err)
		return
	}
	if _, err = io.Copy(w, r); err != nil {
		err = fmt.Errorf("compression failed: %w", err)
		return
	}
	if err = w.Close(); err != nil {
		err = fmt.Errorf("finish compression failed: %w", err)
		return
	}
	data := buf.Bytes()

	done := make(chan struct{})
	st.fns <- func(st *Store) {
		defer close(done)
		if st.total >= st.cfg.FilesTotal {
			err = errors.New("too many files total")
			return
		}
		userEntries := st.entries[user]
		if len(userEntries) >= st.cfg.FilesUser {
			err = errors.New("too many files per user")
			return
		}
		handle = -1
		for i := range userEntries {
			handle = max(handle, userEntries[i].Handle)
		}
		handle++
		now := time.Now()
		st.entries[user] = append(userEntries, entry{
			Entry: Entry{
				Inserted: now,
				Accessed: now,
				Filename: filename,
				Handle:   handle,
			},
			data: data,
		})
		st.total++
	}
	<-done
	return
}

// cleanup removes files from store which were idle for too long.
func (st *Store) cleanup(now time.Time) {
	best := now.Add(-st.cfg.StorageDuration)
	for user, userEntries := range st.entries {
		entries := slices.DeleteFunc(userEntries, func(e entry) bool {
			return e.Accessed.Before(best)
		})
		switch {
		case len(entries) == 0:
			delete(st.entries, user)
		case len(entries) != len(userEntries):
			st.entries[user] = entries
		}
	}
}
