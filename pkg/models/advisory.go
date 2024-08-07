// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

// AdvisoryKey specifies a unique advisory.
type AdvisoryKey struct {
	Publisher  string `uri:"publisher" binding:"required" json:"publisher"`
	TrackingID string `uri:"trackingid" binding:"required" json:"tracking_id"`
}

// AdvisoryState specifies a unique advisory with its state.
type AdvisoryState struct {
	Publisher  string   `uri:"publisher" binding:"required" json:"publisher"`
	TrackingID string   `uri:"trackingid" binding:"required" json:"tracking_id"`
	State      Workflow `uri:"state" binding:"required" json:"state"`
}
