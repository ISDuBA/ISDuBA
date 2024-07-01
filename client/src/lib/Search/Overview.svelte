<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount, tick } from "svelte";
  import { Button, ButtonGroup, Search, Toggle } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import AdvisoryTable from "$lib/Table/Table.svelte";
  import { SEARCHPAGECOLUMNS } from "$lib/Queries/query";
  import Queries from "./Queries.svelte";

  let searchTerm: string | null;
  let advisoryTable: any;
  let advancedSearch = false;
  let selectedCustomQuery: any;

  const resetQuery = () => {
    return {
      columns: [...SEARCHPAGECOLUMNS.ADVISORY],
      advisories: true,
      orders: ["cvss_v3_score"],
      query: ""
    };
  };

  let query = resetQuery();

  const setQueryBack = async () => {
    query = resetQuery();
    searchTerm = "";
    sessionStorage.setItem("documentSearchTerm", "");
    await tick();
    advisoryTable.fetchData();
  };

  $: if (selectedCustomQuery === -1) {
    setQueryBack();
  }

  const triggerSearch = async () => {
    if (!advancedSearch) {
      query.query = searchTerm ? `"${searchTerm}" german search msg as` : "";
      if (
        searchTerm &&
        !query.columns.find((c) => {
          return c === "msg";
        })
      )
        query.columns.push("msg");
    }
    await tick();
    advisoryTable.fetchData();
  };

  onMount(async () => {});
</script>

<svelte:head>
  <title>Search</title>
</svelte:head>

<div class="flex flex-row">
  <SectionHeader title="Search"></SectionHeader>
</div>
<hr class="mb-6" />
<Queries
  on:querySelected={async (e) => {
    let { detail } = e;
    query = {
      query: detail.query,
      columns: detail.columns,
      advisories: detail.advisories,
      orders: detail.orders || []
    };
    searchTerm = "";
    await tick();
    advisoryTable.fetchData();
  }}
  bind:selectedIndex={selectedCustomQuery}
></Queries>
<div class="mb-3 flex">
  <div class="flex w-2/3 flex-row">
    <Search
      size="sm"
      placeholder={advancedSearch ? "Enter a query" : "Enter a searchterm"}
      bind:value={searchTerm}
      on:keyup={(e) => {
        sessionStorage.setItem("documentSearchTerm", searchTerm ?? "");
        if (e.key === "Enter") triggerSearch();
      }}
    >
      {#if searchTerm}
        <button
          class="mr-3"
          on:click={async () => {
            searchTerm = "";
            query.query = "";
            query.columns = query.columns.filter((c) => {
              return c !== "msg";
            });
            await tick();
            advisoryTable.fetchData();
            sessionStorage.setItem("documentSearchTerm", "");
          }}>x</button
        >
      {/if}
      <Button
        size="xs"
        class="h-7 py-3.5"
        on:click={() => {
          triggerSearch();
        }}>{advancedSearch ? "Apply" : "Search"}</Button
      >
    </Search>
    <Toggle bind:checked={advancedSearch} class="ml-3">Advanced</Toggle>
  </div>
  {#if selectedCustomQuery === -1}
    <ButtonGroup class="ml-auto h-7">
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${query.advisories ? "bg-gray-200 hover:bg-gray-100" : ""}`}
        on:click={() => {
          query.advisories = true;
          query.columns = SEARCHPAGECOLUMNS.ADVISORY;
        }}>Advisories</Button
      >
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${!query.advisories ? "bg-gray-200 hover:bg-gray-100" : ""}`}
        on:click={() => {
          query.advisories = false;
          query.columns = SEARCHPAGECOLUMNS.DOCUMENT;
        }}>Documents</Button
      >
    </ButtonGroup>
  {/if}
</div>
{#if searchTerm !== null}
  <AdvisoryTable
    defaultOrderBy={query.orders[0]}
    columns={query.columns}
    loadAdvisories={query.advisories}
    query={`${query.query}`}
    orderBy={query.orders?.join(" ") ?? ""}
    bind:this={advisoryTable}
  ></AdvisoryTable>
{/if}
