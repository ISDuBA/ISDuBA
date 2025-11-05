<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount, tick, untrack } from "svelte";
  import { Toggle } from "flowbite-svelte";
  import AdvisoryTable from "$lib/Table/Table.svelte";
  import { searchColumnName } from "$lib/Table/defaults";
  import { SEARCHPAGECOLUMNS, SEARCHTYPES } from "$lib/Queries/query";
  import Queries from "./Queries.svelte";
  import { appStore } from "$lib/store.svelte";
  import { querystring } from "svelte-spa-router";
  import { parse } from "qs";
  import Toolbox from "./Toolbox.svelte";
  import CSearch from "$lib/Components/CSearch.svelte";
  import TypeToggle from "$lib/Search/TypeToggle.svelte";

  let searchTerm: string = $state("");
  let advisoryTable: any = $state(null);
  let advancedSearch = $state(false);
  let searchResults = $state(true);
  let selectedCustomQuery: boolean = $state(false);
  let queryString: any = $state();
  let defaultQuery: any = $state(null);
  // let searchqueryTimer: any = null;

  $effect(() => {
    untrack(() => selectedCustomQuery);
    untrack(() => query);
    if (defaultQuery) {
      if (!selectedCustomQuery) {
        query = getDefaultQuery();
      }
    }
  });

  const getDefaultQuery = () => {
    let searchType = SEARCHTYPES.ADVISORY;
    if (defaultQuery) {
      searchType = defaultQuery.kind;
    }
    let position = sessionStorage.getItem("tablePosition" + "" + searchType);
    const orderBy = position ? JSON.parse(position)[3] : undefined;
    if (defaultQuery) {
      return {
        columns: defaultQuery.columns,
        queryType: defaultQuery.kind,
        query: defaultQuery.query,
        queryReset: "",
        orders: orderBy ?? defaultQuery.orders
      };
    } else {
      return {
        columns: [...SEARCHPAGECOLUMNS.ADVISORY],
        queryType: SEARCHTYPES.ADVISORY,
        orders: orderBy ?? ["-critical"],
        query: "",
        queryReset: ""
      };
    }
  };

  let query = $state(getDefaultQuery());

  const setQueryBack = async () => {
    query = getDefaultQuery();
    searchTerm = "";
    sessionStorage.setItem("documentSearchTerm", "");
    await tick();
    advisoryTable.fetchData();
  };

  $effect(() => {
    untrack(() => query);
    untrack(() => searchTerm);
    untrack(() => advisoryTable);
    if (!selectedCustomQuery) {
      setQueryBack();
    }
  });

  const triggerSearch = async () => {
    if (!advancedSearch) {
      if (!selectedCustomQuery) {
        query.query = searchTerm ? `"${searchTerm}" search ${searchColumnName} as` : "";
      } else {
        query.query = `${query.queryReset} ${searchTerm ? `"${searchTerm}" search ${searchColumnName} as and` : ""}`;
      }
      if (
        searchTerm &&
        !query.columns.find((c: any) => {
          return c === searchColumnName;
        })
      ) {
        query.columns.push(searchColumnName);
      }
      if (!searchTerm)
        query.columns = query.columns.filter((c: any) => {
          return c !== searchColumnName;
        });
    } else {
      query.columns = query.columns.filter((c: any) => {
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
    query.columns = query.columns.filter((c: any) => {
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

  const filterOrderCriteria = (orders: string[], possibleOrders: string[]) => {
    return orders.filter((criterium) => {
      if (criterium.charAt(0) === "-") {
        criterium = criterium.slice(1);
      }
      return possibleOrders.indexOf(criterium) != -1;
    });
  };
</script>

<svelte:head>
  <title>Search</title>
</svelte:head>

<div class="mb-8 flex flex-wrap justify-between gap-4">
  <Queries
    onQuerySelected={async (detail: any) => {
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
    bind:defaultQuery
  ></Queries>
  {#if !selectedCustomQuery}
    <TypeToggle
      selected={query.queryType}
      eventButtonVisible={appStore.isEditor() ||
        appStore.isReviewer() ||
        appStore.isAdmin() ||
        appStore.isAuditor()}
      onSelect={(newType: SEARCHTYPES) => {
        query.queryType = newType;
        if (newType === SEARCHTYPES.ADVISORY) {
          query.columns = SEARCHPAGECOLUMNS.ADVISORY;
          query.orders = filterOrderCriteria(query.orders, SEARCHPAGECOLUMNS.ADVISORY);
        } else if (newType === SEARCHTYPES.DOCUMENT) {
          query.columns = SEARCHPAGECOLUMNS.DOCUMENT;
          query.orders = filterOrderCriteria(query.orders, SEARCHPAGECOLUMNS.DOCUMENT);
        } else if (newType === SEARCHTYPES.EVENT) {
          query.columns = SEARCHPAGECOLUMNS.EVENT;
        }
        if (
          (newType === SEARCHTYPES.ADVISORY || newType === SEARCHTYPES.DOCUMENT) &&
          query.orders.length === 0
        ) {
          query.orders = ["-critical"];
        } else if (newType === SEARCHTYPES.EVENT) {
          query.orders = ["-time"];
        }
        clearSearch();
      }}
    ></TypeToggle>
  {/if}
</div>
<div class="mb-3 flex flex-row flex-wrap gap-2">
  <CSearch
    buttonText={advancedSearch ? "Apply" : "Search"}
    placeholder={advancedSearch ? "Enter a query" : "Enter a search term"}
    search={() => {
      triggerSearch();
    }}
    onKeyup={(e) => {
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
    bind:searchTerm
  ></CSearch>
  <div class="mt-1" title="Define finer grained search queries">
    <Toggle bind:checked={advancedSearch} class="ml-3">Advanced</Toggle>
  </div>
  <div class="mt-1" title="Show every single time the search term was found">
    <Toggle
      onchange={() => {
        advisoryTable.fetchData();
      }}
      bind:checked={searchResults}
      class="ml-3">Detailed</Toggle
    >
  </div>
</div>
{#if searchTerm !== undefined}
  <AdvisoryTable
    defaultOrderBy={query.orders}
    columns={query.columns}
    tableType={query.queryType}
    query={`${query.query}`}
    orderBy={query.orders}
    bind:this={advisoryTable}
    {searchResults}
  ></AdvisoryTable>
{/if}

{#if appStore.isEditor() || appStore.isReviewer()}
  <Toolbox></Toolbox>
{/if}
