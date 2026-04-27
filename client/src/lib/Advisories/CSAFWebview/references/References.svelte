<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Link from "$lib/Components/Link.svelte";
  import { Table, TableBody, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  import SearchableText from "../SearchableText.svelte";

  interface Props {
    path: string;
    references: any;
  }
  let { path, references }: Props = $props();

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
            <TableBodyCell class={cellStyle}>
              <SearchableText
                textPath={`${path}/references[${i}]/category`}
                text={reference.category}
              />
            </TableBodyCell>
            <TableBodyCell class={cellStyle}
              ><p class="mb-2">
                <SearchableText
                  textPath={`${path}/references[${i}]/summary`}
                  text={reference.summary}
                ></SearchableText>
              </p>
              <Link class="underline" href={reference.url}
                ><i class="bx bx-link"></i>
                <SearchableText textPath={`${path}/references[${i}]/url`} text={reference.url}
                ></SearchableText>
              </Link></TableBodyCell
            >
            <TableBodyCell></TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
{/if}
