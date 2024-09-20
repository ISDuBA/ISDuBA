<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell, Spinner, Label, Select, PaginationItem } from "flowbite-svelte";
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
  let count = 0;
  let currentPage = 1;
  let numberOfPages = 1000;

  $: numberOfPages = Math.ceil(count / limit);

  const paginationItemClass =
    "text-gray-500 bg-white hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white";
  const paginationItemDeactivatedClass =
    "text-gray-400 bg-gray-50 dark:border-gray-700 dark:text-gray-400 dark:bg-gray-700 cursor-not-allowed";

  const loadFeed = async (id: number) => {
    let result = await fetchFeed(id);
    if (result.ok) {
      feed = result.value;
    } else {
      loadFeedError = result.error;
    }
  };

  const previous = async () => {
    if (offset - limit >= 0) {
      offset = offset - limit > 0 ? offset - limit : 0;
      currentPage -= 1;
    }
    await loadLogs();
  };

  const next = async () => {
    if (offset + limit <= count) {
      offset = offset + limit;
      currentPage += 1;
    }
    await loadLogs();
  };

  const first = async () => {
    offset = 0;
    currentPage = 1;
    await loadLogs();
  };

  const last = async () => {
    offset = (numberOfPages - 1) * limit;
    currentPage = numberOfPages;
    await loadLogs();
  };

  const loadLogs = async () => {
    if (!feed) {
      return;
    }
    if (!feed.id) {
      return;
    }
    loadingLogs = true;
    let result = await fetchFeedLogs(feed.id, offset, limit, true);
    loadingLogs = false;
    if (result.ok) {
      [logs, count] = result.value;
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

  <div class="mx-3 flex w-full flex-row flex-wrap items-center gap-3">
    <Select
      size="sm"
      id="pagecount"
      class="h-7 w-24 p-1 leading-3"
      items={[
        { name: "10", value: 10 },
        { name: "25", value: 25 },
        { name: "50", value: 50 },
        { name: "100", value: 100 }
      ]}
      bind:value={limit}
      on:change={async () => {
        offset = 0;
        currentPage = 1;
        await loadLogs();
      }}
    ></Select>
    <div class="flex flex-row flex-wrap items-center">
      <div class:flex={true} class:mr-3={true}>
        <PaginationItem
          normalClass={currentPage === 1 ? paginationItemDeactivatedClass : paginationItemClass}
          on:click={first}
        >
          <i class="bx bx-arrow-to-left"></i>
        </PaginationItem>
        <PaginationItem
          normalClass={currentPage === 1 ? paginationItemDeactivatedClass : paginationItemClass}
          on:click={previous}
        >
          <i class="bx bx-chevrons-left"></i>
        </PaginationItem>
      </div>

      <div class="flex flex-row flex-wrap items-center">
        <input
          class={`${numberOfPages < 10000 ? "w-16" : "w-20"} cursor-pointer border pr-1 text-right`}
          on:change={() => {
            if (!parseInt("" + currentPage)) currentPage = 1;
            currentPage = Math.floor(currentPage);
            if (currentPage < 1) currentPage = 1;
            if (currentPage > numberOfPages) currentPage = numberOfPages;
            offset = (currentPage - 1) * limit;
            loadLogs();
          }}
          bind:value={currentPage}
        />
        <span class="ml-2 mr-3 w-max text-nowrap">of {numberOfPages} pages</span>
      </div>

      <div class:flex={true}>
        <PaginationItem
          normalClass={currentPage === numberOfPages
            ? paginationItemDeactivatedClass
            : paginationItemClass}
          on:click={next}
        >
          <i class="bx bx-chevrons-right"></i>
        </PaginationItem>
        <PaginationItem
          normalClass={currentPage === numberOfPages
            ? paginationItemDeactivatedClass
            : paginationItemClass}
          on:click={last}
        >
          <i class="bx bx-arrow-to-right"></i>
        </PaginationItem>
      </div>
    </div>
  </div>

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
      <div
        class:invisible={!loadingLogs}
        class={loadingLogs ? "loadingFadeIn" : ""}
        class:mb-4={true}
      >
        Loading ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
    </div>
  </CustomTable>
{/if}

<ErrorMessage error={loadLogsError}></ErrorMessage>
<ErrorMessage error={loadFeedError}></ErrorMessage>
