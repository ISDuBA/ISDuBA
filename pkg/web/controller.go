// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package web implements the endpoints of the web server.
package web

import (
	"database/sql"
	"log/slog"
	"net/http"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gocsaf/csaf/v3/csaf"
	sloggin "github.com/samber/slog-gin"

	"github.com/ISDuBA/ISDuBA/pkg/aggregators"
	"github.com/ISDuBA/ISDuBA/pkg/config"
	"github.com/ISDuBA/ISDuBA/pkg/database"
	"github.com/ISDuBA/ISDuBA/pkg/forwarder"
	"github.com/ISDuBA/ISDuBA/pkg/ginkeycloak"
	"github.com/ISDuBA/ISDuBA/pkg/models"
	"github.com/ISDuBA/ISDuBA/pkg/sources"
	"github.com/ISDuBA/ISDuBA/pkg/tempstore"
)

// Controller binds the endpoints to the internal logic.
type Controller struct {
	cfg *config.Config
	db  *database.DB
	fm  *forwarder.ForwardManager
	ts  *tempstore.Store
	sm  *sources.Manager
	am  *aggregators.Manager
	val csaf.RemoteValidator
}

// NewController returns a new Controller.
func NewController(
	cfg *config.Config,
	db *database.DB,
	fm *forwarder.ForwardManager,
	ts *tempstore.Store,
	dl *sources.Manager,
	am *aggregators.Manager,
	val csaf.RemoteValidator,
) *Controller {
	return &Controller{
		cfg: cfg,
		db:  db,
		fm:  fm,
		ts:  ts,
		sm:  dl,
		am:  am,
		val: val,
	}
}

