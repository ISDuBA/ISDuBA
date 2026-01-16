// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const selectCVERelatedSQL = `
WITH related AS (
   SELECT
     dc.documents_id,
     uc.cve,
     dc.cve_id
   FROM documents_cves dc
     JOIN unique_cves uc ON dc.cve_id          = uc.id
     JOIN documents docs ON dc.documents_id    = docs.id
     JOIN advisories ads ON docs.advisories_id = ads.id
   WHERE documents_id = :doc_placeholder AND %[1]s
),
others AS (
  SELECT
    dc2.documents_id,
    related.cve
  FROM documents_cves dc2
    JOIN related ON dc2.cve_id = related.cve_id
  WHERE dc2.documents_id <> :doc_placeholder
)
SELECT
  others.documents_id,
  others.cve,
  ads.state,
  d.ssvc,
  d.title,
  ads.tracking_id,
  ads.publisher
FROM others
  JOIN documents  d   ON others.documents_id = d.id
  JOIN advisories ads ON d.advisories_id     = ads.id
WHERE %[1]s
`

// cveRelatedDocuments is an endpoint that returns the documents
// that are related with the given document by CVEs.
//
//	@Summary        Returns CVE related documents.
//	@Description    Returns the documents related to this document by CVE.
//	@Produce        json
//	@Success        200 {array}     models.RelatedDocument
//	@Failure        400 {object}    models.Error
//	@Failure        401
//	@Failure        500 {object}    models.Error
//	@Router         /documents/{id}/cve_related [get]
func (c *Controller) cveRelatedDocuments(ctx *gin.Context) {
	// Get an ID from context
	docID, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}

	var (
		sb = query.SQLBuilder{
			Replacements: []any{},
		}
		tlps       = c.tlps(ctx)
		allowedDoc = tlps.AsExprPublisher("ads.publisher")
		tlpCheck   = sb.CreateWhere(allowedDoc)
	)

	sb.Replacements = append(sb.Replacements, docID)
	docIndex := len(sb.Replacements)
	selectSQL := fmt.Sprintf(selectCVERelatedSQL, tlpCheck)

	var relatedDocuments []*models.RelatedDocument

	selectSQL = strings.ReplaceAll(selectSQL, ":doc_placeholder", fmt.Sprintf("$%d", docIndex))

	if err := c.db.Run(
		ctx.Request.Context(),
		func(rctx context.Context, conn *pgxpool.Conn) error {
			rows, _ := conn.Query(rctx, selectSQL, sb.Replacements...)
			var err error
			relatedDocuments, err = pgx.CollectRows(rows,
				func(row pgx.CollectableRow) (*models.RelatedDocument, error) {
					var related models.RelatedDocument
					if err := row.Scan(
						&related.DocumentID,
						&related.CVE,
						&related.State,
						&related.SSVC,
						&related.Title,
						&related.TrackingID,
						&related.Publisher,
					); err != nil {
						return nil, err
					}
					return &related, nil
				})
			return err
		}, 0,
	); err != nil {
		slog.Error("database error", "err", err)
		models.SendError(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, relatedDocuments)
}
