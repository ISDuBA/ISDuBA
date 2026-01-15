<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass } from "$lib/Table/defaults";
  import { Button, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  import { onMount } from "svelte";

  interface Related {
    [key: string]: string[];
  }

  let documents: any[] | undefined = $state(undefined);
  let cves: Related | undefined = $state(undefined);

  let headers = $derived.by(() => {
    const tmpHeaders = [
      {
        label: "CVE",
        attribute: "cve"
      }
    ];
    if (documents) {
      tmpHeaders.push(
        ...(documents as string[]).map((document: string) => {
          return {
            label: document,
            attribute: document
          };
        })
      );
    }
    return tmpHeaders;
  });

  onMount(() => {
    //const response = request("ENDPOINT", "GET");
    const docs = [
      "ESA-2024-0001",
      "ESA-2024-0002",
      "ESA-2024-0003",
      "ESA-2024-0004",
      "ESA-2024-0005"
    ];
    cves = {
      "CVE-1970-0001": [docs[3], docs[4]],
      "CVE-1970-0002": [docs[0], docs[1], docs[2], docs[3]],
      "CVE-1970-0003": [docs[0], docs[3], docs[4]],
      "CVE-1970-0004": [docs[1], docs[2], docs[3], docs[4]]
    };
    documents = [];
    Object.keys(cves).forEach((key) => {
      cves?.[key].forEach((doc) => {
        if (!documents?.includes(doc)) {
          documents?.push(doc);
        }
      });
    });
    documents.sort();
  });
</script>

<div>
  {#if documents && cves}
    <CustomTable title="Documents having the same CVEs as document ESA-2024-0001" {headers}>
      {#snippet mainSlot()}
        {#each Object.keys(cves as Related) as string[] as cve, index (index)}
          <TableBodyRow>
            <TableBodyCell class={`${tdClass}`}>{cve}</TableBodyCell>
            {#each documents as _doc}
              <TableBodyCell class={`${tdClass}`}>
                {@const i = Math.random()}
                {#if i < 0.5}
                  <i class="bx bx-check"></i>
                {/if}
              </TableBodyCell>
            {/each}
          </TableBodyRow>
        {/each}
        <TableBodyRow>
          <TableBodyCell class={`${tdClass}`}>Status</TableBodyCell>
          {#each documents as _doc}
            <TableBodyCell class={`${tdClass}`}></TableBodyCell>
          {/each}
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell class={`${tdClass}`}>SSVC</TableBodyCell>
          {#each documents as _doc}
            <TableBodyCell class={`${tdClass}`}></TableBodyCell>
          {/each}
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell class={`${tdClass}`}></TableBodyCell>
          {#each documents as _doc, index}
            <TableBodyCell class={`${tdClass}`}>
              {#if index > 0}
                <Button color="light">Compare</Button>
              {/if}
            </TableBodyCell>
          {/each}
        </TableBodyRow>
      {/snippet}
    </CustomTable>
  {/if}
</div>
