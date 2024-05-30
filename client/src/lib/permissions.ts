/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

import { appStore } from "./store";

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

const WORKFLOW_TRANSITIONS: WorkflowStateTransition[] = [
  {
    from: NEW,
    to: READ,
    roles: [EDITOR]
  },
  { from: READ, to: NEW, roles: [EDITOR] },
  { from: READ, to: ASSESSING, roles: [EDITOR] },
  { from: READ, to: REVIEW, roles: [] },
  { from: READ, to: ARCHIVED, roles: [] },
  { from: READ, to: DELETE, roles: [EDITOR, REVIEWER] },
  { from: ASSESSING, to: NEW, roles: [EDITOR] },
  { from: ASSESSING, to: READ, roles: [EDITOR] },
  { from: ASSESSING, to: REVIEW, roles: [EDITOR] },
  { from: ASSESSING, to: ARCHIVED, roles: [] },
  { from: ASSESSING, to: DELETE, roles: [EDITOR, REVIEWER] },
  { from: REVIEW, to: NEW, roles: [REVIEWER] },
  { from: REVIEW, to: READ, roles: [REVIEWER] },
  { from: REVIEW, to: ASSESSING, roles: [REVIEWER] },
  { from: REVIEW, to: ARCHIVED, roles: [REVIEWER] },
  { from: REVIEW, to: DELETE, roles: [REVIEWER] },
  { from: ARCHIVED, to: NEW, roles: [] },
  { from: ARCHIVED, to: READ, roles: [ADMIN] },
  { from: ARCHIVED, to: ASSESSING, roles: [ADMIN] },
  { from: ARCHIVED, to: REVIEW, roles: [ADMIN] },
  { from: ARCHIVED, to: DELETE, roles: [EDITOR, REVIEWER] },
  { from: DELETE, to: NEW, roles: [] },
  { from: DELETE, to: READ, roles: [ADMIN] },
  { from: DELETE, to: ASSESSING, roles: [ADMIN] },
  { from: DELETE, to: REVIEW, roles: [ADMIN] },
  { from: DELETE, to: ARCHIVED, roles: [ADMIN] }
];

export function isRoleIncluded(roles: Role[], rolesToCheck: Role[]) {
  for (let i = 0; i < rolesToCheck.length; i++) {
    if (roles.includes(rolesToCheck[i])) {
      return true;
    }
  }
  return false;
}

export function canSetStateNew(currentState: WorkflowState) {
  return allowedToChangeWorkflow(appStore.getRoles(), currentState, NEW);
}

export function canSetStateRead(currentState: WorkflowState) {
  return allowedToChangeWorkflow(appStore.getRoles(), currentState, READ);
}

export function canSetStateReview(currentState: WorkflowState) {
  return allowedToChangeWorkflow(appStore.getRoles(), currentState, REVIEW);
}

export function canSetStateAssessing(currentState: WorkflowState) {
  return allowedToChangeWorkflow(appStore.getRoles(), currentState, ASSESSING);
}

export function canSetStateArchived(currentState: WorkflowState) {
  return allowedToChangeWorkflow(appStore.getRoles(), currentState, ARCHIVED);
}

export function canSetStateDeleted(currentState: WorkflowState) {
  return allowedToChangeWorkflow(appStore.getRoles(), currentState, DELETE);
}

export function allowedToChangeWorkflow(
  roles: Role[],
  oldState: WorkflowState,
  newState: WorkflowState
) {
  for (let i = 0; i < WORKFLOW_TRANSITIONS.length; i++) {
    const change = WORKFLOW_TRANSITIONS[i];
    if (change.from === oldState && change.to === newState && isRoleIncluded(change.roles, roles)) {
      return true;
    }
  }
  return false;
}

export function getAllowedWorkflowChanges(currentState: WorkflowState) {
  return WORKFLOW_TRANSITIONS.filter(
    (transition) =>
      isRoleIncluded(transition.roles, appStore.getRoles()) && transition.from === currentState
  );
}
