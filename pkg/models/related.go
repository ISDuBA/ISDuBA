// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

package models

// RelatedDocument is a document that relates to another document by an CVE.
type RelatedDocument struct {
	DocumentID int64   `json:"document_id"`
	CVE        string  `json:"cve"`
	State      string  `json:"state"`
	SSVC       *string `json:"ssvc,omitempty"`
	Title      *string `json:"title,omitempty"`
	TrackingID string  `json:"tracking_id"`
	Publisher  string  `json:"publisher"`
}
