<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { querystring as qs } from "svelte-spa-router";
  import { parse } from "qs";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { Button, ButtonGroup } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "$routes/router.svelte";
  import { truncate } from "$lib/utils";
  import {
    defaultQuery as defaultQ,
    queryState,
    invisibleQueries,
    sortedQueries
  } from "./search.svelte";
  import type { Query, SEARCHTYPES } from "$lib/Queries/query";

  interface Props {
    onQuerySelected: (id: number | undefined) => void;
    selectedType: SEARCHTYPES;
  }

  let { onQuerySelected, selectedType }: Props = $props();

  const uid = $props.id();

  let queryString: any = $derived($qs ? parse($qs) : undefined);

  let defaultQuery = $derived(defaultQ());

  let queryID: number | undefined = $derived(
    queryString?.queryID ? Number(queryString.queryID) : undefined
  );
  let selectedQuery: Query | null = $derived.by(() => {
    const state = $state.snapshot(queryState);
    const queries = state.queries;
    if (queries === null) return null;
    const queryByID = $state.snapshot(queries).find((q) => q.id === queryID);
    if (queryByID) return queryByID;
    return defaultQuery ?? null;
  });

  let selectedInvisibleQuery = $derived(invisibleQueries().find((q) => q.id === queryID));
  let errorMessage: ErrorDetails | null = $state(null);
  let advancedQueryErrorMessage: ErrorDetails | null = null;
  const globalQueryButtonColor = "primary";
  const defaultQueryButtonClass = "flex flex-col py-0 text-xs";
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
      queryState.ignoredQueries = response.content;
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load queries.`, response);
    }
  };

  onMount(async () => {
    fetchIgnored();
    queryState.queries = [];
    queryState.ignoredQueries = [];
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queryState.queries = response.content;
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load user defined queries.`, response);
    }
  });

  const selectQuery = (id: number | undefined) => {
    onQuerySelected(id);
  };
</script>

<div class="flex flex-col flex-wrap gap-4">
  <div class="flex items-center gap-x-4">
    <ButtonGroup class="h-7 flex-wrap">
      {#each sortedQueries() as query, index (`queries-${uid}-${index}`)}
        <Button
          color="light"
          onclick={() => selectQuery(query.id === queryID ? undefined : query.id)}
          class={getClass(
            query.global,
            query.id === selectedQuery?.id && selectedType === query.kind
          )}
        >
          <span title={query.description}>{truncate(query.name, 30)}</span>
        </Button>
      {/each}
      {#if selectedInvisibleQuery}
        <Button color="light" onclick={() => selectQuery(undefined)} class={getClass(true, true)}>
          <span>{truncate(selectedInvisibleQuery.name, 30)}</span>
        </Button>
      {/if}
      <Button
        class="py-1"
        color="light"
        title="Configure queries"
        size="xs"
        onclick={() => {
          push("/queries");
        }}
      >
        <i class="bx bx-cog"></i>
      </Button>
    </ButtonGroup>
    {#if sortedQueries().length === 0}
      <span class="text-xs text-gray-400">No queries defined yet</span>
    {/if}
  </div>
  <ErrorMessage error={advancedQueryErrorMessage}></ErrorMessage>
</div>
<ErrorMessage error={errorMessage}></ErrorMessage>
