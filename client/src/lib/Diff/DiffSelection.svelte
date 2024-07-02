<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { Button, Img } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
</script>

{#if $appStore.app.diff.docA || $appStore.app.diff.docB}
  <div class="left-50 fixed bottom-0 flex justify-center">
    <div class="flex gap-6 bg-gray-100 px-4 py-2 pb-5 shadow-md shadow-gray-800">
      <div class="flex items-center gap-1">
        <div>
          <span>First:</span>
          <span>{$appStore.app.diff.docA?.title ?? "-"}</span>
        </div>
        <Button on:click={() => appStore.setDiffDocA(null)} color="light" class="p-0">
          <i class="bx bx-x text-lg"></i>
        </Button>
      </div>
      <div class="flex items-center gap-1">
        <div>
          <span>Second:</span>
          <span>{$appStore.app.diff.docB?.title ?? "-"}</span>
        </div>
        <Button on:click={() => appStore.setDiffDocB(null)} color="light" class="p-0">
          <i class="bx bx-x text-lg"></i>
        </Button>
      </div>
      <Button
        on:click={() => push("/diff")}
        disabled={$appStore.app.diff.docA === null || $appStore.app.diff.docB === null}
        size="sm"
        class="flex gap-x-2"
      >
        <Img src="plus-minus.svg" class="w-5 invert" />
        <span>Compare</span>
      </Button>
    </div>
  </div>
{/if}
