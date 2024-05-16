<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { onMount } from "svelte";
  import { request } from "$lib/utils";
  import { Button, ButtonGroup } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import AdvisoryTable from "$lib/table/AdvisoryTable.svelte";

  let queries: any[] = [];
  let selectedIndex = 0;
  let pressedButtonClass = "bg-gray-200 hover:bg-gray-100";
  let errorMessage = "";

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
</script>

{#if $appStore.app.isUserLoggedIn}
  {#if queries.length > 0}
    <SectionHeader title="Queries"></SectionHeader>
    <ButtonGroup>
      {#each queries as query, index}
        <Button
          on:click={() => selectQuery(index)}
          class={`${index === selectedIndex ? pressedButtonClass : ""} flex flex-col p-0`}
        >
          <span title={query.description} class="m-2 h-full w-full">{query.name}</span>
        </Button>
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
