<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch
  } from "flowbite-svelte";
  let documents: any = [];
  let searchTerm: string = "";
  const sortState: any = {
    id: "",
    publisher: "",
    title: "",
    trackingID: "",
    version: "",
    activeSortColumn: ""
  };
  import { tdClass, tablePadding } from "$lib/table/defaults";
  import SectionHeader from "$lib/SectionHeader.svelte";

  const defaultSortFunction = (attribute: string) => {
    return {
      asc: (ad1: any, ad2: any) => {
        if (ad1[attribute] < ad2[attribute]) return -1;
        if (ad1[attribute] > ad2[attribute]) return 1;
        return 0;
      },
      desc: (ad2: any, ad1: any) => {
        if (ad1[attribute] < ad2[attribute]) return -1;
        if (ad1[attribute] > ad2[attribute]) return 1;
        return 0;
      }
    };
  };
  const sortFunctionsByColumn: any = {
    id: defaultSortFunction("id"),
    publisher: defaultSortFunction("publisher"),
    title: defaultSortFunction("title"),
    trackingID: defaultSortFunction("tracking_id"),
    version: defaultSortFunction("version"),
    state: defaultSortFunction("state")
  };
  const sortDocuments = (column: string) => {
    sortState["activeSortColumn"] = column;
    if (sortState[column] === "asc") {
      documents = [...documents.sort(sortFunctionsByColumn[column]["desc"])];
      sortState[column] = "desc";
    } else {
      documents = [...documents.sort(sortFunctionsByColumn[column]["asc"])];
      sortState[column] = "asc";
    }
  };
  const allUnreadDocuments = encodeURI(
    "/api/documents?columns=id publisher title tracking_id version"
  );
  $: filteredItems = documents;
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      $appStore.app.keycloak.updateToken(5).then(async () => {
        const response = await fetch(allUnreadDocuments, {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          }
        });
        if (response.ok) {
          ({ documents } = await response.json());
          sortDocuments("id");
        } else {
          // Do errorhandling
        }
      });
    }
  });
</script>

<SectionHeader title="New Events"></SectionHeader>
<Table>
  <TableHead>
    <TableHeadCell padding={tablePadding}>Description</TableHeadCell>
    <TableHeadCell padding={tablePadding}>Advisory</TableHeadCell>
  </TableHead>
  <TableBody>
    <TableBodyRow>
      <TableBodyCell {tdClass}>Comment added</TableBodyCell>
      <TableBodyCell {tdClass}>Sick PSIRT SCA-2022-0032</TableBodyCell>
    </TableBodyRow>
  </TableBody>
</Table>
<SectionHeader title="New Advisories"></SectionHeader>
{#if documents}
  <div style="width: 100%;overflow-y: auto">
    <TableSearch placeholder="Search" hoverable={true} bind:inputValue={searchTerm}>
      <TableHead class="cursor-pointer">
        <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("id")}
          >ID<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "id" && sortState["id"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "id" &&
              sortState["id"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("publisher")}
          >Publisher<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "publisher" &&
              sortState["publisher"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "publisher" &&
              sortState["publisher"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("title")}
          >Title<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "title" &&
              sortState["title"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "title" &&
              sortState["title"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("trackingID")}
          >Tracking ID<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "trackingID" &&
              sortState["trackingID"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "trackingID" &&
              sortState["trackingID"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("version")}
          >Version<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "version" &&
              sortState["version"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "version" &&
              sortState["version"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => sortDocuments("state")}
          >State<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "state" &&
              sortState["state"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "state" &&
              sortState["state"] === "desc"}
          ></i></TableHeadCell
        >
      </TableHead>
      <TableBody>
        {#each filteredItems as item}
          <TableBodyRow
            class="cursor-pointer"
            on:click={() => {
              push(`/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`);
            }}
          >
            <TableBodyCell {tdClass}>{item.id}</TableBodyCell>
            <TableBodyCell {tdClass}>{item.publisher}</TableBodyCell>
            <TableBodyCell {tdClass}>{item.title}</TableBodyCell>
            <TableBodyCell {tdClass}>{item.tracking_id}</TableBodyCell>
            <TableBodyCell {tdClass}>{item.version}</TableBodyCell>
            <TableBodyCell {tdClass}>{item.state}</TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </TableSearch>
  </div>
{/if}
