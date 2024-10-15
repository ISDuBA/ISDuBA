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
  import { tablePadding } from "$lib/Table/defaults";
  const baseCellStyle = "py-0 px-2";
  const cellStyle = "whitespace-nowrap " + baseCellStyle;
</script>

{#if $appStore.webview.doc?.isRevisionHistoryPresent}
  <div class="mt-1 w-fit pl-5">
    <Table noborder striped={true}>
      <TableHead>
        <TableHeadCell padding={tablePadding}>Number</TableHeadCell>
        <TableHeadCell padding={tablePadding}>Date</TableHeadCell>
        <TableHeadCell padding={tablePadding}>Summary</TableHeadCell>
        <TableHeadCell padding={tablePadding}>Legacy_version</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each $appStore.webview.doc?.revisionHistory as entry}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyle}>{entry.number}</TableBodyCell>
            <TableBodyCell tdClass={cellStyle}>{entry.date}</TableBodyCell>
            <TableBodyCell tdClass={baseCellStyle + " min-w-52"}>{entry.summary}</TableBodyCell>
            <TableBodyCell
              >{#if entry.legacyVersion}{entry.legacyVersion}{/if}</TableBodyCell
            >
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
{/if}
