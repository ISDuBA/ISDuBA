// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package results represents the models used at the HTTP interface.
package results

import "github.com/gin-gonic/gin"

// ID represents a database id.
type ID struct {
	ID int64 `json:"id"`
}

// Error represents an error.
type Error struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// Success represents a success message.
type Success struct {
	Message string `json:"message"`
}

// SendSuccess sends a success message to a gin context.
func SendSuccess(ctx *gin.Context, status int, msg string) {
	ctx.JSON(status, Success{Message: msg})
}

// SendErrorMessage sends an error message to a gin context.
func SendErrorMessage(ctx *gin.Context, status int, msg string) {
	e := Error{
		Error: msg,
		Code:  status,
	}
	ctx.JSON(status, e)
}

// SendError sends a Go error to a gin context.
func SendError(ctx *gin.Context, status int, err error) {
	SendErrorMessage(ctx, status, err.Error())
}
