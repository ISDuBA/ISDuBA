<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  /* eslint-disable svelte/no-at-html-tags */
  import { onMount } from "svelte";
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
  import { tdClass, tablePadding, title, publisher } from "$lib/table/defaults";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Spinner } from "flowbite-svelte";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";
  import { request } from "$lib/utils";

  let openRow: number | null;

  const toggleRow = (i: number) => {
    openRow = openRow === i ? null : i;
  };
  let limit = 10;
  let offset = 0;
  let count = 0;
  let currentPage = 1;
  let documents: any = [];
  let searchTerm: string = "";
  let columns = [
    "id",
    "tracking_id",
    "version",
    "publisher",
    "current_release_date",
    "initial_release_date",
    "title",
    "tlp",
    "cvss_v2_score",
    "cvss_v3_score",
    "four_cves"
  ];
  let orderBy = "-cvss_v3_score";
  let loading = false;
  let error: string;

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

  const fetchData = async () => {
    const searchSuffix = searchTerm ? `query="${searchTerm}" german search msg as &` : "";
    const searchColumn = searchTerm ? " msg" : "";
    const documentURL = encodeURI(
      `/api/documents?${searchSuffix}count=1&order=${orderBy}&limit=${limit}&offset=${offset}&columns=${columns.join(" ")}${searchColumn}`
    );
    loading = true;
    error = "";
    const response = await request(documentURL, "GET");
    if (response.ok) {
      ({ count, documents } = response.content);
      documents = documents || [];
    } else if (response.error) {
      error = response.error;
    }
    loading = false;
  };

  onMount(async () => {
    if ($appStore.app.keycloak.authenticated) {
      fetchData();
    }
  });
</script>

<svelte:head>
  <title>Documents</title>
</svelte:head>

<SectionHeader title="Documents"></SectionHeader>
{#if documents}
  <div class="mb-3 w-2/3">
    <Search
      bind:value={searchTerm}
      on:keyup={(e) => {
        if (e.key === "Enter") fetchData();
      }}
    >
      {#if searchTerm}
        <button
          class="mr-3"
          on:click={() => {
            searchTerm = "";
            fetchData();
          }}>x</button
        >
      {/if}
      <Button
        on:click={() => {
          fetchData();
        }}>Search</Button
      >
    </Search>
  </div>
  <div class="mb-2 mt-8 flex items-center justify-between">
    {#if documents.length > 0}
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
                currentPage = Math.floor(parseInt(currentPage));
                if (currentPage < 1) currentPage = 1;
                if (currentPage > numberOfPages) currentPage = numberOfPages;
                offset = (currentPage - 1) * limit;
                fetchData();
              }}
              bind:value={currentPage}
            />
            <span>of {numberOfPages} pages</span>
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
    {/if}
    <div class="mr-3">
      {#if searchTerm}
        {count} entries found
      {:else}
        {count} entries in total
      {/if}
    </div>
  </div>
  <ErrorMessage message={error}></ErrorMessage>
  <div class:invisible={!loading} class:mb-4={true}>
    Loading ...
    <Spinner color="gray" size="4"></Spinner>
  </div>
  <Table hoverable={true} noborder={true}>
    <TableHead class="cursor-pointer">
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("cvss_v3_score")}
        >CVSS3<i
          class:bx={true}
          class:bx-caret-up={orderBy == "cvss_v3_score"}
          class:bx-caret-down={orderBy == "-cvss_v3_score"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("cvss_v2_score")}
        >CVSS2<i
          class:bx={true}
          class:bx-caret-up={orderBy == "cvss_v2_score"}
          class:bx-caret-down={orderBy == "-cvss_v2_score"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding}>CVEs</TableHeadCell>
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("publisher")}
        >Publisher<i
          class:bx={true}
          class:bx-caret-up={orderBy == "publisher"}
          class:bx-caret-down={orderBy == "-publisher"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("title")}
        >Title<i
          class:bx={true}
          class:bx-caret-up={orderBy == "title"}
          class:bx-caret-down={orderBy == "-title"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("tracking_id")}
        >Tracking ID<i
          class:bx={true}
          class:bx-caret-up={orderBy == "tracking_id"}
          class:bx-caret-down={orderBy == "tracking_id"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("initial_release_date")}
        >Initial Release<i
          class:bx={true}
          class:bx-caret-up={orderBy == "initial_release_date"}
          class:bx-caret-down={orderBy == "-initial_release_date"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("current_release_date")}
        >Current Release<i
          class:bx={true}
          class:bx-caret-up={orderBy == "current_release_date"}
          class:bx-caret-down={orderBy == "-current_release_date"}
        ></i></TableHeadCell
      >
      <TableHeadCell padding={tablePadding} on:click={() => switchSort("version")}
        >Version<i
          class:bx={true}
          class:bx-caret-up={orderBy == "version"}
          class:bx-caret-down={orderBy == "-version"}
        ></i></TableHeadCell
      >
    </TableHead>
    <TableBody>
      {#each documents as item, i}
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
            ><span class:text-red-500={Number(item.cvss_v2_score) > 5.0}
              >{item.cvss_v2_score == null ? "" : item.cvss_v2_score}</span
            ></TableBodyCell
          >
          <TableBodyCell {tdClass}
            >{#if item.four_cves[0]}
              <!-- svelte-ignore a11y-click-events-have-key-events -->
              <!-- svelte-ignore a11y-no-static-element-interactions -->
              {#if item.four_cves.length > 1}
                <div class="mr-2 flex">
                  <div class="flex-grow">
                    {item.four_cves[0]}
                  </div>
                  <span on:click|stopPropagation={() => toggleRow(i)}>
                    {#if openRow === i}
                      <i class="bx bx-minus"></i>
                    {:else}
                      <i class="bx bx-plus"></i>
                    {/if}
                  </span>
                </div>
              {:else}
                <span>{item.four_cves[0]}</span>
              {/if}
            {/if}</TableBodyCell
          >
          <TableBodyCell tdClass={publisher}
            ><span title={item.publisher}>{item.publisher}</span></TableBodyCell
          >
          <TableBodyCell tdClass={title}><span title={item.title}>{item.title}</span></TableBodyCell
          >
          <TableBodyCell {tdClass}>{item.tracking_id}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.initial_release_date.split("T")[0]}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.current_release_date.split("T")[0]}</TableBodyCell>
          <TableBodyCell {tdClass}>{item.version}</TableBodyCell>
        </TableBodyRow>
        {#if openRow === i}
          <TableBodyRow>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}>
              <div>
                {#each item.four_cves as cve, i}
                  {#if i !== 0}
                    <div>{cve}</div>
                  {/if}
                {/each}
              </div>
            </TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
          </TableBodyRow>
        {/if}
        {#if item.msg}
          <TableBodyRow class="border border-indigo-500/100 bg-slate-100">
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}>{@html item.msg}</TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
            <TableBodyCell {tdClass}></TableBodyCell>
          </TableBodyRow>
        {/if}
      {/each}
    </TableBody>
  </Table>
{/if}
