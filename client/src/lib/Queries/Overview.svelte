<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { tablePadding, tdClass } from "$lib/Table/defaults";
  import {
    Button,
    Table,
    TableHead,
    TableHeadCell,
    TableBodyCell,
    Spinner,
    Checkbox
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { push } from "svelte-spa-router";
  import { Modal, Img } from "flowbite-svelte";
  import { ADMIN } from "$lib/workflow";
  import { isRoleIncluded } from "$lib/permissions";
  import { appStore } from "$lib/store";
  import Sortable from "sortablejs";
  import { saveStoredQuery, type Query } from "./query";
  let deleteModalOpen = false;

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  let queries: Query[] = [];
  let ignoredQueries: number[] = [];
  let orderBy = "";
  let errorMessage: ErrorDetails | null;
  let ignoreErrorMessage: ErrorDetails | null;
  let cloneErrorMessage: ErrorDetails | null;
  let querytoDelete: any = resetQueryToDelete();
  let loading = false;
  let columnList: any;
  let columnListAdmin: any;
  let clonedQueriesAlready = true;

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

  const fetchIgnored = async () => {
    const response = await request(`/api/queries/ignore`, "GET");
    if (response.ok) {
      ignoredQueries = response.content;
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load queries.`, response);
    }
  };

  const changeIgnored = async (id: number, isChecked: boolean) => {
    ignoreErrorMessage = null;
    const method = isChecked ? "DELETE" : "POST";
    const response = await request(`/api/queries/ignore/${id}`, method);
    if (response.ok) {
      if (isChecked) {
        ignoredQueries = ignoredQueries.filter((i) => i !== id);
      } else {
        ignoredQueries.push(id);
      }
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not change option.`, response);
    }
  };

  const deleteQuery = async () => {
    const response = await request(`/api/queries/${querytoDelete.id}`, "DELETE");
    if (response.error) {
      errorMessage = getErrorDetails(`Could not delete query ${querytoDelete.name}.`, response);
      querytoDelete = resetQueryToDelete();
      deleteModalOpen = false;
    }
    fetchData();
  };

  const updateQueryOrder = async (queries: Query[]) => {
    let nodes = columnList.querySelectorAll(".columnName");
    type Order = {
      id: number;
      order: number;
    };
    let orders: Order[] = [];
    let i = 0;
    for (const node of nodes) {
      let columnName = node.innerText;
      let query = queries.find((q) => q.name === columnName);
      if (query) {
        orders.push({ id: query.id, order: i });
      }
      i++;
    }

    let response = await request(`/api/queries/orders`, "POST", JSON.stringify(orders));
    if (!response.ok && response.error) {
      errorMessage = getErrorDetails(`Could not update query order.`, response);
    }
    if (response.ok) {
      push(`/queries/`);
    }
  };

  const checkIfCloned = () => {
    const personalQueries = queries.filter((q) => !q.global);
    const globalDashboardQueries = queries.filter((q) => q.dashboard && q.global);
    const firstTwoQueries = globalDashboardQueries.slice(0, 2);
    clonedQueriesAlready = false;
    firstTwoQueries.forEach((globalQuery: Query) => {
      const foundQuery = personalQueries.find((persQuery: Query) => {
        const keys = Object.keys(globalQuery);
        keys.forEach((key) => {
          if (persQuery[key] !== globalQuery[key]) {
            return false;
          }
        });
        return true;
      });
      if (foundQuery) {
        clonedQueriesAlready = true;
      }
    });
  };

  const fetchData = async () => {
    await fetchQueries();
    fetchIgnored();
    checkIfCloned();
  };

  onMount(() => {
    fetchData();
  });
  $: userQueries = queries.filter((q: Query) => {
    return !q.global;
  });
  $: adminQueries = queries.filter((q: Query) => {
    return q.global;
  });

  const cloneDashboardQueries = async () => {
    cloneErrorMessage = null;
    const globalDashboardQueries = queries.filter((q) => q.dashboard && q.global);
    const firstTwoQueries = globalDashboardQueries.slice(0, 2);
    for (let i = 0; i < 2; i++) {
      const queryToClone = firstTwoQueries[i];
      if (queryToClone) {
        queryToClone.global = false;
        const response = await saveStoredQuery(queryToClone);
        if (!response.ok && response.error) {
          cloneErrorMessage = getErrorDetails(`Failed to clone queries.`, response);
        }
      }
    }
    fetchData();
  };

  const elementDragEventUserQuery = () => {
    updateQueryOrder(userQueries);
  };

  const elementDragEventAdminQuery = () => {
    updateQueryOrder(adminQueries);
  };

  $: if (columnList) {
    Sortable.create(columnList, {
      animation: 150,
      onEnd: elementDragEventUserQuery
    });
  }
  $: if (columnListAdmin) {
    Sortable.create(columnList, {
      animation: 150,
      onEnd: elementDragEventAdminQuery
    });
  }
