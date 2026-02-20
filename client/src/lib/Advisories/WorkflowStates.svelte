<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { ASSESSING, ARCHIVED, DELETE, NEW, READ, REVIEW, EDITOR, REVIEWER } from "$lib/workflow";
  import { allowedToChangeWorkflow, isRoleIncluded } from "$lib/permissions";
  import { appStore } from "$lib/store.svelte";
  import WorkflowButton from "./WorkflowButton.svelte";

  interface Props {
    advisoryState: string;
    updateStateFn: (state: string) => Promise<void>;
  }
  let { advisoryState = "", updateStateFn }: Props = $props();

  const updateStateIfAllowed = async (state: string) => {
    if (allowedToChangeWorkflow(appStore.getRoles(), advisoryState, state)) {
      await updateStateFn(state);
    }
  };

  const getBadgeColor = (state: string, currentState: string) => {
    if (state === currentState) {
      return "green";
    } else if (allowedToChangeWorkflow(appStore.getRoles(), currentState, state)) {
      return "dark";
    } else {
      return "none";
    }
  };
</script>

{#if advisoryState}
  <WorkflowButton
    onClick={() => updateStateIfAllowed(NEW)}
    color={getBadgeColor(NEW, advisoryState)}
    label={NEW}
    tooltip="Mark as read"
  >
    {#snippet icon()}
      <i class="bx bxs-certification"></i>
    {/snippet}
  </WorkflowButton>
  <WorkflowButton
    onClick={() => updateStateIfAllowed(READ)}
    color={getBadgeColor(READ, advisoryState)}
    label={READ}
    tooltip="Mark as read"
  >
    {#snippet icon()}
      <i class="bx bx-show"></i>
    {/snippet}
  </WorkflowButton>
  {#if (isRoleIncluded( appStore.getRoles(), [EDITOR, REVIEWER] ) && advisoryState === REVIEW) || (isRoleIncluded( appStore.getRoles(), [EDITOR] ) && advisoryState === ARCHIVED)}
    <WorkflowButton
      onClick={() => {
        document.getElementById("comment-textarea")?.focus();
      }}
      color="dark"
      label={ASSESSING}
      tooltip="Mark as assesing"
    >
      {#snippet icon()}
        <i class="bx bxs-analyse"></i>
      {/snippet}
    </WorkflowButton>
  {:else}
    <WorkflowButton
      onClick={() => updateStateIfAllowed(ASSESSING)}
      color={getBadgeColor(ASSESSING, advisoryState)}
      label={ASSESSING}
      tooltip="Mark as assesing"
    >
      {#snippet icon()}
        <i class="bx bxs-analyse"></i>
      {/snippet}
    </WorkflowButton>
  {/if}
  {#if advisoryState === ARCHIVED && isRoleIncluded(appStore.getRoles(), [EDITOR])}
    <WorkflowButton
      onClick={() => {
        document.getElementById("comment-textarea")?.focus();
      }}
      color="dark"
      label={REVIEW}
      tooltip="Release for review"
    >
      {#snippet icon()}
        <i class="bx bx-book-open"></i>
      {/snippet}
    </WorkflowButton>
  {:else}
    <WorkflowButton
      onClick={() => updateStateIfAllowed(REVIEW)}
      color={getBadgeColor(REVIEW, advisoryState)}
      label={REVIEW}
      tooltip="Release for review"
    >
      {#snippet icon()}
        <i class="bx bx-book-open"></i>
      {/snippet}
    </WorkflowButton>
  {/if}
  <WorkflowButton
    onClick={() => updateStateIfAllowed(ARCHIVED)}
    color={getBadgeColor(ARCHIVED, advisoryState)}
    label={ARCHIVED}
    tooltip="Archive"
  >
    {#snippet icon()}
      <i class="bx bx-archive"></i>
    {/snippet}
  </WorkflowButton>
  <WorkflowButton
    onClick={() => updateStateFn(DELETE)}
    color={getBadgeColor(DELETE, advisoryState)}
    label={DELETE}
    tooltip="Mark for deletion"
  >
    {#snippet icon()}
      <i class="bx bx-trash"></i>
    {/snippet}
  </WorkflowButton>
{/if}
