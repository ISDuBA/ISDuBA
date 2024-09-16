// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"github.com/gin-gonic/gin"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
	"github.com/ISDuBA/ISDuBA/pkg/ginkeycloak"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// extractTLPs extracts the TLP from the JWT token.
func extractTLPs(claims func(any) error, kc *ginkeycloak.KeycloakToken) error {
	var wrapper struct {
		TLP []models.PublishersTLPs `json:"TLP"`
	}
	if err := claims(&wrapper); err != nil {
		return err
	}
	// Merge multivalued attributes
	tlps := models.PublishersTLPs{}
	for _, tlp := range wrapper.TLP {
		for key, value := range tlp {
			_, ok := tlps[key]
			if ok {
				tlps[key] = append(tlps[key], value...)
			} else {
				tlps[key] = value
			}
		}
	}
	kc.CustomClaims = tlps
	return nil
}

// tlps fetches the TLPs from the given Gin context.
func (c *Controller) tlps(ctx *gin.Context) models.PublishersTLPs {
	token, ok := ctx.Get("token")
	if !ok {
		return c.cfg.PublishersTLPs
	}
	kct, ok := token.(*ginkeycloak.KeycloakToken)
	if !ok || kct == nil {
		return c.cfg.PublishersTLPs
	}
	tlps, ok := kct.CustomClaims.(models.PublishersTLPs)
	if !ok || len(tlps) == 0 {
		return c.cfg.PublishersTLPs
	}
	return tlps
}

// andTLPExpr adds a filter expressin to only fetch the permitted documents.
func (c *Controller) andTLPExpr(ctx *gin.Context, expr *query.Expr) *query.Expr {
	tlps := c.tlps(ctx)
	return expr.And(tlps.AsExpr())
}

// rolesAsStrings converts a slice of roles to a slice of strings.
func rolesAsStrings(roles []models.WorkflowRole) []string {
	s := make([]string, len(roles))
	for i, r := range roles {
		s[i] = string(r)
	}
	return s
}

// hasAnyRole checks if at least one of the roles is fullfilled.
func (c *Controller) hasAnyRole(ctx *gin.Context, roles ...models.WorkflowRole) bool {
	token, ok := ctx.Get("token")
	if !ok {
		return false
	}
	kct, ok := token.(*ginkeycloak.KeycloakToken)
	if !ok || kct == nil {
		return false
	}
	return kct.RealmAccess.ContainsAny(rolesAsStrings(roles))
}