</script>

<svelte:head>
  <title>User defined queries</title>
</svelte:head>

<Modal size="xs" title={querytoDelete.name} bind:open={deleteModalOpen} autoclose outsideclose>
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
<ErrorMessage error={ignoreErrorMessage}></ErrorMessage>
<Button class="mb-6 mt-3" href="/#/queries/new"><i class="bx bx-plus"></i>New query</Button>
{#if queries.length > 0}
  <div class="flex flex-row flex-wrap gap-12">
    <div class="mb-12 w-fit">
      <span class="text-2xl">Personal</span>
      <hr class="mb-6" />
      <div class="max-h-[66vh] overflow-auto">
        <Table hoverable={true} noborder={true}>
          <TableHead>
            <TableHeadCell padding={tablePadding}></TableHeadCell>
            <TableHeadCell padding={tablePadding} on:click={() => {}}
              >Name<i
                class:bx={true}
                class:bx-caret-up={orderBy == "name"}
                class:bx-caret-down={orderBy == "-name"}
              ></i></TableHeadCell
            >
            <TableHeadCell padding={tablePadding} on:click={() => {}}
              >Description<i
                class:bx={true}
                class:bx-caret-up={orderBy == "description"}
                class:bx-caret-down={orderBy == "-description"}
              ></i>
            </TableHeadCell>
            <TableHeadCell padding={tablePadding} on:click={() => {}}>
              <div title={"Show on your personal dashboard"}>Pers. Dashb.</div>
            </TableHeadCell>
            <TableHeadCell></TableHeadCell>
          </TableHead>
          <tbody bind:this={columnList}>
            {#each userQueries as query, index (index)}
              <tr
                on:click={() => {
                  push(`/queries/${query.id}`);
                }}
                class="cursor-pointer"
                ><TableBodyCell {tdClass}>
                  <Img
                    src="grid-dots-vertical-rounded.svg"
                    class="h-4 min-h-2 min-w-2 invert-[.5]"
                  />
                </TableBodyCell>
                <TableBodyCell {tdClass}>
                  <span class="columnName">{query.name ?? "-"}</span>
                </TableBodyCell>
                <TableBodyCell {tdClass}>{query.description ?? "-"}</TableBodyCell>
                <TableBodyCell {tdClass}>
                  <Checkbox
                    on:click={(event) => {
                      event.stopPropagation();
                      // @ts-expect-error Cannot use TS:
                      // https://github.com/sveltejs/language-tools/blob/master/docs/preprocessors/typescript.md#can-i-use-typescript-syntax-inside-the-templatemustache-tags
                      // But without ignore we would get an error.
                      changeIgnored(query.id, event.target?.checked);
                    }}
                    disabled={!ignoredQueries}
                    checked={!ignoredQueries.includes(query.id)}
                  ></Checkbox>
                </TableBodyCell>
                <td>
                  <button
                    title={`clone ${query.name}`}
                    on:click|stopPropagation={() => {
                      push(`/queries/new?clone=${query.id}`);
                    }}><i class="bx bx-copy"></i></button
                  >
                  <button
                    on:click|stopPropagation={() => {
                      querytoDelete = {
                        name: query.name,
                        id: query.id
                      };
                      deleteModalOpen = true;
                    }}
                    title={`delete ${query.name}`}><i class="bx bx-trash text-red-500"></i></button
                  >
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </div>
    </div>
    <div class="mb-12 w-fit">
      <span class="text-2xl">Global</span>
      <hr class="mb-6" />
      <div class="mb-2 max-h-[66vh] overflow-auto">
        <Table hoverable={true} noborder={true}>
          <TableHead>
            <TableHeadCell padding={tablePadding}></TableHeadCell>
            <TableHeadCell padding={tablePadding} on:click={() => {}}
              >Name<i
                class:bx={true}
                class:bx-caret-up={orderBy == "name"}
                class:bx-caret-down={orderBy == "-name"}
              ></i></TableHeadCell
            >
            <TableHeadCell padding={tablePadding} on:click={() => {}}
              >Description<i
                class:bx={true}
                class:bx-caret-up={orderBy == "description"}
                class:bx-caret-down={orderBy == "-description"}
              ></i>
            </TableHeadCell>
            <TableHeadCell padding={tablePadding} on:click={() => {}}>
              <div title={"Show on your personal dashboard"}>Pers. Dashb.</div>
            </TableHeadCell>
            <TableHeadCell></TableHeadCell>
          </TableHead>
          <tbody bind:this={columnListAdmin}>
            {#each adminQueries as query, index (index)}
              <tr
                on:click={() => {
                  if (query.global && isRoleIncluded(appStore.getRoles(), [ADMIN])) {
                    push(`/queries/${query.id}`);
                  }
                }}
                class={!(query.global && !isRoleIncluded(appStore.getRoles(), [ADMIN]))
                  ? "cursor-pointer"
                  : ""}
                ><TableBodyCell {tdClass}>
                  {#if !(query.global && !isRoleIncluded(appStore.getRoles(), [ADMIN]))}
                    <Img
                      src="grid-dots-vertical-rounded.svg"
                      class="h-4 min-h-2 min-w-2 invert-[.5]"
                    />
                  {/if}
                </TableBodyCell>
                <TableBodyCell {tdClass}>
                  <span>{query.name ?? "-"}</span>
                </TableBodyCell>
                <TableBodyCell {tdClass}>{query.description ?? "-"}</TableBodyCell>
                <TableBodyCell {tdClass}>
                  <Checkbox
                    on:click={(event) => {
                      event.stopPropagation();
                      // @ts-expect-error Cannot use TS:
                      // https://github.com/sveltejs/language-tools/blob/master/docs/preprocessors/typescript.md#can-i-use-typescript-syntax-inside-the-templatemustache-tags
                      // But without ignore we would get an error.
                      changeIgnored(query.id, event.target?.checked);
                    }}
                    disabled={!ignoredQueries}
                    checked={!ignoredQueries.includes(query.id)}
                  ></Checkbox>
                </TableBodyCell>
                <td>
                  <button
                    title={`clone ${query.name}`}
                    on:click|stopPropagation={() => {
                      push(`/queries/new?clone=${query.id}`);
                    }}><i class="bx bx-copy"></i></button
                  >
                  {#if !(query.global && !isRoleIncluded(appStore.getRoles(), [ADMIN]))}
                    <button
                      on:click|stopPropagation={() => {
                        querytoDelete = {
                          name: query.name,
                          id: query.id
                        };
                        deleteModalOpen = true;
                      }}
                      title={`delete ${query.name}`}
                      ><i class="bx bx-trash text-red-500"></i></button
                    >
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </div>
      <Button on:click={cloneDashboardQueries} disabled={clonedQueriesAlready} color="light"
        >Clone the global dashboard queries for me</Button
      >
      <ErrorMessage error={cloneErrorMessage}></ErrorMessage>
    </div>
  </div>
{/if}
