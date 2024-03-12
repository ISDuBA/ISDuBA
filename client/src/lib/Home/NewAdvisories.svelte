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
    Label,
    Select,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch,
    Pagination
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";

  const previous = () => {
    if (offset - limit >= 0) {
      offset = offset - limit > 0 ? offset - limit : 0;
    }
    fetchData();
  };
  const next = () => {
    if (offset + limit <= count) {
      offset = offset + limit;
    }
    fetchData();
  };
  const handleClick = () => {
    alert("Page clicked");
  };

  let orderBy = "title";
  let limit = 10;
  let offset = 0;
  let count = 0;
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

  $: pages = Array.from({ length: Math.ceil(count / limit) }, (_, v) => v + 1).map((x) => {
    return { name: x };
  });
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
        fetchData();
      }}
    ></Select>
  </div>
  <div class="mr-3 flex-grow">
    <Pagination {pages} on:previous={previous} on:next={next} on:click={handleClick} />
  </div>
  <div class="mr-3">
    {count} entries in total
  </div>
</div>
