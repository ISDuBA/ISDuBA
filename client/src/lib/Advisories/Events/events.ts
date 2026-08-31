// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import type { WorkflowState } from "$lib/workflow";

type CommentEventType = "add_comment" | "change_comment" | "delete_comment";
type SSVCEventType = "add_sscv" | "change_sscv" | "delete_sscv";
type OtherEventType = "import_document" | "delete_document" | "state_change";
type GeneralEventType = CommentEventType | SSVCEventType | OtherEventType;

type BaseEvent = {
  document_id: number;
  comment_id?: number;
  time: string;
  documentVersion?: string; // Added by client
};

type GeneralEvent = BaseEvent & {
  event_type: GeneralEventType;
  actor: string;
  state: WorkflowState;
};

type OtherEvent = BaseEvent & {
  event_type: OtherEventType;
  actor: string;
  state: WorkflowState;
};

type CommentEvent = BaseEvent & {
  event_type: CommentEventType;
  commentator: string;
  id: number;
  message: string;
  times: any[]; // Added by client
};

type SSVCEvent = {
  event_type: SSVCEventType;
  prev_ssvc?: string;
  ssvc?: string;
  actor?: string;
  changedate: string;
  change_number: number;
  documents_id: number;
  documents_version: number;
  time: string;
  documentVersion?: string; // Added by client
};

export type {
  GeneralEventType,
  CommentEventType,
  BaseEvent,
  GeneralEvent,
  CommentEvent,
  SSVCEvent,
  OtherEvent
};
