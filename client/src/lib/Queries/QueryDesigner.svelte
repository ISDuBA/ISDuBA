<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import {
    Input,
    Spinner,
    Button,
    Checkbox,
    Img,
    RadioButton,
    ButtonGroup,
    Label,
    Select
  } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import {
    COLUMNS,
    ORDERDIRECTIONS,
    SEARCHTYPES,
    generateQueryString,
    type Search,
    type Column
  } from "$lib/Queries/query";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { onMount } from "svelte";
  import { push, querystring } from "svelte-spa-router";
  import { parse } from "qs";
  import { appStore } from "$lib/store";
  import { ADMIN, AUDITOR, EDITOR, IMPORTER, REVIEWER, SOURCE_MANAGER } from "$lib/workflow";
  import { isRoleIncluded } from "$lib/permissions";
  import Sortable from "sortablejs";

  export let params: any = null;
  let editName = false;
  let editDescription = false;
  let queryCount: any = null;
  let loading = false;
  let errorMessage: ErrorDetails | null;
  let saveErrorMessage: ErrorDetails | null;
  let loadQueryError: ErrorDetails | null;
  let loadedData: any = null;
  let abortController: AbortController;
  let placeholder = "";

  let columnList: any;

  // Prop items of MultiSelect doesn't accept simple strings
  const roles = [EDITOR, REVIEWER, AUDITOR, IMPORTER, SOURCE_MANAGER, ADMIN].map((r) => {
    return {
      name: r,
      value: r
    };
  });

  const unsetMessages = () => {
    queryCount = null;
    errorMessage = null;
    saveErrorMessage = null;
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
        errorMessage = getErrorDetails(`An error occured.`, response);
      }
    }
    loading = false;
  };

  const columnsFromNames = (columns: string[]): Column[] => {
    const result = [...columns];
    return result.map((name) => {
      return {
        name: name,
        visible: false
      };
    });
  };

  const newQuery = (): Search => {
    const columns = columnsFromNames(COLUMNS.ADVISORY);
    return {
      searchType: SEARCHTYPES.ADVISORY,
      columns: columns,
      orderBy: [],
      name: "New Query",
      query: "",
      description: "",
      global: false,
      dashboard: false,
      roles: []
    };
  };

  let currentSearch = newQuery();

  const sortColumns = (columns: Column[]): Column[] => {
    let nodes = columnList.querySelectorAll(".columnName");
    let sortedColumns: Column[] = [];
    for (const node of nodes) {
      let columnName = node.innerText;
      let found = columns.find((c) => c.name === columnName);
      if (found) {
        sortedColumns.push(found);
      }
    }
    return sortedColumns;
  };

  const saveQuery = async () => {
    unsetMessages();
    const formData = new FormData();
    formData.append("kind", currentSearch.searchType);
    formData.append("name", currentSearch.name);
    formData.append("global", `${currentSearch.global}`);
    formData.append("dashboard", `${currentSearch.dashboard}`);
    formData.append("role", `${currentSearch.roles}`);
    if (currentSearch.description.length > 0) {
      formData.append("description", currentSearch.description);
    }
    if (currentSearch.query.length > 0) {
      formData.append("query", currentSearch.query);
    }
    let sortedColumns = sortColumns(currentSearch.columns);
    const columns = sortedColumns.filter((c) => c.visible).map((c) => c.name);
    formData.append("columns", columns.join(" "));
    const columnsForOrder = currentSearch.orderBy.filter((order) => order[0] !== "");
    const orderBy = columnsForOrder.map(
      (c) => `${c[1] === ORDERDIRECTIONS.DESC ? "-" : ""}${c[0]}`
    );
    formData.append("orders", orderBy.join(" "));
    let response;
    if (loadedData) {
      response = await request(`/api/queries/${loadedData.id}`, "PUT", formData);
    } else {
      response = await request("/api/queries", "POST", formData);
    }
    if (!response.ok && response.error) {
      saveErrorMessage = getErrorDetails(`Failed to save query.`, response);
      if (response.error === "409")
        saveErrorMessage = getErrorDetails(
          `A query with the name "${currentSearch.name}" already exists.`
        );
    }
    if (response.ok) {
      push(`/queries/`);
    }
  };

  const switchOrderDirection = (name: string) => {
    let order = currentSearch.orderBy.find((o) => o[0] === name);
    if (order) {
      if (order[1] === ORDERDIRECTIONS.ASC) {
        order[1] = ORDERDIRECTIONS.DESC;
      } else {
        currentSearch.orderBy = currentSearch.orderBy.filter((o) => o[0] !== name);
      }
    } else {
      currentSearch.orderBy.push([name, ORDERDIRECTIONS.ASC]);
    }
    currentSearch.columns = currentSearch.columns;
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
    if (currentSearch.searchType === SEARCHTYPES.EVENT) {
      currentSearch.columns = columnsFromNames(COLUMNS.EVENT);
    }
  };

  const shorten = (text: string) => {
    if (!text) return "";
    if (text.length < 20) return text;
    return `${text.substring(0, 20)}...`;
  };

  const generateQueryFrom = (result: any): Search => {
    let searchType = SEARCHTYPES.DOCUMENT;
    let columns = [];
    if (result.kind === SEARCHTYPES.ADVISORY) {
      searchType = SEARCHTYPES.ADVISORY;
      columns = COLUMNS.ADVISORY;
    } else if (result.kind === SEARCHTYPES.DOCUMENT) {
      searchType = SEARCHTYPES.DOCUMENT;
      columns = COLUMNS.DOCUMENT;
    } else {
      searchType = SEARCHTYPES.EVENT;
      columns = COLUMNS.EVENT;
    }
    columns = result.columns.concat(
      columns.filter((c: string) => {
        if (!result.columns.includes(c)) return c;
      })
    );
    columns = columnsFromNames(columns);
    columns = columns.map((c) => {
      if (result.columns.includes(c.name)) c.visible = true;
      return c;
    });

    let orderBy: [string, ORDERDIRECTIONS][] = [];
    if (result.orders) {
      for (let order of result.orders) {
        let direction = ORDERDIRECTIONS.ASC;
        if (order.startsWith("-")) {
          direction = ORDERDIRECTIONS.DESC;
          order = order.substring(1);
        }
        orderBy.push([order, direction]);
      }
    }

    return {
      searchType: searchType,
      columns: columns,
      orderBy: orderBy,
      name: result.name,
      query: result.query,
      description: result.description || "",
      global: result.global,
      dashboard: result.dashboard,
      roles: result.roles
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
    let queryString;
    if ($querystring) {
      queryString = parse($querystring);
    }
    let id;
    if (queryString?.clone) {
      id = queryString?.clone;
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
        if (queryString?.clone) {
          currentSearch.name = proposeName(result, currentSearch.name);
          if (!isRoleIncluded(appStore.getRoles(), [ADMIN])) {
            currentSearch.global = false;
          }
        }
      } else if (response.error) {
        loadQueryError = getErrorDetails(`Could not load query.`, response);
      }
    }
  });

  const getOrderDirection = (name: string): [number, ORDERDIRECTIONS] | undefined => {
    let index = currentSearch.orderBy.findIndex((o) => o[0] === name);
    if (index >= 0) {
      return [index, currentSearch.orderBy[index][1]];
    }
    return undefined;
  };

  $: if (columnList) {
    Sortable.create(columnList, {
      animation: 150
    });
  }

  $: noColumnSelected = currentSearch.columns.every((c) => c.visible == false);
  $: disableSave = noColumnSelected || currentSearch.name == "";
