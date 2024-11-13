// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package sources

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/gocsaf/csaf/v3/csaf"
	"github.com/gocsaf/csaf/v3/util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// dlStatus tracks the results of the different validation checks per download.
type dlStatus int

// Defines a set of status codes for the download operation using bit flags.
const (
	allSucceeded   dlStatus = 0
	downloadFailed dlStatus = 1 << iota
	filenameFailed
	schemaValidationFailed
	remoteValidationFailed
	checksumFailed
	signatureFailed
	duplicateFailed
)

func (ds *dlStatus) set(mask dlStatus) { *ds |= mask }

func (ds dlStatus) has(mask dlStatus) bool { return ds&mask == mask }

func (ds dlStatus) toInserter(i *inserter) {
	i.add("download_failed", ds.has(downloadFailed))
	i.add("filename_failed", ds.has(filenameFailed))
	i.add("schema_failed", ds.has(schemaValidationFailed))
	i.add("remote_failed", ds.has(remoteValidationFailed))
	i.add("checksum_failed", ds.has(checksumFailed))
	i.add("signature_failed", ds.has(signatureFailed))
	i.add("duplicate_failed", ds.has(duplicateFailed))
}

type inserter struct {
	keys   []string
	values []any
}

func (i *inserter) add(key string, value any) {
	i.keys = append(i.keys, key)
	i.values = append(i.values, value)
}

func (i *inserter) sql(table string) string {
	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table, strings.Join(i.keys, ","), placeholders(len(i.values)))
}

