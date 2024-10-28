<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Spinner, TableBodyCell } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { push } from "svelte-spa-router";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails, getErrorDetails } from "$lib/Errors/error";
  import { tdClass } from "$lib/Table/defaults";
  import { request } from "$lib/request";
  import { onDestroy, onMount } from "svelte";
  import { type CSAFAggregator } from "$lib/aggregatorTypes";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import {
    type Source,
    type Aggregator,
    fetchSources,
    fetchAggregators,
    fetchAggregatorData
  } from "$lib/Sources/source";
  import { appStore } from "$lib/store";

  let messageError: ErrorDetails | null;
  let sourcesError: ErrorDetails | null;
  let aggregatorError: ErrorDetails | null;

  let loadingSources: boolean = false;
  let loadingAggregators: boolean = false;

  let sources: Source[] = [];
  let aggregators: Aggregator[] = [];
  let aggregatorData = new Map<number, [string, string][]>();

  async function getMessage() {
    const response = await request("api/sources/message", "GET");
    if (response.ok) {
      return response.content;
    } else {
      messageError = getErrorDetails(`Couldn't load default message`, response);
    }
    return new Map<string, [string]>();
  }

  let sourceUpdate = setInterval(async () => {
    if (appStore.isEditor() || appStore.isSourceManager()) {
      getSources();
    }
  }, 30 * 1000);

  const getSources = async () => {
    loadingSources = true;
    const result = await fetchSources(true);
    loadingSources = false;
    if (result.ok) {
      sources = result.value;
    } else {
      sources = [];
      sourcesError = result.error;
    }
  };

  const getAggregators = async () => {
    loadingAggregators = true;
    const result = await fetchAggregators();
    loadingAggregators = false;
    if (result.ok) {
      aggregators = result.value;
    } else {
      aggregatorError = result.error;
    }
  };

  const parseAggregatorData = (data: CSAFAggregator): [string, string][] => {
    return data.csaf_providers.map((i: any) => [i.metadata.publisher.name, i.metadata.url]);
  };

  const toggleAggregatorView = async (aggregator: Aggregator) => {
    if (aggregatorData.get(aggregator.id ?? -1)) {
      aggregatorData.delete(aggregator.id ?? -1);
      aggregatorData = aggregatorData;
      return;
    }
    const resp = await fetchAggregatorData(aggregator.url);
    if (resp.ok) {
      aggregatorData.set(aggregator.id ?? -1, parseAggregatorData(resp.value.aggregator));
      aggregatorData = aggregatorData;
    } else {
      aggregatorError = resp.error;
    }
  };

  onMount(async () => {
    if (appStore.isEditor() || appStore.isSourceManager()) {
      await getSources();
      await getAggregators();
    }
  });

  onDestroy(() => {
    clearInterval(sourceUpdate);
  });
</script>

<svelte:head>
  <title>Sources</title>
</svelte:head>

<div>
  {#if appStore.isEditor() || appStore.isSourceManager()}
    <SectionHeader title="Sources"></SectionHeader>
    <CustomTable
      title="CSAF Provider"
      headers={[
        {
          label: "Name",
          attribute: "name"
        },
        {
          label: "Domain/PMD",
          attribute: "url"
        },
        {
          label: "Active",
          attribute: "active"
        },
        {
          label: "Loading/Queued",
          attribute: "stats"
        }
      ]}
    >
      {#each sources as source, index (index)}
        <tr
          on:click={() => {
            if (appStore.isSourceManager()) {
              push(`/sources/${source.id}`);
            }
          }}
          on:blur={() => {}}
          on:focus={() => {}}
          class={appStore.isSourceManager() ? "cursor-pointer" : ""}
        >
          <TableBodyCell {tdClass}>{source.name}</TableBodyCell>
          <TableBodyCell {tdClass}>{source.url}</TableBodyCell>
          <TableBodyCell {tdClass}
            ><i class={"bx " + (source.active ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
          >
          <TableBodyCell {tdClass}
            >{source.stats?.downloading}/{source.stats?.waiting}</TableBodyCell
          >
        </tr>
      {/each}
      <div slot="bottom">
        <div
          class:invisible={!loadingSources}
          class={loadingSources ? "loadingFadeIn" : ""}
          class:mb-4={true}
        >
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
        {#if appStore.isSourceManager()}
          <Button href="/#/sources/new" class="mb-2" color="primary" size="xs">
            <i class="bx bx-plus"></i>
            <span>Add source</span>
          </Button>
        {/if}
        <ErrorMessage error={sourcesError}></ErrorMessage>
      </div>
    </CustomTable>
    <CustomTable
      title="Aggregator"
      headers={[
        {
          label: "Name",
          attribute: "name"
        },
        {
          label: "URL",
          attribute: "url"
        }
      ]}
    >
      {#each aggregators as aggregator, index (index)}
        <tr
          on:click={async () => {
            if (appStore.isSourceManager()) {
              await toggleAggregatorView(aggregator);
            }
          }}
          on:blur={() => {}}
          on:focus={() => {}}
          class={appStore.isSourceManager() ? "cursor-pointer" : ""}
        >
          <TableBodyCell {tdClass}>{aggregator.name}</TableBodyCell>
          <TableBodyCell {tdClass}>{aggregator.url}</TableBodyCell>
        </tr>
        {@const list = aggregatorData.get(aggregator.id ?? -1) ?? []}
        {#each list as entry}
          <tr class="bg-slate-200">
            <TableBodyCell {tdClass}>{entry[0]}</TableBodyCell>
            <TableBodyCell {tdClass}>{entry[1]}</TableBodyCell>
          </tr>
        {/each}
      {/each}
      <div slot="bottom">
        <div
          class:invisible={!loadingAggregators}
          class={loadingAggregators ? "loadingFadeIn" : ""}
          class:mb-4={true}
        >
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
        {#if appStore.isSourceManager()}
          <Button href="/#/aggregator/new" class="mb-2" color="primary" size="xs">
            <i class="bx bx-plus"></i>
            <span>Add aggregator</span>
          </Button>
        {/if}
        <ErrorMessage error={aggregatorError}></ErrorMessage>
      </div>
    </CustomTable>
  {/if}
  {#await getMessage() then resp}
    {#if resp.message}
      {resp.message}
    {/if}
  {/await}
  <ErrorMessage error={sourcesError}></ErrorMessage>
  <ErrorMessage error={messageError}></ErrorMessage>

  <br />
  {#if appStore.isImporter()}
    <Button href="/#/sources/upload" class="my-2" color="primary" size="xs">
      <i class="bx bx-upload"></i>
      <span>Upload documents</span>
    </Button>
  {/if}
</div>
