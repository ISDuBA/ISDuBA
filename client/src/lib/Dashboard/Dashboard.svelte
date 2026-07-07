<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { onMount } from "svelte";
  import AdvisoryQuery from "./AdvisoryQuery.svelte";
  import EventQuery from "./EventQuery.svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import ImportStats from "$lib/Statistics/ImportStats.svelte";
  import SourceEvents from "./SourceEvents.svelte";

  const uid = $props.id();

  let filteredQueries: any[] = $state([]);
  let loadIgnoredError: ErrorDetails | null = $state(null);
  let loadQueryError: ErrorDetails | null = $state(null);

  const fetchStoredQueries = async (): Promise<any[]> => {
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      const result = response.content;
      return result.sort((q1: any, q2: any) => {
        return q1.num > q2.num;
      });
    } else if (response.error) {
      loadQueryError = getErrorDetails(`Could not load queries.`, response);
    }
    return [];
  };

  const fetchIgnored = async () => {
    loadIgnoredError = null;
    const response = await request(`/api/queries/ignore`, "GET");
    if (response.ok) {
      return response.content;
    } else if (response.error) {
      loadIgnoredError = getErrorDetails(`Could not load queries.`, response);
      return undefined;
    }
  };

  onMount(async () => {
    const allQueries = await fetchStoredQueries();
    const ignoredQueries = await fetchIgnored();
    const userDashboardQueries = allQueries.filter(
      (query) =>
        query.dashboard &&
        query.definer === appStore.state.app.tokenParsed?.preferred_username &&
        !query.global &&
        (!ignoredQueries || !ignoredQueries.includes(query.id))
    );
    const globalDashboardQueries = allQueries.filter(
      (query) =>
        query.dashboard &&
        query.global &&
        query.definer === "system-default" &&
        (!ignoredQueries || !ignoredQueries.includes(query.id))
    );
    filteredQueries = [...userDashboardQueries, ...globalDashboardQueries.slice(0, 2)];
  });
</script>

<svelte:head>
  <title>Dashboard</title>
</svelte:head>

{#if appStore.state.app.isUserLoggedIn}
  <div class="mt-8 mb-8 flex flex-row flex-wrap gap-10">
    {#each filteredQueries as query, i (`dashboard-${uid}-${i}`)}
      {#if [SEARCHTYPES.ADVISORY, SEARCHTYPES.DOCUMENT].includes(query.kind)}
        <AdvisoryQuery storedQuery={query}></AdvisoryQuery>
      {:else}
        <EventQuery storedQuery={query}></EventQuery>
      {/if}
    {/each}
    {#if appStore.state.app.isUserLoggedIn && appStore.isSourceManager()}
      <SourceEvents></SourceEvents>
    {/if}
  </div>
  <ErrorMessage error={loadQueryError}></ErrorMessage>
  <ErrorMessage error={loadIgnoredError}></ErrorMessage>
  <div class="mb-8 flex w-full max-w-[96%] flex-col gap-4 2xl:w-[46%]">
    <ImportStats updateIntervalInMinutes={10}></ImportStats>
  </div>
{/if}