// download fetches the files of a document and stores
// them into the database.
func (l *location) download(m *Manager, f *feed) {

	var (
		strictMode     bool                     // All checks have to be fulfilled.
		signatureCheck bool                     // Take signature check seriously.
		filename       string                   // We need it later to check it against the tracking id.
		writers        []io.Writer              // Enables to decode JSON and calculating the checksum at once.
		checks         []func(*dlStatus, *feed) // List of checks to pass.
		data           bytes.Buffer             // The raw data will be stored in the database.
		signatureData  []byte                   // The signature will be stored in the database.
		client         *http.Client
	)

	// The manager owns the configuration so extract the parameters beforehand.
	m.inManager(func(m *Manager, _ context.Context) {
		strictMode = f.source.useStrictMode(m)
		signatureCheck = f.source.checkSignature(m)
		client = f.source.httpClient(m)
	})
	defer client.CloseIdleConnections()

	// checks is a list of checks to have to be passed in strict mode.
	checks = []func(ds *dlStatus, f *feed){
		// Ignore advisories with none conforming file names.
		func(ds *dlStatus, f *feed) {
			if filename = filepath.Base(l.doc.String()); !util.ConformingFileName(filename) {
				ds.set(filenameFailed)
				f.log(m, config.WarnFeedLogLevel, "File name %q is not conforming", filename)
			}
		},
	}

	// Loading the hash
	if l.hash != nil { // ROLIE gave us an URL to hash file.
		var checksum hash.Hash
		hashFile := l.hash.String()
		switch lc := strings.ToLower(hashFile); {
		case strings.HasSuffix(lc, ".sha512"):
			checksum = sha512.New()
		case strings.HasSuffix(lc, ".sha256"):
			checksum = sha256.New()
		}
		if checksum != nil {
			var check func(*dlStatus, *feed)
			if remoteChecksum, err := f.source.loadHash(client, m, hashFile); err != nil {
				check = func(ds *dlStatus, f *feed) {
					ds.set(checksumFailed)
					f.log(m, config.WarnFeedLogLevel, "Fetching hash %q failed: %v", hashFile, err)
				}
			} else {
				writers = append(writers, checksum)
				check = func(ds *dlStatus, f *feed) {
					if !bytes.Equal(checksum.Sum(nil), remoteChecksum) {
						ds.set(checksumFailed)
						f.log(m, config.ErrorFeedLogLevel, "Checksum mismatch for document %q", l.doc)
					}
				}
			}
			checks = append(checks, check)
		}
	} else if !f.rolie { // If we are directory based, do some guessing
		var checksum hash.Hash
		var remoteChecksum []byte
		for _, h := range []struct {
			ext  string
			cstr func() hash.Hash
		}{
			{".sha512", sha512.New},
			{".sha256", sha256.New},
		} {
			guess := l.doc.String() + h.ext
			if rc, err := f.source.loadHash(client, m, guess); err == nil {
				remoteChecksum, checksum = rc, h.cstr()
				break
			}
		}
		var check func(*dlStatus, *feed)
		if checksum != nil { // We found a hash
			writers = append(writers, checksum)
			check = func(ds *dlStatus, f *feed) {
				if !bytes.Equal(checksum.Sum(nil), remoteChecksum) {
					ds.set(checksumFailed)
					f.log(m, config.ErrorFeedLogLevel, "Checksum mismatch for document %q", l.doc)
				}
			}
		} else { // We didn't found a hash.
			check = func(ds *dlStatus, f *feed) {
				ds.set(checksumFailed)
				f.log(m, config.WarnFeedLogLevel, "Fetching hash for %q failed", l.doc)
			}
		}
		checks = append(checks, check)
	}

	// Keep the raw data.
	writers = append(writers, &data)

	// Download the CSAF document.
	resp, err := f.source.httpGet(client, m, l.doc.String())
	if err != nil {
		f.log(m, config.ErrorFeedLogLevel, "downloading %q failed: %v", l.doc, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		f.log(m, config.ErrorFeedLogLevel, "downloading %q failed: %s (%d)",
			l.doc, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}

	// Decode document into JSON.
	var doc any
	if err := func() error {
		defer resp.Body.Close()
		// Prevent over-sized downloads.
		limited := io.LimitReader(resp.Body, int64(m.cfg.General.AdvisoryUploadLimit))
		tee := io.TeeReader(limited, io.MultiWriter(writers...))
		return json.NewDecoder(tee).Decode(&doc)
	}(); err != nil {
		// If it is not JSON there is no way to carry on.
		f.log(m, config.ErrorFeedLogLevel, "decoding document %q failed: %v", l.doc, err)
		return
	}

	// Check if the tracking id matches the filename.
	checks = append(checks, func(ds *dlStatus, f *feed) {
		expr := util.NewPathEval()
		if err := util.IDMatchesFilename(expr, doc, filename); err != nil {
			ds.set(filenameFailed)
			f.log(m, config.ErrorFeedLogLevel, "Tracking ID in %q is not conforming: %v", l.doc, err)
		}
	})

	// Check document against schema.
	checks = append(checks, func(ds *dlStatus, f *feed) {
		if errors, err := csaf.ValidateCSAF(doc); err != nil || len(errors) > 0 {
			ds.set(schemaValidationFailed)
			if err != nil {
				f.log(m, config.ErrorFeedLogLevel,
					"Schema validation of document %q failed: %v", l.doc, err)
			} else {
				f.log(m, config.ErrorFeedLogLevel,
					"Schema validation of document %q has %d errors", l.doc, len(errors))
			}
			return
		}
	})

	// Check against remote validator if configured.
	if m.val != nil {
		checks = append(checks, func(ds *dlStatus, f *feed) {
			switch rvr, err := m.val.Validate(doc); {
			case err != nil:
				ds.set(remoteValidationFailed)
				slog.Error("Remote validation failed", "err", err, "url", l.doc)
				f.log(m, config.ErrorFeedLogLevel,
					"Remote validation of document %q failed: %v", l.doc, err)
			case !rvr.Valid:
				// XXX: Maybe we should tell more details here?!
				ds.set(remoteValidationFailed)
				f.log(m, config.ErrorFeedLogLevel,
					"Remote validator classifies document %q as invalid", l.doc)
			}
		})
	}

	// Check signatures
	keys, err := m.openPGPKeys(f.source)
	if err != nil {
		f.log(m, config.ErrorFeedLogLevel, "Loading OpenPGP keys failed: %v", err)
	} else if keys.CountEntities() > 0 {
		// Only check signature if we have something in the key ring.
		checks = append(checks, func(ds *dlStatus, f *feed) {
			var sign *url.URL
			switch {
			case l.signature != nil: // from ROLIE feed.
				sign = l.signature
			case !f.rolie: // If we are directory based, do some guessing:
				guess := l.doc.String() + ".asc"
				sign, _ = url.Parse(guess)
			default:
				// XXX: Should not happen.
				return
			}
			var err error
			var signature *crypto.PGPSignature
			if signature, signatureData, err = f.source.loadSignature(client, m, sign); err != nil {
				if signatureCheck {
					ds.set(signatureFailed)
					f.log(m, config.ErrorFeedLogLevel,
						"Loading OpenPGP signature for %q failed: %v", l.doc, err)
				}
			} else {
				pm := crypto.NewPlainMessage(data.Bytes())
				if err := keys.VerifyDetached(pm, signature, crypto.GetUnixTime()); err != nil {
					if signatureCheck {
						ds.set(signatureFailed)
						f.log(m, config.ErrorFeedLogLevel,
							"Verifying OpenPGP signature of %q failed: %v", l.doc, err)
					}
				}
			}
		})
	}

	// Run the checks.
	status := allSucceeded
	for _, check := range checks {
		check(&status, f)
	}

	if strictMode && status != allSucceeded {
		// Don't import, only write the stats.
		if err := m.db.Run(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
			var i inserter
			status.toInserter(&i)
			sql := i.sql("downloads")
			_, err := conn.Exec(ctx, sql, i.values...)
			return err
		}, 0); err != nil {
			f.log(m, config.ErrorFeedLogLevel, "storing stats of %q failed: %v", l.doc, err)
		}
		return
	}

	// Store stats in database.
	storeStats := func(ctx context.Context, tx pgx.Tx, docID int64, duplicate bool) error {
		var i inserter
		if !duplicate {
			i.add("documents_id", docID)
		} else {
			status.set(duplicateFailed)
		}
		if !f.invalid.Load() {
			i.add("feeds_id", f.id)
		}
		status.toInserter(&i)
		sql := i.sql("downloads")
		_, err := tx.Exec(ctx, sql, i.values...)
		return err
	}

	// Store signature data in database.
	storeSignature := func(ctx context.Context, tx pgx.Tx, docID int64, duplicate bool) error {
		if duplicate {
			return nil
		}
		const insertSQL = `UPDATE documents ` +
			`SET (signature, filename) = ($1, $2)` +
			`WHERE id = $3`
		_, err := tx.Exec(ctx, insertSQL, signatureData, filename, docID)
		return err
	}

	var importer *string
	if !m.cfg.General.AnonymousEventLogging {
		importer = &m.cfg.Sources.FeedImporter
	}

	switch err := m.db.Run(context.Background(), func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := models.ImportDocumentData(
			ctx, conn,
			doc, data.Bytes(),
			importer,
			m.cfg.Sources.PublishersTLPs,
			models.ChainInTx(storeStats, storeSignature, f.storeLastChanges(l)),
			false)
		return err
	}, 0); {
	case errors.Is(err, models.ErrAlreadyInDatabase):
		f.log(m, config.InfoFeedLogLevel, "not storing duplicate %q: %v", l.doc, err)
	case err != nil:
		f.log(m, config.ErrorFeedLogLevel, "storing %q failed: %v", l.doc, err)
		return
	}

	f.log(m, config.InfoFeedLogLevel, "downloading %q done", l.doc)
}
