<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { tablePadding } from "$lib/table/defaults";
  import {
    Button,
    Table,
    TableHead,
    TableBody,
    TableHeadCell,
    TableBodyRow,
    TableBodyCell
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";

  let queries: any[] = [];
  let orderBy = "";
  let errorMessage = "";

  onMount(async () => {
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queries = response.content;
    } else if (response.error) {
      errorMessage = response.error;
    }
  });
</script>

<h2 class="mb-3 text-lg">User defined queries</h2>
<Button class="mb-6 mt-3" href="/#/configuration/userqueries"
  ><i class="bx bx-plus"></i>New query</Button
>
<div class="flex flex-row">
  <div class="mb-12 w-1/3">
    <Table hoverable={true} noborder={true}>
      <TableHead class="cursor-pointer">
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Name<i
            class:bx={true}
            class:bx-caret-up={orderBy == "name"}
            class:bx-caret-down={orderBy == "-name"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Description<i
            class:bx={true}
            class:bx-caret-up={orderBy == "description"}
            class:bx-caret-down={orderBy == "-description"}
          ></i>
        </TableHeadCell>
      </TableHead>
      <TableBody>
        {#each queries as query}
          <TableBodyRow class="cursor-pointer">
            <TableBodyCell>{query.name ?? "-"}</TableBodyCell>
            <TableBodyCell>{query.description ?? "-"}</TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
    <ErrorMessage message={errorMessage}></ErrorMessage>
  </div>
</div>
