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

	"github.com/ISDuBA/ISDuBA/pkg/version"
)

// about return the backend version number.
func (c *Controller) about(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"version": version.SemVersion})
}

// view Returns the publisher and tlp levels that are visible.
func (c *Controller) view(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, c.tlps(ctx))
}

func (c *Controller) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"status": "up"})
}
