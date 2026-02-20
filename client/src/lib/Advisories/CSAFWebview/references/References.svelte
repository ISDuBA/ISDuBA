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
    references: any;
  }
  let { references }: Props = $props();

  const uid = $props.id();

  const baseCellStyle = "py-2 px-2";
  const cellStyle = "" + baseCellStyle;
</script>

{#if references}
  <div class="mt-1 w-full pl-5">
    <Table border={false} striped={true}>
      <TableBody>
        {#each references as reference, i (`references-${uid}-${i}`)}
          <TableBodyRow>
            <TableBodyCell class={cellStyle}>{reference.category}</TableBodyCell>
            <TableBodyCell class={cellStyle}
              ><p class="mb-2">{reference.summary}</p>
              <a class="underline" href={reference.url}><i class="bx bx-link"></i>{reference.url}</a
              ></TableBodyCell
            >
            <TableBodyCell></TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
{/if}
