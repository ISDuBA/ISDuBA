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
	"github.com/csaf-poc/csaf_distribution/v3/util"
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
	Length   int64     `json:"length"`
	ID       int64     `json:"id"`
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

// Total returns the total number of entries in the store.
func (st *Store) Total() int {
	result := make(chan int)
	st.fns <- func(st *Store) { result <- st.total }
	return <-result
}

// List lists the entries for a given user.
func (st *Store) List(user string) (entries []Entry) {
	done := make(chan struct{})
	st.fns <- func(st *Store) {
		defer close(done)
		now := time.Now()
		best := now.Add(-st.cfg.StorageDuration)
		userEntries := st.entries[user]
		entries = make([]Entry, 0, len(userEntries))
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
func (st *Store) Delete(user string, id int64) bool {
	result := make(chan bool)
	st.fns <- func(st *Store) {
		userEntries := st.entries[user]
		if len(userEntries) == 0 {
			result <- false
			return
		}
		deleted := false
		best := time.Now().Add(-st.cfg.StorageDuration)
		entries := slices.DeleteFunc(userEntries, func(e entry) bool {
			found := e.ID == id
			deleted = deleted || found
			return found || e.Accessed.Before(best)
		})
		if diff := len(userEntries) - len(entries); diff > 0 {
			st.total -= diff
			if len(entries) > 0 {
				st.entries[user] = entries
			} else {
				delete(st.entries, user)
			}
		}
		result <- deleted
	}
	return <-result
}

// Fetch fetches a stored file for a given user and id.
func (st *Store) Fetch(user string, id int64) (r io.Reader, result Entry, err error) {
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
			if entry := &userEntries[i]; entry.ID == id {
				if entry.Accessed.Before(best) {
					err = errors.New("expired")
				} else {
					entry.Accessed = now
					result = entry.Entry
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
// Returns a unique id to fetch it afterwards.
func (st *Store) Store(user, filename string, store func(io.Writer) error) (id int64, err error) {
	var buf bytes.Buffer
	var w *gzip.Writer
	if w, err = gzip.NewWriterLevel(&buf, gzip.BestSpeed); err != nil {
		err = fmt.Errorf("init compression failed: %w", err)
		return
	}
	nw := util.NWriter{Writer: w, N: 0}
	if err = store(&nw); err != nil {
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
		id = -1
		for i := range userEntries {
			id = max(id, userEntries[i].ID)
		}
		id++
		now := time.Now()
		st.entries[user] = append(userEntries, entry{
			Entry: Entry{
				Inserted: now,
				Accessed: now,
				Filename: filename,
				Length:   nw.N,
				ID:       id,
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
		if diff := len(userEntries) - len(entries); diff > 0 {
			st.total -= diff
			if len(entries) > 0 {
				st.entries[user] = entries
			} else {
				delete(st.entries, user)
			}
		}
	}
}
