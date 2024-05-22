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
  import { Button, ButtonGroup } from "flowbite-svelte";
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
</script>

{#if $appStore.app.isUserLoggedIn}
  {#if queries.length > 0}
    <ButtonGroup class="me-2">
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
    {@const query = queries[selectedIndex]}
    <AdvisoryTable
      columns={query.columns}
      loadAdvisories={query.advisories}
      query={query.query}
      orderBy={query.orders?.join(" ") ?? ""}
    ></AdvisoryTable>
  {/if}
  <ErrorMessage message={errorMessage}></ErrorMessage>
{/if}
