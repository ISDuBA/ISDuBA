// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

package models

// Event is an event in the event log.
type Event string

// Documents to export
const (
	ImportDocumentEvent Event = "import_document" // ImportDocumentEvent represents a document import.
	DeleteDocumentEvent Event = "delete_document" // DeleteDocumentEvent represents a document deletion.
	StateChangeEvent    Event = "state_change"    // StateChangeEvent represents changing the advisory state.
	AddSSVCEvent        Event = "add_sscv"        // AddSSVCEvent represents the addtion of a SSVC score.
	ChangeSSVCEvent     Event = "change_sscv"     // ChangeSSVCEvent represents the change of a SSVC score.
	DeleteSSVCEvent     Event = "delete_sscv"     // DeleteSSVCEvent represents the deletion of a SSVC score.
	AddCommentEvent     Event = "add_comment"     // AddCommentEvent represents the addition of a comment.
	ChangeCommentEvent  Event = "change_comment"  // ChangeCommentEvent represents the change of a comment.
	DeleteCommentEvent  Event = "delete_comment"  // DeleteCommentEvent represents the deletion of a comment.
)
