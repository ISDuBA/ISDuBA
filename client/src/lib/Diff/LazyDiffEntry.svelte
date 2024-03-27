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
  import { A } from "flowbite-svelte";

  export let operation: string;
  export let path: string;
  export let urlPath: string;
  let result: any;

  const loadEntry = async (event: any) => {
    event.preventDefault();
    const requestPath = encodeURI(`${urlPath}&item_op=${operation}&item_path=${path}`);
    const response = await request(requestPath, "GET");
    result = await response.json();
  };
</script>

<div>
  <div class="mb-1">
    <b class="me-4">
      <code>
        {path}
      </code>
    </b>
    {#if !result}
      <A on:click={loadEntry}>Load Entry</A>
    {/if}
  </div>
  {#if result}
    <DiffEntry content={result} {operation}></DiffEntry>
  {/if}
</div>
