<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { Button, Search } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import AdvisoryTable from "$lib/Advisories/AdvisoryTable.svelte";

  let searchTerm: string | null;
  let advisoryTable: any;
  let columns = [
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
    "version"
  ];
  onMount(async () => {
    let savedSearch = sessionStorage.getItem("advisorySearchTerm");
    searchTerm = savedSearch ?? "";
  });
</script>

<svelte:head>
  <title>Advisories</title>
</svelte:head>

<SectionHeader title="Advisories"></SectionHeader>
<div class="mb-3 w-2/3">
  <Search
    bind:value={searchTerm}
    on:keyup={(e) => {
      sessionStorage.setItem("advisorySearchTerm", searchTerm ?? "");
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
      class="py-3.5"
      on:click={() => {
        advisoryTable.fetchData();
      }}>Search</Button
    >
  </Search>
</div>
{#if searchTerm !== null}
  <AdvisoryTable {searchTerm} bind:this={advisoryTable} loadAdvisories={true} {columns}
  ></AdvisoryTable>
{/if}
