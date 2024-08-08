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

	"github.com/ISDuBA/ISDuBA/pkg/sources"
	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/csaf-poc/csaf_distribution/v3/util"
	"github.com/gin-gonic/gin"
)

func (c *Controller) pmd(ctx *gin.Context) {
	var input struct {
		URL string `form:"url" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindQuery(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	header := http.Header{}
	header.Add("User-Agent", sources.UserAgent)
	client := util.HeaderClient{
		Client: &http.Client{},
		Header: header,
	}
	pmdLoader := csaf.NewProviderMetadataLoader(&client)
	lpmd := pmdLoader.Load(input.URL)
	if !lpmd.Valid() {
		h := gin.H{}
		if n := len(lpmd.Messages); n > 0 {
			msgs := make([]string, 0, n)
			for i := range lpmd.Messages {
				msgs = append(msgs, lpmd.Messages[i].Message)
			}
			h["messages"] = msgs
		}
		ctx.JSON(http.StatusBadGateway, h)
		return
	}
	ctx.JSON(http.StatusOK, lpmd.Document)
}
