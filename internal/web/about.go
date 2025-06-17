// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ISDuBA/ISDuBA/internal/version"
)

// about returns the backend version number.
//
//	@Summary		Returns application information.
//	@Description	Returns general information about the application, like version.
//	@Produce		json
//	@Success		200	{object}	web.about.info
//	@Failure		401
//	@Router			/about [get]
func (c *Controller) about(ctx *gin.Context) {
	type info struct {
		Version string `json:"version"`
	}
	ctx.JSON(http.StatusOK, info{Version: version.SemVersion})
}

// view returns the publisher and tlp levels that are visible.
//
//	@Summary		Returns publisher and levels that are visible.
//	@Description	Returns information what documents the user can view and comment.
//	@Produce		json
//	@Success		200	{object}	models.PublishersTLPs
//	@Failure		401
//	@Router			/view [get]
func (c *Controller) view(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.tlps(ctx))
}

// clientConfig returns client configuration.
//
//	@Summary		Returns client configuration.
//	@Description	Returns information that the client needs to operate.
//	@Produce		json
//	@Success		200	{object}	config.Client
//	@Router			/client-config [get]
func (c *Controller) clientConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.cfg.Client)
}
