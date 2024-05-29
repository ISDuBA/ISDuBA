<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";
  import { request } from "$lib/utils";
  import { Button, ButtonGroup, Input } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import AdvisoryTable from "$lib/Documents/Table.svelte";
  import { push } from "svelte-spa-router";

  let queries: any[] = [];
  $: sortedQueries = queries.sort((a: any, b: any) => {
    if (a.global && !b.global) {
      return -1;
    } else if (!a.global && b.global) {
      return 1;
    }
    return 0;
  });
  let selectedIndex = 0;
  let errorMessage = "";
  let advancedQuery = "";
  let appliedAdvancedQuery = "";
  let isAdvancedQueryValid = true;
  let advancedQueryErrorMessage = "";
  let globalQueryButtonColor = "primary";
  let defaultQueryButtonClass = "flex flex-col p-0 focus:text-black hover:text-black";
  let queryButtonClass = "bg-white hover:bg-gray-100";
  let pressedQueryButtonClass = "bg-gray-200 text-black hover:!bg-gray-100";
  let globalQueryButtonClass = `border-${globalQueryButtonColor}-500 hover:!text-black`;
  let pressedGlobalQueryButtonClass = `border-${globalQueryButtonColor}-500 bg-${globalQueryButtonColor}-600 focus:text-white text-white hover:text-black`;

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

  onMount(async () => {
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queries = response.content;
    } else if (response.error) {
      errorMessage = `Could not load user defined queries. ${getErrorMessage(response.error)}`;
    }
  });

  const selectQuery = (index: number) => {
    selectedIndex = index;
  };

  const testAdvancedQuery = async () => {
    advancedQueryErrorMessage = "";
    const selectedQuery = queries[selectedIndex];
    const query = `${selectedQuery.query} ${advancedQuery.length > 0 ? advancedQuery.concat(" and") : ""}`;
    const documentURL = encodeURI(
      `/api/documents?query=${query}&advisories=${selectedQuery.advisories}&count=1&order=${selectedQuery.orders?.join(" ") ?? ""}&limit=0&columns=${selectedQuery.columns.join(" ")}`
    );
    const result = await request(documentURL, "GET");
    isAdvancedQueryValid = result.ok;
    if (!result.ok) {
      advancedQueryErrorMessage = result.content;
    }
  };

  const applyAdvancedQueries = () => {
    appliedAdvancedQuery = advancedQuery;
  };
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="mb-8 flex flex-col flex-wrap gap-4">
    <div class="flex gap-x-4">
      <ButtonGroup class="flex-wrap">
        {#each sortedQueries as query, index}
          <Button
            size="xs"
            on:click={() => selectQuery(index)}
            class={getClass(query.global, index === selectedIndex)}
          >
            <span title={query.description} class="p-2">{query.name}</span>
          </Button>
        {/each}
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
    </div>
    {#if queries.length > 0}
      <div class="flex flex-row flex-wrap items-center gap-2 leading-3">
        <div class="mt-3 flex flex-wrap gap-x-2">
          <Input
            class="h-8 w-96"
            size="sm"
            bind:value={advancedQuery}
            on:input={testAdvancedQuery}
            id="advanced-parameters"
            type="text"
            placeholder="Advanced search"
          />
          <Button
            size="xs"
            color="light"
            on:click={applyAdvancedQueries}
            disabled={!isAdvancedQueryValid}>Apply</Button
          >
        </div>
      </div>
    {/if}
    <ErrorMessage message={advancedQueryErrorMessage}></ErrorMessage>
  </div>
  {#if queries.length > 0}
    {@const query = queries[selectedIndex]}
    <AdvisoryTable
      columns={query.columns}
      loadAdvisories={query.advisories}
      query={`${query.query} ${appliedAdvancedQuery.length > 0 ? appliedAdvancedQuery.concat(" and") : ""}`}
      orderBy={query.orders?.join(" ") ?? ""}
    ></AdvisoryTable>
  {/if}
  <ErrorMessage message={errorMessage}></ErrorMessage>
{/if}
