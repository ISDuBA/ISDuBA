// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
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
		r.StaticFS("/web", http.Dir(c.cfg.Web.Static))
	}

	kcCfg := c.cfg.Keycloak.Config(extractTLPs)

	authRoles := func(roles ...string) gin.HandlerFunc {
		return ginkeycloak.Auth(ginkeycloak.RoleCheck(roles...), kcCfg)
	}

	var (
		authIm     = authRoles("importer")
		authBeRe   = authRoles("bearbeiter", "reviewer")
		authBeReAu = authRoles("bearbeiter", "reviewer", "auditor")
	)

	api := r.Group("/api")

	// Documents
	api.POST("/documents", authIm, c.importDocument)
	api.GET("/documents", authBeReAu, c.overviewDocuments)
	api.GET("/documents/:id", authBeReAu, c.viewDocument)

	// Comments
	api.POST("/comments/:document", authBeRe, c.createComment)
	api.PUT("/comments/:id", authBeRe, c.updateComment)
	api.GET("/comments/:document", authBeReAu, c.viewComments)

	return r
}
