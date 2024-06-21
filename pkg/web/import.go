// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/worker"
	"github.com/csaf-poc/csaf_distribution/v3/lib/downloader"
	"github.com/gin-gonic/gin"
)

const (
	defaultWorker         = 2
	defaultPreset         = "mandatory"
	defaultForwardQueue   = 5
	defaultValidationMode = downloader.ValidationStrict
)

// importProvider downloads the advisories from the specified source
func (c *Controller) importProvider(ctx *gin.Context) {
	domainsQuery, ok := ctx.GetQuery("domains")
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing domains query parameter"})
		return
	}
	domains := strings.Split(domainsQuery, ",")

	c.downloadWorker.Enqueue(worker.DownloadJob{
		Domains:        domains,
		Worker:         defaultWorker,
		ForwardQueue:   defaultForwardQueue,
		Preset:         defaultPreset,
		ValidationMode: defaultValidationMode,
		Db:             c.db,
	})
	slog.Info("Queued download for domains", "domains", domains)

	ctx.JSON(http.StatusOK, gin.H{"msg": "queued import job"})
}
