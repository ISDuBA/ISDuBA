<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { tablePadding } from "$lib/Table/defaults";
  import { getReadableDateString } from "../helpers";
  const baseCellStyle = "py-0 px-2";
  const cellStyle = "whitespace-nowrap " + baseCellStyle;
</script>

{#if appStore.state.webview.doc?.isRevisionHistoryPresent}
  <div class="mt-1 w-fit pl-5">
    <Table noborder striped={true}>
      <TableHead>
        <TableHeadCell padding={tablePadding}>#</TableHeadCell>
        <TableHeadCell padding={tablePadding}>Date</TableHeadCell>
        <TableHeadCell padding={tablePadding}>Summary</TableHeadCell>
        <TableHeadCell padding={tablePadding}>Legacy_version</TableHeadCell>
      </TableHead>
      <TableBody>
        {#each appStore.state.webview.doc?.revisionHistory as entry}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyle}>{entry.number}</TableBodyCell>
            <TableBodyCell tdClass={cellStyle}>{getReadableDateString(entry.date)}</TableBodyCell>
            <TableBodyCell tdClass={baseCellStyle + " min-w-52"}>{entry.summary}</TableBodyCell>
            <TableBodyCell
              >{#if entry.legacy_version}{entry.legacy_version}{/if}</TableBodyCell
            >
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
{/if}
