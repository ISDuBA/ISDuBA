<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import {
    Button,
    Label,
    PaginationItem,
    Select,
    Search,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Table
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";
  import { onMount } from "svelte";
  import { Spinner } from "flowbite-svelte";

  let limit = 10;
  let offset = 0;
  let count = 0;
  let currentPage = 1;
  let documents: any = [];
  let searchTerm: string = "";
  let loading = false;
  let columns = [
    "id",
    "publisher",
    "title",
    "tracking_id",
    "version",
    "cvss_v2_score",
    "cvss_v3_score"
  ];
  let orderBy = "title";
  const fetchData = () => {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      loading = true;
      const response = await fetch(documentURL, {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      });
      if (response.ok) {
        ({ count, documents } = await response.json());
        documents = documents || [];
      } else {
        appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
      }
      loading = false;
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

  const switchSort = (column: string) => {
    if (column === orderBy) {
      orderBy[0] === "-" ? (orderBy = column) : (orderBy = `-${column}`);
    } else {
      orderBy = column;
    }
    fetchData();
  };

  $: searchSuffix = searchTerm ? ` "${searchTerm}" german search msg as and` : "";
  $: numberOfPages = Math.ceil(count / limit);
  $: documentURL = encodeURI(
    `/api/documents?query=$state new workflow =${searchSuffix}&advisories=true&count=1&order=${orderBy}&limit=${limit}&offset=${offset}&columns=${columns.join(" ")}`
  );
  onMount(async () => {
    fetchData();
  });
</script>

<div style="width: 100%;overflow-y: auto">
  {#if loading}
    <Spinner color="gray" />
  {/if}
  <Table hoverable={true} noborder={true}>
    <TableHead class="cursor-pointer">
      <TableHeadCell
        padding={tablePadding}
        on:click={() => {
          switchSort("cvss_v3_score");
        }}
        >CVSS<i
          class:bx={true}
          class:bx-caret-up={orderBy === "cvss_v3_score"}
          class:bx-caret-down={orderBy === "-cvss_v3_score"}
        ></i></TableHeadCell
      >
      <TableHeadCell
        padding={tablePadding}
        on:click={() => {
          switchSort("publisher");
        }}
        >Publisher<i
          class:bx={true}
          class:bx-caret-up={orderBy === "publisher"}
          class:bx-caret-down={orderBy === "-publisher"}
        ></i></TableHeadCell
      >
      <TableHeadCell
        padding={tablePadding}
        on:click={() => {
          switchSort("title");
        }}
        >Title<i
          class:bx={true}
          class:bx-caret-up={orderBy === "title"}
          class:bx-caret-down={orderBy === "-title"}
        ></i></TableHeadCell
      >
      <TableHeadCell
        padding={tablePadding}
        on:click={() => {
          switchSort("tracking_id");
        }}
        >Tracking ID<i
          class:bx={true}
          class:bx-caret-up={orderBy === "tracking_id"}
          class:bx-caret-down={orderBy === "-tracking_id"}
        ></i></TableHeadCell
      >
      <TableHeadCell
        padding={tablePadding}
        on:click={() => {
          switchSort("version");
        }}
        >Version<i
          class:bx={true}
          class:bx-caret-up={orderBy === "version"}
          class:bx-caret-down={orderBy === "-version"}
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
          <TableBodyCell {tdClass}
            ><span class:text-red-500={Number(item.cvss_v3_score) > 5.0}
              >{item.cvss_v3_score == null ? "" : item.cvss_v3_score}</span
            ></TableBodyCell
          >
          <TableBodyCell {tdClass}>{item.publisher}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.title}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.tracking_id}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.version}</TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
</div>
<div class="mb-12 mt-3 flex items-center">
  {#if documents.length > 0}
    <div class="flex flex-grow items-center">
      <Label class="mr-3">Items per page</Label>
      <Select
        id="pagecount"
        class="mt-2 w-24"
        items={[
          { name: "10", value: 10 },
          { name: "25", value: 25 },
          { name: "50", value: 50 },
          { name: "100", value: 100 }
        ]}
        bind:value={limit}
        on:change={() => {
          offset = 0;
          currentPage = 1;
          fetchData();
        }}
      ></Select>
    </div>
    <div class="mr-3 flex-grow">
      <div class="flex">
        {#if currentPage > 1}
          <PaginationItem on:click={first}>
            <i class="bx bx-arrow-to-left"></i>
          </PaginationItem>
          <PaginationItem on:click={previous}>
            <i class="bx bx-chevrons-left"></i>
          </PaginationItem>
        {/if}
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
          <span>of {numberOfPages} Pages</span>
        </div>
        {#if currentPage !== numberOfPages}
          <PaginationItem on:click={next}>
            <i class="bx bx-chevrons-right"></i>
          </PaginationItem>
          <PaginationItem on:click={last}>
            <i class="bx bx-arrow-to-right"></i>
          </PaginationItem>
        {/if}
      </div>
    </div>
  {/if}
  <div class="mr-3">
    {#if searchTerm}
      {count} entries found
    {:else}
      {count} entries in total
    {/if}
  </div>
</div>
