// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/csaf-poc/csaf_distribution/v3/csaf"
	"github.com/gin-gonic/gin"
)

func (c *Controller) importTempDocument(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	f, err := file.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limited := http.MaxBytesReader(
		ctx.Writer, f, int64(c.cfg.General.AdvisoryUploadLimit))
	defer limited.Close()

	user := ctx.GetString("uid")
	id, err := c.tmpStore.Store(user, file.Filename, func(w io.Writer) error {
		// Check if the uploaded document is a valid CSAF document.
		var document any
		r := io.TeeReader(limited, w)
		if err := json.NewDecoder(r).Decode(&document); err != nil {
			return fmt.Errorf("decoding JSON failed: %w", err)
		}
		msgs, err := csaf.ValidateCSAF(document)
		if err != nil {
			return fmt.Errorf("schema validation failed: %w", err)
		}
		if len(msgs) > 0 {
			return errors.New("schema validation failed: " + strings.Join(msgs, ", "))
		}
		return nil
	})
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

func (c *Controller) overviewTempDocuments(ctx *gin.Context) {
	user := ctx.GetString("uid")
	files := c.tmpStore.List(user)
	free := max(0, min(
		c.cfg.TempStore.FilesTotal-c.tmpStore.Total(),
		c.cfg.TempStore.FilesUser-len(files)))
	ctx.JSON(http.StatusOK, gin.H{
		"files": files,
		"free":  free,
	})
}

func (c *Controller) viewTempDocument(ctx *gin.Context) {
	idS := ctx.Param("id")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := ctx.GetString("uid")
	r, entry, err := c.tmpStore.Fetch(user, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf("\"attachment; filename=\"%s\"",
			strings.ReplaceAll(entry.Filename, `"`, ``)),
	}
	ctx.DataFromReader(
		http.StatusOK, entry.Length,
		"application/json",
		r,
		extraHeaders)
}

func (c *Controller) deleteTempDocument(ctx *gin.Context) {
	idS := ctx.Param("id")
	id, err := strconv.ParseInt(idS, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := ctx.GetString("uid")
	if c.tmpStore.Delete(user, id) {
		ctx.JSON(http.StatusOK, gin.H{"message": "deleted"})
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	}
}