// currentUser returns the current user to be used in database queries.
func (c *Controller) currentUser(ctx *gin.Context) sql.NullString {
	var user sql.NullString
	if !c.cfg.General.AnonymousEventLogging {
		user.String = ctx.GetString("uid")
		user.Valid = true
	}
	return user
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

	authRoles := func(roles ...models.WorkflowRole) gin.HandlerFunc {
		return ginkeycloak.Auth(ginkeycloak.RoleCheck(rolesAsStrings(roles)...), kcCfg)
	}

	var (
		authAdm      = authRoles(models.Admin)
		authIm       = authRoles(models.Importer)
		authEdRe     = authRoles(models.Editor, models.Reviewer)
		authEdReAdAu = authRoles(models.Editor, models.Reviewer, models.Admin, models.Auditor)
		authEdReAu   = authRoles(models.Editor, models.Reviewer, models.Auditor)
		authEdReAd   = authRoles(models.Editor, models.Reviewer, models.Admin)
		authSM       = authRoles(models.SourceManager)
		authEdSM     = authRoles(models.Editor, models.SourceManager)
		authAll      = authRoles(models.Admin, models.Importer, models.Editor,
			models.Reviewer, models.Auditor, models.SourceManager)
	)

	api := r.Group("/api")

	// Documents
	// Importer can import (POST) documents
	api.POST("/documents", authIm, c.importDocument)
	// Everyone can view (GET) overviewDocuments and viewDocuments?
	api.GET("/documents", authAll /* authEdReAu */, c.overviewDocuments)
	api.GET("/documents/:id", authAll /* authEdReAu */, c.viewDocument)
	api.GET("/documents/forward", authAll /* authEdReAu */, c.viewForwardTargets)
	api.POST("/documents/forward/:id/:target", authAll /* authEdReAu */, c.forwardDocument)
	// Admin can delete documents
	api.DELETE("/documents/:id", authAdm, c.deleteDocument)

	// Advisories
	api.DELETE("/advisory/:publisher/:trackingid", authAdm, c.deleteAdvisory)

	// Comments
	api.POST("/comments/:document", authEdReAd, c.createComment)
	api.GET("/comments/:publisher/:trackingid", authEdReAdAu, c.viewComments)
	api.PUT("/comments/post/:id", authEdReAd, c.updateComment)
	api.GET("/comments/post/:id", authEdReAdAu, c.viewComment)

	// Stored queries
	api.POST("/queries", authAll, c.createStoredQuery)
	api.POST("/queries/orders", authAll, c.updateOrder)
	api.GET("/queries", authAll, c.listStoredQueries)
	api.GET("/queries/:query", authAll, c.fetchStoredQuery)
	api.PUT("/queries/:query", authAll, c.updateStoredQuery)
	api.DELETE("/queries/:query", authAll, c.deleteStoredQuery)
	api.GET("/queries/ignore", authAll, c.getDefaultQueryExclusion)
	api.POST("/queries/ignore/:query", authAll, c.insertDefaultQueryExclusion)
	api.DELETE("/queries/ignore/:query", authAll, c.deleteDefaultQueryExclusion)

	// Events
	api.GET("/events", authEdReAdAu, c.overviewEvents)
	api.GET("/events/:publisher/:trackingid", authEdReAdAu, c.viewEvents)

	// State change
	api.PUT("/status/:publisher/:trackingid/:state", authEdReAd, c.changeStatus)
	api.PUT("/status", authEdReAd, c.changeStatusBulk)

	// SSVC change
	api.PUT("/ssvc/:document", authEdRe, c.changeSSVC)

	// Calculate diff
	api.GET("/diff/:document1/:document2", authEdRe, c.viewDiff)

	// Manage temporary documents
	api.POST("/tempdocuments", authEdReAu, c.importTempDocument)
	api.GET("/tempdocuments", authEdReAu, c.overviewTempDocuments)
	api.GET("/tempdocuments/:id", authEdReAu, c.viewTempDocument)
	api.DELETE("/tempdocuments/:id", authEdReAu, c.deleteTempDocument)

	// Backend information
	api.GET("/about", authAll, c.about)

	// Visibility information
	api.GET("/view", authAll, c.view)

	// Client configuration
	api.GET("/client-config", c.clientConfig)

	// PMD proxy
	api.GET("/pmd", authSM, c.pmd)

	// Source manager
	api.GET("/sources", authEdSM, c.viewSources)
	api.POST("/sources", authSM, c.createSource)
	api.GET("/sources/message", authAll, c.defaultMessage)
	api.GET("/sources/attention", authSM, c.attentionSources)
	api.GET("/sources/default", authSM, c.defaultSourceConfig)
	api.DELETE("/sources/:id", authSM, c.deleteSource)
	api.GET("/sources/:id", authSM, c.viewSource)
	api.PUT("/sources/:id", authSM, c.updateSource)

	// Source feeds
	api.GET("/sources/:id/feeds", authEdSM, c.viewFeeds)
	api.POST("/sources/:id/feeds", authSM, c.createFeed)
	api.GET("/sources/feeds/:id", authEdSM, c.viewFeed)
	api.PUT("/sources/feeds/:id", authSM, c.updateFeed)
	api.DELETE("/sources/feeds/:id", authSM, c.deleteFeed)
	api.GET("/sources/feeds/log", authSM, c.allFeedsLog)
	api.GET("/sources/feeds/:id/log", authSM, c.feedLog)

	// Import stats
	api.GET("/stats/imports/source/:id", authAll, c.importStatsSource)
	api.GET("/stats/imports/feed/:id", authAll, c.importStatsFeed)
	api.GET("/stats/imports", authAll, c.importStatsAllSources)
	api.GET("/stats/cve/source/:id", authAll, c.cveStatsSource)
	api.GET("/stats/cve/feed/:id", authAll, c.cveStatsFeed)
	api.GET("/stats/cve", authAll, c.cveStatsAllSources)
	api.GET("/stats/critical/source/:id", authAll, c.criticalStatsSource)
	api.GET("/stats/critical/feed/:id", authAll, c.criticalStatsFeed)
	api.GET("/stats/critical", authAll, c.criticalStatsAllSources)
	api.GET("/stats/totals", authAll, c.statsTotal)

	// Aggregators
	api.GET("/aggregator", authEdSM, c.aggregatorProxy)
	api.GET("/aggregators", authEdSM, c.viewAggregators)
	api.GET("/aggregators/:id", authEdSM, c.viewAggregator)
	api.PUT("/aggregators/:id", authSM, c.updateAggregator)
	api.GET("/aggregators/attention", authSM, c.attentionAggregators)
	api.POST("/aggregators", authSM, c.createAggregator)
	api.DELETE("/aggregators/:id", authSM, c.deleteAggregator)

	return r
}
