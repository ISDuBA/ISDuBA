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
	"hash"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/ProtonMail/gopenpgp/v2/crypto"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/jackc/pgx/v5/pgxpool"
)

// download fetches the files of a document and stores
// them into the database.
func (l location) download(m *Manager, f *feed, done func()) {
	defer done()

	var (
		writers        []io.Writer
		checksum       hash.Hash
		remoteChecksum []byte
	)

	// Loading the hash
	if l.hash != nil { // ROLIE gave us an URL to hash file.
		hashFile := l.hash.String()
		switch lc := strings.ToLower(hashFile); {
		case strings.HasSuffix(lc, ".sha256"):
			checksum = sha256.New()
		case strings.HasSuffix(lc, ".sha512"):
			checksum = sha512.New()
		}
		if checksum != nil {
			var err error
			if remoteChecksum, err = f.source.loadHash(hashFile); err != nil {
				f.log(m, config.WarnFeedLogLevel, "fetching hash %q failed: %v", hashFile, err)
			} else {
				writers = append(writers, checksum)
			}
		}
	} else if !f.rolie { // If we are directory based, do some guessing:
		for _, h := range []struct {
			ext  string
			cstr func() hash.Hash
		}{
			{".sha512", sha512.New},
			{".sha256", sha256.New},
		} {
			guess := l.doc.String() + h.ext
			if rc, err := f.source.loadHash(guess); err == nil {
				remoteChecksum, checksum = rc, h.cstr()
				writers = append(writers, checksum)
				break
			}
		}
	}

	// Download the CSAF document.
	resp, err := f.source.httpGet(l.doc.String())
	if err != nil {
		f.log(m, config.ErrorFeedLogLevel, "downloading %q failed: %v", l.doc, err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		f.log(m, config.ErrorFeedLogLevel, "downloading %q failed: %s (%d)",
			l.doc, http.StatusText(resp.StatusCode), resp.StatusCode)
		return
	}
	defer resp.Body.Close()

	var data bytes.Buffer
	writers = append(writers, &data)

	// Prevent over-sized downloads.
	limited := io.LimitReader(resp.Body, int64(m.cfg.General.AdvisoryUploadLimit))

	tee := io.TeeReader(limited, io.MultiWriter(writers...))

	// Decode document into JSON.
	var doc any
	if err := json.NewDecoder(tee).Decode(&doc); err != nil {
		f.log(m, config.ErrorFeedLogLevel, "decoding document %q failed: %v", l.doc, err)
		return
	}

	// Compare checksums.
	if remoteChecksum != nil {
		if !bytes.Equal(checksum.Sum(nil), remoteChecksum) {
			f.log(m, config.ErrorFeedLogLevel, "Checksum mismatch for document %q", l.doc)
			return
		}
	}

	// Check document against schema.
	if errors, err := csaf.ValidateCSAF(doc); err != nil || len(errors) > 0 {
		if err != nil {
			f.log(m, config.ErrorFeedLogLevel,
				"Schema validation of document %q failed: %v", l.doc, err)
		} else {
			f.log(m, config.ErrorFeedLogLevel,
				"Schema validation of document %q has %d errors", l.doc, len(errors))
		}
		return
	}

	// Check against remote validator if configured.
	if m.val != nil {
		rvr, err := m.val.Validate(doc)
		if err != nil {
			slog.Error("Remote validation failed", "err", err, "url", l.doc)
			f.log(m, config.ErrorFeedLogLevel,
				"Remote validation of document %q failed: %v", l.doc, err)
			return
		}
		if !rvr.Valid {
			// XXX: Maybe we should tell more details here?!
			f.log(m, config.ErrorFeedLogLevel,
				"Remote validator classifies document %q as invalid", l.doc)
			return
		}
	}

	// Check signatures
	var signature *crypto.PGPSignature
	var signatureData []byte

	keys, err := m.openPGPKeys(f.source)
	if err != nil {
		f.log(m, config.ErrorFeedLogLevel,
			"OpenPGP signature for %q failed: %v", l.doc, err)
	} else if keys.CountEntities() > 0 {
		// Only check signature if we have something in the key ring.
		var sign *url.URL
		switch {
		case l.signature != nil: // from ROLIE feed.
			sign = l.signature
		case !f.rolie: // If we are directory based, do some guessing:
			guess := l.doc.String() + ".asc"
			sign, _ = url.Parse(guess)
		default:
			goto skipSignatureCheck
		}
		var err error
		if signature, signatureData, err = f.source.loadSignature(sign); err != nil {
			f.log(m, config.ErrorFeedLogLevel,
				"Loading OpenPGP signature for %q failed: %v", l.doc, err)
		} else {
			pm := crypto.NewPlainMessage(data.Bytes())
			if err := keys.VerifyDetached(pm, signature, crypto.GetUnixTime()); err != nil {
				f.log(m, config.ErrorFeedLogLevel,
					"Verifying OpenPGP signature of %q failed: %v", l.doc, err)
				// TODO: Survive failed signature check.
				return
			}
		}
	}
skipSignatureCheck:

	// TODO: store signature data in database.
	_ = signatureData

	// TODO: Filename check. (???)
	// TODO: Statistics

	var importer *string
	if !m.cfg.General.AnonymousEventLogging {
		importer = &m.cfg.Sources.FeedImporter
	}

	ctx := context.Background()
	if err := m.db.Run(ctx, func(ctx context.Context, conn *pgxpool.Conn) error {
		_, err := models.ImportDocumentData(
			ctx, conn,
			doc, data.Bytes(),
			importer,
			m.cfg.Sources.PublishersTLPs,
			f.storeLastChanges(&l),
			false)
		return err
	}, 0); err != nil {
		f.log(m, config.ErrorFeedLogLevel, "storing %q failed: %v", l.doc, err)
		return
	}

	f.log(m, config.InfoFeedLogLevel, "downloading %q done", l.doc)
}