</script>

<SectionHeader title="Queries"></SectionHeader>
<hr class="mb-6" />

{#if loadQueryError !== null}
  <div class="w-3/4">
    <div class="flex flex-col">
      <div class="flex flex-row">
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
      </div>
      <div class="mb-4">
        <small class={currentSearch.name === "" ? "text-red-500" : "text-gray-400"}>Required</small>
      </div>
    </div>
    <div class="mb-4 flex gap-4">
      {#if isRoleIncluded(appStore.getRoles(), [ADMIN])}
        <div class="flex flex-row items-center gap-x-2">
          <span>Global:</span>
          <Checkbox
            checked={currentSearch.global}
            on:change={() => {
              currentSearch.global = !currentSearch.global;
            }}
          ></Checkbox>
        </div>
      {/if}
      <div class="flex flex-row items-center gap-x-2">
        <span>Show on dashboard:</span>
        <Checkbox
          checked={currentSearch.dashboard}
          on:change={() => {
            currentSearch.dashboard = !currentSearch.dashboard;
          }}
        ></Checkbox>
      </div>
    </div>
    <div class="mb-6">
      {#if isRoleIncluded(appStore.getRoles(), [ADMIN])}
        <Label class="mb-1" for="roles">Roles:</Label>
        <Select id="roles" items={roles} bind:value={currentSearch.roles}></Select>
      {/if}
    </div>
    <div class="flex flex-row">
      <div class="flex w-1/3 min-w-40 -flex-row items-baseline gap-x-3">
        <h5 class="text-lg font-medium text-gray-500 dark:text-gray-400">Searching</h5>
        <small class:text-red-500={noColumnSelected} class:text-gray-400={!noColumnSelected}
          >Select at least 1 column</small
        >
      </div>
      <ButtonGroup class="ml-6">
        <RadioButton
          class="h-8"
          on:change={toggleSearchType}
          value={SEARCHTYPES.ADVISORY}
          bind:group={currentSearch.searchType}
        >
          Advisories</RadioButton
        >
        <RadioButton
          class="h-8"
          on:change={toggleSearchType}
          value={SEARCHTYPES.DOCUMENT}
          bind:group={currentSearch.searchType}>Documents</RadioButton
        >
        <RadioButton
          class="h-8"
          on:change={toggleSearchType}
          value={SEARCHTYPES.EVENT}
          bind:group={currentSearch.searchType}>Events</RadioButton
        >
      </ButtonGroup>
    </div>
    <div class="mt-4">
      <div class="mb-2 flex flex-row">
        <div class="ml-6 w-1/3 min-w-40">Column</div>
        <div class="w-1/4 min-w-28">Visible</div>
        <div class="text-nowrap">Order</div>
      </div>
      <section bind:this={columnList}>
        {#each currentSearch.columns as col, index (index)}
          {@const order = getOrderDirection(currentSearch.columns[index].name)}
          <div
            role="presentation"
            class="mb-1 flex cursor-pointer flex-row items-center"
            on:blur={() => {}}
            on:focus={() => {}}
          >
            <div class:w-6={true} class:flex={true} class:flex-col={true}>
              <button class="h-4">
                <Img src="grid-dots-vertical-rounded.svg" class="h-4 min-h-2 min-w-2 invert-[.5]" />
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
                switchOrderDirection(col.name);
              }}
            >
              {#if order}
                {order[0] + 1}
                {#if order[1] === ORDERDIRECTIONS.ASC}
                  <i class="bx bx-sort-a-z"></i>
                {/if}
                {#if order[1] === ORDERDIRECTIONS.DESC}
                  <i class="bx bx-sort-z-a"></i>
                {/if}
              {:else}
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
      {#if saveErrorMessage}
        <ErrorMessage error={saveErrorMessage}></ErrorMessage>
      {/if}
    </div>
  </div>
{:else}
  <ErrorMessage error={loadQueryError}></ErrorMessage>
{/if}
