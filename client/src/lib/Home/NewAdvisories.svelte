<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import {
    Input,
    Label,
    Select,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch,
    PaginationItem
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";

  let orderBy = "title";
  let limit = 10;
  let offset = 0;
  let count = 0;
  let currentPage = 1;
  let documents: any = [];
  let searchTerm: string = "";
  const sortState: any = {
    publisher: "",
    title: "",
    trackingID: "",
    version: "",
    activeSortColumn: ""
  };

  const fetchData = () => {
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
  };

  const previous = () => {
    if (offset - limit >= 0) {
      offset = offset - limit > 0 ? offset - limit : 0;
      currentPage -= 1;
    }
    fetchData();
  };
  const next = () => {
    if (offset + limit <= count) {
      offset = offset + limit;
      currentPage += 1;
    }
    fetchData();
  };

  const first = () => {
    offset = 0;
    currentPage = 1;
    fetchData();
  };

  const last = () => {
    offset = count - (count % limit);
    currentPage = numberOfPages;
    fetchData();
  };
  $: numberOfPages = Math.ceil(count / limit);
  $: documentURL = encodeURI(
    `/api/documents?&count=1&order=${orderBy}&limit=${limit}&offset=${offset}&columns=id publisher title tracking_id version`
  );
  $: if ($appStore.app.keycloak.authenticated) {
    fetchData();
  }
</script>

{#if documents}
  <div style="width: 100%;overflow-y: auto">
    <TableSearch placeholder="Search" hoverable={true} bind:inputValue={searchTerm}>
      <TableHead class="cursor-pointer">
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Publisher<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "publisher" &&
              sortState["publisher"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "publisher" &&
              sortState["publisher"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Title<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "title" &&
              sortState["title"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "title" &&
              sortState["title"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Tracking ID<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "trackingID" &&
              sortState["trackingID"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "trackingID" &&
              sortState["trackingID"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Version<i
            class:bx={true}
            class:bx-caret-up={sortState["activeSortColumn"] == "version" &&
              sortState["version"] === "asc"}
            class:bx-caret-down={sortState["activeSortColumn"] == "version" &&
              sortState["version"] === "desc"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => {}}
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
<div class="mb-12 mt-3 flex items-center">
  <div class="flex flex-grow items-center">
    <Label class="mr-3">Items per page</Label>
    <Select
      id="pagecount"
      class="mt-2 w-24"
      items={[
        { name: "10", value: 10 },
        { name: "20", value: 20 },
        { name: "50", value: 50 },
        { name: "100", value: 100 }
      ]}
      bind:value={limit}
      on:change={() => {
        offset = 0;
        fetchData();
      }}
    ></Select>
  </div>
  <div class="mr-3 flex-grow">
    <div class="flex">
      <PaginationItem on:click={first}>
        <i class="bx bx-arrow-to-left"></i>
      </PaginationItem>
      <PaginationItem on:click={previous}>
        <i class="bx bx-chevrons-left"></i>
      </PaginationItem>
      <div class="mx-3 flex items-center">
        <input
          class="w-16 cursor-pointer border pr-1 text-right"
          on:change={() => {
            if (!parseInt("" + currentPage)) currentPage = 1;
            currentPage = Math.floor(currentPage);
            if (currentPage < 1) currentPage = 1;
            if (currentPage > numberOfPages) currentPage = numberOfPages;
            offset = (currentPage - 1) * limit;
            fetchData();
          }}
          bind:value={currentPage}
        />
        <span class="mr-9">of {numberOfPages} Pages</span>
      </div>
      <PaginationItem on:click={next}>
        <i class="bx bx-chevrons-right"></i>
      </PaginationItem>
      <PaginationItem on:click={last}>
        <i class="bx bx-arrow-to-right"></i>
      </PaginationItem>
    </div>
  </div>
  <div class="mr-3">
    {count} entries in total
  </div>
</div>
