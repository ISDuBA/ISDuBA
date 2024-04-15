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
	"strings"
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
	TLPWhite TLP = "WHITE" // TLPWhite represents TLP:WHITE
	TLPGreen TLP = "GREEN" // TLPGreen represents TLP:GREEN
	TLPAmber TLP = "AMBER" // TLPAmber represents TLP:AMBER
	TLPRed   TLP = "RED"   // bearbeiterTLPRed   represents TLP:RED
)

// UnmarshalText implements [encoding.TextUnmarshaler].
func (tlp *TLP) UnmarshalText(text []byte) error {
	s := TLP(text)
	switch s {
	case TLPWhite, TLPGreen, TLPAmber, TLPRed:
		*tlp = s
		return nil
	default:
		return fmt.Errorf("unknown TLP value: %q", s)
	}
}

// Allowed checks if a pair of publisher/tlp is allowed.
func (ptlps PublishersTLPs) Allowed(publisher string, tlp TLP) bool {
	if p, ok := ptlps[Publisher(publisher)]; ok {
		return slices.Contains(p, tlp)
	}
	wildcard, ok := ptlps["*"]
	return ok && slices.Contains(wildcard, tlp)
}

// AsConditions returns the list of TLP rules as a postfix expression.
func (ptlps PublishersTLPs) AsConditions() string {

	// As map iteration order is random we sort to simplify test
	publisherOrder := make([]Publisher, 0, len(ptlps))
	for publisher := range ptlps {
		if publisher != "*" {
			publisherOrder = append(publisherOrder, publisher)
		}
	}
	slices.Sort(publisherOrder)

	var b strings.Builder
	var noneWildcards int
	for _, publisher := range publisherOrder {
		noneWildcards++
		tlps := ptlps[publisher]
		for j, t := range tlps {
			b.WriteString(" $tlp ")
			b.WriteString(string(t))
			b.WriteString(" =")
			if j > 0 {
				b.WriteString(" or")
			}
		}
		b.WriteString(` $publisher "`)
		publisher := strings.ReplaceAll(string(publisher), `"`, `\"`)
		b.WriteString(publisher)
		b.WriteString(`" =`)
		if len(tlps) > 0 {
			b.WriteString(" and")
		}
		if noneWildcards > 1 {
			b.WriteString(" or")
		}
	}

	if wildcard, ok := ptlps["*"]; ok {
		for j, t := range wildcard {
			b.WriteString(" $tlp ")
			b.WriteString(string(t))
			b.WriteString(" =")
			if j > 0 {
				b.WriteString(" or")
			}
		}
		first := true
		for _, publisher := range publisherOrder {
			b.WriteString(` $publisher "`)
			publisher := strings.ReplaceAll(string(publisher), `"`, `\"`)
			b.WriteString(publisher)
			b.WriteString(`" !=`)
			if !first {
				b.WriteString(" and")
			} else {
				first = false
			}
		}
		if noneWildcards > 0 {
			b.WriteString(" and or")
		}
	}
	return b.String()
}
