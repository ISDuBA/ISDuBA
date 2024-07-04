/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

// THIS FILE IS MACHINE GENERATED!
// Use "go generate ./..." in the root folder to regenerate it.

export type WorkflowState = string;
export const NEW: WorkflowState = "new";
export const READ: WorkflowState = "read";
export const ASSESSING: WorkflowState = "assessing";
export const REVIEW: WorkflowState = "review";
export const ARCHIVED: WorkflowState = "archived";
export const DELETE: WorkflowState = "delete";

export const WORKFLOW_STATES = [NEW, READ, ASSESSING, REVIEW, ARCHIVED, DELETE];

export type Role = string;
export const ADMIN: Role = "admin";
export const IMPORTER: Role = "importer";
export const EDITOR: Role = "editor";
export const REVIEWER: Role = "reviewer";
export const AUDITOR: Role = "auditor";
export const SOURCE_MANAGER: Role = "source-manager";

export type WorkflowStateTransition = {
  from: WorkflowState;
  to: WorkflowState;
  roles: Role[];
};

export const WORKFLOW_TRANSITIONS: WorkflowStateTransition[] = [
  { from: ARCHIVED, to: ASSESSING, roles: [ADMIN, EDITOR] },
  { from: ARCHIVED, to: DELETE, roles: [EDITOR, REVIEWER] },
  { from: ARCHIVED, to: NEW, roles: [IMPORTER] },
  { from: ARCHIVED, to: READ, roles: [ADMIN] },
  { from: ARCHIVED, to: REVIEW, roles: [ADMIN, EDITOR] },
  { from: ASSESSING, to: DELETE, roles: [EDITOR, REVIEWER] },
  { from: ASSESSING, to: NEW, roles: [IMPORTER] },
  { from: ASSESSING, to: READ, roles: [EDITOR] },
  { from: ASSESSING, to: REVIEW, roles: [EDITOR] },
  { from: DELETE, to: ARCHIVED, roles: [ADMIN] },
  { from: DELETE, to: ASSESSING, roles: [ADMIN] },
  { from: DELETE, to: READ, roles: [ADMIN] },
  { from: DELETE, to: REVIEW, roles: [ADMIN] },
  { from: NEW, to: READ, roles: [EDITOR] },
  { from: READ, to: ASSESSING, roles: [EDITOR] },
  { from: READ, to: DELETE, roles: [EDITOR, REVIEWER] },
  { from: READ, to: NEW, roles: [EDITOR] },
  { from: REVIEW, to: ARCHIVED, roles: [REVIEWER] },
  { from: REVIEW, to: ASSESSING, roles: [REVIEWER, EDITOR] },
  { from: REVIEW, to: DELETE, roles: [REVIEWER] },
  { from: REVIEW, to: NEW, roles: [IMPORTER] },
  { from: REVIEW, to: READ, roles: [REVIEWER] }
];
