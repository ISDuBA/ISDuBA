// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package web implements the endpoints of the web server.
package web

import (
	"log/slog"
	"net/http"

	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/ginkeycloak"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

// Controller binds the endpoints to the internal logic.
type Controller struct {
	cfg *config.Config
	db  *database.DB
}

// NewController returns a new Controller.
func NewController(cfg *config.Config, db *database.DB) *Controller {
	return &Controller{cfg: cfg, db: db}
}

// Bind return a http handler to be used in a web server.
func (c *Controller) Bind() http.Handler {
	r := gin.New()
	r.Use(sloggin.New(slog.Default()))
	r.Use(gin.Recovery())

	if c.cfg.Web.Static != "" {
		r.Use(static.Serve("/", static.LocalFile(c.cfg.Web.Static, false)))
	}

	kcCfg := c.cfg.Keycloak.Config(extractTLPs)

	authRoles := func(roles ...string) gin.HandlerFunc {
		return ginkeycloak.Auth(ginkeycloak.RoleCheck(roles...), kcCfg)
	}

	var (
		authIm     = authRoles(models.Importer)
		authEdRe   = authRoles(models.Editor, models.Reviewer)
		authEdReAu = authRoles(models.Editor, models.Reviewer, models.Auditor)
		authEdReAd = authRoles(models.Editor, models.Reviewer, models.Admin)
	)

	api := r.Group("/api")

	// Documents
	api.POST("/documents", authIm, c.importDocument)
	api.GET("/documents", authEdReAu, c.overviewDocuments)
	api.GET("/documents/:id", authEdReAu, c.viewDocument)

	// Comments
	api.POST("/comments/:document", authEdRe, c.createComment)
	api.PUT("/comments/:id", authEdRe, c.updateComment)
	api.GET("/comments/:document", authEdReAu, c.viewComments)

	// State change
	api.PUT("/status/:id/:state", authEdReAd, c.changeStatus)

	return r
}
