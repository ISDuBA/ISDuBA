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
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import type { PaginationParameters } from "./search";

  let { params } = $props();

  let searchTerm: string = $state("");
  let advisoryTable: any = $state(null);
  let advancedSearch = $state(false);
  let searchResults = $state(true);
  let loading = $state(false);
  let selectedCustomQuery: boolean = $state(false);
  let queryString: any = $state();
  let defaultQuery: any = $state(null);
  let openRow: number | null = $state(null);
  let count = $state(0);
  let offset = $state(0);
  let limit = $state(10);
  let currentPage = $state(1);
  let error: ErrorDetails | null = $state(null);
  let prevQuery = "";
  let abortController: AbortController;
  let requestOngoing = false;
  let documents: any = $state(null);
  let postitionRestored: boolean = $state(false);
  // let searchqueryTimer: any = null;

  let numberOfPages = $derived(Math.ceil(count / limit));

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
    fetchData();
  };

  $effect(() => {
    untrack(() => query);
    untrack(() => searchTerm);
    untrack(() => advisoryTable);
    if (!selectedCustomQuery && !params?.searchTerm) {
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
    fetchData();
  };

  const clearSearch = async () => {
    searchTerm = "";
    query.query = query.queryReset;
    query.columns = query.columns.filter((c: any) => {
      return c !== searchColumnName;
    });
    await tick();
    fetchData();
    sessionStorage.setItem("documentSearchTerm", "");
  };

  const savePosition = () => {
    let position = [
      $state.snapshot(offset),
      $state.snapshot(currentPage),
      $state.snapshot(limit),
      $state.snapshot(query).orders
    ];
    sessionStorage.setItem(
      "tablePosition" + query.query + query.queryType,
      JSON.stringify(position)
    );
  };

  const setPaginationParameters = (paginationParameters: PaginationParameters, fetch = true) => {
    if (paginationParameters.offset !== undefined) {
      offset = paginationParameters.offset;
    }
    if (paginationParameters.currentPage !== undefined) {
      currentPage = paginationParameters.currentPage;
    }
    if (paginationParameters.limit !== undefined) {
      limit = paginationParameters.limit;
    }
    if (paginationParameters.orderBy !== undefined) {
      query.orders = paginationParameters.orderBy;
    }
    if (fetch) fetchData();
  };

  const restorePosition = () => {
    let position = sessionStorage.getItem("tablePosition" + query.query + query.queryType);
    if (position) {
      const [offset, currentPage, limit, orderBy] = JSON.parse(position);
      setPaginationParameters(
        {
          offset,
          currentPage,
          limit,
          orderBy
        },
        false
      );
    } else {
      setPaginationParameters(
        {
          offset: 0,
          currentPage: 1,
          limit: 10,
          orderBy: ["-critical"]
        },
        false
      );
    }
  };

  const last = async () => {
    setPaginationParameters({
      currentPage: numberOfPages,
      offset: (numberOfPages - 1) * limit
    });
  };

  async function fetchData(): Promise<void> {
    appStore.setDocuments([]);
    appStore.clearSelectedDocumentIDs();
    openRow = null;
    if (query.query !== prevQuery) {
      restorePosition();
      savePosition();
      prevQuery = query.query;
    }
    const searchColumn = searchTerm ? ` ${searchColumnName}` : "";
    let queryParam = "";
    if (query.query) {
      queryParam = `query=${query.query}`;
    }
    let fetchColumns = [...query.columns];
    let requiredColumns = ["id", "tracking_id", "publisher"];
    for (let c of requiredColumns) {
      if (!fetchColumns.includes(c)) {
        fetchColumns.push(c);
      }
    }
    let documentURL = "";

    if (query.queryType === SEARCHTYPES.EVENT) {
      documentURL = encodeURI(
        `/api/events?${queryParam}&count=1&orders=${query.orders.join(" ")}&limit=${limit}&offset=${offset}&columns=${fetchColumns.join(" ")}${searchColumn}`
      );
    } else {
      const loadAdvisories = query.queryType === SEARCHTYPES.ADVISORY;
      documentURL = encodeURI(
        `/api/documents?${queryParam}&advisories=${loadAdvisories}&count=1&orders=${query.orders.join(" ")}&limit=${limit}&offset=${offset}&results=${searchResults}&columns=${fetchColumns.join(" ")}${searchColumn}`
      );
    }

    error = null;
    loading = true;
    if (!requestOngoing) {
      requestOngoing = true;
      abortController = new AbortController();
    } else {
      abortController.abort();
    }
    const response = await request(documentURL, "GET");
    if (response.ok) {
      ({ count, documents } = response.content);
      if (query.queryType === SEARCHTYPES.EVENT) {
        count = response.content.count;
        documents = response.content.events;
      } else {
        ({ count, documents } = response.content);
      }
      appStore.setDocuments(documents);
      // We are outside the range of available documents,
      // try the last page
      if (offset >= count) {
        await last();
      }
    } else if (response.error) {
      error =
        response.error === "400"
          ? getErrorDetails(`Please check your search syntax.`, response)
          : response.content.includes("deadline exceeded")
            ? getErrorDetails(`The server wasn't able to answer your request in time.`)
            : getErrorDetails(`Could not load query.`, response);
    }
    loading = false;
    requestOngoing = false;
  }

  onMount(async () => {
    if (!postitionRestored) {
      restorePosition();
      postitionRestored = true;
    }
    if ($querystring) {
      queryString = parse($querystring);
    }
    if (params?.searchTerm) {
      searchTerm = params.searchTerm;
      triggerSearch();
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
      fetchData();
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
        savePosition();
        query.queryType = newType;
        restorePosition();
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
        fetchData();
      }}
      bind:checked={searchResults}
      class="ml-3">Detailed</Toggle
    >
  </div>
</div>
{#if searchTerm !== undefined}
  <AdvisoryTable
    columns={query.columns}
    {documents}
    {error}
    {loading}
    {numberOfPages}
    dataChanged={fetchData}
    tableType={query.queryType}
    query={`${query.query}`}
    bind:currentPage
    bind:count
    bind:limit
    bind:offset
    bind:openRow
    bind:orderBy={query.orders}
    {last}
    {setPaginationParameters}
    {searchResults}
  ></AdvisoryTable>
{/if}

{#if appStore.isEditor() || appStore.isReviewer()}
  <Toolbox></Toolbox>
{/if}
