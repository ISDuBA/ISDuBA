<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Radio, Input, Spinner, Button, Checkbox } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { COLUMNS, ORDERDIRECTIONS, SEARCHTYPES, generateQueryString } from "$lib/query/query";

  const unsetMessages = () => {
    queryCount = null;
    errorMessage = "";
  };

  const testQuery = async () => {
    loading = true;
    unsetMessages();
    const query = generateQueryString(currentSearch);
    const response = await request(query, "GET");
    if (response.ok) {
      queryCount = response.content.count;
    } else if (response.error) {
      errorMessage = response.error;
    }
    loading = false;
  };

  const columnsFromNames = (columns: string[]) => {
    const result = [...columns];
    return result.map((name) => {
      return {
        name: name,
        visible: false,
        orderBy: null
      };
    });
  };

  const newQuery = () => {
    const columns: any[] = columnsFromNames(COLUMNS.ADVISORY);
    return {
      searchType: SEARCHTYPES.ADVISORY,
      columns: columns,
      name: "New Query",
      query: "",
      description: "",
      global: false
    };
  };

  const switchOrderDirection = (index: number) => {
    if (currentSearch.columns[index].orderBy == null) {
      currentSearch.columns[index].orderBy = ORDERDIRECTIONS.ASC;
      return;
    }
    if (currentSearch.columns[index].orderBy === ORDERDIRECTIONS.ASC) {
      currentSearch.columns[index].orderBy = ORDERDIRECTIONS.DESC;
    } else {
      currentSearch.columns[index].orderBy = null;
    }
  };

  const setVisible = (index: number) => {
    currentSearch.columns[index].visible = !currentSearch.columns[index].visible;
  };

  let currentSearch = newQuery();
  let editName = false;
  let editDescription = false;
  // let hoveredLine = "";
  let queryCount: any = null;
  let loading = false;
  let errorMessage = "";

  const toggleSearchType = () => {};

  const shorten = (text: string) => {
    if (text.length < 10) return text;
    return `${text.substring(0, 10)}...`;
  };
</script>

<SectionHeader title="Configuration"></SectionHeader>
<hr class="mb-6" />
<h2 class="mb-6 text-lg">User defined queries</h2>

<div class="w-2/3">
  <div class="flex h-1 flex-row">
    <div class="flex w-1/3 flex-row items-center gap-x-2">
      <span>Name:</span>
      <button
        on:click={() => {
          editName = !editName;
        }}
      >
        {#if editName}
          <Input
            autofocus
            bind:value={currentSearch.name}
            on:keyup={(e) => {
              if (e.key === "Enter") editName = false;
              if (e.key === "Escape") editName = false;
              e.preventDefault();
            }}
            on:blur={() => {
              editName = false;
            }}
          />
        {:else}
          <div class="flex flex-row items-center" title={currentSearch.name}>
            <h5 class="text-xl font-medium text-gray-500 dark:text-gray-400">
              {shorten(currentSearch.name)}
            </h5>
            <i class="bx bx-edit-alt ml-1"></i>
          </div>
        {/if}
      </button>
    </div>
    <div class="flex w-1/3 flex-row items-center gap-x-2">
      <span>Description:</span>
      <button
        on:click={() => {
          editDescription = !editDescription;
        }}
      >
        {#if editDescription}
          <Input
            autofocus
            bind:value={currentSearch.description}
            on:keyup={(e) => {
              if (e.key === "Enter") editDescription = false;
              if (e.key === "Escape") editDescription = false;
              e.preventDefault();
            }}
            on:blur={() => {
              editDescription = false;
            }}
          />
        {:else}
          <div class="flex flex-row items-center" title={currentSearch.description}>
            <h5 class="text-xl font-medium text-gray-500 dark:text-gray-400">
              {shorten(currentSearch.description)}
            </h5>
            <i class="bx bx-edit-alt ml-1"></i>
          </div>
        {/if}
      </button>
    </div>
    <div class="w-1/3">
      <div class="flex h-1 flex-row items-center gap-x-3">
        <span>Global:</span>
        <Checkbox checked={currentSearch.global}></Checkbox>
      </div>
    </div>
  </div>
  <hr class="mb-4 mt-4 w-4/5" />
  <div class="flex w-1/2 flex-row">
    <div class="w-1/3">
      <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Searching</h5>
    </div>
    <div class="w-1/3">
      <Radio
        name="queryType"
        on:change={toggleSearchType}
        value={SEARCHTYPES.ADVISORY}
        bind:group={currentSearch.searchType}>Advisories</Radio
      >
    </div>
    <div>
      <Radio
        name="queryType"
        on:change={toggleSearchType}
        value={SEARCHTYPES.DOCUMENT}
        bind:group={currentSearch.searchType}>Documents</Radio
      >
    </div>
  </div>
  <div class="mt-4 w-1/2">
    <div class="mb-2 flex flex-row">
      <div class="w-1/3">Column</div>
      <div class="flex w-1/3 flex-row">
        <div>Visible</div>
      </div>
      <div>Order direction</div>
    </div>
    {#each currentSearch.columns as col, index (index)}
      <div class="mb-1 flex flex-row">
        <div class="w-1/3">{col.name}</div>
        <div class="w-1/3">
          <Checkbox
            on:change={() => {
              setVisible(index);
            }}
          ></Checkbox>
        </div>
        <div class="">
          <button
            on:click={() => {
              switchOrderDirection(index);
            }}
          >
            {#if col.orderBy === ORDERDIRECTIONS.ASC}
              <i class="bx bx-sort-a-z"></i>
            {/if}
            {#if col.orderBy === ORDERDIRECTIONS.DESC}
              <i class="bx bx-sort-z-a"></i>
            {/if}
            {#if col.orderBy === null}
              <i class="bx bx-minus"></i>
            {/if}
          </button>
        </div>
      </div>
    {/each}
  </div>
  <div class="mt-6 w-4/5">
    <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Query criteria</h5>
    <div class="flex flex-row">
      <div class="w-full">
        <Input bind:value={currentSearch.query} />
      </div>
    </div>
    <div class="mt-3 flex flex-row">
      {#if loading}
        <div class="mr-4 mt-3">
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
      {/if}
      {#if queryCount !== null}
        <div class:mt-3={true}>
          The query found {queryCount} results.
        </div>
      {/if}
      {#if errorMessage}
        <span class="text-red-600">{errorMessage}</span>
      {/if}
      <div class="my-2 ml-auto flex flex-row gap-3">
        <Button on:click={testQuery} color="light"
          ><i class="bx bx-test-tube me-2"></i> Test query</Button
        >
        <Button
          on:click={() => {
            currentSearch = newQuery();
            queryCount = null;
          }}
          color="light"><i class="bx bx-undo me-2 text-xl"></i> Reset</Button
        >
        <Button color="light"><i class="bx bxs-save me-2"></i> Save</Button>
      </div>
    </div>
  </div>
</div>
