// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package web

import (
	"github.com/gin-gonic/gin"
	"github.com/ISDuBA/ISDuBA/pkg/ginkeycloak"
	"github.com/ISDuBA/ISDuBA/pkg/models"
)

// extractTLPs extracts the TLP from the JWT token.
func extractTLPs(claims func(any) error, kc *ginkeycloak.KeycloakToken) error {
	var wrapper struct {
		TLP models.PuplishersTLPs `json:"TLP"`
	}
	if err := claims(&wrapper); err != nil {
		return err
	}
	kc.CustomClaims = wrapper.TLP
	return nil
}

// tlps fetches the TLPs from the given Gin context.
func (c *Controller) tlps(ctx *gin.Context) models.PuplishersTLPs {
	token, ok := ctx.Get("token")
	if !ok {
		return c.cfg.PublishersTLPs
	}
	kct, ok := token.(*ginkeycloak.KeycloakToken)
	if !ok || kct == nil {
		return c.cfg.PublishersTLPs
	}
	tlps, ok := kct.CustomClaims.(models.PuplishersTLPs)
	if !ok {
		return c.cfg.PublishersTLPs
	}
	return tlps
}
