<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { type Feed, logLevels } from "$lib/Sources/source";
  import { Checkbox, Select, Input, TableBodyCell, Button } from "flowbite-svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass } from "$lib/Table/defaults";

  export let feeds: Feed[] = [];
  export let edit: boolean = false;
</script>

<CustomTable
  title="Feeds"
  headers={[
    {
      label: "Active",
      attribute: "enable"
    },
    {
      label: "URL",
      attribute: "url"
    },
    {
      label: "Log level",
      attribute: "log_level"
    },
    {
      label: "Label",
      attribute: "label"
    }
  ]}
>
  {#each feeds as feed, index (index)}
    <tr>
      <TableBodyCell {tdClass}><Checkbox bind:checked={feed.enable}></Checkbox></TableBodyCell>
      <TableBodyCell {tdClass}>{feed.url}</TableBodyCell>
      <TableBodyCell {tdClass}
        ><Select items={logLevels} bind:value={feed.log_level} /></TableBodyCell
      >
      <TableBodyCell {tdClass}><Input bind:value={feed.label}></Input></TableBodyCell>
      {#if edit}
        <td>
          <Button
            on:click={() => {
              console.log("TODO");
            }}
            title={`Edit feed "${feed.label}"`}
            class="border-0 p-2"
            color="light"
          >
            <i class="bx bx-edit text-xl"></i>
          </Button>
        </td>
        <td>
          <Button
            on:click={(event) => {
              console.log(event);
            }}
            title={`Delete feed "${feed.label}"`}
            class="border-0 p-2"
            color="light"
          >
            <i class="bx bx-trash text-xl text-red-500"></i>
          </Button>
        </td>
      {/if}
    </tr>
  {/each}
</CustomTable>
