<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { tablePadding } from "$lib/table/defaults";
  import {
    Button,
    Table,
    TableHead,
    TableBody,
    TableHeadCell,
    TableBodyRow,
    TableBodyCell
  } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { push } from "svelte-spa-router";
  import { Modal } from "flowbite-svelte";
  let deleteModalOpen = false;

  const resetQueryToDelete = () => {
    return { name: "", id: -1 };
  };

  let queries: any[] = [];
  let orderBy = "";
  let errorMessage = "";
  let querytoDelete: any = resetQueryToDelete();

  const fetchQueries = async () => {
    const response = await request("/api/queries", "GET");
    if (response.ok) {
      queries = response.content;
    } else if (response.error) {
      errorMessage = `Could not load queries. ${getErrorMessage(response.error)}`;
    }
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

  onMount(() => {
    fetchQueries();
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

<h2 class="mb-3 text-lg">User defined queries</h2>
<Button class="mb-6 mt-3" href="/#/configuration/userqueries"
  ><i class="bx bx-plus"></i>New query</Button
>
{#if queries.length > 0}
  <div class="flex flex-row">
    <div class="mb-12 w-1/3">
      <Table hoverable={true} noborder={true}>
        <TableHead class="cursor-pointer">
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
          {#each queries as query}
            <TableBodyRow
              on:click={() => {
                push(`/configuration/userqueries/${query.id}`);
              }}
              class="cursor-pointer"
            >
              <TableBodyCell>{query.name ?? "-"}</TableBodyCell>
              <TableBodyCell>{query.description ?? "-"}</TableBodyCell>
              <td>
                <button
                  title={`clone ${query.name}`}
                  on:click|stopPropagation={() => {
                    push(`/configuration/userqueries/?clone=${query.id}`);
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
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
      <ErrorMessage message={errorMessage}></ErrorMessage>
    </div>
  </div>
{/if}
