<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Table, TableBody, TableBodyCell, TableBodyRow } from "flowbite-svelte";

  interface Props {
    keys: Array<string>;
    values: any;
  }
  let { keys, values }: Props = $props();

  const cellStyle = "px-6 py-0";
</script>

<div class="ml-2 w-fit">
  <Table border={false}>
    <TableBody>
      {#each keys as key, index}
        {#if key === "text" || key === "Text"}
          <TableBodyRow color="default">
            <TableBodyCell class={cellStyle}>{key}</TableBodyCell>
            <TableBodyCell class={cellStyle}>
              <div class="markdown-text">
                <div class="display-markdown max-w-2/3">
                  {index}
                </div>
              </div>
            </TableBodyCell>
          </TableBodyRow>
        {:else}
          <TableBodyRow color="default"
            ><TableBodyCell class={cellStyle}>{key}</TableBodyCell>
            <TableBodyCell class={cellStyle}>
              {#if typeof values[index] === "string" && values[index].startsWith && values[index].startsWith("https://")}
                <a class="underline" href={values[index]}>
                  <i class="bx bx-link"></i>{values[index]}
                </a>
              {:else}
                {values[index]}
              {/if}
            </TableBodyCell>
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
  }
</style>
