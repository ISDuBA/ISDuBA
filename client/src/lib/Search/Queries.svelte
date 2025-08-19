<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { Button, ButtonGroup } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "svelte-spa-router";
  import { createEventDispatcher } from "svelte";
  import { type Query } from "$lib/Queries/query";
  import { truncate } from "$lib/utils";

  const dispatch = createEventDispatcher();

  let queries: any[] = [];
  $: sortedQueries = queries.sort((a: any, b: any) => {
    if (a.global && !b.global) {
      return -1;
    } else if (!a.global && b.global) {
      return 1;
    }
    return 0;
  });
  let selectedIndex = -2;
  export let selectedQuery: boolean = false;
  export let defaultQuery: any;
  export let queryString: any;
  let ignoredQueries: Query[] = [];
  let errorMessage: ErrorDetails | null = null;
  let advancedQueryErrorMessage: ErrorDetails | null = null;
  const globalQueryButtonColor = "primary";
  const defaultQueryButtonClass = "flex flex-col p-0";
  const queryButtonClass = "bg-white hover:bg-gray-100";
  const pressedQueryButtonClass =
    "bg-gray-200 text-black hover:text-black hover:!bg-gray-100 dark:bg-gray-600 dark:hover:!bg-gray-700";
  const globalQueryButtonClass = `border-${globalQueryButtonColor}-500 dark:border-${globalQueryButtonColor}-500 dark:hover:border-${globalQueryButtonColor}-500 hover:!text-black dark:hover:!text-white`;
  const pressedGlobalQueryButtonClass = `border-${globalQueryButtonColor}-500 bg-${globalQueryButtonColor}-600 hover:bg-${globalQueryButtonColor}-700 hover:text-white dark:bg-${globalQueryButtonColor}-600 dark:hover:bg-${globalQueryButtonColor}-500 dark:border-${globalQueryButtonColor}-700 dark:hover:border-${globalQueryButtonColor}-700 text-white focus:text-white`;

  const getClass = (isGlobal: boolean, isPressed: boolean) => {
    const addition = isGlobal
      ? isPressed
        ? pressedGlobalQueryButtonClass
        : globalQueryButtonClass
      : isPressed
        ? pressedQueryButtonClass
        : queryButtonClass;
    return `${defaultQueryButtonClass} ${addition}`;
  };

  const fetchIgnored = async () => {
    const response = await request(`/api/queries/ignore`, "GET");
    if (response.ok) {
      ignoredQueries = response.content;
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load queries.`, response);
    }
  };

  onMount(async () => {
    fetchIgnored();
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queries = response.content.filter((q: Query) => !q.dashboard && !q.default_query);
      let defaultQueries = response.content.filter((q: Query) => q.default_query);
      if (defaultQueries) {
        defaultQuery = defaultQueries[0];
      }
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load user defined queries.`, response);
    }
    if (queryString?.query) {
      // Need to wait until sortedQueries is filled.
      setTimeout(() => {
        const index = sortedQueries.findIndex((q) => q.id === Number(queryString.query));

        // Probably a dashboard query
        if (index === -1) {
          const query = response.content.filter((q: any) => `${q.id}` === queryString.query)?.[0];
          selectedIndex = -2;
          currentQueryTitle = query.name;
          if (query) {
            dispatch("querySelected", query);
            selectedQuery = true;
          }
        } else {
          selectQuery(index);
        }
      }, 100);
    }
  });

  let currentQueryTitle: string | undefined;

  const selectQuery = (index: number) => {
    if (selectedIndex === index || index === -1) {
      selectedIndex = -1;
      currentQueryTitle = undefined;
      selectedQuery = false;
    } else {
      selectedIndex = index;
      dispatch("querySelected", sortedQueries[selectedIndex]);
      currentQueryTitle = sortedQueries[selectedIndex].name;
      selectedQuery = true;
    }
  };
</script>

<div class="mb-8 flex flex-col flex-wrap gap-4">
  <div class="flex items-center gap-x-4">
    <ButtonGroup class="flex-wrap">
      {#each sortedQueries as query, index}
        {#if !ignoredQueries.includes(query.id)}
          <Button
            size="xs"
            on:click={() => selectQuery(index)}
            class={getClass(query.global, index === selectedIndex)}
          >
            <span title={query.description} class="p-2">{truncate(query.name, 30)}</span>
          </Button>
        {/if}
      {/each}
      {#if currentQueryTitle && selectedIndex < 0}
        <Button size="xs" on:click={() => selectQuery(-1)} class={getClass(true, true)}>
          <span class="p-2">{truncate(currentQueryTitle, 30)}</span>
        </Button>
      {/if}
      <Button
        title="Configure queries"
        size="xs"
        on:click={() => {
          push("/queries");
        }}
      >
        <i class="bx bx-cog"></i>
      </Button>
    </ButtonGroup>
    {#if sortedQueries.length === 0}
      <span class="text-xs text-gray-400">No queries defined yet</span>
    {/if}
  </div>
  <ErrorMessage error={advancedQueryErrorMessage}></ErrorMessage>
</div>
<ErrorMessage error={errorMessage}></ErrorMessage>
