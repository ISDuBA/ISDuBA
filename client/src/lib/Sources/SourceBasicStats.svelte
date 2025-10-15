<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { fetchBasicStatistic, type StatisticGroup } from "$lib/Statistics/statistics";
  import { Spinner } from "flowbite-svelte";
  import type { Result } from "$lib/types";
  import type { ErrorDetails } from "$lib/Errors/error";

  interface Props {
    from?: Date;
    sourceID: number;
  }

  let { from = new Date(0), sourceID }: Props = $props();

  let oldResponse: Result<StatisticGroup, ErrorDetails> | null = $state(null);
  export const reload = async () => {
    promise = fetch();
  };
  const fetch = async () => {
    oldResponse = await fetchBasicStatistic(
      from,
      new Date(),
      Date.now() - from.getTime(),
      "imports",
      sourceID
    );
    return oldResponse;
  };
  let promise = $state(fetch());
</script>

{#await promise}
  {#if oldResponse}
    {#if oldResponse.ok}
      {oldResponse.value.imports?.[0][1] ?? 0}
    {:else}
      <span class="text-red-700">Couldn't load value.</span>
    {/if}
  {:else}
    <Spinner color="gray" size="4"></Spinner>
  {/if}
{:then response}
  {#if response.ok}
    {response.value.imports?.[0][1] ?? 0}
  {:else}
    <span class="text-red-700">Couldn't load value.</span>
  {/if}
{/await}
