// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package ginkeycloak

import (
	"slices"

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

// NewAccessBuilder returns a new AccessBuilder for a given Keycloak configuration.
func NewAccessBuilder(cfg *Config) *AccessBuilder {
	return &AccessBuilder{cfg: cfg}
}

// Clone returns a copy of the given AccessBuilder.
func (ab *AccessBuilder) Clone() *AccessBuilder {
	return &AccessBuilder{
		cfg:      ab.cfg,
		accesses: slices.Clone(ab.accesses),
		uids:     slices.Clone(ab.uids),
		roles:    slices.Clone(ab.roles),
	}
}

// With configures an AccessBuilder.
func (ab *AccessBuilder) With(options ...AccessBuilderOption) *AccessBuilder {
	for _, option := range options {
		option(ab)
	}
	return ab
}

// Allow return a custom check for a AccessBuilder.
func Allow(acfn AccessCheckFunction) AccessBuilderOption {
	return func(ab *AccessBuilder) {
		ab.accesses = append(ab.accesses, acfn)
	}
}

// AllowUID returns a check for a given user id.
func AllowUID(uid string) AccessBuilderOption {
	return func(ab *AccessBuilder) {
		ab.uids = append(ab.uids, uid)
	}
}

// AllowRole return a check for a given role.
func AllowRole(realm string) AccessBuilderOption {
	return func(ab *AccessBuilder) {
		ab.roles = append(ab.roles, realm)
	}
}

func (ab *AccessBuilder) build(cond func(...AccessCheckFunction) AccessCheckFunction) gin.HandlerFunc {
	var acfns []AccessCheckFunction
	if len(ab.accesses) > 0 {
		acfns = append(acfns, ab.accesses...)
	}
	if len(ab.uids) > 0 {
		acfns = append(acfns, UIDCheck(ab.uids...))
	}
	if len(ab.roles) > 0 {
		acfns = append(acfns, RoleCheck(ab.roles...))
	}
	return Auth(cond(acfns...), ab.cfg)
}

// Any returns a Gin middleware which grants access if any of the
// configured checks passes.
func (ab *AccessBuilder) Any() gin.HandlerFunc {
	return ab.build(anyCheck)
}

// All returns a Gin middleware which grants access if all of the
// configured checks pass.
func (ab *AccessBuilder) All() gin.HandlerFunc {
	return ab.build(allCheck)
}

// RoleCheck returns a check which passes if any of the given roles
// is part of the claims.
func RoleCheck(roles ...string) AccessCheckFunction {
	return func(tc *TokenContainer, _ *gin.Context) bool {
		return tc.KeycloakToken.RealmAccess.ContainsAny(roles)
	}
}

// UIDCheck returns a check which passes if any of the given userids
// is part of the claims.
func UIDCheck(uids ...string) AccessCheckFunction {
	return func(tc *TokenContainer, _ *gin.Context) bool {
		user := tc.KeycloakToken.PreferredUsername
		for _, uid := range uids {
			if uid == user {
				return true
			}
		}
		return false
	}
}

// AuthCheck always grants access.
func AuthCheck() AccessCheckFunction {
	return func(_ *TokenContainer, _ *gin.Context) bool {
		return true
	}
}

func allCheck(acfns ...AccessCheckFunction) AccessCheckFunction {
	return func(tc *TokenContainer, ctx *gin.Context) bool {
		for _, acfn := range acfns {
			if !acfn(tc, ctx) {
				return false
			}
		}
		return true
	}
}

func anyCheck(acfns ...AccessCheckFunction) AccessCheckFunction {
	return func(tc *TokenContainer, ctx *gin.Context) bool {
		for _, acfn := range acfns {
			if acfn(tc, ctx) {
				return true
			}
		}
		return false
	}
}
