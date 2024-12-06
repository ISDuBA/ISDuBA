<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Img, Table, TableHead, TableHeadCell, TableBodyCell } from "flowbite-svelte";
  import { tablePadding, tdClass } from "$lib/Table/defaults";
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
  import CIconButton from "$lib/Components/CIconButton.svelte";
  import { getContext } from "svelte";
  import { setIgnored, updateStoredQuery, type Query } from "./query";
  import { push } from "svelte-spa-router";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import Sortable from "sortablejs";
  import { request } from "$lib/request";

  export let tableContainerID: string | null = null;
  export let title = "";
  export let queries: Query[] | undefined = [];
  export let newQueries: Query[] = [];
  export let ignoredQueries: number[] | null = null;
  export let isAllowedToEdit = false;
  export let isAllowedToClone = true;

  const queryContext: any = getContext("queryContext");

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  let orderBy = "";
  let querytoDelete: any = resetQueryToDelete();
  let ignoreErrorMessage: ErrorDetails | null = null;
  let cloneErrorMessage: ErrorDetails | null = null;
  let orderErrorMessage: ErrorDetails | null = null;
  let columnList: any;
  let isLoading = false;

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
      orderErrorMessage = getErrorDetails(`Could not update query order.`, response);
    }
    if (response.ok) {
      push(`/queries/`);
    }
  };

  const elementDragEventUserQuery = () => {
    if (queries) {
      updateQueryOrder(queries);
    }
  };

  $: if (columnList && isAllowedToEdit) {
    Sortable.create(columnList, {
      animation: 150,
      onEnd: elementDragEventUserQuery
    });
  }

  const unsetErrors = () => {
    ignoreErrorMessage = null;
    cloneErrorMessage = null;
  };

  const changeIgnored = async (id: number, isChecked: boolean) => {
    isLoading = true;
    if (ignoredQueries) {
      unsetErrors();
      let newIgnored: number[] | null = null;
      ({ ignoredQueries: newIgnored, errorMessage: ignoreErrorMessage } = await setIgnored(
        id,
        isChecked
      ));
      if (ignoreErrorMessage === null) {
        type UpdateIgnoredFunction = (newIgnored: number[]) => void;
        const updateIgnored: UpdateIgnoredFunction = queryContext["updateIgnored"];
        updateIgnored(newIgnored);
      }
    }
    isLoading = false;
  };

  const changeDashboard = async (id: number, isChecked: boolean) => {
    isLoading = true;
    unsetErrors();
    if (queries) {
      const queryToUpdate = queries.filter((q) => q.id === id)[0];
      if (queryToUpdate) {
        queryToUpdate.dashboard = isChecked;
        const response = await updateStoredQuery(queryToUpdate);
        if (!response.ok && response.error) {
          ignoreErrorMessage = getErrorDetails(`Could not change option.`, response);
        }
      }
    }
    isLoading = false;
  };

  const clone = async (query: Query) => {
    if (!queries) return;
    const cloneQuery = queryContext["cloneQuery"];
    const queryToClone = query;
    if (!isAllowedToEdit) {
      queryToClone.global = false;
    }
    await cloneQuery(query);
  };
</script>

<div class="w-fit">
  <div class="mb-1 flex items-center gap-4">
    <span class="text-2xl">{title}</span>
  </div>
  <hr class="mb-6" />
  <div id={tableContainerID} class="mb-2 max-h-[66vh] overflow-auto">
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
      {#if queries !== undefined && queries.length > 0}
        <tbody bind:this={columnList}>
          {#each queries as query, index (index)}
            <tr
              on:click={() => {
                push(`/queries/${query.id}`);
              }}
              class:cursor-pointer={true}
              class:motion-safe:animate-pulse-fast={newQueries.map((q) => q.id).includes(query.id)}
              ><TableBodyCell {tdClass}>
                {#if isAllowedToEdit}
                  <Img
                    src="grid-dots-vertical-rounded.svg"
                    class="h-4 min-h-2 min-w-2 invert-[.5]"
                  />
                {/if}
              </TableBodyCell>
              <TableBodyCell {tdClass}>
                <div class="text-wrap break-all">
                  <span class="columnName">{query.name ?? "-"}</span>
                </div>
              </TableBodyCell>
              <TableBodyCell {tdClass}
                ><div class="text-wrap break-all">{query.description ?? "-"}</div></TableBodyCell
              >
              <TableBodyCell {tdClass}>
                <CCheckbox
                  on:change={() => {
                    changeDashboard(query.id, query.dashboard);
                  }}
                  bind:checked={query.dashboard}
                  class={isAllowedToEdit ? "" : "text-gray-300"}
                  disabled={!isAllowedToEdit || isLoading}
                ></CCheckbox>
              </TableBodyCell>
              <TableBodyCell {tdClass}>
                <CCheckbox
                  on:click={() => {
                    changeIgnored(query.id, !ignoredQueries?.includes(query.id));
                  }}
                  disabled={!ignoredQueries || isLoading}
                  checked={ignoredQueries !== null && ignoredQueries.includes(query.id)}
                ></CCheckbox>
              </TableBodyCell>
              <td>
                {#if !queries.find((q) => q.id === query.id) || isAllowedToClone}
                  <CIconButton
                    title={`clone ${query.name}`}
                    icon="copy"
                    on:click={async () => {
                      clone(query);
                    }}
                  ></CIconButton>
                {/if}
                {#if isAllowedToEdit}
                  <CIconButton
                    on:click={() => {
                      querytoDelete = {
                        name: query.name,
                        id: query.id
                      };
                      const openDeleteModal = queryContext["openDeleteModal"];
                      openDeleteModal(querytoDelete);
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
      {:else if queries !== undefined && queries.length === 0}
        <tbody>
          <span>No queries found</span>
        </tbody>
      {/if}
    </Table>
  </div>
  <div class="flex flex-col">
    <slot></slot>
    <ErrorMessage error={ignoreErrorMessage}></ErrorMessage>
    <ErrorMessage error={cloneErrorMessage}></ErrorMessage>
    <ErrorMessage error={orderErrorMessage}></ErrorMessage>
  </div>
</div>
