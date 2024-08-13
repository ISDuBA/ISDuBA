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

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/gin-gonic/gin"
)

// parserMode parses parser mode from a given string.
func parserMode(s string) (query.ParserMode, error) {
	var pm query.ParserMode
	if err := pm.UnmarshalText([]byte(s)); err != nil {
		return 0, err
	}
	return pm, nil
}

// toInt64 parses a given string to a 64bit integer.
func toInt64(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }

// parse parses a string with a given function to a value.
// If that fails a bad request status code is set in the gin context.
func parse[T any](ctx *gin.Context, conv func(string) (T, error), s string) (T, bool) {
	v, err := conv(s)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return v, false
	}
	return v, true
}
