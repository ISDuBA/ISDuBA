<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Table, TableBody, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  export let references: any;
  const baseCellStyle = "py-2 px-2";
  const cellStyle = "" + baseCellStyle;
</script>

{#if references}
  <div class="mt-1 w-full pl-5">
    <Table noborder striped={true}>
      <TableBody>
        {#each references as reference}
          <TableBodyRow>
            <TableBodyCell tdClass={cellStyle}>{reference.category}</TableBodyCell>
            <TableBodyCell tdClass={cellStyle}
              ><p class="mb-2">{reference.summary}</p>
              <p><i class="bx bx-link"></i>{reference.url}</p></TableBodyCell
            >
            <TableBodyCell></TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
{/if}
