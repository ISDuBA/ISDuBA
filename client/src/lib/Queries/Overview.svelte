<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { tablePadding, tdClass } from "$lib/table/defaults";
  import {
    Button,
    Table,
    TableHead,
    TableBody,
    TableHeadCell,
    TableBodyCell,
    Spinner
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { push } from "svelte-spa-router";
  import { Modal } from "flowbite-svelte";
  import { ADMIN } from "$lib/workflow";
  import { isRoleIncluded } from "$lib/permissions";
  import { appStore } from "$lib/store";
  let deleteModalOpen = false;

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  let queries: any[] = [];
  let orderBy = "";
  let errorMessage = "";
  let querytoDelete: any = resetQueryToDelete();
  let hoveredAdminQuery: any;
  let hoveredUserQuery: any;
  let loading = false;

  const fetchQueries = async () => {
    loading = true;
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      const result = response.content;
      queries = result.sort((q1: any, q2: any) => {
        return q1.num > q2.num;
      });
    } else if (response.error) {
      errorMessage = `Could not load queries. ${getErrorMessage(response.error)}`;
    }
    loading = false;
  };

  const deleteQuery = async () => {
    const response = await request(`/api/queries/${querytoDelete.id}`, "DELETE");
    if (response.error) {
      errorMessage = `Could not delete query ${querytoDelete.name}. ${getErrorMessage(response.error)}`;
      querytoDelete = resetQueryToDelete();
      deleteModalOpen = false;
    }
    fetchQueries();
  };

  const swapQueryNum = async (query1: any, query2: any) => {
    loading = true;
    let formData = new FormData();
    let TEMP_NUM = 1_000_000;
    formData.append("num", `${TEMP_NUM}`);
    const response1 = await request(`/api/queries/${query1.id}`, "PUT", formData);
    if (response1.ok) {
      formData = new FormData();
      formData.append("num", `${query1.num}`);
      const response2 = await request(`/api/queries/${query2.id}`, "PUT", formData);
      if (response2.ok) {
        formData = new FormData();
        formData.append("num", `${query2.num}`);
        const response3 = await request(`/api/queries/${query1.id}`, "PUT", formData);
        if (response3.error) {
          errorMessage = `An error occured while swapping order of queries`;
        }
      }
      if (response2.error) {
        errorMessage = `An error occured while swapping order of queries`;
      }
    }
    if (response1.error) {
      errorMessage = `An error occured while swapping order of queries`;
    }
    loading = false;
    fetchQueries();
  };

  const promoteUserQuery = () => {
    if (hoveredUserQuery === 0) return;
    const first = useryQueries[hoveredUserQuery];
    const second = useryQueries[hoveredUserQuery - 1];
    swapQueryNum(second, first);
  };

  const demoteUserQuery = () => {
    if (hoveredUserQuery === useryQueries.length - 1) return;
    const first = useryQueries[hoveredUserQuery];
    const second = useryQueries[hoveredUserQuery + 1];
    swapQueryNum(first, second);
  };

  const promoteAdminQuery = () => {
    if (hoveredAdminQuery === 0) return;
    const first = adminQueries[hoveredAdminQuery];
    const second = adminQueries[hoveredAdminQuery - 1];
    swapQueryNum(second, first);
  };

  const demoteAdminQuery = () => {
    if (hoveredAdminQuery === adminQueries.length - 1) return;
    const first = adminQueries[hoveredAdminQuery];
    const second = adminQueries[hoveredAdminQuery + 1];
    swapQueryNum(first, second);
  };

  onMount(() => {
    fetchQueries();
  });
  $: useryQueries = queries.filter((q: any) => {
    return !q.global;
  });
  $: adminQueries = queries.filter((q: any) => {
    return q.global;
  });
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
<div class:invisible={!loading}>
  Loading ...
  <Spinner color="gray" size="4"></Spinner>
</div>
<ErrorMessage message={errorMessage}></ErrorMessage>
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
          <TableBody>
            {#each useryQueries as query, index (index)}
              <tr
                on:click={() => {
                  push(`/queries/${query.id}`);
                }}
                on:mouseover={() => {
                  hoveredUserQuery = index;
                }}
                on:mouseout={() => {
                  hoveredUserQuery = -1;
                }}
                on:blur={() => {}}
                on:focus={() => {}}
                class="cursor-pointer"
                ><TableBodyCell {tdClass}>
                  <div
                    class:invisible={hoveredUserQuery !== index}
                    class:w-1={true}
                    class:flex={true}
                    class:flex-col={true}
                  >
                    <button
                      class="h-4"
                      on:click|stopPropagation={() => {
                        promoteUserQuery();
                      }}
                    >
                      <i class="bx bxs-up-arrow-circle"></i>
                    </button>
                    <button
                      on:click|stopPropagation={() => {
                        demoteUserQuery();
                      }}
                      class="h-4"
                    >
                      <i class="bx bxs-down-arrow-circle"></i>
                    </button>
                  </div>
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
          </TableBody>
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
          <TableBody>
            {#each adminQueries as query, index (index)}
              <tr
                on:click={() => {
                  if (query.global && isRoleIncluded(appStore.getRoles(), [ADMIN])) {
                    push(`/queries/${query.id}`);
                  }
                }}
                on:mouseover={() => {
                  hoveredAdminQuery = index;
                }}
                on:mouseout={() => {
                  hoveredAdminQuery = -1;
                }}
                on:blur={() => {}}
                on:focus={() => {}}
                class={!(query.global && !isRoleIncluded(appStore.getRoles(), [ADMIN]))
                  ? "cursor-pointer"
                  : ""}
                ><TableBodyCell {tdClass}>
                  {#if !(query.global && !isRoleIncluded(appStore.getRoles(), [ADMIN]))}
                    <div
                      class:invisible={hoveredAdminQuery !== index}
                      class:w-1={true}
                      class:flex={true}
                      class:flex-col={true}
                    >
                      <button
                        class="h-4"
                        on:click|stopPropagation={() => {
                          promoteAdminQuery();
                        }}
                      >
                        <i class="bx bxs-up-arrow-circle"></i>
                      </button>
                      <button
                        on:click|stopPropagation={() => {
                          demoteAdminQuery();
                        }}
                        class="h-4"
                      >
                        <i class="bx bxs-down-arrow-circle"></i>
                      </button>
                    </div>
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
          </TableBody>
        </Table>
      </div>
    </div>
  </div>
{/if}
