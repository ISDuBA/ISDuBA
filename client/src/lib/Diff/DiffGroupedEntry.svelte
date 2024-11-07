<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Table, TableBody, TableBodyCell, TableHead, TableHeadCell } from "flowbite-svelte";
  import { getKeysOfAllObjects } from "$lib/utils";
  import type { JsonDiffResult } from "./Diff";
  import { onMount } from "svelte";

  export let change: any;

  let headers: string[] = [];

  onMount(() => {
    headers = getKeysOfAllObjects(change.result.map((r: JsonDiffResult) => r.value));
  });
</script>

<div class="mb-1 text-sm font-bold">
  <code>
    {change.result[0].path.split("/").slice(0, -1).join("/")}
  </code>
</div>
<Table>
  <TableHead class="bg-transparent">
    <TableHeadCell padding="px-2 py-2" class="text-gray-500 dark:text-white">Index</TableHeadCell>
    {#each headers as header}
      <TableHeadCell padding="px-2 py-2" class="text-gray-500 dark:text-white"
        >{header}</TableHeadCell
      >
    {/each}
  </TableHead>
  <TableBody>
    {#each change.result as member}
      <tr>
        <TableBodyCell tdClass="px-2 py-2 align-top font-medium !text-gray-500 dark:text-white">
          {member.path.split("/").slice(-1)}
        </TableBodyCell>
        {#each headers as header}
          <TableBodyCell tdClass="px-2 py-2 align-top font-medium !text-gray-500 dark:text-white">
            {#if Array.isArray(member.value[header])}
              {member.value[header].join(", ")}
            {:else}
              {member.value[header] ?? ""}
            {/if}
          </TableBodyCell>
        {/each}
      </tr>
    {/each}
  </TableBody>
</Table>
