<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell, Spinner, Label, PaginationItem, Select } from "flowbite-svelte";
  import { tdClass } from "$lib/Table/defaults";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { onMount } from "svelte";
  import {
    fetchFeedLogs,
    fetchAllFeedLogs,
    fetchFeed,
    type Feed,
    logLevels,
    LogLevel
  } from "./source";
  import ImportStats from "$lib/Statistics/ImportStats.svelte";
  import { DAY_MS } from "$lib/time";
  import Button from "flowbite-svelte/Button.svelte";
  import CSearch from "$lib/Components/CSearch.svelte";
  import DateRange from "$lib/Components/DateRange.svelte";
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
  import debounce from "debounce";

  export let params: any = null;

  let logs: any[] = [];
  let loadingLogs: boolean = false;
  let isDownloadingLogs: boolean = false;
  let downloadAbortController = new AbortController();
  let abortController: AbortController | undefined = undefined;
  let loadFeedError: ErrorDetails | null = null;
  let loadLogsError: ErrorDetails | null = null;
  let loadConfigError: ErrorDetails | null = null;

  let feed: Feed | null = null;

  let offset = 0;
  let limit = 10;
  let count: number | undefined = undefined;
  let currentPage = 1;
  let numberOfPages = 1000;
  let searchTerm = "";
  let selectedLogLevels: LogLevel[] = [];
  let from: Date | undefined;
  let to: Date | undefined = new Date();

  $: numberOfPages = Math.max(1, Math.ceil((count ?? 0) / limit));
  $: realLogLevels = logLevels.filter((l) => l.value !== LogLevel.default);

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
    if (offset + limit <= (count ?? 0)) {
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
    if (!feed || !feed.id) {
      return;
    }
    loadingLogs = true;
    let result = await fetchFeedLogs(
      feed.id,
      offset,
      limit,
      from ? new Date(from) : undefined,
      to ? new Date(to) : undefined,
      searchTerm,
      selectedLogLevels,
      true,
      abortController
    );
    loadingLogs = false;
    if (result.ok) {
      [logs, count] = result.value;
    } else {
      loadLogsError = result.error;
    }
  };

  const delayedLoadLogs = debounce(() => {
    abortController = new AbortController();
    loadLogs();
  }, 600);

  const downloadFeedLogs = async () => {
    if (!feed || !feed.id) {
      return;
    }
    downloadAbortController = new AbortController();
    isDownloadingLogs = true;
    let result = await fetchAllFeedLogs(feed.id, false, downloadAbortController);
    isDownloadingLogs = false;
    if (result.ok) {
      let resultstring = JSON.stringify(result.value);
      let file = new Blob([resultstring], { type: "application/json" });
      let a = document.createElement("a"),
        url = URL.createObjectURL(file);
      a.href = url;
      a.download = feed.label;
      document.body.appendChild(a);
      a.click();
      setTimeout(() => {
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
      }, 0);
    } else if (result.error.message !== "AbortError") {
      loadLogsError = result.error;
    }
  };

  onMount(async () => {
    selectedLogLevels = realLogLevels.map((l) => l.value);
    let id = params?.id;
    if (id) {
      await loadFeed(id);
      await loadLogs();
    }
  });

  const onRangeChanged = () => {
    resetPagination();
    resetPagination();
    delayedLoadLogs();
  };

  const resetPagination = () => {
    offset = 0;
    currentPage = 1;
  };

  const toggleLevel = (level: LogLevel) => {
    if (selectedLogLevels.includes(level)) {
      const index = selectedLogLevels.findIndex((l) => l === level);
      if (index !== -1 && selectedLogLevels.length > 1) {
        selectedLogLevels = selectedLogLevels.toSpliced(index, 1);
      } else if (selectedLogLevels.length === 1) {
        selectedLogLevels = realLogLevels.map((l) => l.value);
      }
    } else {
      selectedLogLevels.push(level);
    }
    // Necessary in case the last checkbox was deselected. Otherwise, the last checkbox would appear as if it
    // was deselected although this isn't the case.
    setTimeout(() => {
      selectedLogLevels = selectedLogLevels;
    }, 0);
    resetPagination();
    loadLogs();
  };
</script>

