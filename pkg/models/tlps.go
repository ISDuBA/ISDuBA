// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

import (
	"fmt"
	"slices"

	"github.com/ISDuBA/ISDuBA/pkg/database/query"
)

type (
	// TLP represents a Traffic Light Protocol 1 value.
	TLP string
	// Publisher represents the publisher name.
	Publisher string
	// PublishersTLPs is a list of TLPs per publisher.
	PublishersTLPs map[Publisher][]TLP
)

// The different TLP levels
const (
	TLPWhite       TLP = "WHITE"        // TLPWhite represents TLP:WHITE
	TLPClear       TLP = "CLEAR"        // TLPClear represents TLP:CLEAR (replaces WHITE in CSAF 2.1)
	TLPGreen       TLP = "GREEN"        // TLPGreen represents TLP:GREEN
	TLPAmber       TLP = "AMBER"        // TLPAmber represents TLP:AMBER
	TLPAmberStrict TLP = "AMBER+STRICT" // TLPAmberStrict represents TLP:AMBER+STRICT (only in CSAF 2.1)
	TLPRed         TLP = "RED"          // TLPRed   represents TLP:RED
)

// UnmarshalText implements [encoding.TextUnmarshaler].
func (tlp *TLP) UnmarshalText(text []byte) error {
	s := TLP(text)
	switch s {
	case TLPWhite, TLPClear, TLPGreen, TLPAmber, TLPAmberStrict, TLPRed:
		*tlp = s
		return nil
	default:
		return fmt.Errorf("unknown TLP value: %q", s)
	}
}

// temporaryTLPMapping temporarily maps (WHITE -> CLEAR) and (AMBER -> AMBER+STRICT) rules
// TODO: This will need to be migrated to proper TLP 2.0 handling later
var temporaryTLPMapping = map[TLP][]TLP{
	TLPWhite: {TLPClear},
	TLPAmber: {TLPAmberStrict},
}

// mappedTLPs returns list of tlps for which rights exist, as well as mapped ones added temporarily
// TODO: This is a temporary change until TLP 2 migration
func mappedTLPs(tlps []TLP) []TLP {
	mapped := make([]TLP, 0, len(tlps)+2)
	add := func(tlp TLP) {
		if !slices.Contains(mapped, tlp) {
			mapped = append(mapped, tlp)
		}
	}
	for _, tlp := range tlps {
		add(tlp)
		for _, tmpTLP := range temporaryTLPMapping[tlp] {
			add(tmpTLP)
		}
	}
	return mapped
}

// Allowed checks if a pair of publisher/tlp is allowed.
func (ptlps PublishersTLPs) Allowed(publisher string, tlp TLP) bool {
	if p, ok := ptlps[Publisher(publisher)]; ok {
		return slices.Contains(mappedTLPs(p), tlp)
	}
	wildcard, ok := ptlps["*"]
	return ok && slices.Contains(mappedTLPs(wildcard), tlp)
}

// or transforms a slice of string kinds into list of or-ed string field accesses.
func or[T ~string](field string, tlps []T) *query.Expr {
	var ts *query.Expr
	for _, t := range tlps {
		if tlp := query.FieldEqString(field, string(t)); ts == nil {
			ts = tlp
		} else {
			ts = ts.Or(tlp)
		}
	}
	return ts
}

// AsExpr returns the list of TLP rules as an expression tree.
func (ptlps PublishersTLPs) AsExpr() *query.Expr {
	return ptlps.AsExprPublisher("publisher")
}

// AsExprPublisher returns the list of TLP rules as an expression tree with a given
// publisher field name.
func (ptlps PublishersTLPs) AsExprPublisher(publisher string) *query.Expr {
	// Make build process deterministic.
	pubs := make([]Publisher, 0, len(ptlps))
	for pub := range ptlps {
		if pub != "*" {
			pubs = append(pubs, pub)
		}
	}
	slices.Sort(pubs)

	var root *query.Expr
	for _, pub := range pubs {
		ts := or("tlp", mappedTLPs(ptlps[pub]))
		if ts == nil {
			// List is empty.
			continue
		}
		curr := query.FieldEqString(publisher, string(pub)).
			And(ts)

		if root == nil {
			root = curr
		} else {
			root = root.Or(curr)
		}
	}

	// Do we have a wildcard?
	if tlps, ok := ptlps["*"]; ok {
		if ts := or("tlp", mappedTLPs(tlps)); ts != nil {
			// If we have other publishers,
			// don't apply wildcard in these cases.
			if len(pubs) > 0 {
				ts = ts.And(or(publisher, pubs).Not())
			}
			if root == nil {
				root = ts
			} else {
				root = root.Or(ts)
			}
		}
	}

	if root == nil {
		return query.False()
	}
	return root
}
