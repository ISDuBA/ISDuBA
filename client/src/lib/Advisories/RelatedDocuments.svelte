<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass, type TableHeader } from "$lib/Table/defaults";
  import { Button, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  import { onMount } from "svelte";

  interface Related {
    [key: string]: string[];
  }

  let documents: any[] | undefined = $state(undefined);
  let cves: Related | undefined = $state(undefined);

  const bodyCellClass = "border-x-2 border-gray-400 dark:border-gray-600 text-center";
  const topBodyCellClass = `${bodyCellClass} rounded-t-lg border-t-2`;
  const bottomBodyCellClass = `${bodyCellClass} rounded-b-lg border-b-2`;

  let headers = $derived.by(() => {
    const tmpHeaders: TableHeader[] = [
      {
        label: "CVE",
        attribute: "cve"
      }
    ];
    if (documents) {
      tmpHeaders.push(
        ...(documents as string[]).map((document: string, index: number) => {
          const header: TableHeader = { label: document, attribute: document };
          if (index === 0) header["class"] = topBodyCellClass;
          else header["class"] = "text-center";
          return header;
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
            {#each documents as _doc, index}
              <TableBodyCell class={`${tdClass} text-center ${index === 0 ? bodyCellClass : ""}`}>
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
          {#each documents as _doc, index}
            <TableBodyCell class={`${tdClass} text-center ${index === 0 ? bodyCellClass : ""}`}
            ></TableBodyCell>
          {/each}
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell class={`${tdClass}`}>SSVC</TableBodyCell>
          {#each documents as _doc, index}
            <TableBodyCell
              class={`${tdClass} text-center ${index === 0 ? bottomBodyCellClass : ""}`}
            ></TableBodyCell>
          {/each}
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell class={`${tdClass}`}></TableBodyCell>
          {#each documents as _doc, index}
            <TableBodyCell class={`${tdClass} text-center`}>
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
