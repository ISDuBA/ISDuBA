<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    saveAggregator,
    fetchAggregatorData,
    fetchAggregators,
    type Aggregator
  } from "$lib/Sources/source";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Input, Spinner, Label, Button, TableBodyCell } from "flowbite-svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass } from "$lib/Table/defaults";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";
  import { type CSAFAggregator } from "$lib/aggregatorTypes";
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";

  let loadingAggregators: boolean = false;
  let aggregators: Aggregator[] = [];
  let aggregatorData = new Map<number, [string, string][]>();

  let aggregatorError: ErrorDetails | null;
  let aggregatorSaveError: ErrorDetails | null;

  let validUrl: boolean | null = null;
  let urlColor: "red" | "green" | "base" = "base";
  $: if (validUrl !== undefined) {
    if (validUrl === null) {
      urlColor = "base";
    } else if (validUrl) {
      urlColor = "green";
    } else {
      urlColor = "red";
    }
  }

  let aggregator: Aggregator = {
    name: "",
    url: ""
  };

  let formClass = "max-w-[800pt]";

  const checkUrl = async () => {
    if (aggregator.url === "") {
      validUrl = null;
      return;
    }
    if (aggregator.url.startsWith("https://") && aggregator.url.endsWith("aggregator.json")) {
      validUrl = null;
      return;
    }
    validUrl = false;
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

  const submitAggregator = async () => {
    let result = await saveAggregator(aggregator);
    if (!result.ok) {
      aggregatorSaveError = result.error;
    } else {
      await getAggregators();
    }
  };
  onMount(async () => {
    await getAggregators();
  });
</script>

<svelte:head>
  <title>Sources - Aggregator</title>
</svelte:head>

<div>
  <SectionHeader title="Aggregator"></SectionHeader>
  <CustomTable
    title="Aggregator"
    headers={[
      { label: "", attribute: "expand" },
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
      {@const list = aggregatorData.get(aggregator.id ?? -1) ?? []}
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
        <TableBodyCell {tdClass}>
          {#if list.length === 0}
            <i class="bx bx-plus"></i>
          {:else}
            <i class="bx bx-minus"></i>
          {/if}
        </TableBodyCell>
        <TableBodyCell {tdClass}>{aggregator.name}</TableBodyCell>
        <TableBodyCell {tdClass}>{aggregator.url}</TableBodyCell>
      </tr>
      {#each list as entry}
        <tr class="bg-slate-200">
          <TableBodyCell {tdClass}></TableBodyCell>
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
      <ErrorMessage error={aggregatorError}></ErrorMessage>
    </div>
  </CustomTable>
  {#if appStore.isSourceManager()}
    <form on:submit={submitAggregator} class={formClass}>
      <Label>Name</Label>
      <Input bind:value={aggregator.name}></Input>
      <Label>URL</Label>
      <Input bind:value={aggregator.url} on:input={checkUrl} color={urlColor}></Input>
      <br />
      <Button
        type="submit"
        color="light"
        disabled={validUrl === false || aggregator.name === "" || aggregator.url === ""}
      >
        <i class="bx bx-check me-2"></i>
        <span>Save aggregator</span>
      </Button>
    </form>
  {/if}
  <ErrorMessage error={aggregatorSaveError}></ErrorMessage>

  <br />
</div>
