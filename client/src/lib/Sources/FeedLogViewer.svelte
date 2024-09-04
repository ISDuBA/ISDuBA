<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell, Spinner, Label, Select } from "flowbite-svelte";
  import { tdClass } from "$lib/Table/defaults";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { onMount } from "svelte";
  import { fetchFeedLogs, fetchFeed, type Feed } from "./source";

  export let params: any = null;

  let logs: any[] = [];
  let loadingLogs: boolean = false;
  let loadFeedError: ErrorDetails | null = null;
  let loadLogsError: ErrorDetails | null = null;

  let feed: Feed | null = null;

  let offset = 0;
  let limit = 10;

  const loadFeed = async (id: number) => {
    let result = await fetchFeed(id);
    if (result.ok) {
      feed = result.value;
    } else {
      loadFeedError = result.error;
    }
  };

  const loadLogs = async () => {
    if (!feed) {
      return;
    }
    if (!feed.id) {
      return;
    }
    loadingLogs = true;
    let result = await fetchFeedLogs(feed.id, offset, limit);
    loadingLogs = false;
    if (result.ok) {
      logs = result.value;
    } else {
      loadLogsError = result.error;
    }
  };

  onMount(async () => {
    let id = params?.id;
    if (id) {
      await loadFeed(id);
      await loadLogs();
    }
  });
</script>

{#if feed}
  <SectionHeader title={feed.label}></SectionHeader>
  <Label class="mr-3 text-nowrap">Logs per page</Label>
  <Select
    size="sm"
    id="pagecount"
    class="mt-2 h-7 w-24 p-1 leading-3"
    items={[
      { name: "10", value: 10 },
      { name: "25", value: 25 },
      { name: "50", value: 50 },
      { name: "100", value: 100 }
    ]}
    bind:value={limit}
    on:change={async () => {
      await loadLogs();
    }}
  ></Select>
  <CustomTable
    title="Logs"
    headers={[
      {
        label: "Time",
        attribute: "time"
      },
      {
        label: "level",
        attribute: "level"
      },
      {
        label: "Message",
        attribute: "msg"
      }
    ]}
  >
    {#each logs as log, index (index)}
      <tr>
        <TableBodyCell {tdClass}>{log.time}</TableBodyCell>
        <TableBodyCell {tdClass}>{log.level}</TableBodyCell>
        <TableBodyCell {tdClass}>{log.msg}</TableBodyCell>
      </tr>
    {/each}
    <div slot="bottom">
      <div class:hidden={!loadingLogs} class:mb-4={true}>
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
    </div>
  </CustomTable>
{/if}

<ErrorMessage error={loadLogsError}></ErrorMessage>
<ErrorMessage error={loadFeedError}></ErrorMessage>
