// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package forwarder

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// forwarderWakeupInterval is a saftey net wakeup interval for each
// forwarder if a ping from the manager is missed some how.
const forwarderWakeupInterval = 2 * time.Minute

type forwarder struct {
	cfg     *config.ForwardTarget
	db      *database.DB
	fns     chan (func(*forwarder))
	done    bool
	client  *http.Client
	headers http.Header
}

func newForwarder(cfg *config.ForwardTarget, db *database.DB) (*forwarder, error) {
	// Init http clients
	var tlsConfig tls.Config
	if cfg.ClientPrivateCert != "" && cfg.ClientPublicCert != "" {
		clientCert, err := tls.LoadX509KeyPair(
			cfg.ClientPublicCert,
			cfg.ClientPrivateCert)
		if err != nil {
			return nil, fmt.Errorf(
				"cannot load client cert for forward target %q: %w",
				cfg.URL, err)
		}
		tlsConfig.Certificates = []tls.Certificate{clientCert}
	}
	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
		},
	}
	headers := make(http.Header, len(cfg.Header))
	for _, header := range cfg.Header {
		if k, v, ok := strings.Cut(header, ":"); ok {
			headers.Add(k, v)
			continue
		}
		return nil, fmt.Errorf(
			"header %q of forwarder target %q is missing ':'",
			header, cfg.URL)
	}
	return &forwarder{
		cfg:     cfg,
		db:      db,
		fns:     make(chan func(*forwarder)),
		client:  client,
		headers: headers,
	}, nil
}

func (f *forwarder) run(ctx context.Context) {
	ticker := time.NewTicker(forwarderWakeupInterval)
	defer ticker.Stop()
	for !f.done {
		if err := f.forward(ctx); err != nil {
			slog.Error("forwarder has issues", "error", err, "forwarder", f.cfg.URL)
		}
		select {
		case fn := <-f.fns:
			fn(f)
		case <-ctx.Done():
			return
		}
	}
}

func (f *forwarder) forward(ctx context.Context) error {
	for {
		docIDs, err := f.loadDocIDs(ctx)
		if err != nil {
			return fmt.Errorf(
				"loading document ids failed: %w", err)
		}
		if len(docIDs) == 0 {
			break
		}
		if err := f.loadForwardDocuments(ctx, docIDs); err != nil {
			return fmt.Errorf("load and forwarding docs failed: %w", err)
		}
	}
	return nil
}

func (f *forwarder) loadDocIDs(ctx context.Context) ([]int64, error) {
	const docIDsSQL = `` +
		`SELECT` +
		` documents_id ` +
		`FROM forwarders_queue fwq ` +
		`JOIN forwarders fw ON fwq.forwarders_id = fw.id ` +
		`WHERE` +
		` fwq.state IN ('pending', 'failed')` +
		` AND fw.url = $1 ` +
		`ORDER BY upload_order DESC ` +
		`LIMIT 20` // Poll in smaller batches.
	var docIDs []int64
	if err := f.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := conn.Query(rctx, docIDsSQL, f.cfg.URL)
			var err error
			docIDs, err = pgx.CollectRows(
				rows,
				func(row pgx.CollectableRow) (int64, error) {
					var docID int64
					err := row.Scan(&docID)
					return docID, err
				})
			return err
		}, 0,
	); err != nil {
		return nil, fmt.Errorf(
			"scanning for pending/failed documents failed: %w", err)
	}
	return docIDs, nil
}

