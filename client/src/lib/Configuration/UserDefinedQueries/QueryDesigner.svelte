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
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { onMount } from "svelte";
  import { push, querystring } from "svelte-spa-router";
  // @ts-expect-error ignore complaining qs has not type declaration
  import { parse } from "qs";
  import { appStore } from "$lib/store";
  import { ADMIN, isRoleIncluded } from "$lib/permissions";

  export let params: any = null;
  let editName = false;
  let editDescription = false;
  let hoveredLine = -1;
  let queryCount: any = null;
  let loading = false;
  let errorMessage = "";
  let saveErrorMessage = "";
  let loadQueryError = "";
  let loadedData: any = null;
  let abortController: AbortController;
  let placeholder = "";

  const unsetMessages = () => {
    queryCount = null;
    errorMessage = "";
    saveErrorMessage = "";
  };

  const testQuery = async () => {
    loading = true;
    unsetMessages();
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

  const saveQuery = async () => {
    unsetMessages();
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
      push(`/configuration/`);
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
      currentSearch.columns = currentSearch.columns.filter((c) => {
        if (c.name !== "ssvc" && c.name !== "state") return c;
      });
    }
    if (currentSearch.searchType === SEARCHTYPES.ADVISORY) {
      const newCols = columnsFromNames(["ssvc", "state"]);
      currentSearch.columns = [...currentSearch.columns, ...newCols];
    }
  };

  const promoteColumn = (index: number) => {
    if (index === 0) return;
    let tmp = currentSearch.columns[index - 1];
    currentSearch.columns[index - 1] = currentSearch.columns[index];
    currentSearch.columns[index] = tmp;
  };

  const demoteColumn = (index: number) => {
    if (index === currentSearch.columns.length - 1) return;
    let tmp = currentSearch.columns[index + 1];
    currentSearch.columns[index + 1] = currentSearch.columns[index];
    currentSearch.columns[index] = tmp;
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

  onMount(async () => {
    const queryString = parse($querystring);
    let id;
    if (queryString.clone) {
      id = queryString.clone;
    }
    if (params) id = params.id;
    if (id) {
      const response = await request(`/api/queries/${id}`, "GET");
      if (response.ok) {
        const result = await response.content;
        if (params && params.id) {
          loadedData = result;
        }
        currentSearch = generateQueryFrom(result);
        if (queryString.clone) {
          currentSearch.name = ``;
          editName = true;
          placeholder = "Please enter a name";
        }
      } else if (response.error) {
        loadQueryError = `Could not load query. ${getErrorMessage(response.error)}`;
      }
    }
  });

  $: disableSave =
    currentSearch.columns.every((c) => c.visible == false) || currentSearch.name == "";
</script>

<SectionHeader title="Configuration"></SectionHeader>
<hr class="mb-6" />
<h2 class="mb-6 text-lg">User defined queries</h2>

{#if loadQueryError === ""}
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
              <h5 class="text-xl font-medium text-gray-500 dark:text-gray-400">
                {shorten(currentSearch.name)}
              </h5>
              <i class="bx bx-edit-alt ml-1"></i>
            </div>
          {/if}
        </button>
      </div>
      <div class="ml-6 flex w-1/3 flex-row items-center gap-x-2">
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
              <h5 class="text-xl font-medium text-gray-500 dark:text-gray-400">
                {shorten(currentSearch.description)}
              </h5>
              <i class="bx bx-edit-alt ml-1"></i>
            </div>
          {/if}
        </button>
      </div>
      <div class="w-1/3">
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
    <hr class="mb-4 mt-4 w-4/5" />
    <div class="flex w-1/2 flex-row">
      <div class="w-1/3">
        <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Searching</h5>
      </div>
      <div class="ml-6 w-1/3">
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
        <div class="ml-6 w-1/3">Column</div>
        <div class="flex w-1/3 flex-row">
          <div>Visible</div>
        </div>
        <div>Order direction</div>
      </div>
      {#each currentSearch.columns as col, index (index)}
        <div
          role="presentation"
          class="mb-1 flex cursor-pointer flex-row items-center"
          on:mouseover={() => {
            hoveredLine = index;
          }}
          on:mouseout={() => {
            hoveredLine = -1;
          }}
          on:blur={() => {}}
          on:focus={() => {}}
        >
          <div
            class:w-6={true}
            class:flex={true}
            class:flex-col={true}
            class:invisible={hoveredLine !== index}
          >
            <button
              class="h-4"
              on:click={() => {
                promoteColumn(index);
              }}
            >
              <i class="bx bxs-up-arrow-circle"></i>
            </button>
            <button
              on:click={() => {
                demoteColumn(index);
              }}
              class="h-4"
            >
              <i class="bx bxs-down-arrow-circle"></i>
            </button>
          </div>
          <div class="w-1/3">{col.name}</div>
          <div class="w-1/3">
            <Checkbox
              on:change={() => {
                setVisible(index);
              }}
              checked={currentSearch.columns[index].visible}
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
              {#if col.orderBy === ""}
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
