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
  import { Button, ButtonGroup, Input, Label } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import AdvisoryTable from "$lib/Advisories/AdvisoryTable.svelte";

  let defaultQueries = [
    {
      name: "New advisories",
      query: "$state new workflow =",
      advisories: true,
      columns: [
        "id",
        "publisher",
        "title",
        "tracking_id",
        "version",
        "cvss_v2_score",
        "cvss_v3_score"
      ]
    }
  ];
  let queries: any[] = [];
  let selectedIndex = 0;
  let pressedButtonClass = "bg-gray-200 hover:bg-gray-100";
  let errorMessage = "";
  let isAdvancedParametersEnabled = false;
  let advancedQuery = "";
  let appliedAdvancedQuery = ";";
  let isAdvancedQueryValid = true;
  let advancedQueryErrorMessage = "";

  onMount(async () => {
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queries = [...defaultQueries, ...response.content];
    } else if (response.error) {
      errorMessage = `Could not load user defined queries. ${getErrorMessage(response.error)}`;
    }
  });

  const selectQuery = (index: number) => {
    selectedIndex = index;
  };

  const toggleAdvancedParameters = () => {
    isAdvancedParametersEnabled = !isAdvancedParametersEnabled;
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
  {#if queries.length > 0}
    <div class="mb-4 flex gap-x-4">
      <ButtonGroup>
        <Button
          on:click={() => selectQuery(0)}
          class={`${0 === selectedIndex ? pressedButtonClass : ""} flex flex-col p-0`}
        >
          <span title="New advisories" class="m-2 h-full w-full">New advisories</span>
        </Button>
      </ButtonGroup>
      <ButtonGroup>
        {#each queries as query, index}
          {#if index > defaultQueries.length - 1}
            <Button
              on:click={() => selectQuery(index)}
              class={`${index + defaultQueries.length - 1 === selectedIndex ? pressedButtonClass : ""} flex flex-col p-0`}
            >
              <span title={query.description} class="m-2 h-full w-full">{query.name}</span>
            </Button>
          {/if}
        {/each}
      </ButtonGroup>
      <Button on:click={toggleAdvancedParameters} color="light">
        <span>Advanced</span>
        {#if isAdvancedParametersEnabled}
          <i class="bx bx-chevron-up text-xl"></i>
        {:else}
          <i class="bx bx-chevron-down text-xl"></i>
        {/if}
      </Button>
    </div>
    {#if isAdvancedParametersEnabled}
      <div class="flex items-end gap-x-2">
        <div>
          <Label for="advanced-parameters">Parameters:</Label>
          <div class="flex gap-x-2">
            <Input
              bind:value={advancedQuery}
              on:input={testAdvancedQuery}
              id="advanced-parameters"
              type="text"
            />
            <Button on:click={applyAdvancedQueries} disabled={!isAdvancedQueryValid}>Apply</Button>
          </div>
        </div>
      </div>
      {#if advancedQueryErrorMessage.length > 0}
        <ErrorMessage message={advancedQueryErrorMessage}></ErrorMessage>
      {/if}
    {/if}
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
