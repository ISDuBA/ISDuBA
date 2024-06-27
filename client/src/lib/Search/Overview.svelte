<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { Button, ButtonGroup, Search, Toggle } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import AdvisoryTable from "$lib/table/Table.svelte";
  import { SEARCHPAGECOLUMNS } from "$lib/Queries/query";

  let searchTerm: string | null;
  let advisoryTable: any;
  let advisoriesOnly = true;
  let advancedSearch = false;

  let defaultOrderBy = "cvss_v3_score";

  $: columns = advisoriesOnly ? SEARCHPAGECOLUMNS.ADVISORY : SEARCHPAGECOLUMNS.DOCUMENT;

  onMount(async () => {
    let savedSearch = sessionStorage.getItem("documentSearchTerm");
    searchTerm = savedSearch ?? "";
  });
</script>

<svelte:head>
  <title>Search</title>
</svelte:head>

<div class="flex flex-row">
  <SectionHeader title="Search"></SectionHeader>
</div>
<hr class="mb-6" />

<div class="mb-3 flex">
  <div class="flex w-2/3 flex-row">
    <Search
      size="sm"
      placeholder={advancedSearch ? "Enter a query" : "Enter a searchterm"}
      bind:value={searchTerm}
      on:keyup={(e) => {
        sessionStorage.setItem("documentSearchTerm", searchTerm ?? "");
        if (e.key === "Enter") advisoryTable.fetchData();
      }}
    >
      {#if searchTerm}
        <button
          class="mr-3"
          on:click={() => {
            searchTerm = "";
            advisoryTable.fetchData();
            sessionStorage.setItem("documentSearchTerm", "");
          }}>x</button
        >
      {/if}
      <Button
        size="xs"
        class="h-7 py-3.5"
        on:click={() => {
          advisoryTable.fetchData();
        }}>{advancedSearch ? "Apply" : "Search"}</Button
      >
    </Search>
    <Toggle bind:checked={advancedSearch} class="ml-3">Advanced</Toggle>
  </div>
  <ButtonGroup class="ml-auto h-7">
    <Button
      size="xs"
      color="light"
      class={`h-7 py-1 text-xs ${advisoriesOnly ? "bg-gray-200 hover:bg-gray-100" : ""}`}
      on:click={() => {
        advisoriesOnly = true;
      }}>Advisories</Button
    >
    <Button
      size="xs"
      color="light"
      class={`h-7 py-1 text-xs ${!advisoriesOnly ? "bg-gray-200 hover:bg-gray-100" : ""}`}
      on:click={() => {
        advisoriesOnly = false;
      }}>Documents</Button
    >
  </ButtonGroup>
</div>
{#if searchTerm !== null}
  <AdvisoryTable
    {defaultOrderBy}
    {searchTerm}
    bind:this={advisoryTable}
    loadAdvisories={advisoriesOnly}
    {columns}
  ></AdvisoryTable>
{/if}
