/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

import { appStore } from "./store.svelte";
import type { Role, WorkflowState, WorkflowStateTransition } from "./workflow";
import { NEW, READ, REVIEW, ASSESSING, ARCHIVED, DELETE, WORKFLOW_TRANSITIONS } from "./workflow";

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

export function getAllowedWorkflowChanges(
  currentStates: WorkflowState[]
): WorkflowStateTransition[] {
  const workflowTransitions: WorkflowStateTransition[] = [];
  if (currentStates.length === 0) return workflowTransitions;
  WORKFLOW_TRANSITIONS.forEach((transition: WorkflowStateTransition) => {
    if (
      !workflowTransitions.includes(transition) &&
      isRoleIncluded(transition.roles, appStore.getRoles()) &&
      currentStates.length ===
        currentStates.filter((s: WorkflowState) => s === transition.from).length
    ) {
      workflowTransitions.push(transition);
    }
  });
  return workflowTransitions;
}
