<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Radio, Badge, Input, Spinner, Button, Checkbox } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import {
    COLUMNS,
    ORDERDIRECTIONS,
    SEARCHTYPES,
    generateQueryString,
    newQuery
  } from "$lib/query/query";

  const chooseColumn = (col: any) => {
    unsetMessages();
    currentSearch.chosenColumns = [
      ...currentSearch.chosenColumns,
      { name: col, searchOrder: ORDERDIRECTIONS.ASC }
    ];
    currentSearch.activeColumns = currentSearch.activeColumns.filter((column: any) => {
      return column !== col;
    });
  };

  const undoColumn = (col: any) => {
    unsetMessages();
    currentSearch.chosenColumns = currentSearch.chosenColumns.filter((column: any) => {
      return column.name !== col.name;
    });
    const activeColumns =
      currentSearch.searchType === SEARCHTYPES.ADVISORY ? COLUMNS.ADVISORY : COLUMNS.DOCUMENT;
    const activeColumnsToChoose = new Set([...currentSearch.activeColumns, col.name]);
    currentSearch.activeColumns = activeColumns.reduce((acc: string[], column) => {
      if (activeColumnsToChoose.has(column)) acc.push(column);
      return acc;
    }, []);
  };

  const hoverLine = (col: any) => {
    hoveredLine = col;
  };

  const indexOfCol = (col: any) => {
    return currentSearch.chosenColumns.map((col: any) => col.name).indexOf(col);
  };

  const changeSearchType = () => {
    if (currentSearch.searchType === SEARCHTYPES.ADVISORY) {
      currentSearch.activeColumns = [...COLUMNS.ADVISORY];
    } else {
      currentSearch.activeColumns = [...COLUMNS.DOCUMENT];
    }
    currentSearch.chosenColumns = [];
  };

  const promoteColumn = (col: any) => {
    unsetMessages();
    const index = indexOfCol(col);
    if (index === 0) return;
    const tmp = currentSearch.chosenColumns[index - 1];
    currentSearch.chosenColumns[index - 1] = currentSearch.chosenColumns[index];
    currentSearch.chosenColumns[index] = tmp;
  };

  const demoteColumn = (col: any) => {
    unsetMessages();
    const index = indexOfCol(col);
    if (index === currentSearch.chosenColumns.length - 1) return;
    const tmp = currentSearch.chosenColumns[index + 1];
    currentSearch.chosenColumns[index + 1] = currentSearch.chosenColumns[index];
    currentSearch.chosenColumns[index] = tmp;
  };

  const switchOrder = (col: any) => {
    const index = indexOfCol(col);
    const selectedCol = currentSearch.chosenColumns[index];
    let searchOrder = ORDERDIRECTIONS.DESC;
    if (selectedCol.searchOrder === ORDERDIRECTIONS.DESC) {
      searchOrder = ORDERDIRECTIONS.ASC;
    }
    currentSearch.chosenColumns[index] = {
      name: selectedCol.name,
      searchOrder: searchOrder
    };
  };

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

  let currentSearch = newQuery();
  let editDescription = false;
  let hoveredLine = "";
  let queryCount: any = null;
  let loading = false;
  let errorMessage = "";
</script>

<SectionHeader title="Configuration"></SectionHeader>
<hr class="mb-6" />
<h2 class="mb-6 text-lg">User defined queries</h2>

