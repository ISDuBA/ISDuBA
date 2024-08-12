<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Radio, Input, Spinner, Button, Checkbox, Img } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { COLUMNS, ORDERDIRECTIONS, SEARCHTYPES, generateQueryString } from "$lib/Queries/query";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { onMount } from "svelte";
  import { push, querystring } from "svelte-spa-router";
  // @ts-expect-error ignore complaining qs has not type declaration
  import { parse } from "qs";
  import { appStore } from "$lib/store";
  import { ADMIN } from "$lib/workflow";
  import { isRoleIncluded } from "$lib/permissions";
  import Sortable from "sortablejs";

  export let params: any = null;
  let editName = false;
  let editDescription = false;
  let queryCount: any = null;
  let loading = false;
  let errorMessage = "";
  let saveErrorMessage = "";
  let loadQueryError = "";
  let loadedData: any = null;
  let abortController: AbortController;
  let placeholder = "";

  let columnList: any;

  const unsetMessages = () => {
    queryCount = null;
    errorMessage = "";
    saveErrorMessage = "";
  };

  const testQuery = async () => {
    loading = true;
    unsetMessages();
    sortColumns();
    abortController = new AbortController();
    const query = generateQueryString(currentSearch);
    const response = await request(query, "GET", undefined, abortController);
    if (response.ok) {
      queryCount = response.content.count;
    } else if (response.error) {
      if (/Error/.test(response.error)) {
        // Intentionally ignore errors induced by aborting the request
      } else {
        errorMessage = `An error occured: ${response.content}`;
      }
    }
    loading = false;
  };

  const columnsFromNames = (columns: string[]) => {
    const result = [...columns];
    return result.map((name) => {
      return {
        name: name,
        visible: false,
        orderBy: ""
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

  let currentSearch = newQuery();

  const sortColumns = () => {
    let oldColumns = [...currentSearch.columns];
    currentSearch.columns = [];
    let columns = columnList.querySelectorAll(".columnName");
    for (const column of columns) {
      let columnName = column.innerText;
      let newColumn = oldColumns.find((c) => c.name === columnName);
      currentSearch.columns.push(newColumn);
    }
  };

  const saveQuery = async () => {
    unsetMessages();
    sortColumns();
    const formData = new FormData();
    formData.append("advisories", `${currentSearch.searchType === SEARCHTYPES.ADVISORY}`);
    formData.append("name", currentSearch.name);
    formData.append("global", `${currentSearch.global}`);
    if (currentSearch.description.length > 0) {
      formData.append("description", currentSearch.description);
    }
    if (currentSearch.query.length > 0) {
      formData.append("query", currentSearch.query);
    }
    const columns = currentSearch.columns.filter((c) => c.visible).map((c) => c.name);
    formData.append("columns", columns.join(" "));
    const columnsForOrder = currentSearch.columns.filter((c) => c.orderBy);
    const orderBy = columnsForOrder.map((c) => `${c.orderBy === "desc" ? "-" : ""}${c.name}`);
    formData.append("orders", orderBy.join(" "));
    let response;
    if (loadedData) {
      response = await request(`/api/queries/${loadedData.id}`, "PUT", formData);
    } else {
      response = await request("/api/queries", "POST", formData);
    }
    if (!response.ok && response.error) {
      saveErrorMessage = `${getErrorMessage(response.error)} Reason: "${response.content}"`;
      if (response.error === "409")
        saveErrorMessage = `A query with the name "${currentSearch.name}" already exists.`;
    }
    if (response.ok) {
      push(`/queries/`);
    }
  };

  const switchOrderDirection = (index: number) => {
    if (currentSearch.columns[index].orderBy == "") {
      currentSearch.columns[index].orderBy = ORDERDIRECTIONS.ASC;
      return;
    }
    if (currentSearch.columns[index].orderBy === ORDERDIRECTIONS.ASC) {
      currentSearch.columns[index].orderBy = ORDERDIRECTIONS.DESC;
    } else {
      currentSearch.columns[index].orderBy = "";
    }
  };

  const setVisible = (index: number) => {
    currentSearch.columns[index].visible = !currentSearch.columns[index].visible;
  };

  const toggleSearchType = () => {
    if (currentSearch.searchType === SEARCHTYPES.DOCUMENT) {
      currentSearch.columns = columnsFromNames(COLUMNS.DOCUMENT);
    }
    if (currentSearch.searchType === SEARCHTYPES.ADVISORY) {
      currentSearch.columns = columnsFromNames(COLUMNS.ADVISORY);
    }
  };

  const shorten = (text: string) => {
    if (!text) return "";
    if (text.length < 20) return text;
    return `${text.substring(0, 20)}...`;
  };

  const generateQueryFrom = (result: any) => {
    let searchType = "";
    let columns = [];
    if (result.advisories) {
      searchType = SEARCHTYPES.ADVISORY;
      columns = COLUMNS.ADVISORY;
    } else {
      searchType = SEARCHTYPES.DOCUMENT;
      columns = COLUMNS.DOCUMENT;
    }
    columns = result.columns.concat(
      columns.filter((c: string) => {
        if (!result.columns.includes(c)) return c;
      })
    );
    columns = columnsFromNames(columns);
    columns = columns.map((c) => {
      if (result.columns.includes(c.name)) c.visible = true;
      if (result.orders?.includes(c.name)) c.orderBy = ORDERDIRECTIONS.ASC;
      if (result.orders?.includes(`-${c.name}`)) c.orderBy = ORDERDIRECTIONS.DESC;
      return c;
    });
    return {
      searchType: searchType,
      columns: columns,
      name: result.name,
      query: result.query,
      description: result.description || "",
      global: result.global
    };
  };

  /**
   * Takes the list of existing queries, looks for already given clones and returns a proper name.
   * Expamples:
   *
   * For non existing clones
   *
   * Monat -> Monat (1)
   * Monat (1) -> Monat (1) (1)
   *
   * Say there is already a clone
   *
   * Monat and Monat (1) -> Monat (2)
   * Monat (1) and Monat (1) (1) -> Monat (1) (2)
   * Monat (1) (2) and Monat (1) (1) -> Monat (1) (3)
   *
   * And so on.
   *
   * @param result list of queries
   * @param name name of the query
   */
  const proposeName = (result: any, name: string) => {
    const clones = result
      .filter((r: any) => {
        const re = new RegExp(name.replaceAll("(", "\\(").replaceAll(")", "\\)") + " \\(\\d+\\)");
        return re.test(r.name);
      })
      .map((r: any) => {
        return r.name;
      })
      .sort((a: string, b: string) => a.localeCompare(b, "en", { numeric: true }));
    if (clones.length === 0) return `${name} (1)`;
    const highestIndex = parseInt(clones[clones.length - 1].split(name + " (")[1]);
    return `${name} (${highestIndex + 1})`;
  };

  onMount(async () => {
    const queryString = parse($querystring);
    let id;
    if (queryString.clone) {
      id = queryString.clone;
    }
    if (params) id = params.id;
    if (id) {
      const response = await request(`/api/queries/`, "GET");
      if (response.ok) {
        const result = await response.content;
        const thisQuery = result.find((q: any) => {
          return q.id == id;
        });
        if (params && params.id) {
          loadedData = thisQuery;
        }
        currentSearch = generateQueryFrom(thisQuery);
        if (queryString.clone) {
          currentSearch.name = proposeName(result, currentSearch.name);
          if (!isRoleIncluded(appStore.getRoles(), [ADMIN])) {
            currentSearch.global = false;
          }
        }
      } else if (response.error) {
        loadQueryError = `Could not load query. ${getErrorMessage(response.error)}`;
      }
    }
  });

  $: if (columnList) {
    Sortable.create(columnList, {
      handle: "#handle",
      animation: 150
    });
  }

  $: noColumnSelected = currentSearch.columns.every((c) => c.visible == false);
  $: disableSave = noColumnSelected || currentSearch.name == "";
</script>

<SectionHeader title="Queries"></SectionHeader>
<hr class="mb-6" />

{#if loadQueryError === ""}
  <div class="w-3/4">
    <div class="flex h-1 flex-row">
      <div class="flex w-1/3 min-w-40 flex-row items-center gap-x-2">
        <span class={currentSearch.name === "" ? "text-red-500" : ""}>Name:</span>
        <button
          on:click={() => {
            editName = !editName;
          }}
        >
          {#if editName}
            <Input
              autofocus
              {placeholder}
              bind:value={currentSearch.name}
              on:keyup={(e) => {
                if (e.key === "Enter") editName = false;
                if (e.key === "Escape") editName = false;
                e.preventDefault();
              }}
              on:blur={() => {
                editName = false;
              }}
              on:click={(e) => e.stopPropagation()}
            />
          {:else}
            <div class="flex flex-row items-center" title={currentSearch.name}>
              <h5 class="font-medium text-gray-500 dark:text-gray-400">
                {shorten(currentSearch.name)}
              </h5>
              <i class="bx bx-edit-alt ml-1"></i>
            </div>
          {/if}
        </button>
      </div>
      <div class="ml-6 flex w-1/3 min-w-96 flex-row items-center gap-x-2">
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
              on:click={(e) => e.stopPropagation()}
            />
          {:else}
            <div class="flex flex-row items-center" title={currentSearch.description}>
              <h5 class="font-medium text-gray-500 dark:text-gray-400">
                {shorten(currentSearch.description)}
              </h5>
              <i class="bx bx-edit-alt ml-1"></i>
            </div>
          {/if}
        </button>
      </div>
      <div>
        {#if isRoleIncluded(appStore.getRoles(), [ADMIN])}
          <div class="flex h-1 flex-row items-center gap-x-3">
            <span>Global:</span>
            <Checkbox
              checked={currentSearch.global}
              on:change={() => {
                currentSearch.global = !currentSearch.global;
              }}
            ></Checkbox>
          </div>
        {/if}
      </div>
    </div>
    <div class="my-2">
      <small class={currentSearch.name === "" ? "text-red-500" : "text-gray-400"}>Required</small>
    </div>
    <hr class="mb-4 w-4/5 min-w-96" />
    <div class="flex flex-row">
      <div class="flex w-1/3 min-w-40 -flex-row items-baseline gap-x-3">
        <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Searching</h5>
        <small class:text-red-500={noColumnSelected} class:text-gray-400={!noColumnSelected}
          >Select at least 1 column</small
        >
      </div>
      <div class="ml-6 w-1/4 min-w-28">
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
    <div class="mt-4">
      <div class="mb-2 flex flex-row">
        <div class="ml-6 w-1/3 min-w-40">Column</div>
        <div class="w-1/4 min-w-28">Visible</div>
        <div class="text-nowrap">Orderdirection</div>
      </div>
      <section bind:this={columnList}>
        {#each currentSearch.columns as col, index (index)}
          <div
            role="presentation"
            class="mb-1 flex cursor-pointer flex-row items-center"
            on:blur={() => {}}
            on:focus={() => {}}
          >
            <div class:w-6={true} class:flex={true} class:flex-col={true}>
              <button class="h-4">
                <Img
                  id="handle"
                  src="grid-dots-vertical-rounded.svg"
                  class="h-4 min-h-2 min-w-2 invert-[.5]"
                />
              </button>
            </div>
            <div class="columnName w-1/3 min-w-40">{col.name}</div>
            <div class="w-1/4 min-w-28">
              <Checkbox
                on:change={() => {
                  setVisible(index);
                }}
                checked={currentSearch.columns[index].visible}
              ></Checkbox>
            </div>
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
              {#if col.orderBy === ""}
                <i class="bx bx-minus"></i>
              {/if}
            </button>
          </div>
        {/each}
      </section>
    </div>
    <div class="mt-6 w-full min-w-96">
      <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Query criteria</h5>
      <div class="flex flex-row">
        <Input bind:value={currentSearch.query} />
      </div>
      <div class="mt-3 flex flex-row">
        {#if loading}
          <div class="loadingFadeIn mr-4 mt-3">
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
          {#if !loading}
            <Button on:click={testQuery} color="light"
              ><i class="bx bx-test-tube me-2"></i> Test query</Button
            >
          {/if}
          {#if loading}
            <Button
              on:click={() => {
                if (abortController) abortController.abort();
                loading = false;
                unsetMessages();
              }}
              color="light"><i class="bx bx-stop-circle"></i> Abort query</Button
            >
          {/if}
          <Button
            on:click={() => {
              if (loadedData) {
                currentSearch = generateQueryFrom(loadedData);
              } else {
                currentSearch = newQuery();
              }
              queryCount = null;
            }}
            color="light"><i class="bx bx-undo me-2 text-xl"></i> Reset</Button
          >
          <Button disabled={disableSave} on:click={saveQuery} color="light"
            ><i class="bx bxs-save me-2"></i> Save</Button
          >
        </div>
      </div>
      {#if saveErrorMessage.length > 0}
        <ErrorMessage message={saveErrorMessage}></ErrorMessage>
      {/if}
    </div>
  </div>
{:else}
  <ErrorMessage message={loadQueryError}></ErrorMessage>
{/if}
