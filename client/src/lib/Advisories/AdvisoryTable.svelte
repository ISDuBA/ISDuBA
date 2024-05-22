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
    TableHead,
    TableHeadCell,
    Table
  } from "flowbite-svelte";
  import { tdClass, tablePadding, title, publisher } from "$lib/table/defaults";
  import { onMount } from "svelte";
  import { Spinner } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";

  let openRow: number | null;

  const toggleRow = (i: number) => {
    openRow = openRow === i ? null : i;
  };
  let limit = 10;
  let offset = 0;
  let count = 0;
  let currentPage = 1;
  let documents: any = null;
  let loading = false;
  let error: string;
  export let columns: string[];
  export let query: string = "";
  export let searchTerm: string = "";
  export let loadAdvisories: boolean;
  export let orderBy = "title";

  let anchorLink: string | null;

  const getColumnDisplayName = (column: string): string => {
    let names: { [key: string]: string } = {
      id: "ID",
      tracking_id: "TRACKING ID",
      version: "VERSION",
      publisher: "PUBLISHER",
      current_release_date: "CURRENT RELEASE",
      initial_release_date: "INITIAL RELEASE",
      title: "TITLE",
      tlp: "TLP",
      cvss_v2_score: "CVSS2",
      cvss_v3_score: "CVSS3",
      ssvc: "SSVC",
      four_cves: "CVES",
      state: "STATE"
    };

    return names[column] ?? column;
  };

  const calcSSVC = (documents: any) => {
    if (!documents) return [];
    documents.map((d: any) => {
      if (d["ssvc"]) d["ssvc"] = convertVectorToLabel(d["ssvc"]);
    });
    return documents;
  };

  const savePosition = () => {
    let position = [offset, currentPage, limit, orderBy];
    sessionStorage.setItem("tablePosition" + query + loadAdvisories, JSON.stringify(position));
  };

  let postitionRestored: boolean = false;
  const restorePosition = () => {
    let position = sessionStorage.getItem("tablePosition" + query + loadAdvisories);
    if (position) {
      [offset, currentPage, limit, orderBy] = JSON.parse(position);
    }
  };

  let searchTimeout: any = null;

  $: if (searchTerm !== undefined) {
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }
    searchTimeout = setTimeout(() => {
      fetchData();
    }, 500);
  }

  $: if (columns || query || loadAdvisories || !loadAdvisories) {
    fetchData();
  }

  $: if (offset || currentPage || limit || orderBy) {
    if (!postitionRestored) {
      restorePosition();
      postitionRestored = true;
    }
    savePosition();
  }

  export async function fetchData(): Promise<void> {
    const searchSuffix = searchTerm ? `"${searchTerm}" german search msg as ` : "";
    const searchColumn = searchTerm ? " msg" : "";

    let queryParam = "";
    if (query || searchSuffix) {
      queryParam = `query=${query}${searchSuffix}`;
    }
    let fetchColumns = [...columns];
    let requiredColumns = ["id", "tracking_id", "publisher"];
    for (let c of requiredColumns) {
      if (!fetchColumns.includes(c)) {
        fetchColumns.push(c);
      }
    }

    const documentURL = encodeURI(
      `/api/documents?${queryParam}&advisories=${loadAdvisories}&count=1&order=${orderBy}&limit=${limit}&offset=${offset}&columns=${fetchColumns.join(" ")}${searchColumn}`
    );
    error = "";
    loading = true;
    const response = await request(documentURL, "GET");
    if (response.ok) {
      ({ count, documents } = response.content);
      documents = calcSSVC(documents) || [];
    } else if (response.error) {
      error = getErrorMessage(response.error);
    }
    loading = false;
  }

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
    restorePosition();
    postitionRestored = true;
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
    <div class="w-auto">
      <a href={anchorLink}>
        <Table style="w-auto" hoverable={true} noborder={true}>
          <TableHead class="cursor-pointer">
            {#each columns as column}
              <TableHeadCell
                padding={tablePadding}
                on:click={() => {
                  switchSort(column);
                }}
                >{getColumnDisplayName(column)}<i
                  class:bx={true}
                  class:bx-caret-up={orderBy === column}
                  class:bx-caret-down={orderBy === "-" + column}
                ></i></TableHeadCell
              >
            {/each}
          </TableHead>
          <TableBody>
            {#each documents as item, i}
              <tr
                class="cursor-pointer bg-white hover:bg-gray-50 dark:border-gray-700 dark:bg-gray-800 dark:hover:bg-gray-600"
                on:click={() => {
                  push(`/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`);
                }}
                on:mouseenter={() => {
                  anchorLink = `#/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`;
                }}
                on:mouseleave={() => {
                  anchorLink = null;
                }}
              >
                {#each columns as column}
                  {#if column === "cvss_v3_score" || column === "cvss_v2_score"}
                    <TableBodyCell {tdClass}
                      ><span class:text-red-500={Number(item[column]) > 5.0}
                        >{item[column] == null ? "" : item[column]}</span
                      ></TableBodyCell
                    >
                  {:else if column === "ssvc"}
                    <TableBodyCell {tdClass}
                      ><span style={item[column] ? `color:${item[column].color}` : ""}
                        >{item[column]?.label || ""}</span
                      ></TableBodyCell
                    >
                  {:else if column === "state"}
                    <TableBodyCell {tdClass}
                      ><i
                        title={item[column]}
                        class:bx={true}
                        class:bxs-star={item[column] === "new"}
                        class:bx-show={item[column] === "read"}
                        class:bxs-analyse={item[column] === "assessing"}
                        class:bx-book-open={item[column] === "review"}
                        class:bx-archive={item[column] === "archived"}
                        class:bx-trash={item[column] === "delete"}
                      ></i>
                    </TableBodyCell>
                  {:else if column === "initial_release_date" || column === "current_release_date"}
                    <TableBodyCell {tdClass}
                      >{item.initial_release_date?.split("T")[0]}</TableBodyCell
                    >
                  {:else if column === "title"}
                    <TableBodyCell tdClass={title}
                      ><span title={item[column]}>{item[column]}</span></TableBodyCell
                    >
                  {:else if column === "publisher"}
                    <TableBodyCell tdClass={publisher}
                      ><span title={item[column]}>{item[column]}</span></TableBodyCell
                    >
                  {:else if column === "four_cves"}
                    <TableBodyCell {tdClass}
                      >{#if item[column][0]}
                        <!-- svelte-ignore a11y-click-events-have-key-events -->
                        <!-- svelte-ignore a11y-no-static-element-interactions -->
                        {#if item[column] > 1}
                          <div class="mr-2 flex">
                            <div class="flex-grow">
                              {item[column][0]}
                            </div>
                            <span
                              on:mouseenter={() => (anchorLink = null)}
                              on:click|stopPropagation={() => toggleRow(i)}
                            >
                              {#if openRow === i}
                                <i class="bx bx-minus"></i>
                              {:else}
                                <i class="bx bx-plus"></i>
                              {/if}
                            </span>
                          </div>
                        {:else}
                          <span>{item[column][0]}</span>
                        {/if}
                      {/if}</TableBodyCell
                    >
                  {:else}
                    <TableBodyCell {tdClass}>{item[column]}</TableBodyCell>
                  {/if}
                {/each}
              </tr>
            {/each}
          </TableBody>
        </Table>
      </a>
    </div>
  {/if}
</div>
