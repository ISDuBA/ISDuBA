<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { tablePadding, tdClass } from "$lib/Table/defaults";
  import { Button, Table, TableHead, TableHeadCell, TableBodyCell, Spinner } from "flowbite-svelte";
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
  let deleteModalOpen = false;

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  type Query = {
    advisories: boolean;
    columns: string[];
    definer: string;
    global: boolean;
    id: number;
    name: string;
    num: number;
    orders: string[] | undefined;
    query: string;
    description: string | undefined;
  };

  let queries: Query[] = [];
  let orderBy = "";
  let errorMessage: ErrorDetails | null;
  let querytoDelete: any = resetQueryToDelete();
  let loading = false;
  let columnList: any;
  let columnListAdmin: any;

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

  const deleteQuery = async () => {
    const response = await request(`/api/queries/${querytoDelete.id}`, "DELETE");
    if (response.error) {
      errorMessage = getErrorDetails(`Could not delete query ${querytoDelete.name}.`, response);
      querytoDelete = resetQueryToDelete();
      deleteModalOpen = false;
    }
    fetchQueries();
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

  onMount(() => {
    fetchQueries();
  });
  $: userQueries = queries.filter((q: Query) => {
    return !q.global;
  });
  $: adminQueries = queries.filter((q: Query) => {
    return q.global;
  });

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
    </div>
  </div>
{/if}
