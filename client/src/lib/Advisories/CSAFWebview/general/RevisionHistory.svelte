<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
</script>

{#if $appStore.webview.doc?.isRevisionHistoryPresent}
  <div class="w-max">
    <Table noborder>
      <TableHead>
        <TableHeadCell>Date</TableHeadCell>
        <TableHeadCell>Number</TableHeadCell>
        <TableHeadCell>Summary</TableHeadCell>
        <TableHeadCell>Legacy_version</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each $appStore.webview.doc?.revisionHistory as entry}
          <TableBodyRow>
            <TableBodyCell>{entry.date}</TableBodyCell>
            <TableBodyCell>{entry.number}</TableBodyCell>
            <TableBodyCell>{entry.summary}</TableBodyCell>
            <TableBodyCell
              >{#if entry.legacyVersion}{entry.legacyVersion}{/if}</TableBodyCell
            >
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
{/if}
