<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { ASSESSING, ARCHIVED, DELETE, NEW, READ, REVIEW, EDITOR, REVIEWER } from "$lib/workflow";
  import { isRoleIncluded } from "$lib/permissions";
  import { appStore } from "$lib/store.svelte";
  import WorkflowButton from "./WorkflowButton.svelte";
  import { setContext } from "svelte";

  interface Props {
    advisoryState: string;
    updateStateFn: (state: string) => Promise<void>;
  }
  let { advisoryState = "", updateStateFn }: Props = $props();

  setContext("currentState", () => advisoryState);
  setContext("updateStateFn", () => updateStateFn);
</script>

<WorkflowButton state={NEW} tooltip="Mark as read">
  {#snippet icon()}
    <i class="bx bxs-certification"></i>
  {/snippet}
</WorkflowButton>
<WorkflowButton state={READ} tooltip="Mark as read">
  {#snippet icon()}
    <i class="bx bx-show"></i>
  {/snippet}
</WorkflowButton>
{#if (isRoleIncluded( appStore.getRoles(), [EDITOR, REVIEWER] ) && advisoryState === REVIEW) || (isRoleIncluded( appStore.getRoles(), [EDITOR] ) && advisoryState === ARCHIVED)}
  <WorkflowButton
    onClick={() => {
      document.getElementById("comment-textarea")?.focus();
    }}
    state={ASSESSING}
    tooltip="Mark as assesing"
  >
    {#snippet icon()}
      <i class="bx bxs-analyse"></i>
    {/snippet}
  </WorkflowButton>
{:else}
  <WorkflowButton state={ASSESSING} tooltip="Mark as assesing">
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
    state={REVIEW}
    tooltip="Release for review"
  >
    {#snippet icon()}
      <i class="bx bx-book-open"></i>
    {/snippet}
  </WorkflowButton>
{:else}
  <WorkflowButton state={REVIEW} tooltip="Release for review">
    {#snippet icon()}
      <i class="bx bx-book-open"></i>
    {/snippet}
  </WorkflowButton>
{/if}
<WorkflowButton state={ARCHIVED} tooltip="Archive">
  {#snippet icon()}
    <i class="bx bx-archive"></i>
  {/snippet}
</WorkflowButton>
<WorkflowButton state={DELETE} tooltip="Mark for deletion">
  {#snippet icon()}
    <i class="bx bx-trash"></i>
  {/snippet}
</WorkflowButton>
