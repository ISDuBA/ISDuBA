// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
)

func (c *Controller) changeStatus(ctx *gin.Context) {
	idS := ctx.Param("id")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	state := models.Workflow(ctx.Param("state"))
	if !state.Valid() {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("invalid state %q", state)})
		return
	}

	// TODO: Implement me!

	_ = id
	_ = state
	_ = c.hasAnyRole
}
