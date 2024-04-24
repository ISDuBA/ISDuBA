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

  let currentSearch = reset();
  let orderBy = "";
  let edit = false;

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

<div class="mb-12 w-1/2">
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
    <div class="w-1/2">
      <h5 class="my-4 text-lg font-medium text-gray-500 dark:text-gray-400">
        Choose columns to include
      </h5>
    </div>
    <div class="w-1/2">
      <h5 class="my-4 text-lg font-medium text-gray-500 dark:text-gray-400">Choosen columns</h5>
    </div>
  </div>
  <div class="flex flex-row">
    <div class="my-3 ml-3 w-1/2">
      <div class="flex flex-col gap-3">
        {#each currentSearch.activeColumns as col}
          <div class="flex items-center">
            <Checkbox bind:group={currentSearch.chosenColumns} value={col}></Checkbox>
            <Badge>{col}</Badge>
          </div>
        {/each}
      </div>
    </div>
    <div class="my-3 ml-3 w-1/2">
      <div class="flex flex-col gap-3">
        {#each currentSearch.chosenColumns as col}
          <div class="flex items-center">
            <Badge>{col}</Badge>
          </div>
        {/each}
      </div>
    </div>
  </div>
</Card>
