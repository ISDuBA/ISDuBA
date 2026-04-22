// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package models

type (
	// TextPath is JSON path into a document and
	// the positions which needed to be highlighted when showing
	// search matches.
	// The original text can queried with the JSON path
	// so it is not included here.
	TextPath struct {
		Path      string   `json:"path"`
		Text      *string  `json:"text,omitempty"`
		Positions [][2]int `json:"positions"`
	}
	// TextPaths is a list of matches into a document.
	TextPaths []TextPath
)
