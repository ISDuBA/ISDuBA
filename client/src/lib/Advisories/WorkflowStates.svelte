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
  import CBadge from "$lib/Components/CBadge.svelte";

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
  <a href={"javascript:void(0);"} class="inline-flex" onclick={() => updateStateIfAllowed(NEW)}>
    <CBadge title="Mark as new" class="flex w-fit gap-1" color={getBadgeColor(NEW, advisoryState)}>
      <i class="bx bxs-certification"></i>
      <span>{NEW}</span>
    </CBadge>
  </a>
  <a href={"javascript:void(0);"} class="inline-flex" onclick={() => updateStateIfAllowed(READ)}>
    <CBadge
      title="Mark as read"
      class="flex w-fit gap-1"
      color={getBadgeColor(READ, advisoryState)}
    >
      <i class="bx bx-show"></i>
      <span>{READ}</span></CBadge
    >
  </a>
  {#if isRoleIncluded(appStore.getRoles(), [EDITOR, REVIEWER]) && advisoryState === REVIEW}
    <a
      href={"javascript:void(0);"}
      class="inline-flex"
      onclick={() => {
        document.getElementById("comment-textarea")?.focus();
      }}
    >
      <CBadge title="Mark as assesing" class="flex w-fit gap-1" color="dark">
        <i class="bx bxs-analyse"></i>
        <span>{ASSESSING}</span></CBadge
      >
    </a>
  {:else if isRoleIncluded(appStore.getRoles(), [EDITOR]) && advisoryState === ARCHIVED}
    <a
      href={"javascript:void(0);"}
      class="inline-flex"
      onclick={() => {
        document.getElementById("comment-textarea")?.focus();
      }}
    >
      <CBadge title="Mark as assesing" class="flex w-fit gap-1" color="dark">
        <i class="bx bxs-analyse"></i>
        <span>{ASSESSING}</span>
      </CBadge>
    </a>
  {:else}
    <a
      href={"javascript:void(0);"}
      class="inline-flex"
      onclick={() => updateStateIfAllowed(ASSESSING)}
    >
      <CBadge
        title="Mark as assesing"
        class="flex w-fit gap-1"
        color={getBadgeColor(ASSESSING, advisoryState)}
      >
        <i class="bx bxs-analyse"></i>
        <span>{ASSESSING}</span>
      </CBadge>
    </a>
  {/if}
  {#if advisoryState === ARCHIVED && isRoleIncluded(appStore.getRoles(), [EDITOR])}
    <a
      href={"javascript:void(0);"}
      class="inline-flex"
      onclick={() => {
        document.getElementById("comment-textarea")?.focus();
      }}
    >
      <CBadge title="Release for review" class="flex w-fit gap-1" color="dark">
        <i class="bx bx-book-open"></i>
        <span>{REVIEW}</span>
      </CBadge>
    </a>
  {:else}
    <a
      href={"javascript:void(0);"}
      class="inline-flex"
      onclick={() => updateStateIfAllowed(REVIEW)}
    >
      <CBadge
        title="Release for review"
        class="flex w-fit gap-1"
        color={getBadgeColor(REVIEW, advisoryState)}
      >
        <i class="bx bx-book-open"></i>
        <span>{REVIEW}</span>
      </CBadge>
    </a>
  {/if}
  <a
    href={"javascript:void(0);"}
    class="inline-flex"
    onclick={() => updateStateIfAllowed(ARCHIVED)}
  >
    <CBadge title="Archive" class="flex w-fit gap-1" color={getBadgeColor(ARCHIVED, advisoryState)}>
      <i class="bx bx-archive"></i>
      <span>{ARCHIVED}</span>
    </CBadge>
  </a>
  <a href={"javascript:void(0);"} class="inline-flex" onclick={() => updateStateIfAllowed(DELETE)}>
    <CBadge
      title="Mark for deletion"
      onclick={() => updateStateFn(DELETE)}
      class="flex w-fit gap-1"
      color={getBadgeColor(DELETE, advisoryState)}
    >
      <i class="bx bx-trash"></i>
      <span>{DELETE}</span>
    </CBadge>
  </a>
{/if}
