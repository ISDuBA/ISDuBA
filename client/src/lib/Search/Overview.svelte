<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount, tick } from "svelte";
  import { Toggle } from "flowbite-svelte";
  import AdvisoryTable from "$lib/Table/Table.svelte";
  import { searchColumnName } from "$lib/Table/defaults";
  import { SEARCHPAGECOLUMNS, SEARCHTYPES } from "$lib/Queries/query";
  import Queries from "./Queries.svelte";
  import { appStore } from "$lib/store.svelte";
  import { push, querystring as qs } from "svelte-spa-router";
  import { parse } from "qs";
  import Toolbox from "./Toolbox.svelte";
  import CSearch from "$lib/Components/CSearch.svelte";
  import TypeToggle from "$lib/Search/TypeToggle.svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import type { SearchParameters } from "./search";

  const INITIAL_LIMIT = 10;
  const INITIAL_ORDER = ["-critical"];

  let advancedSearch = $state(false);
  let loading = $state(false);
  let selectedCustomQuery: boolean = $state(false);
  let defaultQuery: any = $state(null);
  let openRow: number | null = $state(null);
  let count = $state(0);
  let searchTermInputValue = $state("");
  let error: ErrorDetails | null = $state(null);
  let prevQuery = "";
  let abortController: AbortController;
  let requestOngoing = false;
  let documents: any = $state(null);
  // let searchqueryTimer: any = null;

  // Variables derived from URL parameters
  let queryString: any = $derived($qs ? parse($qs) : undefined);
  let searchTerm: string = $derived(queryString?.searchTerm ? queryString.searchTerm : "");
  let orderBy: string[] = $derived(
    queryString?.orderBy ? queryString.orderBy.split(" ") : INITIAL_ORDER
  );
  let detailed: boolean = $derived(
    queryString?.detailed !== undefined ? (queryString.detailed === "true" ? true : false) : true
  );
  let type = $derived(queryString?.type !== undefined ? queryString.type : SEARCHTYPES.ADVISORY);
  let currentPage: number = $derived(Number(queryString?.page ?? 1));
  let limit: number = $derived(Number(queryString?.limit ?? INITIAL_LIMIT));
  let offset: number = $derived(Number((currentPage - 1) * limit));

  let numberOfPages = $derived(Math.ceil(count / limit));

  $effect(() => {
    if (searchTerm) {
      searchTermInputValue = $state.snapshot(searchTerm);
    }
  });

  const getDefaultQuery = (): SearchQuery => {
    if (defaultQuery) {
      return {
        columns: defaultQuery.columns,
        queryType: defaultQuery.kind,
        query: defaultQuery.query,
        queryReset: "",
        orders: defaultQuery.orders
      };
    } else {
      return {
        columns: [...SEARCHPAGECOLUMNS.ADVISORY],
        queryType: SEARCHTYPES.ADVISORY,
        orders: INITIAL_ORDER,
        query: "",
        queryReset: ""
      };
    }
  };

  interface SearchQuery {
    columns: string[];
    queryType: SEARCHTYPES;
    orders: string[];
    query: string;
    queryReset: string;
  }

  let query: SearchQuery = $state(getDefaultQuery());

  const setQueryBack = async () => {
    query = getDefaultQuery();
  };

  const prepareSearch = async () => {
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
  };

  const clearSearch = async () => {
    query.query = query.queryReset;
    query.columns = query.columns.filter((c: any) => {
      return c !== searchColumnName;
    });
  };

  const getCurrentSearchParameters = (): SearchParameters => {
    return {
      type,
      limit,
      currentPage,
      searchTerm,
      orderBy,
      detailed
    };
  };

  const setSearchParameters = async (searchParameters: SearchParameters, fetch = true) => {
    // Don't save non-default parameters in the URL to keep the URL as short as possible.
    let newURL = "/search?";
    if (searchParameters.searchTerm) {
      newURL = newURL.concat(`&searchTerm=${encodeURIComponent(searchParameters.searchTerm)}`);
    } else if (searchParameters.searchTerm === undefined && searchTerm) {
      newURL = newURL.concat(`&searchTerm=${encodeURIComponent(searchTerm)}`);
    }

    if (searchParameters.type && searchParameters.type !== SEARCHTYPES.ADVISORY) {
      newURL = newURL.concat(`&type=${encodeURIComponent(searchParameters.type)}`);
    } else if (searchParameters.type === undefined && type !== SEARCHTYPES.ADVISORY) {
      newURL = newURL.concat(`&type=${encodeURIComponent(type)}`);
    }

    if (
      searchParameters.orderBy &&
      JSON.stringify(searchParameters.orderBy) !== JSON.stringify(INITIAL_ORDER)
    ) {
      newURL = newURL.concat(`&orderBy=${encodeURIComponent(searchParameters.orderBy.join(" "))}`);
    } else if (
      !searchParameters.orderBy &&
      JSON.stringify(orderBy) !== JSON.stringify(INITIAL_ORDER)
    ) {
      newURL = newURL.concat(`&orderBy=${encodeURIComponent(orderBy.join(" "))}`);
    }

    if (searchParameters.currentPage !== undefined && searchParameters.currentPage !== 1) {
      newURL = newURL.concat(`&page=${searchParameters.currentPage}`);
    } else if (currentPage !== 1 && !searchParameters.currentPage) {
      newURL = newURL.concat(`&page=${currentPage}`);
    }

    if (searchParameters.limit !== undefined && searchParameters.limit !== INITIAL_LIMIT) {
      newURL = newURL.concat(`&limit=${searchParameters.limit}`);
    } else if (limit !== INITIAL_LIMIT && !searchParameters.limit) {
      newURL = newURL.concat(`&limit=${limit}`);
    }

    if (searchParameters.detailed !== undefined && searchParameters.detailed !== true) {
      newURL = newURL.concat(`&detailed=${searchParameters.detailed}`);
    } else if (searchParameters.detailed === undefined && detailed !== true) {
      newURL = newURL.concat(`&detailed=${encodeURIComponent(detailed)}`);
    }
    // Don't extend the history if the new URL would not contain any new information. Otherwise
    // the user had to go back multiple times to get to a page with older search parameters.
    if (newURL === "/search?" && $qs?.length && $qs.length > 0) {
      push(newURL.replace("?", ""));
      appStore.setSearchURL(undefined);
    } else if (newURL !== "/search?") {
      push(newURL);
      appStore.setSearchURL(newURL);
    } else {
      appStore.setSearchURL(undefined);
    }

    // Need to wait for the derived values to be updated
    setTimeout(() => {
      if (fetch) {
        prepareSearch();
        fetchData();
      }
    }, 200);
  };

  const last = async () => {
    setSearchParameters({
      currentPage: Math.max(numberOfPages, 1)
    });
  };

  async function fetchData(): Promise<void> {
    appStore.setDocuments([]);
    appStore.clearSelectedDocumentIDs();
    openRow = null;
    if (query.query !== prevQuery) {
      prevQuery = query.query;
    }
    const searchColumn = searchTerm ? ` ${searchColumnName}` : "";
    let queryParam = "";
    if (query.query) {
      queryParam = `query=${query.query}`;
    }
    const orderByParam = selectedCustomQuery ? (query.orders ?? []) : orderBy;
    let fetchColumns = [...query.columns];
    let requiredColumns = ["id", "tracking_id", "publisher"];
    for (let c of requiredColumns) {
      if (!fetchColumns.includes(c)) {
        fetchColumns.push(c);
      }
    }
    let documentURL = "";

    if (
      (selectedCustomQuery && query.queryType === SEARCHTYPES.EVENT) ||
      type === SEARCHTYPES.EVENT
    ) {
      documentURL = encodeURI(
        `/api/events?${queryParam}&count=1&orders=${orderByParam.join(" ")}&limit=${limit}&offset=${offset}&columns=${fetchColumns.join(" ")}${searchColumn}`
      );
    } else {
      const loadAdvisories = type === SEARCHTYPES.ADVISORY;
      documentURL = encodeURI(
        `/api/documents?${queryParam}&advisories=${loadAdvisories}&count=1&orders=${orderByParam.join(" ")}&limit=${limit}&offset=${offset}&results=${detailed}&columns=${fetchColumns.join(" ")}${searchColumn}`
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

  const filterOrderCriteria = (orders: string[], possibleOrders: string[]) => {
    return orders.filter((criterium) => {
      if (criterium.charAt(0) === "-") {
        criterium = criterium.slice(1);
      }
      return possibleOrders.indexOf(criterium) != -1;
    });
  };

  onMount(() => {
    fetchData();
  });
</script>

<svelte:head>
  <title>Search</title>
</svelte:head>

<div class="mb-8 flex flex-wrap justify-between gap-4">
  <Queries
    onQuerySelected={(detail: any) => {
      if (detail) {
        query = {
          query: detail.query,
          queryReset: detail.query,
          columns: [...detail.columns],
          queryType: detail.kind,
          orders: detail.orders || []
        };
      } else {
        setQueryBack();
      }
      setSearchParameters({
        searchTerm: ""
      });
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
        appStore.setSearchParametersForType(type, getCurrentSearchParameters());
        query.queryType = newType;
        const newParameters = $state.snapshot(appStore.getSearchParametersForType(newType));
        if (newType === SEARCHTYPES.ADVISORY) {
          query.columns = SEARCHPAGECOLUMNS.ADVISORY;
          query.orders = filterOrderCriteria(query.orders, SEARCHPAGECOLUMNS.ADVISORY);
        } else if (newType === SEARCHTYPES.DOCUMENT) {
          query.columns = SEARCHPAGECOLUMNS.DOCUMENT;
          query.orders = filterOrderCriteria(query.orders, SEARCHPAGECOLUMNS.DOCUMENT);
        } else if (newType === SEARCHTYPES.EVENT) {
          query.columns = SEARCHPAGECOLUMNS.EVENT;
        }
        clearSearch();
        searchTermInputValue = "";
        if (newParameters) {
          searchTermInputValue = "";
          setSearchParameters(newParameters);
        } else {
          setSearchParameters({ searchTerm: "", type: newType });
        }
      }}
    ></TypeToggle>
  {/if}
</div>
<div class="mb-3 flex flex-row flex-wrap gap-2">
  <CSearch
    buttonText={advancedSearch ? "Apply" : "Search"}
    placeholder={advancedSearch ? "Enter a query" : "Enter a search term"}
    search={(term) => {
      if (term === "") {
        setSearchParameters({
          searchTerm: ""
        });
      } else {
        setSearchParameters({
          searchTerm: term
        });
      }
    }}
    searchTerm={searchTermInputValue}
  ></CSearch>
  <div class="mt-1" title="Define finer grained search queries">
    <Toggle bind:checked={advancedSearch} class="ml-3">Advanced</Toggle>
  </div>
  <div class="mt-1" title="Show every single time the search term was found">
    <Toggle
      onclick={() => {
        setSearchParameters({ detailed: !$state.snapshot(detailed) });
      }}
      checked={detailed}
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
    {limit}
    {offset}
    {currentPage}
    {orderBy}
    dataChanged={fetchData}
    tableType={query.queryType}
    query={`${query.query}`}
    bind:count
    bind:openRow
    {last}
    {setSearchParameters}
  ></AdvisoryTable>
{/if}

{#if appStore.isEditor() || appStore.isReviewer()}
  <Toolbox></Toolbox>
{/if}
