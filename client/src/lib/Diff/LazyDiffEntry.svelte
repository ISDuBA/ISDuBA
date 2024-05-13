<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { request } from "$lib/utils";
  import DiffEntry from "./DiffEntry.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { Button } from "flowbite-svelte";

  export let operation: string;
  export let path: string;
  export let urlPath: string;
  let result: any;
  let error: string;
  let isOpen = false;

  const loadEntry = async () => {
    isOpen = !isOpen;
    if (result) return;
    error = "";
    const requestPath = encodeURI(`${urlPath}&item_op=${operation}&item_path=${path}`);
    const response = await request(requestPath, "GET");
    if (response.ok) {
      result = response.content;
    } else if (response.error) {
      error = response.error;
    }
  };
</script>

<div>
  <Button
    class="flex items-end gap-x-2 bg-inherit pl-1 text-gray-500 hover:bg-inherit"
    on:click={loadEntry}
  >
    {#if isOpen}
      <i class="bx bx-chevron-up text-2xl"></i>
    {:else}
      <i class="bx bx-chevron-down text-2xl"></i>
    {/if}
    <code class="text-md font-bold">
      {path}
    </code>
  </Button>
  {#if result && isOpen}
    <div class="mt-2">
      <DiffEntry content={result} {operation}></DiffEntry>
    </div>
  {/if}
  <ErrorMessage message={error}></ErrorMessage>
</div>
