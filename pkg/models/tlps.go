// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
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
	// PuplisherTLPs is a list of allowed TLP values per publisher.
	PuplisherTLPs struct {
		Publisher string `json:"publisher" toml:"publisher"`
		TLPs      []TLP  `json:"tlps" toml:"tlps"`
	}
	// PuplishersTLPs is a list of TLPs per publisher.
	PuplishersTLPs []PuplisherTLPs
)

const (
	TLPWhite TLP = "WHITE" // TLPWhite represents TLP:WHITE
	TLPGreen TLP = "GREEN" // TLPGreen represents TLP:GREEN
	TLPAmber TLP = "AMBER" // TLPAmber represents TLP:AMBER
	TLPRed   TLP = "RED"   // TLPRed   represents TLP:RED
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
func (ptlps PuplishersTLPs) Allowed(publisher string, tlp TLP) bool {
	var wildcard *PuplisherTLPs
	for i := range ptlps {
		ptlp := &ptlps[i]
		if ptlp.Publisher == "" {
			wildcard = ptlp
			continue
		}
		if ptlp.Publisher == publisher {
			return slices.Contains(ptlp.TLPs, tlp)
		}
	}
	if wildcard != nil {
		return slices.Contains(wildcard.TLPs, tlp)
	}
	return false
}

// AsConditions returns the list of TLP rules as a postfix expression.
func (ptlps PuplishersTLPs) AsConditions() string {
	var b strings.Builder
	var noneWildcards int
	var wildcard *PuplisherTLPs
	for i := range ptlps {
		p := &ptlps[i]
		if p.Publisher == "" {
			wildcard = p
			continue
		}
		noneWildcards++
		for j, t := range p.TLPs {
			b.WriteString(" $tlp ")
			b.WriteString(string(t))
			b.WriteString(" =")
			if j > 0 {
				b.WriteString(" or")
			}
		}
		b.WriteString(` $publisher "`)
		publisher := strings.ReplaceAll(p.Publisher, `"`, `\"`)
		b.WriteString(publisher)
		b.WriteString(`" =`)
		if len(ptlps[i].TLPs) > 0 {
			b.WriteString(" and")
		}
		if noneWildcards > 1 {
			b.WriteString(" or")
		}
	}
	if wildcard != nil {
		for j, t := range wildcard.TLPs {
			b.WriteString(" $tlp ")
			b.WriteString(string(t))
			b.WriteString(" =")
			if j > 0 {
				b.WriteString(" or")
			}
		}
		for j, k := 0, 0; j < noneWildcards; j, k = j+1, k+1 {
			for ptlps[k].Publisher == "" {
				k++
			}
			b.WriteString(` $publisher "`)
			publisher := strings.ReplaceAll(ptlps[k].Publisher, `"`, `\"`)
			b.WriteString(publisher)
			b.WriteString(`" !=`)
			if j > 0 {
				b.WriteString(" and")
			}
		}
		if noneWildcards > 0 {
			b.WriteString(" and or")
		}
	}
	return b.String()
}
