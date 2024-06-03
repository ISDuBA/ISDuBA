<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { ASSESSING, ARCHIVED, DELETE, NEW, READ, REVIEW } from "$lib/workflow";
  import { allowedToChangeWorkflow } from "$lib/permissions";
  import { appStore } from "$lib/store";
  import { Badge } from "flowbite-svelte";

  export let advisoryState = "";
  export let updateStateFn;

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
  <a href={"javascript:void(0);"} class="inline-flex" on:click={() => updateStateIfAllowed(NEW)}>
    <Badge title="Mark as new" class="w-fit" color={getBadgeColor(NEW, advisoryState)}>{NEW}</Badge>
  </a>
  <a href={"javascript:void(0);"} class="inline-flex" on:click={() => updateStateIfAllowed(READ)}>
    <Badge title="Mark as read" class="w-fit" color={getBadgeColor(READ, advisoryState)}
      >{READ}</Badge
    >
  </a>
  <a
    href={"javascript:void(0);"}
    class="inline-flex"
    on:click={() => updateStateIfAllowed(ASSESSING)}
  >
    <Badge title="Mark as assesing" class="w-fit" color={getBadgeColor(ASSESSING, advisoryState)}
      >{ASSESSING}</Badge
    >
  </a>
  <a href={"javascript:void(0);"} class="inline-flex" on:click={() => updateStateIfAllowed(REVIEW)}>
    <Badge title="Release for review" class="w-fit" color={getBadgeColor(REVIEW, advisoryState)}
      >{REVIEW}</Badge
    >
  </a>
  <a
    href={"javascript:void(0);"}
    class="inline-flex"
    on:click={() => updateStateIfAllowed(ARCHIVED)}
  >
    <Badge title="Archive" class="w-fit" color={getBadgeColor(ARCHIVED, advisoryState)}
      >{ARCHIVED}</Badge
    >
  </a>
  <a href={"javascript:void(0);"} class="inline-flex" on:click={() => updateStateIfAllowed(DELETE)}>
    <Badge
      title="Mark for deletion"
      on:click={() => updateStateFn(DELETE)}
      class="w-fit"
      color={getBadgeColor(DELETE, advisoryState)}>{DELETE}</Badge
    >
  </a>
{/if}
