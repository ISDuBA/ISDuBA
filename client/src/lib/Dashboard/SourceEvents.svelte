<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";
  import {
    type Attention,
    fetchSourceAttentionList,
    fetchAggregatorAttentionList
  } from "$lib/Sources/source";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails } from "$lib/Errors/error";
  import Activity from "./Activity.svelte";
  import { Button, Spinner } from "flowbite-svelte";
  import { push } from "svelte-spa-router";

  interface MergedAttention extends Attention {
    isSource: boolean;
  }

  let attentionCount = 0;
  let attentions: MergedAttention[] = [];
  let loadAttentionError: ErrorDetails | null;
  let isLoading = false;

  const loadAttentionList = async () => {
    let sourceResult = await fetchSourceAttentionList();
    if (sourceResult.ok) {
      attentions.push(
        ...sourceResult.value.map((i) => {
          return { ...i, isSource: true };
        })
      );
    } else {
      loadAttentionError = sourceResult.error;
    }
    let aggregatorResult = await fetchAggregatorAttentionList();
    if (aggregatorResult.ok) {
      attentions.push(
        ...aggregatorResult.value.map((i) => {
          return { ...i, isSource: false };
        })
      );
    } else {
      loadAttentionError = aggregatorResult.error;
    }
    attentionCount = attentions.length;
    attentions = attentions;
  };

  onMount(async () => {
    isLoading = true;
    await loadAttentionList();
    isLoading = false;
  });
</script>

{#if $appStore.app.isUserLoggedIn && appStore.isSourceManager()}
  <div class="flex flex-col gap-4 md:w-[46%] md:max-w-[46%]">
    <SectionHeader title="Changed sources"></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if isLoading}
        <div class:invisible={!isLoading} class={isLoading ? "loadingFadeIn" : ""}>
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
      {/if}
      {#if attentions}
        {#if attentions.length > 0}
          {#each attentions as attention}
            <Activity
              on:click={() => {
                if (attention.id) {
                  if (attention.isSource) {
                    push(`/sources/${attention.id}`);
                  } else {
                    push(`/sources/aggregators/${attention.id}`);
                  }
                }
              }}
            >
              <div slot="top-left">
                {attention.isSource ? "Source change" : "Aggregator change"}
              </div>
              <div>{attention.name}</div>
            </Activity>
          {/each}
        {:else}
          <div class="text-gray-600">No source changes found.</div>
        {/if}
      {/if}
      <Button
        on:click={async () => await push(`/sources/`)}
        color="light"
        class="h-fit w-fit rounded-md !px-2 !py-1"
      >
        <i class="bx bx-spreadsheet text-lg"></i>
      </Button>
      {#if attentionCount > 10}<div class="">…There are more events</div>{/if}
    </div>
    <ErrorMessage error={loadAttentionError}></ErrorMessage>
  </div>
{/if}
