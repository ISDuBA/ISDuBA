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
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";

  export let operation: string;
  export let path: string;
  export let urlPath: string;
  let result: any;
  let error: string;

  const loadEntry = async (event: any) => {
    event.preventDefault();
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
  <div class="mb-1">
    <b class="me-4">
      <code>
        {path}
      </code>
    </b>
    {#if !result}
      <A on:click={loadEntry}>Load Entry</A>
    {/if}
    <ErrorMessage message={error} plain={true}></ErrorMessage>
  </div>
  {#if result}
    <DiffEntry content={result} {operation}></DiffEntry>
  {/if}
</div>
