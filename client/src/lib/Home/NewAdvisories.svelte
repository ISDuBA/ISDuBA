<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { push } from "svelte-spa-router";
  import {
    Label,
    PaginationItem,
    Select,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Table
  } from "flowbite-svelte";
  import { tdClass, tablePadding, title } from "$lib/table/defaults";
  import { onMount } from "svelte";
  import { Spinner } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";

  let limit = 10;
  let offset = 0;
  let count = 0;
  let currentPage = 1;
  let documents: any = null;
  let searchTerm: string = "";
  let loading = false;
  let error: string;
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
  const fetchData = async () => {
    const searchSuffix = searchTerm ? ` "${searchTerm}" german search msg as and` : "";
    const documentURL = encodeURI(
      `/api/documents?query=$state new workflow =${searchSuffix}&advisories=true&count=1&order=${orderBy}&limit=${limit}&offset=${offset}&columns=${columns.join(" ")}`
    );
    error = "";
    loading = true;
    const response = await request(documentURL, "GET");
    console.log(response);
    if (response.ok) {
      ({ count, documents } = response.content);
      documents = documents || [];
    } else if (response.error) {
      error = response.error;
    }
    loading = false;
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
    offset = (numberOfPages - 1) * limit;
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

  $: numberOfPages = Math.ceil(count / limit);
  $: onMount(async () => {
    await fetchData();
  });
</script>

<div>
  <div class="mb-2 mt-2 flex items-center justify-between">
    {#if documents?.length > 0}
      <div class="flex items-center">
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
      <div>
        <div class="flex">
          <div class:invisible={currentPage === 1} class:flex={true}>
            <PaginationItem on:click={first}>
              <i class="bx bx-arrow-to-left"></i>
            </PaginationItem>
            <PaginationItem on:click={previous}>
              <i class="bx bx-chevrons-left"></i>
            </PaginationItem>
          </div>
          <div class="mx-3 flex items-center">
            <input
              class="mr-1 w-16 cursor-pointer border pr-1 text-right"
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
          <div class:invisible={currentPage === numberOfPages} class:flex={true}>
            <PaginationItem on:click={next}>
              <i class="bx bx-chevrons-right"></i>
            </PaginationItem>
            <PaginationItem on:click={last}>
              <i class="bx bx-arrow-to-right"></i>
            </PaginationItem>
          </div>
        </div>
      </div>
      <div class="mr-3">
        {#if searchTerm}
          {count} entries found
        {:else}
          {count} entries in total
        {/if}
      </div>
    {/if}
  </div>
  <div class:invisible={!loading} class:mb-4={true}>
    Loading ...
    <Spinner color="gray" size="4"></Spinner>
  </div>
  <ErrorMessage message={error}></ErrorMessage>
  {#if documents?.length > 0}
    <div class="w-fit">
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
              <TableBodyCell {tdClass}
                ><span title={item.publisher}>{item.publisher}</span></TableBodyCell
              >
              <TableBodyCell style="max-width: 48rem;" tdClass={title}
                ><span title={item.title}>{item.title}</span></TableBodyCell
              >
              <TableBodyCell {tdClass}>{item.tracking_id}</TableBodyCell>
              <TableBodyCell {tdClass}>{item.version}</TableBodyCell>
            </TableBodyRow>
          {/each}
        </TableBody>
      </Table>
    </div>
  {/if}
</div>
