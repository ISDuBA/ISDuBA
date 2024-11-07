<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import DiffSelection from "$lib/Diff/DiffSelection.svelte";
  import { appStore } from "$lib/store";
  import { Button, Img } from "flowbite-svelte";

  $: docA = $appStore.app.diff.docA;
  $: docB = $appStore.app.diff.docB;
</script>

<div class={`sticky bottom-0 left-0 flex flex-col items-start justify-center`}>
  <div class="flex">
    <Button
      on:click={appStore.toggleToolbox}
      class="rounded-none rounded-t-md border-b-0"
      color="light"
    >
      <span class="me-2"
        >Diff {docA
          ? `${docA?.document?.title.substring(0, 25)}${docA?.document?.title.length > 25 ? "..." : ""}`
          : ""}
        {docB
          ? ` - ${docB?.document?.title.substring(0, 25)}${docB?.document?.title.length > 25 ? "..." : ""}`
          : ""}</span
      >
      <Img src="plus-minus.svg" class="h-4 min-w-4 dark:invert" />
    </Button>
  </div>
  {#if $appStore.app.isToolboxOpen}
    <div
      class="flex min-h-48 w-full min-w-full max-w-[700pt] items-stretch gap-6 rounded-tr-md border border-solid border-gray-300 bg-white p-4 shadow-gray-800 dark:border-gray-600 dark:bg-gray-800 dark:shadow-gray-200 md:min-w-96 lg:w-auto"
    >
      <DiffSelection></DiffSelection>
    </div>
  {/if}
</div>
