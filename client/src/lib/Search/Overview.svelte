<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount, tick } from "svelte";
  import { Button, ButtonGroup, Search, Toggle } from "flowbite-svelte";
  import AdvisoryTable from "$lib/Table/Table.svelte";
  import { searchColumnName } from "$lib/Table/defaults";
  import { SEARCHPAGECOLUMNS, SEARCHTYPES } from "$lib/Queries/query";
  import Queries from "./Queries.svelte";
  import { appStore } from "$lib/store";
  import { querystring } from "svelte-spa-router";
  import { parse } from "qs";
  import Toolbox from "./Toolbox.svelte";

  let searchTerm: string | null;
  let advisoryTable: any;
  let advancedSearch = false;
  let selectedCustomQuery: boolean;
  let queryString: any;
  // let searchqueryTimer: any = null;

  const resetQuery = () => {
    return {
      columns: [...SEARCHPAGECOLUMNS.ADVISORY],
      queryType: SEARCHTYPES.ADVISORY,
      orders: ["-critical"],
      query: "",
      queryReset: ""
    };
  };

  let query = resetQuery();

  const setQueryBack = async () => {
    query = resetQuery();
    searchTerm = "";
    sessionStorage.setItem("documentSearchTerm", "");
    await tick();
    advisoryTable.fetchData();
  };

  $: if (!selectedCustomQuery) {
    setQueryBack();
  }

  const triggerSearch = async () => {
    if (!advancedSearch) {
      if (!selectedCustomQuery) {
        query.query = searchTerm ? `"${searchTerm}" search ${searchColumnName} as` : "";
      } else {
        query.query = `${query.queryReset} ${searchTerm ? `"${searchTerm}" search ${searchColumnName} as and` : ""}`;
      }
      if (
        searchTerm &&
        !query.columns.find((c) => {
          return c === searchColumnName;
        })
      ) {
        query.columns.push(searchColumnName);
      }
      if (!searchTerm)
        query.columns = query.columns.filter((c) => {
          return c !== searchColumnName;
        });
    } else {
      query.columns = query.columns.filter((c) => {
        return c !== searchColumnName;
      });
      if (!selectedCustomQuery) {
        query.query = searchTerm || "";
      } else {
        query.query = `${query.queryReset} ${searchTerm ? searchTerm + " and" : ""}`;
      }
    }
    await tick();
    advisoryTable.fetchData();
  };

  const clearSearch = async () => {
    searchTerm = "";
    query.query = query.queryReset;
    query.columns = query.columns.filter((c) => {
      return c !== searchColumnName;
    });
    await tick();
    advisoryTable.fetchData();
    sessionStorage.setItem("documentSearchTerm", "");
  };

  onMount(async () => {
    if ($querystring) {
      queryString = parse($querystring);
    }
  });
</script>

<svelte:head>
  <title>Search</title>
</svelte:head>

<Queries
  on:querySelected={async (e) => {
    let { detail } = e;
    query = {
      query: detail.query,
      queryReset: detail.query,
      columns: [...detail.columns],
      queryType: detail.kind,
      orders: detail.orders || []
    };
    searchTerm = "";
    await tick();
    advisoryTable.fetchData();
  }}
  {queryString}
  bind:selectedQuery={selectedCustomQuery}
></Queries>
<div class="mb-3 flex">
  <div class="flex w-2/3 flex-row">
    <Search
      size="sm"
      placeholder={advancedSearch ? "Enter a query" : "Enter a search term"}
      bind:value={searchTerm}
      on:keyup={(e) => {
        sessionStorage.setItem("documentSearchTerm", searchTerm ?? "");
        if (e.key === "Enter") triggerSearch();
        // if (searchTerm && searchTerm.length > 2) {
        //   if (searchqueryTimer) clearTimeout(searchqueryTimer);
        //   searchqueryTimer = setTimeout(() => {
        //     triggerSearch();
        //   }, 500);
        // }
        if (searchTerm === "") clearSearch();
      }}
    >
      {#if searchTerm}
        <button
          class="mr-3"
          on:click={() => {
            clearSearch();
          }}>x</button
        >
      {/if}
      <Button
        size="xs"
        class="h-7 py-3.5"
        on:click={() => {
          triggerSearch();
        }}>{advancedSearch ? "Apply" : "Search"}</Button
      >
    </Search>
    <div class="mt-1" title="Define finer grained search queries">
      <Toggle bind:checked={advancedSearch} class="ml-3">Advanced</Toggle>
    </div>
  </div>
  {#if !selectedCustomQuery}
    <ButtonGroup class="ml-auto h-7">
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${query.queryType === SEARCHTYPES.ADVISORY ? "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700" : ""}`}
        on:click={() => {
          query.queryType = SEARCHTYPES.ADVISORY;
          query.columns = SEARCHPAGECOLUMNS.ADVISORY;
          query.orders = query.orders.filter((criterium) => {
            if (criterium.charAt(0) === "-") {
              criterium = criterium.slice(1);
            }
            return SEARCHPAGECOLUMNS.ADVISORY.indexOf(criterium) != -1;
          });
          if (query.orders.length === 0) {
            query.orders = ["-critical"];
          }
          clearSearch();
        }}>Advisories</Button
      >
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${query.queryType === SEARCHTYPES.DOCUMENT ? "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700" : ""}`}
        on:click={() => {
          query.queryType = SEARCHTYPES.DOCUMENT;
          query.columns = SEARCHPAGECOLUMNS.DOCUMENT;
          query.orders = query.orders.filter((criterium) => {
            if (criterium.charAt(0) === "-") {
              criterium = criterium.slice(1);
            }
            return SEARCHPAGECOLUMNS.DOCUMENT.indexOf(criterium) != -1;
          });
          if (query.orders.length === 0) {
            query.orders = ["-critical"];
          }
          clearSearch();
        }}>Documents</Button
      >
      {#if appStore.isEditor() || appStore.isReviewer() || appStore.isAdmin() || appStore.isAuditor()}
        <Button
          size="xs"
          color="light"
          class={`h-7 py-1 text-xs ${query.queryType === SEARCHTYPES.EVENT ? "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700" : ""}`}
          on:click={() => {
            query.queryType = SEARCHTYPES.EVENT;
            query.columns = SEARCHPAGECOLUMNS.EVENT;
            query.orders = ["-time"];
            clearSearch();
          }}>Events</Button
        >
      {/if}
    </ButtonGroup>
  {/if}
</div>
{#if searchTerm !== null}
  <AdvisoryTable
    defaultOrderBy={query.orders}
    columns={query.columns}
    tableType={query.queryType}
    query={`${query.query}`}
    orderBy={query.orders}
    bind:this={advisoryTable}
  ></AdvisoryTable>
{/if}

{#if appStore.isEditor() || appStore.isReviewer()}
  <Toolbox></Toolbox>
{/if}
