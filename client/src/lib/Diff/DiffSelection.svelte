<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { Button, Dropzone, Img, P } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { getPublisher } from "$lib/utils";

  $: docA = $appStore.app.diff.docA;
  $: docB = $appStore.app.diff.docB;
</script>

{#if docA || docB}
  <div class="left-50 fixed bottom-0 flex justify-center border border-solid border-gray-300">
    <div class="flex items-stretch gap-6 bg-gray-100 px-4 py-2 pb-8 shadow-md shadow-gray-800">
      <div class="flex items-center gap-1">
        <div class="flex min-h-28 justify-between gap-2 rounded-md px-3 py-2">
          {#if docA}
            <div>
              <Button on:click={() => appStore.setDiffDocA(null)} color="light" class="p-1">
                <i class="bx bx-x text-lg"></i>
              </Button>
            </div>
            <div class="flex flex-col">
              <div title={docA.title}>{docA.title}</div>
              <div class="text-gray-600">{getPublisher(docA.publisher)}</div>
              <div class="text-gray-600">Version: {docA.version}</div>
            </div>
          {:else}
            <div class="flex flex-col gap-2">
              <P italic>Select another document or upload a local one.</P>
              <!-- TODO: When upload is implemented remove the background color and the "cursor-not-allowed" -->
              <Dropzone class="h-16 cursor-not-allowed border-dashed bg-gray-200 hover:bg-gray-200">
                <i class="bx bx-upload text-xl text-gray-500"></i>
                <p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
                  <span class="font-semibold">Click to upload</span> or drag and drop
                </p>
              </Dropzone>
            </div>
          {/if}
        </div>
      </div>
      <div class="flex items-center gap-1">
        <div class="flex min-h-28 justify-between gap-2 rounded-md px-3 py-2">
          {#if docB}
            <div>
              <Button on:click={() => appStore.setDiffDocB(null)} color="light" class="p-1">
                <i class="bx bx-x text-lg"></i>
              </Button>
            </div>
            <div class="flex flex-col">
              <div title={docB.title}>{docB.title}</div>
              <div class="text-gray-600">{getPublisher(docB.publisher)}</div>
              <div class="text-gray-600">Version: {docB.version}</div>
            </div>
          {:else}
            <div class="flex flex-col gap-2">
              <P italic>Select another document or upload a local one.</P>
              <!-- TODO: When upload is implemented remove the background color and the "cursor-not-allowed" -->
              <Dropzone class="h-16 cursor-not-allowed border-dashed bg-gray-200 hover:bg-gray-200">
                <i class="bx bx-upload text-xl text-gray-500"></i>
                <p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
                  <span class="font-semibold">Click to upload</span> or drag and drop
                </p>
              </Dropzone>
            </div>
          {/if}
        </div>
      </div>
      <div class="flex h-full items-center">
        <Button
          on:click={() => push("/diff")}
          disabled={docA === null || docB === null}
          size="sm"
          class="flex gap-x-2"
        >
          <Img src="plus-minus.svg" class="w-5 invert" />
          <span>Compare</span>
        </Button>
      </div>
    </div>
  </div>
{/if}