{#if feed}
  <SectionHeader title={feed.label}>
    <div slot="right" class="flex items-center gap-4">
      <Button
        title="Download all logs"
        on:click={() => {
          if (isDownloadingLogs) {
            downloadAbortController.abort();
          } else {
            downloadFeedLogs();
          }
        }}
        color={isDownloadingLogs ? "red" : "light"}
        class={`ml-3 h-8 py-1 text-xs`}
        outline={isDownloadingLogs}
      >
        {#if isDownloadingLogs}
          <i class="bx bx-x text-lg"></i>
        {:else}
          <i class="bx bx-download text-lg"></i>
        {/if}
      </Button>
      {#if isDownloadingLogs}
        <div class="flex items-center gap-2">
          <Spinner color="gray" size="4"></Spinner>
          <span class="text-sm text-gray-400"
            >Preparing logs. Depending on the amount of logs this may take a while.</span
          >
        </div>
      {/if}
    </div>
  </SectionHeader>

  <div class="mb-4 flex flex-col gap-4">
    <div class="flex flex-wrap gap-x-8 gap-y-6">
      <CSearch
        on:search={() => {
          resetPagination();
          loadLogs();
        }}
        bind:searchTerm
      ></CSearch>
      <DateRange clearable showTimeControls on:change={onRangeChanged} bind:from bind:to
      ></DateRange>
      <div class="flex flex-wrap items-center gap-1">
        <Label for="log-level-selection">Log levels:</Label>
        {#each realLogLevels as level}
          <CCheckbox
            checked={selectedLogLevels.includes(level.value)}
            on:click={() => {
              toggleLevel(level.value);
            }}>{level.name}</CCheckbox
          >
        {/each}
      </div>
    </div>
    <div class="flex w-full flex-row flex-wrap items-center justify-between gap-3">
      <div class="flex items-baseline gap-2">
        <Select
          size="sm"
          id="pagecount"
          class="h-7 w-24 p-1 leading-3"
          items={[
            { name: "10", value: 10 },
            { name: "25", value: 25 },
            { name: "50", value: 50 },
            { name: "100", value: 100 },
            { name: "1000", value: 1000 },
            { name: "10000", value: 10000 }
          ]}
          bind:value={limit}
          on:change={async () => {
            resetPagination();
            await loadLogs();
          }}
        ></Select>
        <Label class="mr-3 text-nowrap">Logs per page</Label>
      </div>
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
            class={`${numberOfPages < 10000 ? "w-16" : "w-20"} cursor-pointer border pr-1 text-right dark:bg-gray-800 dark:text-white`}
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
      <div class="mr-3 text-nowrap">
        {#if count !== undefined}
          {count} entries found
        {/if}
      </div>
    </div>
  </div>

  {#if count}
    <div
      class="mb-8 overflow-y-scroll"
      style={limit === 10 ? "min-height: 350pt;" : "min-height: 500pt;"}
    >
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
            <TableBodyCell tdClass={`${tdClass} min-w-[170pt] align-baseline`}
              >{log.time}</TableBodyCell
            >
            <TableBodyCell tdClass={`${tdClass} align-baseline`}>{log.level}</TableBodyCell>
            <TableBodyCell tdClass={`${tdClass} break-all whitespace-normal`}
              >{log.msg}</TableBodyCell
            >
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
    </div>
  {:else}
    <div class="mb-8"></div>
  {/if}
{/if}

<ErrorMessage error={loadLogsError}></ErrorMessage>
<ErrorMessage error={loadFeedError}></ErrorMessage>
<ErrorMessage error={loadConfigError}></ErrorMessage>

{#if feed?.id}
  <ImportStats
    axes={[{ label: "", types: ["imports"] }]}
    height="200pt"
    initialFrom={new Date(Date.now() - DAY_MS)}
    showModeToggle
    showRangeSelection
    source={{ id: feed.id, isFeed: true }}
    title="Imports"
  ></ImportStats>
  <ImportStats
    axes={[{ label: "", types: ["importFailures"] }]}
    height="200pt"
    initialFrom={new Date(Date.now() - DAY_MS)}
    isStacked
    showLegend
    showModeToggle
    showRangeSelection
    source={{ id: feed.id, isFeed: true }}
    title="Import errors"
  ></ImportStats>
{/if}