<div class="w-2/5">
  <div class="flex flex-row">
    <div class="mt-0 w-1/2">
      <span class="mr-3">Description:</span>
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
          <h5 class="text-xl font-medium text-gray-500 dark:text-gray-400">
            {currentSearch.description}
          </h5>
        {/if}
      </button>
      <button
        on:click={() => {
          editDescription = !editDescription;
        }}><i class="bx bx-edit-alt ml-1"></i></button
      >
    </div>
    <div class="mt-1 flex w-1/2 flex-row">
      <span class="mr-3">Global:</span>
      <Checkbox checked={currentSearch.global}></Checkbox>
    </div>
  </div>
  <hr class="mb-4 mt-2" />
  <div class="flex flex-row">
    <div class="flex flex-row gap-3">
      <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Searching</h5>
      <Radio
        name="queryType"
        on:change={changeSearchType}
        value={SEARCHTYPES.ADVISORY}
        bind:group={currentSearch.searchType}>Advisories</Radio
      >
      <Radio
        name="queryType"
        on:change={changeSearchType}
        value={SEARCHTYPES.DOCUMENT}
        bind:group={currentSearch.searchType}>Documents</Radio
      >
    </div>
  </div>
  <div class="mt-4 flex flex-row">
    <div class="w-1/2">
      <h5 class="my-1 text-lg font-medium text-gray-500 dark:text-gray-400">Available columns</h5>
    </div>
    <div class="w-1/2">
      <h5 class="my-1 text-lg font-medium text-gray-500 dark:text-gray-400">Choosen columns</h5>
    </div>
  </div>
  <div class="ml-3 flex flex-row">
    <div class="w-1/2">
      <Button
        outline
        size="xs"
        on:click={() => {
          const columns =
            currentSearch.searchType === SEARCHTYPES.ADVISORY ? COLUMNS.ADVISORY : COLUMNS.DOCUMENT;
          currentSearch.chosenColumns = columns.map((col) => {
            return {
              name: col,
              searchOrder: ORDERDIRECTIONS.ASC
            };
          });
          currentSearch.activeColumns = [];
        }}>All <i class="bx bx-right-arrow-alt"></i></Button
      >
    </div>
    <div class="ml-7 w-1/2">
      <Button
        outline
        size="xs"
        on:click={() => {
          const columns =
            currentSearch.searchType === SEARCHTYPES.ADVISORY ? COLUMNS.ADVISORY : COLUMNS.DOCUMENT;
          currentSearch.activeColumns = columns;
          currentSearch.chosenColumns = [];
        }}>None <i class="bx bx-left-arrow-alt"></i></Button
      >
    </div>
  </div>
  <div style="height:30rem" class="flex flex-row">
    <div class="my-3 ml-3 w-1/2">
      <div class="flex flex-col gap-3">
        {#each currentSearch.activeColumns as col}
          <div class="flex items-center">
            <button
              on:click={() => {
                chooseColumn(col);
              }}
              title={`${col} column`}><Badge>{col}</Badge></button
            >
          </div>
        {/each}
      </div>
    </div>
    <div class="flex w-3 flex-col items-center"></div>
    <div class="ml-2 mr-3"></div>
    <div class="my-3 w-1/2">
      <div class="flex flex-col leading-3">
        {#each currentSearch.chosenColumns as col}
          <div
            role="presentation"
            class="flex items-center"
            on:focus={() => {}}
            on:mouseleave={() => {
              hoveredLine = "";
            }}
            on:mouseover={() => {
              hoverLine(col.name);
            }}
          >
            <div class:invisible={hoveredLine !== col.name} class:flex={true} class:flex-col={true}>
              <button
                on:click={() => {
                  promoteColumn(col.name);
                }}
                title="Promote column"
              >
                <i class="bx bxs-up-arrow-circle"></i>
              </button>
              <button
                on:click={() => {
                  demoteColumn(col.name);
                }}
                title="Demote column"
              >
                <i class="bx bxs-down-arrow-circle"></i>
              </button>
            </div>
            <button
              on:click={() => {
                undoColumn(col);
              }}
              title={`${col.name} column`}><Badge>{col.name}</Badge></button
            >
            <button
              class="ml-1"
              on:click={() => {
                switchOrder(col.name);
              }}
            >
              {#if col.searchOrder === ORDERDIRECTIONS.ASC}
                <i class="bx bx-sort-a-z" title={"Ascending order"}></i>
              {:else}
                <i class="bx bx-sort-z-a" title={"Descending order"}></i>
              {/if}
            </button>
          </div>
        {/each}
      </div>
    </div>
  </div>
  <h5 class="mt-2 text-lg font-medium text-gray-500 dark:text-gray-400">Query criteria</h5>
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
