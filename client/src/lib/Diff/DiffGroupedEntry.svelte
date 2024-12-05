<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    Button,
    Table,
    TableBody,
    TableBodyCell,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { getKeysOfAllObjects } from "$lib/utils";
  import type { JsonDiffResult } from "./Diff";
  import { onMount } from "svelte";

  export let change: any;
  export let isClosable = false;

  let headers: string[] = [];
  let isOpen = !isClosable;

  onMount(() => {
    //if (isClosable) isOpen = false;
    headers = getKeysOfAllObjects(change.result.map((r: JsonDiffResult) => r.value));
  });
</script>

<div>
  {#if isClosable}
    <Button
      class="text-md flex items-end gap-x-2 bg-inherit pl-1 font-bold text-gray-500 hover:bg-inherit"
      on:click={() => {
        isOpen = !isOpen;
      }}
    >
      {#if isOpen}
        <i class="bx bx-chevron-up text-2xl"></i>
      {:else}
        <i class="bx bx-chevron-down text-2xl"></i>
      {/if}
      <code>
        {change.result[0].path.split("/").slice(0, -1).join("/")}
      </code>
    </Button>
  {:else}
    <div class="mb-1 text-sm font-bold">
      <code>
        {change.result[0].path.split("/").slice(0, -1).join("/")}
      </code>
    </div>
  {/if}
  {#if isOpen}
    <Table>
      <TableHead class="bg-transparent">
        <TableHeadCell padding="px-2 py-2" class="text-gray-500 dark:text-white">Pos</TableHeadCell>
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
              {Number(member.path.split("/").slice(-1)) + 1}
            </TableBodyCell>
            {#each headers as header}
              <TableBodyCell
                tdClass="px-2 py-2 align-top font-medium !text-gray-500 dark:text-white"
              >
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
  {/if}
</div>
