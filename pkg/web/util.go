// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-gonic/gin"
)

// parseTime parses a time from a given string.
func parseTime(s string) (time.Time, error) {
	for _, layout := range []string{
		time.RFC3339,
		time.DateTime,
		time.DateOnly,
	} {
		if v, err := time.Parse(layout, s); err == nil {
			return v, nil
		}
	}
	return time.Time{}, fmt.Errorf("cannot parse %q as time", s)
}

// parserMode parses parser mode from a given string.
func parserMode(s string) (query.ParserMode, error) {
	var pm query.ParserMode
	if err := pm.UnmarshalText([]byte(s)); err != nil {
		return 0, err
	}
	return pm, nil
}

// endsWith checks if the input ends with a given suffix.
func endsWith(suffix string) func(string) (string, error) {
	return func(s string) (string, error) {
		if !strings.HasSuffix(s, suffix) {
			return "", fmt.Errorf("parameter has to end with %q", suffix)
		}
		return s, nil
	}
}

// notEmpty checks if the parameter contains any text.
func notEmpty(s string) (string, error) {
	if s == "" {
		return "", errors.New("parameter is empty")
	}
	return s, nil
}

// toInt64 parses a given string to a 64bit integer.
func toInt64(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }

// parse parses a string with a given function to a value.
// If that fails a bad request status code is set in the gin context.
func parse[T any](ctx *gin.Context, conv func(string) (T, error), s string) (T, bool) {
	v, err := conv(s)
	if err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return v, false
	}
	return v, true
}
