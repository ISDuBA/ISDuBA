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
    TableSearch,
    Pagination
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";
  import SectionHeader from "$lib/SectionHeader.svelte";

  $: if ($appStore.app.keycloak.authenticated) {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      const response = await fetch(documentURL, {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      });
      if (response.ok) {
        ({ count, documents } = await response.json());
      } else {
        // Do errorhandling
      }
    });
  }

  let pages = [{ name: "1" }, { name: "2" }, { name: "3" }, { name: "4" }, { name: "5" }];
  let documents: any = [];
  let count = 0;
  let searchTerm: string = "";
  const sortState: any = {
    id: "",
    publisher: "",
    title: "",
    trackingID: "",
    version: "",
    activeSortColumn: ""
  };

  let orderBy = "id";
  let limit = 10;
  const previous = () => {
    alert("Previous btn clicked. Make a call to your server to fetch data.");
  };
  const next = () => {
    alert("Next btn clicked. Make a call to your server to fetch data.");
  };
  const handleClick = () => {
    alert("Page clicked");
  };

  const sortDocuments = (column: string) => {};
  $: documentURL = encodeURI(
    `/api/documents?&count=1&order=${orderBy}&limit=${limit}&columns=id publisher title tracking_id version`
  );
</script>

{#if $appStore.app.keycloak.authenticated}
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
          {#each documents as item}
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
  <div class="mt-3 flex">
    <span class="mr-3 flex-grow">Showing the first 10 of {count} Entries</span>
    <Pagination {pages} on:previous={previous} on:next={next} on:click={handleClick} />
  </div>
{/if}