func (f *forwarder) forwardDocument(ctx context.Context, docID int64) error {
	const documentSQL = `` +
		`SELECT` +
		` original,` +
		` filename,` +
		` publisher,` +
		` (filename_failed OR remote_failed OR checksum_failed OR signature_failed) ` +
		`FROM documents` +
		` JOIN downloads ON documents.id = downloads.documents_id` +
		` JOIN advisories ON documents.advisories_id = advisories.id ` +
		`WHERE` +
		` documents.id = $1`
	var (
		doc              []byte
		filename         *string
		failedValidation *bool
		publisher        string
	)
	switch err := f.db.Run(
		ctx,
		func(rctx context.Context, conn *pgxpool.Conn) error {
			return conn.QueryRow(rctx, documentSQL, docID).Scan(
				&doc,
				&filename,
				&publisher,
				&failedValidation)
		}, 0,
	); {
	case errors.Is(err, pgx.ErrNoRows):
		return errors.New("document not found")
	case err != nil:
		return fmt.Errorf("loading document failed: %w", err)
	}
	if !f.acceptsPublisher(publisher) {
		return errors.New("not allowed to forward to target")
	}
	// Build the request.
	req, err := buildRequest(
		doc,
		filename,
		parseValidationStatus(failedValidation),
		f.cfg.URL,
		f.headers)
	if err != nil {
		return fmt.Errorf("building request failed: %w", err)
	}
	// Try to forward.
	res, err := f.client.Do(req)
	if err != nil {
		return fmt.Errorf("sending request failed: %w", err)
	}
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"forwarding failed: code: %d, status: %q", res.StatusCode, res.Status)
	}
	return nil
}

func (f *forwarder) loadForwardDocuments(
	ctx context.Context,
	docIDs []int64,
) error {
	const (
		documentSQL = `` +
			`SELECT` +
			` original,` +
			` filename,` +
			` (filename_failed OR remote_failed OR checksum_failed OR signature_failed) ` +
			`FROM documents` +
			` JOIN downloads ON documents.id = downloads.documents_id` +
			`WHERE` +
			` documents.id = $1`
		updateQueueSQL = `` +
			`UPDATE forwarders_queue` +
			` SET state = $1::forward_state ` +
			`WHERE` +
			` documents_id = $2 AND` +
			` forwarders_id = (SELECT id FROM forwarders WHERE url = $3)`
	)
	for _, docID := range docIDs {
		var (
			doc              []byte
			filename         *string
			failedValidation *bool
		)
		switch err := f.db.Run(
			ctx,
			func(rctx context.Context, conn *pgxpool.Conn) error {
				return conn.QueryRow(rctx, documentSQL, docID).Scan(
					&doc,
					&filename,
					&failedValidation)
			}, 0,
		); {
		case errors.Is(err, pgx.ErrNoRows):
			// Document may be deleted -> ignore it.
			continue
		case err != nil:
			return fmt.Errorf("loading document failed: %w", err)
		}
		// Build the request.
		req, err := buildRequest(
			doc,
			filename,
			parseValidationStatus(failedValidation),
			f.cfg.URL,
			f.headers)
		if err != nil {
			return fmt.Errorf("building request failed: %w", err)
		}
		res, err := f.client.Do(req)
		if err != nil {
			slog.Warn(
				"forwarder",
				"msg", "sending request request failed",
				"err", err)
			// Try again later.
			continue
		}
		var result string
		if res.StatusCode == http.StatusCreated {
			result = "uploaded"
		} else {
			slog.Warn(
				"forwarder",
				"error", "failed",
				"code", res.StatusCode,
				"status", res.Status)
			result = "failed"
		}
		// Update the queue to the result of the upload.
		if err := f.db.Run(
			ctx,
			func(rctx context.Context, conn *pgxpool.Conn) error {
				_, err := conn.Exec(
					rctx, updateQueueSQL,
					result, docID, f.cfg.URL)
				return err
			}, 0,
		); err != nil {
			return fmt.Errorf("updating queue failed: %w", err)
		}
	}
	return nil
}

func (f *forwarder) kill() {
	f.fns <- func(f *forwarder) { f.done = true }
}

// ping should be called by the manager to signal the fowarder
// that there are documents to forward.
func (f *forwarder) ping() {
	f.fns <- func(*forwarder) {}
}

func (f *forwarder) acceptsPublisher(publisher string) bool {
	return f.cfg.Publisher == nil || publisher == *f.cfg.Publisher
}
