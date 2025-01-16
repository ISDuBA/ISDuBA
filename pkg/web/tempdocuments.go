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
	"log/slog"
	"net/http"
	"strings"

	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/ISDuBA/ISDuBA/pkg/tempstore"
	"github.com/gin-gonic/gin"
	"github.com/gocsaf/csaf/v3/csaf"
)

// importTempDocument is an endpoint that saves a temporary document.
//
//	@Summary		Uploads a temporary document.
//	@Description	Uploads a temporary document, that can be used to create diff views.
//	@Param			file	formData	file	true	"Temporary document"
//	@Accept			multipart/form-data
//	@Produce		json
//	@Success		200	{object}	models.ID
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Router			/tempdocuments [post]
func (c *Controller) importTempDocument(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}
	f, err := file.Open()
	if err != nil {
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}
	limited := http.MaxBytesReader(
		ctx.Writer, f, int64(c.cfg.General.AdvisoryUploadLimit))
	defer limited.Close()

	user := ctx.GetString("uid")
	id, err := c.ts.Store(user, file.Filename, func(w io.Writer) error {
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
		models.SendError(ctx, http.StatusBadRequest, err)
		return
	}
	ctx.JSON(http.StatusCreated, models.ID{ID: id})
}

// overviewTempDocuments is an endpoint that returns an overview over all temporary documents.
//
//	@Summary		Returns an overview of all temporary documents.
//	@Description	An overview of all temporary documents that are uploaded by the user are returned.
//	@Produce		json
//	@Success		200	{object}	web.overviewTempDocuments.tempDocuments
//	@Failure		401
//	@Router			/tempdocuments [get]
func (c *Controller) overviewTempDocuments(ctx *gin.Context) {
	type tempDocuments struct {
		Files []tempstore.Entry `json:"files"`
		Free  int               `json:"free"`
	}
	user := ctx.GetString("uid")
	files := c.ts.List(user)
	free := max(0, min(
		c.cfg.TempStore.FilesTotal-c.ts.Total(),
		c.cfg.TempStore.FilesUser-len(files)))
	ctx.JSON(http.StatusOK, tempDocuments{
		Files: files,
		Free:  free,
	})
}

// viewTempDocument is an endpoint that returns a temporary document with the specified ID.
//
//	@Summary		Returns a temporary document.
//	@Description	Returns a temporary document with the specified ID.
//	@Param			id	path	int	true	"Document ID"
//	@Produce		json
//	@Success		200	{object}	any
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		404	{object}	models.Error
//	@Failure		500	{object}	models.Error
//	@Router			/tempdocuments/{id} [get]
func (c *Controller) viewTempDocument(ctx *gin.Context) {
	id, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}
	user := ctx.GetString("uid")
	r, entry, err := c.ts.Fetch(user, id)
	switch {
	case errors.Is(err, tempstore.ErrFileNotFound):
		models.SendError(ctx, http.StatusNotFound, err)
		return
	case err != nil:
		slog.Error("fetch temp file failed", "err", err, "id", id)
		models.SendError(ctx, http.StatusInternalServerError, err)
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

// deleteTempDocument is an endpoint that deletes a temporary document with the specified ID.
//
//	@Summary		Deletes a temporary document.
//	@Description	Deletes a temporary document with the specified ID.
//	@Param			id	path	int	true	"Document ID"
//	@Produce		json
//	@Success		200	{object}	models.Success	"deleted"
//	@Failure		400	{object}	models.Error
//	@Failure		401
//	@Failure		404	{object}	models.Error
//	@Router			/tempdocuments/{id} [delete]
func (c *Controller) deleteTempDocument(ctx *gin.Context) {
	id, ok := parse(ctx, toInt64, ctx.Param("id"))
	if !ok {
		return
	}
	user := ctx.GetString("uid")
	if c.ts.Delete(user, id) {
		models.SendSuccess(ctx, http.StatusOK, "deleted")
	} else {
		models.SendErrorMessage(ctx, http.StatusNotFound, "not found")
	}
}
