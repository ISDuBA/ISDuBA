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

<WorkflowButton state={NEW} tooltip="Mark as read" />
<WorkflowButton state={READ} tooltip="Mark as read" />
{#if (isRoleIncluded( appStore.getRoles(), [EDITOR, REVIEWER] ) && advisoryState === REVIEW) || (isRoleIncluded( appStore.getRoles(), [EDITOR] ) && advisoryState === ARCHIVED)}
  <WorkflowButton
    onClick={() => {
      document.getElementById("comment-textarea")?.focus();
    }}
    state={ASSESSING}
    tooltip="Mark as assesing"
  />
{:else}
  <WorkflowButton state={ASSESSING} tooltip="Mark as assesing" />
{/if}
{#if advisoryState === ARCHIVED && isRoleIncluded(appStore.getRoles(), [EDITOR])}
  <WorkflowButton
    onClick={() => {
      document.getElementById("comment-textarea")?.focus();
    }}
    state={REVIEW}
    tooltip="Release for review"
  />
{:else}
  <WorkflowButton state={REVIEW} tooltip="Release for review" />
{/if}
<WorkflowButton state={ARCHIVED} tooltip="Archive" />
<WorkflowButton state={DELETE} tooltip="Mark for deletion" />
