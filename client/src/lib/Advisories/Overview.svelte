<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { Button, ButtonGroup, Search } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import AdvisoryTable from "$lib/Advisories/Table.svelte";

  let searchTerm: string | null;
  let advisoryTable: any;
  let advisoriesOnly = true;

  $: columns = advisoriesOnly
    ? [
        "critical",
        "cvss_v3_score",
        "cvss_v2_score",
        "ssvc",
        "state",
        "four_cves",
        "publisher",
        "title",
        "tracking_id",
        "initial_release_date",
        "current_release_date",
        "version",
        "comments",
        "recent",
        "versions"
      ]
    : [
        "critical",
        "cvss_v3_score",
        "cvss_v2_score",
        "ssvc",
        "four_cves",
        "publisher",
        "title",
        "tracking_id",
        "initial_release_date",
        "current_release_date",
        "version"
      ];

  onMount(async () => {
    let savedSearch = sessionStorage.getItem("documentSearchTerm");
    searchTerm = savedSearch ?? "";
  });
</script>

<svelte:head>
  <title>Advisories</title>
</svelte:head>

<div class="flex flex-row">
  <SectionHeader title="Advisories"></SectionHeader>
</div>
<hr class="mb-6" />

<div class="mb-3 flex">
  <div class="w-2/3">
    <Search
      size="sm"
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
          }}>x</button
        >
      {/if}
      <Button
        size="xs"
        class="h-7 py-3.5"
        on:click={() => {
          advisoryTable.fetchData();
        }}>Search</Button
      >
    </Search>
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
  <AdvisoryTable {searchTerm} bind:this={advisoryTable} loadAdvisories={advisoriesOnly} {columns}
  ></AdvisoryTable>
{/if}
