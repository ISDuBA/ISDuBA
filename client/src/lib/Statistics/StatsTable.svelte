<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getLabelForKey, type StatisticGroup } from "$lib/Statistics/statistics";
  import { toLocaleISOString } from "$lib/time";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass } from "$lib/Table/defaults";
  import { TableBodyCell } from "flowbite-svelte";

  export let stats: StatisticGroup | null = null;

  $: dates = stats ? (stats[Object.keys(stats)[0]]?.map((s) => s[0]) ?? []) : [];
</script>

{#if stats}
  <div class="max-h-[400pt] overflow-auto">
    <CustomTable
      headers={[
        { label: "Date", attribute: "date" },
        ...Object.keys(stats).map((s) => {
          return {
            label: getLabelForKey(s),
            attribute: s
          };
        })
      ]}
      stickyHeaders
    >
      {#each dates as date, index}
        <tr class="odd:bg-white even:bg-gray-100 dark:odd:bg-gray-800 dark:even:bg-gray-700">
          {#if date instanceof Date}
            <TableBodyCell {tdClass}>{toLocaleISOString(date)}</TableBodyCell>
          {/if}
          {#each Object.keys(stats) as key}
            {@const count = stats[key]?.[index][1]}
            <TableBodyCell
              tdClass={`${tdClass} ${typeof count === "number" && count !== 0 ? "" : "!text-gray-400"}`}
              >{count ?? 0}</TableBodyCell
            >
          {/each}
        </tr>
      {/each}
    </CustomTable>
  </div>
{/if}
