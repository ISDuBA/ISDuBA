<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "#lib/store.svelte.ts";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { tablePadding } from "#lib/Table/defaults.ts";
  import { getReadableDateString } from "../helpers";
  import SearchableText from "../SearchableText.svelte";
  const baseCellStyle = "py-4 px-2";
  const cellStyle = "whitespace-nowrap " + baseCellStyle;

  let revisionHistory = $derived(appStore.state.webview.doc?.revisionHistory);
</script>

{#if appStore.state.webview.doc?.isRevisionHistoryPresent}
  <div class="mt-1 w-fit pl-5">
    <Table border={false} striped={true}>
      <TableHead>
        <TableHeadCell class={tablePadding}>#</TableHeadCell>
        <TableHeadCell class={tablePadding}>Date</TableHeadCell>
        <TableHeadCell class={tablePadding}>Summary</TableHeadCell>
        <TableHeadCell class={tablePadding}>Legacy_version</TableHeadCell>
      </TableHead>
      <TableBody>
        {#if revisionHistory}
          {#each revisionHistory as entry, index (`ref-history${index}`)}
            {@const reversedIndex = revisionHistory.length - 1 - index}
            <TableBodyRow>
              <TableBodyCell class={cellStyle}>{entry.number}</TableBodyCell>
              <TableBodyCell class={cellStyle}>
                <SearchableText
                  text={getReadableDateString(entry.date)}
                  textPath={`/document/tracking/revision_history[${reversedIndex}]/date`}
                />
              </TableBodyCell>
              <TableBodyCell class={baseCellStyle + " min-w-52"}>
                <SearchableText
                  text={entry.summary}
                  textPath={`/document/tracking/revision_history[${reversedIndex}]/summary`}
                />
              </TableBodyCell>
              <TableBodyCell class={baseCellStyle}
                >{#if entry.legacy_version}
                  <SearchableText
                    text={entry.legacy_version}
                    textPath={`/document/tracking/revision_history[${reversedIndex}]/legacy_version`}
                  />
                {/if}</TableBodyCell
              >
            </TableBodyRow>
          {/each}
        {/if}
      </TableBody>
    </Table>
  </div>
{/if}
