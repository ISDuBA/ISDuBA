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
	"sort"
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
	wildcard, ok := ptlps[Publisher("*")]
	if ok {
		if slices.Contains(wildcard, tlp) {
			return true
		}
	}

	p, ok := ptlps[Publisher(publisher)]
	if ok {
		return slices.Contains(p, tlp)
	}
	return false
}

// AsConditions returns the list of TLP rules as a postfix expression.
func (ptlps PublishersTLPs) AsConditions() string {
	var b strings.Builder
	var noneWildcards int
	publisherOrder := make([]string, 0, len(ptlps))

	// As map iteration order is random we sort to simplify test
	for publisher := range ptlps {
		publisherOrder = append(publisherOrder, string(publisher))
	}
	sort.Strings(publisherOrder)

	for _, publisher := range publisherOrder {
		if publisher == "*" {
			continue
		}
		noneWildcards++
		tlps := ptlps[Publisher(publisher)]
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

	wildcard, ok := ptlps[Publisher("*")]
	if ok {
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
			if publisher == "*" {
				continue
			}
			b.WriteString(` $publisher "`)
			publisher := strings.ReplaceAll(string(publisher), `"`, `\"`)
			b.WriteString(publisher)
			b.WriteString(`" !=`)
			if !first {
				b.WriteString(" and")
			}

			first = false
		}
		if noneWildcards > 0 {
			b.WriteString(" and or")
		}
	}
	return b.String()
}
