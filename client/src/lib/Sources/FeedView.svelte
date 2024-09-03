<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { type Feed, logLevels } from "$lib/Sources/source";
  import { Checkbox, Select, Input, TableBodyCell } from "flowbite-svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass } from "$lib/Table/defaults";

  export let feeds: Feed[] = [];
  export let edit: boolean = false;

  export let updateFeed = async (_feed: Feed) => {};
  export let clickFeed = async (_feed: Feed) => {};

  let headers = [
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
  ];

  let headersEdit = [
    ...headers,
    { label: "Downloading", attribute: "downloading" },
    { label: "Waiting", attribute: "waiting" }
  ];
</script>

<CustomTable title="Feeds" headers={edit ? headersEdit : headers}>
  {#each feeds as feed, index (index)}
    <tr on:click={async () => await clickFeed(feed)}>
      <TableBodyCell {tdClass}
        ><Checkbox
          bind:checked={feed.enable}
          on:change={async () => {
            await updateFeed(feed);
            if (!feed.enable) {
              feed.id = undefined;
            }
          }}
        ></Checkbox></TableBodyCell
      >
      <TableBodyCell {tdClass}>{feed.url}</TableBodyCell>
      <TableBodyCell {tdClass}
        ><Select
          items={logLevels}
          bind:value={feed.log_level}
          on:input={async () => await updateFeed(feed)}
        /></TableBodyCell
      >
      {#if edit && !feed.enable}
        <TableBodyCell {tdClass}>N/A</TableBodyCell>
      {:else}
        <TableBodyCell {tdClass}
          ><Input bind:value={feed.label} on:input={async () => await updateFeed(feed)}
          ></Input></TableBodyCell
        >
      {/if}
      {#if edit}
        <TableBodyCell {tdClass}>{feed.stats?.downloading ?? 0}</TableBodyCell>
        <TableBodyCell {tdClass}>{feed.stats?.waiting ?? 0}</TableBodyCell>
      {/if}
    </tr>
  {/each}
</CustomTable>
