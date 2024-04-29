<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { tablePadding } from "$lib/table/defaults";
  import {
    Card,
    Radio,
    Badge,
    Input,
    Table,
    TableHead,
    TableBody,
    TableHeadCell,
    TableBodyRow,
    TableBodyCell,
    Checkbox
  } from "flowbite-svelte";

  const COLUMNS = {
    ADVISORY: [
      "id",
      "tracking_id",
      "version",
      "publisher",
      "current_release_date",
      "initial_release_date",
      "title",
      "tlp",
      "cvss_v2_score",
      "cvss_v3_score",
      "ssvc",
      "four_cves",
      "state"
    ],
    DOCUMENT: [
      "id",
      "tracking_id",
      "version",
      "publisher",
      "current_release_date",
      "initial_release_date",
      "title",
      "tlp",
      "cvss_v2_score",
      "cvss_v3_score",
      "four_cves"
    ]
  };

  const ORDERDIRECTIONS = {
    ASC: "asc",
    DESC: "desc"
  };

  const SEARCHTYPES = {
    ADVISORY: "Advisory",
    DOCUMENT: "Document"
  };

  const STUBQUERIES = [
    { name: "Redhat", description: "Show all RedHat advisories" },
    { name: "Sick", description: "Show all Sick advisories" }
  ];

  const reset = () => {
    return {
      currentStep: 1,
      searchType: SEARCHTYPES.ADVISORY,
      chosenColumns: [],
      activeColumns: [],
      name: "New Query",
      description: ""
    };
  };

  const chooseColumn = (e, col) => {
    if (e.currentTarget.checked) {
      currentSearch.chosenColumns = [
        ...currentSearch.chosenColumns,
        { name: col, searchOrder: ORDERDIRECTIONS.ASC }
      ];
    } else {
      currentSearch.chosenColumns = currentSearch.chosenColumns.filter((column) => {
        return column.name !== col;
      });
    }
  };

  const hoverLine = (col) => {
    hoveredLine = col;
  };

  const indexOfCol = (col) => {
    return currentSearch.chosenColumns.map((col) => col.name).indexOf(col);
  };

  const promoteColumn = (col) => {
    const index = indexOfCol(col);
    if (index === 0) return;
    const tmp = currentSearch.chosenColumns[index - 1];
    currentSearch.chosenColumns[index - 1] = currentSearch.chosenColumns[index];
    currentSearch.chosenColumns[index] = tmp;
  };

  const demoteColumn = (col) => {
    const index = indexOfCol(col);
    if (index === currentSearch.chosenColumns.length - 1) return;
    const tmp = currentSearch.chosenColumns[index + 1];
    currentSearch.chosenColumns[index + 1] = currentSearch.chosenColumns[index];
    currentSearch.chosenColumns[index] = tmp;
  };

  const switchOrder = (col) => {
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

  let currentSearch = reset();
  let orderBy = "";
  let edit = false;
  let hoveredLine = "";

  $: {
    if (currentSearch.searchType === SEARCHTYPES.ADVISORY) {
      currentSearch.activeColumns = [...COLUMNS.ADVISORY];
    }
    if (currentSearch.searchType === SEARCHTYPES.DOCUMENT) {
      currentSearch.activeColumns = [...COLUMNS.DOCUMENT];
    }
  }
</script>

<h2 class="mb-3 text-lg">User defined queries</h2>

<div class="flex flex-row">
  <div class="mb-12 w-1/3">
    <Table hoverable={true} noborder={true}>
      <TableHead class="cursor-pointer">
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Name<i
            class:bx={true}
            class:bx-caret-up={orderBy == "name"}
            class:bx-caret-down={orderBy == "-name"}
          ></i>
        </TableHeadCell>
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Description<i
            class:bx={true}
            class:bx-caret-up={orderBy == "name"}
            class:bx-caret-down={orderBy == "-name"}
          ></i>
        </TableHeadCell>
      </TableHead>
      <TableBody>
        {#each STUBQUERIES as query}
          <TableBodyRow class="cursor-pointer">
            <TableBodyCell>{query.name}</TableBodyCell>
            <TableBodyCell>{query.description}</TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>

  <Card size="lg">
    <div class="my-3">
      <span class="mr-3">Name:</span>
      <button
        on:click={() => {
          edit = !edit;
        }}
      >
        {#if edit}
          <Input
            autofocus
            bind:value={currentSearch.name}
            on:keyup={(e) => {
              if (e.key === "Enter") edit = false;
              if (e.key === "Escape") edit = false;
              e.preventDefault();
            }}
            on:blur={() => {
              edit = false;
            }}
          />
        {:else}
          <h5 class="mb-4 text-xl font-medium text-gray-500 dark:text-gray-400">
            {currentSearch.name}
          </h5>
        {/if}
      </button>
      <button
        on:click={() => {
          edit = !edit;
        }}><i class="bx bx-edit-alt ml-1"></i></button
      >
    </div>
    <h5 class="mb-4 text-lg font-medium text-gray-500 dark:text-gray-400">Choose type of search</h5>
    <div class="ml-3 flex flex-row gap-3">
      <Radio name="queryType" value={SEARCHTYPES.ADVISORY} bind:group={currentSearch.searchType}
        >Advisories</Radio
      >
      <Radio name="queryType" value={SEARCHTYPES.DOCUMENT} bind:group={currentSearch.searchType}
        >Documents</Radio
      >
    </div>
    <div class="flex flex-row">
      <div class="w-1/3">
        <h5 class="my-4 text-lg font-medium text-gray-500 dark:text-gray-400">Available columns</h5>
      </div>
      <div class="w-2/3">
        <h5 class="my-4 text-lg font-medium text-gray-500 dark:text-gray-400">Choosen columns</h5>
      </div>
    </div>
    <div class="flex flex-row">
      <div class="my-3 ml-3 w-1/3">
        <div class="flex flex-col gap-3">
          {#each currentSearch.activeColumns as col}
            <div class="flex items-center">
              <Checkbox
                on:click={(e) => {
                  chooseColumn(e, col);
                }}
              ></Checkbox>
              <Badge>{col}</Badge>
            </div>
          {/each}
        </div>
      </div>
      <div class="my-3 ml-3 w-2/3">
        <div class="flex flex-col leading-3">
          {#each currentSearch.chosenColumns as col}
            <div
              class="flex items-center"
              on:mouseleave={() => {
                hoveredLine = "";
              }}
              on:mouseover={hoverLine(col.name)}
            >
              <div
                class:invisible={hoveredLine !== col.name}
                class:flex={true}
                class:flex-col={true}
              >
                <button on:click={promoteColumn(col.name)}>
                  <i class="bx bxs-up-arrow-circle"></i>
                </button>
                <button on:click={demoteColumn(col.name)}>
                  <i class="bx bxs-down-arrow-circle"></i>
                </button>
              </div>
              <Badge>{col.name}</Badge>
              <button class="ml-1" on:click={switchOrder(col.name)}>
                {#if col.searchOrder === ORDERDIRECTIONS.ASC}
                  <i class="bx bx-sort-a-z"></i>
                {:else}
                  <i class="bx bx-sort-z-a"></i>
                {/if}
              </button>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </Card>
</div>
