// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"fmt"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/gin-gonic/gin"
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
   WHERE documents_id = $1 AND %[1]s
),
others AS (
  SELECT
    dc2.documents_id,
	related.cve
  FROM documents_cves dc2
    JOIN related ON dc2.cve_id = related.cve_id
  WHERE dc2.documents_id <> $1
),
  SELECT
    others.documents_id,
	others.cve,
    ads.state,
    d.ssvc,
    d.title,
    ads.tracking_id
	ads.publisher
  FROM others
    JOIN documents  d   ON others.documents_id = d.id
	JOIN advisories ads ON d.advisories_id = ads.id
  WHERE %[1]s
) select * from info
`

func (c *Controller) cveRelatedDocuments(ctx *gin.Context) {
	// Get an ID from context
	docID, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}
	var (
		sb         = query.SQLBuilder{Replacements: []any{docID}}
		tlps       = c.tlps(ctx)
		allowedDoc = tlps.AsExprPublisher("ads.publisher")
		tlpCheck   = sb.CreateWhere(allowedDoc)
		tlpSQL     = fmt.Sprintf(selectCVERelatedSQL, tlpCheck)
	)

	// TODO: Query and create JSON

	_ = tlpSQL
}
