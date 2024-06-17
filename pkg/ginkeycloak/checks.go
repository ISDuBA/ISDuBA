// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package ginkeycloak

import (
	"github.com/gin-gonic/gin"
)

// AccessCheckFunction is function return true if access should be granted.
type AccessCheckFunction func(*TokenContainer, *gin.Context) bool

// AccessBuilderOption is an option to configure an AccessBuilder
type AccessBuilderOption func(*AccessBuilder)

// AccessBuilder is a helper to produce a Gin middleware based
// on the given constraints.
type AccessBuilder struct {
	cfg      *Config
	accesses []AccessCheckFunction
	uids     []string
	roles    []string
}

// RoleCheck returns a check which passes if any of the given roles
// is part of the claims.
func RoleCheck(roles ...string) AccessCheckFunction {
	return func(tc *TokenContainer, _ *gin.Context) bool {
		return tc.KeycloakToken.RealmAccess.ContainsAny(roles)
	}
}
