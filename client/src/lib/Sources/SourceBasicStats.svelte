<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import {
    fetchBasicStatistic,
    fetchImportFailuresStatistic,
    mergeImportFailureStatistics
  } from "$lib/Statistics/statistics";
  import { Spinner } from "flowbite-svelte";

  export let from: Date = new Date(0);
  export let sourceID: number;
</script>

{#await fetchBasicStatistic(from, new Date(), Date.now() - from.getTime(), "imports", sourceID)}
  <Spinner color="gray" size="4"></Spinner>
{:then response}
  {#if response.ok}
    {response.value.imports?.[0][1] ?? 0}
  {:else}
    <span class="text-red-700">Couldn't load value.</span>
  {/if}
{/await}
/
{#await fetchImportFailuresStatistic(from, new Date(), Date.now() - from.getTime(), sourceID)}
  <Spinner color="gray" size="4"></Spinner>
{:then response}
  {#if response.ok}
    {@const merged = mergeImportFailureStatistics(response.value)}
    {merged.importFailuresCombined?.[0][1]}
  {:else}
    <span class="text-red-700">Couldn't load value.</span>
  {/if}
{/await}
