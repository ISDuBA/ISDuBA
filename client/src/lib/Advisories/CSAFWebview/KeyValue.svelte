<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { marked } from "marked";
  import DOMPurify from "dompurify";
  import { Table, TableBody, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  export let keys: Array<String>;
  export let values: any;
  marked.use({ gfm: true });
  const cellStyle = "px-6 py-2";
</script>

<div class="w-max">
  <Table>
    <TableBody>
      {#each keys as key, index}
        {#if key == "text" || key == "Text"}
          <TableBodyRow>
            <TableBodyCell class={cellStyle}>{key}</TableBodyCell>
            <TableBodyCell class={cellStyle}>
              <div class="markdown-text">
                <div class="display-markdown">
                  {@html DOMPurify.sanitize(
                    marked.parse(
                      values[index].replace(/^[\u200B\u200C\u200D\u200E\u200F\uFEFF]/, "")
                    )
                  )}
                </div>
              </div>
            </TableBodyCell>
          </TableBodyRow>
        {:else}
          <TableBodyRow
            ><TableBodyCell class={cellStyle}>{key}</TableBodyCell>
            <TableBodyCell class={cellStyle}>{values[index]}</TableBodyCell>
          </TableBodyRow>
        {/if}
      {/each}
    </TableBody>
  </Table>
</div>

<style>
  .markdown-text {
    padding: 0.5rem;
    border: 1px solid lightgray;
    width: 100%;
    overflow-x: auto;
    position: relative;
  }
  .display-markdown {
  }
</style>
