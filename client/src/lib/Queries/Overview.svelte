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
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
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
  import { updateStoredQuery, type Query } from "./query";
  import CIconButton from "$lib/Components/CIconButton.svelte";
  let deleteModalOpen = false;

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  let queries: Query[] = [];
  let ignoredQueries: number[] = [];
  let orderBy = "";
  let errorMessage: ErrorDetails | null;
  let ignorePersonalErrorMessage: ErrorDetails | null;
  let ignoreGlobalErrorMessage: ErrorDetails | null;
  let cloneErrorMessage: ErrorDetails | null;
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

  const fetchIgnored = async () => {
    const response = await request(`/api/queries/ignore`, "GET");
    if (response.ok) {
      ignoredQueries = response.content;
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not load queries.`, response);
    }
  };

  const changeIgnored = async (id: number, isChecked: boolean) => {
    unsetErrors();
    const method = isChecked ? "POST" : "DELETE";
    const response = await request(`/api/queries/ignore/${id}`, method);
    if (response.ok) {
      if (isChecked) {
        ignoredQueries.push(id);
      } else {
        ignoredQueries = ignoredQueries.filter((i) => i !== id);
      }
    } else if (response.error) {
      errorMessage = getErrorDetails(`Could not change option.`, response);
    }
  };

  const changeDashboard = async (id: number, isChecked: boolean) => {
    unsetErrors();
    const queryToUpdate = queries.filter((q) => q.id === id)[0];
    if (queryToUpdate) {
      queryToUpdate.dashboard = isChecked;
      const response = await updateStoredQuery(queryToUpdate);
      if (!response.ok && response.error) {
        ignorePersonalErrorMessage = getErrorDetails(`Could not change option.`, response);
      }
    }
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

  const unsetErrors = () => {
    ignoreGlobalErrorMessage = null;
    ignorePersonalErrorMessage = null;
    errorMessage = null;
    cloneErrorMessage = null;
  };

  const fetchData = async () => {
    fetchQueries();
    fetchIgnored();
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
{#if queries.length > 0}
  <div class="flex flex-row flex-wrap gap-12">
    <div class="mb-2 w-fit">
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
              <div>Dashboard</div>
            </TableHeadCell>
            <TableHeadCell padding={tablePadding} on:click={() => {}}>
              <div>Hide</div>
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
                  <CCheckbox
                    on:click={(event) => {
                      // @ts-expect-error Cannot use TS:
                      // https://github.com/sveltejs/language-tools/blob/master/docs/preprocessors/typescript.md#can-i-use-typescript-syntax-inside-the-templatemustache-tags
                      // But without ignore we would get an error.
                      changeDashboard(query.id, event.explicitOriginalTarget?.checked);
                    }}
                    checked={query.dashboard}
                  ></CCheckbox>
                </TableBodyCell>
                <TableBodyCell {tdClass}>
                  <CCheckbox
                    on:click={(event) => {
                      // @ts-expect-error Cannot use TS (see explanation above)
                      changeIgnored(query.id, event.explicitOriginalTarget?.checked);
                    }}
                    disabled={!ignoredQueries}
                    checked={ignoredQueries.includes(query.id)}
                  ></CCheckbox>
                </TableBodyCell>
                <td>
                  <CIconButton
                    title={`clone ${query.name}`}
                    icon="copy"
                    on:click={() => {
                      push(`/queries/new?clone=${query.id}`);
                    }}
                  ></CIconButton>
                  <CIconButton
                    on:click={() => {
                      querytoDelete = {
                        name: query.name,
                        id: query.id
                      };
                      deleteModalOpen = true;
                    }}
                    title={`delete ${query.name}`}
                    icon="trash"
                    color="red"
                  ></CIconButton>
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </div>
      <Button class="mb-6 mt-3" href="/#/queries/new"><i class="bx bx-plus"></i>New query</Button>
      <ErrorMessage error={ignorePersonalErrorMessage}></ErrorMessage>
    </div>
    <div class="mb-2 w-fit">
      <span class="mb-1 text-2xl">Global</span>
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
              <div>Dashboard</div>
            </TableHeadCell>
            <TableHeadCell padding={tablePadding} on:click={() => {}}>
              <div title={"Show on your personal dashboard"}>Hide</div>
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
                  <CCheckbox
                    on:click={(event) => {
                      console.log(event);
                      // @ts-expect-error Cannot use TS (see explanation above)
                      changeDashboard(query.id, event.explicitOriginalTarget?.checked);
                    }}
                    checked={query.dashboard}
                    class={appStore.isAdmin() ? "" : "text-gray-300"}
                    disabled={!appStore.isAdmin()}
                  ></CCheckbox>
                </TableBodyCell>
                <TableBodyCell {tdClass}>
                  <CCheckbox
                    on:click={(event) => {
                      // @ts-expect-error Cannot use TS (see explanation above)
                      changeIgnored(query.id, event.explicitOriginalTarget?.checked);
                    }}
                    disabled={!ignoredQueries}
                    checked={ignoredQueries.includes(query.id)}
                  ></CCheckbox>
                </TableBodyCell>
                <td>
                  <CIconButton
                    title={`clone ${query.name}`}
                    icon="copy"
                    on:click={() => {
                      push(`/queries/new?clone=${query.id}`);
                    }}
                  ></CIconButton>
                  {#if !(query.global && !isRoleIncluded(appStore.getRoles(), [ADMIN]))}
                    <CIconButton
                      on:click={() => {
                        querytoDelete = {
                          name: query.name,
                          id: query.id
                        };
                        deleteModalOpen = true;
                      }}
                      title={`delete ${query.name}`}
                      icon="trash"
                      color="red"
                    ></CIconButton>
                  {/if}
                </td>
              </tr>
            {/each}
          </tbody>
        </Table>
      </div>
      <div class="flex flex-col">
        {#if appStore.isAdmin()}
          <Button class="mb-2 mt-3 w-fit" href="/#/queries/new"
            ><i class="bx bx-plus"></i>New query</Button
          >
        {/if}
        <ErrorMessage error={ignoreGlobalErrorMessage}></ErrorMessage>
        <ErrorMessage error={cloneErrorMessage}></ErrorMessage>
      </div>
    </div>
  </div>
{/if}
