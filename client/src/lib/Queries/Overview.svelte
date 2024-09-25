<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Spinner } from "flowbite-svelte";
  import { onMount, setContext } from "svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { Modal } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import { fetchIgnored, setIgnored, createStoredQuery, proposeName, type Query } from "./query";
  import QueryTable from "./QueryTable.svelte";
  let deleteModalOpen = false;

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  let queries: Query[] | undefined = [];
  let ignoredQueries: number[] = [];
  let errorMessage: ErrorDetails | null;
  let cloneErrorMessage: ErrorDetails | null;
  let querytoDelete: any = resetQueryToDelete();
  let loading = false;
  let isCloning = false;

  $: globalRelevantQueries = queries
    ?.filter((q) => q.definer === "system-default" && q.global && q.dashboard)
    .slice(0, 2);
  $: globalDashboardQueries = queries?.filter(
    (q) => q.dashboard && q.global && !globalRelevantQueries?.map((q) => q.id).includes(q.id)
  );
  $: globalSearchQueries = queries?.filter((q) => !q.dashboard && q.global);
  $: userQueries = queries?.filter((q: Query) => {
    return !q.global;
  });
  $: adminQueries = queries?.filter((q: Query) => {
    return q.global;
  });

  const fetchQueries = async () => {
    loading = true;
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      const result = response.content;
      queries = result.sort((q1: Query, q2: Query) => {
        return q1.num > q2.num;
      });
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load queries.`, response);
    }
    loading = false;
  };

  const onOpenDeleteModal = async (event: any) => {
    querytoDelete = event;
    deleteModalOpen = true;
  };

  const deleteQuery = async () => {
    unsetErrors();
    const response = await request(`/api/queries/${querytoDelete.id}`, "DELETE");
    if (response.error) {
      errorMessage = getErrorDetails(`Could not delete query ${querytoDelete.name}.`, response);
      querytoDelete = resetQueryToDelete();
      deleteModalOpen = false;
    }
    fetchData();
  };

  const unsetErrors = () => {
    errorMessage = null;
    cloneErrorMessage = null;
  };

  const fetchData = async () => {
    fetchQueries();
    ({ ignoredQueries, errorMessage } = await fetchIgnored());
  };

  const updateIgnored = (newIgnoredQueries: any) => {
    ignoredQueries = newIgnoredQueries;
  };

  onMount(() => {
    fetchData();
  });

  const cloneDashboardQueries = async () => {
    if (!globalRelevantQueries || !userQueries || !queries) return;
    cloneErrorMessage = null;
    const idsOfClonedQueries = [];
    let failed = false;
    isCloning = true;
    // Clone the special queries
    for (let i = 0; i < globalRelevantQueries.length; i++) {
      const queryToClone = globalRelevantQueries[i];
      if (queryToClone) {
        queryToClone.global = false;
        queryToClone.name = proposeName(queries, queryToClone.name);
        const response = await createStoredQuery(queryToClone);
        if (!response.ok && response.error) {
          cloneErrorMessage = getErrorDetails(`Failed to clone queries.`, response);
          failed = true;
        } else {
          idsOfClonedQueries.push(response.content.id);
        }
      }
    }
    if (!failed) {
      // Hide the special queries as they are now replaced
      for (let i = 0; i < globalRelevantQueries.length; i++) {
        if (!ignoredQueries.includes(globalRelevantQueries[i].id)) {
          ({ ignoredQueries, errorMessage = cloneErrorMessage } = await setIgnored(
            globalRelevantQueries[i].id,
            true
          ));
        }
      }
      await placeQueriesAtTop(idsOfClonedQueries);
    }
    isCloning = false;
    fetchData();
  };

  const cloneQuery = async (query: Query) => {
    if (!queries) return;
    isCloning = true;
    query.name = proposeName(queries, query.name);
    const response = await createStoredQuery(query);
    if (!response.ok && response.error) {
      cloneErrorMessage = getErrorDetails(`Failed to clone query.`, response);
    } else if (response.ok) {
      await placeQueriesAtTop([response.content.id]);
      fetchData();
    }
    isCloning = false;
  };

  const placeQueriesAtTop = async (queryIDs: number[], global = false) => {
    if (!userQueries || !adminQueries) return;
    type Order = {
      id: number;
      order: number;
    };
    let orders: Order[] = [];
    let count = 0;
    const userQueryIDs = userQueries.map((q) => q.id);
    const adminQueryIDs = adminQueries.map((q) => q.id);
    const newOrderIDs: number[] = global
      ? [...userQueryIDs, ...queryIDs, ...adminQueryIDs]
      : [...queryIDs, ...userQueryIDs, ...adminQueryIDs];
    for (let i = 0; i < newOrderIDs.length; i++) {
      orders.push({ id: newOrderIDs[i], order: count });
      count++;
    }
    let response = await request(`/api/queries/orders`, "POST", JSON.stringify(orders));
    if (!response.ok && response.error) {
      cloneErrorMessage = getErrorDetails(`Could not update query order.`, response);
    }
  };

  setContext("queryContext", {
    cloneQuery: cloneQuery,
    openDeleteModal: onOpenDeleteModal,
    updateIgnored: updateIgnored
  });
</script>

<svelte:head>
  <title>User defined queries</title>
</svelte:head>

<Modal
  size="xs"
  title={querytoDelete.name}
  bind:open={deleteModalOpen}
  autoclose
  outsideclose
  classHeader="flex justify-between items-center p-4 md:p-5 rounded-t-lg break-all"
>
  <div class="text-center">
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete this query?
    </h3>
    <Button
      on:click={() => {
        deleteQuery();
      }}
      color="red"
      class="me-2">Yes, I'm sure</Button
    >
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

<h2 class="text-lg">Queries</h2>
<div class:invisible={!loading} class={loading ? "loadingFadeIn" : ""}>
  Loading ...
  <Spinner color="gray" size="4"></Spinner>
</div>
<ErrorMessage error={errorMessage}></ErrorMessage>
{#if queries && queries.length > 0}
  <div class="flex flex-row flex-wrap gap-12">
    <QueryTable {ignoredQueries} isAllowedToEdit={true} queries={userQueries} title="Personal">
      <Button class="mb-2 mt-3 w-fit" href="/#/queries/new"
        ><i class="bx bx-plus me-2"></i>New query</Button
      >
    </QueryTable>

    {#if !appStore.isAdmin()}
      <QueryTable
        {ignoredQueries}
        isAllowedToClone={false}
        queries={globalRelevantQueries}
        title="Global relevant dashboard queries"
      >
        <ErrorMessage error={cloneErrorMessage}></ErrorMessage>
        <Button class="h-fit w-fit text-sm" on:click={cloneDashboardQueries} disabled={isCloning}>
          <i class="bx bx-copy me-2"></i>
          <span class="me-2">Clone relevant queries and hide cloned queries</span>
          <div class:invisible={!isCloning} class={isCloning ? "loadingFadeIn text-white" : ""}>
            <Spinner color="white" size="4"></Spinner>
          </div>
        </Button>
      </QueryTable>

      <QueryTable
        {ignoredQueries}
        queries={globalDashboardQueries}
        title="Global dashboard queries (not displayed)"
      ></QueryTable>

      <QueryTable {ignoredQueries} queries={globalSearchQueries} title="Global search queries"
      ></QueryTable>
    {:else}
      <QueryTable
        {ignoredQueries}
        queries={adminQueries}
        title="Global"
        isAllowedToEdit={appStore.isAdmin()}
      >
        <Button class="mb-2 mt-3 w-fit" href="/#/queries/new"
          ><i class="bx bx-plus me-2"></i>New query</Button
        >
      </QueryTable>
    {/if}
  </div>
{/if}
