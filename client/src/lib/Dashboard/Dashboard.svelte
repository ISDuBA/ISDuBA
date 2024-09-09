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
  import AdvisoryQuery from "./AdvisoryQuery.svelte";
  import EventQuery from "./EventQuery.svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  let filteredQueries: any[] = [];
  let loadIgnoredError: ErrorDetails | null;

  $: advisoryQueries = filteredQueries.filter((query: any) =>
    [SEARCHTYPES.ADVISORY, SEARCHTYPES.DOCUMENT].includes(query.kind)
  );
  $: eventQueries = filteredQueries.filter((query: any) => query.kind === SEARCHTYPES.EVENT);
  let loadQueryError: ErrorDetails | null;

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
        query.definer === $appStore.app.tokenParsed?.preferred_username &&
        !query.global &&
        (!ignoredQueries || !ignoredQueries.includes(query.id))
    );
    const globalDashboardQueries = allQueries.filter(
      (query) =>
        query.dashboard &&
        query.global &&
        (appStore.getRoles().includes(query.role) || !query.role) &&
        !userDashboardQueries.find((q) => q.id === query.id) &&
        (!ignoredQueries || !ignoredQueries.includes(query.id))
    );
    filteredQueries = [...userDashboardQueries, ...globalDashboardQueries];
  });
</script>

<svelte:head>
  <title>Dashboard</title>
</svelte:head>

{#if $appStore.app.isUserLoggedIn}
  <div class="mb-8 mt-8 flex flex-wrap gap-10">
    {#each advisoryQueries as query}
      <AdvisoryQuery storedQuery={query}></AdvisoryQuery>
    {/each}
    {#each eventQueries as query}
      <EventQuery storedQuery={query}></EventQuery>
    {/each}
  </div>
  <ErrorMessage error={loadQueryError}></ErrorMessage>
  <ErrorMessage error={loadIgnoredError}></ErrorMessage>
{/if}
