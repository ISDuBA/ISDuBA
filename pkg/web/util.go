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
	"strconv"

	"github.com/gin-gonic/gin"
)

// parseInt parses a given string to a 64bit integer.
// If that fails a bad request status code is set in the gin context.
func parseInt(ctx *gin.Context, s string) (int64, bool) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 0, false
	}
	return v, true
}

// parseBool parses a given string to a bool.
// If that fails a bad request status code is set in the gin context.
func parseBool(ctx *gin.Context, s string) (bool, bool) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false, false
	}
	return v, true
}
