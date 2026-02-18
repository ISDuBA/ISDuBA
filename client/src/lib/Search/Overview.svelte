<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Toggle } from "flowbite-svelte";
  import AdvisoryTable from "$lib/Table/Table.svelte";
  import { searchColumnName } from "$lib/Table/defaults";
  import { SEARCHPAGECOLUMNS, SEARCHTYPES } from "$lib/Queries/query";
  import type { Query } from "$lib/Queries/query";
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

  let loading = $state(false);
  let queries: any[] = $state([]);
  let defaultQuery: any = $state(null);
  let openRow: number | null = $state(null);
  let count = $state(0);
  let searchTermInputValue = $state("");
  let error: ErrorDetails | null = $state(null);
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
  let advanced: boolean = $derived(
    queryString?.advanced !== undefined ? (queryString.advanced === "true" ? true : false) : false
  );
  let detailed: boolean = $derived(
    queryString?.detailed !== undefined ? (queryString.detailed === "true" ? true : false) : true
  );
  let type = $derived(queryString?.type !== undefined ? queryString.type : SEARCHTYPES.ADVISORY);
  let currentPage: number = $derived(Number(queryString?.page ?? 1));
  let limit: number = $derived(Number(queryString?.limit ?? INITIAL_LIMIT));
  let offset: number = $derived(Number((currentPage - 1) * limit));
  let queryID: number | undefined = $derived(
    queryString?.queryID ? Number(queryString.queryID) : undefined
  );
  let selectedQuery: Query | null = $derived.by(() => {
    if (!queryID) return null;
    return $state.snapshot(queries).find((q) => q.id === queryID) ?? null;
  });

  let numberOfPages = $derived(Math.ceil(count / limit));

  $effect(() => {
    if (searchTerm) {
      searchTermInputValue = $state.snapshot(searchTerm);
    } else {
      searchTermInputValue = "";
    }
  });

  interface SearchQuery {
    columns: string[];
    queryType: SEARCHTYPES;
    orders: string[];
    query: string;
    queryReset: string;
  }

  let query: SearchQuery = $derived.by(() => {
    if (queryID !== undefined) {
      return queries.find((q) => q.id === queryID);
    } else if (defaultQuery) {
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
  });

  let queryQuery: string = $derived.by(() => {
    if (!advanced) {
      if (!selectedQuery) {
        return searchTerm ? `"${searchTerm}" search ${searchColumnName} as` : "";
      } else {
        return `${selectedQuery.query} ${searchTerm ? `"${searchTerm}" search ${searchColumnName} as and` : ""}`;
      }
    } else {
      if (!selectedQuery) {
        return searchTerm || "";
      } else {
        return `${selectedQuery.queryReset} ${searchTerm ? searchTerm + " and" : ""}`;
      }
    }
  });

  let columns = $derived.by(() => {
    let tmpColumns: string[];
    if (selectedQuery) {
      tmpColumns = selectedQuery.columns;
    } else {
      if (type === SEARCHTYPES.ADVISORY) {
        tmpColumns = SEARCHPAGECOLUMNS.ADVISORY;
      } else if (type === SEARCHTYPES.DOCUMENT) {
        tmpColumns = SEARCHPAGECOLUMNS.DOCUMENT;
      } else {
        tmpColumns = SEARCHPAGECOLUMNS.EVENT;
      }
    }
    if (!advanced) {
      if (
        searchTerm &&
        !tmpColumns.find((c: any) => {
          return c === searchColumnName;
        })
      ) {
        tmpColumns.push(searchColumnName);
      }
      if (!searchTerm)
        tmpColumns = tmpColumns.filter((c: any) => {
          return c !== searchColumnName;
        });
    } else {
      tmpColumns = tmpColumns.filter((c: any) => {
        return c !== searchColumnName;
      });
    }
    return tmpColumns;
  });

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

  const setSearchParameters = async (searchParameters: SearchParameters) => {
    // Don't save non-default parameters in the URL to keep the URL as short as possible.
    let newURL = "/search?";
    if (searchParameters.searchTerm) {
      newURL = newURL.concat(`&searchTerm=${encodeURIComponent(searchParameters.searchTerm)}`);
    } else if (searchParameters.searchTerm === undefined && searchTerm) {
      newURL = newURL.concat(`&searchTerm=${encodeURIComponent(searchTerm)}`);
    }

    if (searchParameters.type && searchParameters.type !== SEARCHTYPES.ADVISORY) {
      newURL = newURL.concat(`&type=${encodeURIComponent(searchParameters.type)}`);
    } else if (!Object.keys(searchParameters).includes("type") && type !== SEARCHTYPES.ADVISORY) {
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
    } else if (!Object.keys(searchParameters).includes("detailed") && detailed !== true) {
      newURL = newURL.concat(`&detailed=${encodeURIComponent(detailed)}`);
    }

    if (searchParameters.advanced !== undefined && searchParameters.advanced !== false) {
      newURL = newURL.concat(`&advanced=${searchParameters.advanced}`);
    } else if (!Object.keys(searchParameters).includes("advanced") && advanced !== false) {
      newURL = newURL.concat(`&advanced=${encodeURIComponent(advanced)}`);
    }

    if (searchParameters.queryID !== undefined) {
      newURL = newURL.concat(`&queryID=${searchParameters.queryID}`);
    } else if (!Object.keys(searchParameters).includes("queryID") && queryID !== undefined) {
      newURL = newURL.concat(`&queryID=${encodeURIComponent(queryID)}`);
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
  };

  const last = async () => {
    setSearchParameters({
      currentPage: Math.max(numberOfPages, 1)
    });
  };

  $effect(() => {
    if ($qs !== undefined) {
      fetchData();
    }
  });

  async function fetchData(): Promise<void> {
    appStore.setSearchResults([]);
    appStore.clearSelectedDocumentIDs();
    openRow = null;

    const searchColumn = searchTerm ? ` ${searchColumnName}` : "";
    let queryParam = "";
    if (queryQuery) {
      queryParam = `query=${queryQuery}`;
    }
    const orderByParam = selectedQuery ? (query.orders ?? []) : orderBy;
    let fetchColumns = [...$state.snapshot(columns)];
    let requiredColumns = ["id", "tracking_id", "publisher"];
    for (let c of requiredColumns) {
      if (!fetchColumns.includes(c)) {
        fetchColumns.push(c);
      }
    }
    let URLWithoutOffsetAndLimit = "";
    const columnsParam = `columns=${fetchColumns.join(" ")}${searchColumn}`;
    appStore.setSearchOffset(offset);

    if ((selectedQuery && selectedQuery.kind === SEARCHTYPES.EVENT) || type === SEARCHTYPES.EVENT) {
      URLWithoutOffsetAndLimit = encodeURI(
        `/api/events?${queryParam}&count=1&orders=${orderByParam.join(" ")}&${columnsParam}`
      );
    } else {
      const loadAdvisories = type === SEARCHTYPES.ADVISORY;
      URLWithoutOffsetAndLimit = encodeURI(
        `/api/documents?${queryParam}&advisories=${loadAdvisories}&aggregate=true&count=1&orders=${orderByParam.join(" ")}&results=${detailed}&${columnsParam}`
      );
    }
    appStore.setSearchRequestURL(URLWithoutOffsetAndLimit);
    const documentURL = URLWithoutOffsetAndLimit + `&offset=${offset}&limit=${limit}`;

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
      if (
        (selectedQuery && selectedQuery.kind === SEARCHTYPES.EVENT) ||
        type === SEARCHTYPES.EVENT
      ) {
        count = response.content.count;
        documents = response.content.events;
      } else {
        ({ count, documents } = JSON.parse(response.content));
      }
      appStore.setSearchResults($state.snapshot(documents));
      appStore.setSearchResultCount($state.snapshot(count));
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
</script>

<svelte:head>
  <title>Search</title>
</svelte:head>

<div class="mb-8 flex flex-wrap justify-between gap-4">
  <Queries
    onQuerySelected={(id: number | undefined) => {
      const newParameters: SearchParameters = {
        currentPage: 1,
        limit: INITIAL_LIMIT,
        queryID: id,
        searchTerm: "",
        type: undefined
      };
      setSearchParameters(newParameters);
    }}
    selectedQueryID={queryID}
    bind:defaultQuery
    bind:queries
  ></Queries>
  {#if !selectedQuery}
    <TypeToggle
      selected={type}
      eventButtonVisible={appStore.isEditor() ||
        appStore.isReviewer() ||
        appStore.isAdmin() ||
        appStore.isAuditor()}
      onSelect={(newType: SEARCHTYPES) => {
        appStore.setSearchType(newType);
        appStore.setSearchParametersForType(type, getCurrentSearchParameters());
        query.queryType = newType;
        const newParameters = $state.snapshot(appStore.getSearchParametersForType(newType));
        if (newType === SEARCHTYPES.ADVISORY) {
          query.orders = filterOrderCriteria(query.orders, SEARCHPAGECOLUMNS.ADVISORY);
        } else if (newType === SEARCHTYPES.DOCUMENT) {
          query.orders = filterOrderCriteria(query.orders, SEARCHPAGECOLUMNS.DOCUMENT);
        }
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
    buttonText={advanced ? "Apply" : "Search"}
    placeholder={advanced ? "Enter a query" : "Enter a search term"}
    search={(term) => {
      setSearchParameters({
        currentPage: 1,
        searchTerm: term
      });
    }}
    searchTerm={searchTermInputValue}
  ></CSearch>
  <div class="mt-1" title="Define finer grained search queries">
    <Toggle
      onclick={() => {
        setSearchParameters({ advanced: !$state.snapshot(advanced) });
      }}
      checked={advanced}
      class="ml-3">Advanced</Toggle
    >
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
    {columns}
    {documents}
    {error}
    {loading}
    {numberOfPages}
    {limit}
    {offset}
    {currentPage}
    {orderBy}
    dataChanged={fetchData}
    tableType={type}
    query={`${queryQuery}`}
    bind:count
    bind:openRow
    {last}
    {setSearchParameters}
  ></AdvisoryTable>
{/if}

{#if appStore.isEditor() || appStore.isReviewer()}
  <Toolbox></Toolbox>
{/if}
