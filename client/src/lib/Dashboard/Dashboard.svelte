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
  import New from "./New.svelte";
  import RecentActivities from "./RecentActivities.svelte";
  import { request } from "$lib/utils";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  let queries: any[];
  $: filteredQueries = queries
    ? queries.filter(
        (query) => query.dashboard === true && appStore.getRoles().includes(query.role)
      )
    : [];
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

  onMount(async () => {
    queries = await fetchStoredQueries();
  });
</script>

<svelte:head>
  <title>Dashboard</title>
</svelte:head>

{#if $appStore.app.isUserLoggedIn}
  <div class="mb-8 mt-8 flex gap-x-10 gap-y-4">
    {#each advisoryQueries as query}
      <New storedQuery={query}></New>
    {/each}
    {#each eventQueries as query}
      <RecentActivities storedQuery={query}></RecentActivities>
    {/each}
  </div>
  <ErrorMessage error={loadQueryError}></ErrorMessage>
{/if}
