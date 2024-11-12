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
    deleteAggregator,
    type Aggregator,
    resetAggregatorAttention
  } from "$lib/Sources/source";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Input, Spinner, Label, Button, TableBodyCell } from "flowbite-svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { tdClass } from "$lib/Table/defaults";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";
  import { type AggregatorMetadata } from "$lib/aggregatorTypes";
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";

  type AggregatorInfo = {
    name: string;
    url: string;
    availableFeeds: string[];
    publisher: boolean;
    expand: boolean;
  };

  let loadingAggregators: boolean = false;
  let aggregators: Aggregator[] = [];
  let aggregatorData = new Map<number, AggregatorInfo[]>();

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
  let validName: boolean | null = null;
  let nameColor: "red" | "green" | "base" = "base";
  $: if (validName !== undefined) {
    if (validName === null) {
      nameColor = "base";
    } else if (validName) {
      nameColor = "green";
    } else {
      nameColor = "red";
    }
  }

  let aggregator: Aggregator = {
    name: "",
    url: ""
  };

  let formClass = "max-w-[800pt]";
  const extraSmallColumnClass = "w-7 max-w-7 min-w-7";
  const smallColumnClass = "w-10 max-w-10 min-w-10";

  const checkUrl = () => {
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

  const checkName = () => {
    if (aggregators.find((i) => i.name === aggregator.name)) {
      validName = false;
      return;
    }
    validName = null;
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

  const parseAggregatorData = (data: AggregatorMetadata): AggregatorInfo[] => {
    const csafProviders = data.aggregator.csaf_providers.map(
      (i) =>
        <AggregatorInfo>{
          name: i.metadata.publisher.name,
          url: i.metadata.url,
          publisher: false
        }
    );
    const csafPublisher =
      data.aggregator.csaf_publishers?.map(
        (i) =>
          <AggregatorInfo>{
            name: i.metadata.publisher.name,
            url: i.metadata.url,
            publisher: true
          }
      ) ?? [];

    const list = [...csafProviders, ...csafPublisher];
    list.forEach((i) => {
      let found = data.custom.subscriptions.find((s) => s.url === i.url);
      if (found) {
        i.availableFeeds = found.available;
      }
    });
    return list;
  };

  const toggleAggregatorView = async (aggregator: Aggregator) => {
    if (aggregatorData.get(aggregator.id ?? -1)) {
      aggregatorData.delete(aggregator.id ?? -1);
      aggregatorData = aggregatorData;
      return;
    }
    loadingAggregators = true;
    const resp = await fetchAggregatorData(aggregator.url);
    loadingAggregators = false;
    if (resp.ok) {
      aggregatorData.set(aggregator.id ?? -1, parseAggregatorData(resp.value));
      let resetResult = await resetAggregatorAttention(aggregator);
      if (resetResult.ok) {
        aggregator.attention = false;
      } else {
        aggregatorError = resetResult.error;
      }
      aggregatorData = aggregatorData;
    } else {
      aggregatorError = resp.error;
    }
  };

  const removeAggregator = async (id: number) => {
    let result = await deleteAggregator(id);
    if (!result.ok) {
      aggregatorError = result.error;
    }
    aggregatorData.delete(id);
    await getAggregators();
  };

  const submitAggregator = async () => {
    let result = await saveAggregator(aggregator);
    if (!result.ok) {
      aggregatorSaveError = result.error;
    } else {
      aggregator.name = "";
      aggregator.url = "";
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
    headers={[
      {
        label: "",
        attribute: "expand",
        class: smallColumnClass
      },
      {
        label: "",
        attribute: "attention",
        class: extraSmallColumnClass
      },
      { label: "", attribute: "delete", class: smallColumnClass },
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
        class={appStore.isSourceManager() ? "cursor-pointer" : ""}
      >
        <TableBodyCell tdClass={`${tdClass} ${smallColumnClass}`}>
          {#if list.length === 0}
            <i class="bx bx-plus"></i>
          {:else}
            <i class="bx bx-minus"></i>
          {/if}
        </TableBodyCell>
        {#if aggregator.attention}
          <TableBodyCell tdClass={`${tdClass} ${extraSmallColumnClass}`}>
            <i class="bx bx-info-square text-lg"></i>
          </TableBodyCell>
        {:else}
          <TableBodyCell tdClass={`${tdClass} ${extraSmallColumnClass}`}></TableBodyCell>
        {/if}
        <TableBodyCell tdClass={`${tdClass} ${smallColumnClass}`}
          ><Button
            on:click={async () => {
              if (aggregator.id) {
                await removeAggregator(aggregator.id);
              }
            }}
            class="!p-2"
            color="light"
          >
            <i class="bx bx-trash text-red-600"></i>
          </Button>
        </TableBodyCell>
        <TableBodyCell {tdClass}>{aggregator.name}</TableBodyCell>
        <TableBodyCell {tdClass}>{aggregator.url}</TableBodyCell>
      </tr>
      {#each list as entry}
        <tr class="bg-slate-100 dark:bg-gray-700" on:click={() => (entry.expand = !entry.expand)}>
          <TableBodyCell {tdClass}>
            {#if entry.expand}
              <i class="bx bx-minus"></i>
            {:else}
              <i class="bx bx-plus"></i>
            {/if}
          </TableBodyCell>

          <TableBodyCell {tdClass}></TableBodyCell>
          <TableBodyCell {tdClass}>
            <Button
              on:click={async () => {
                push(`/sources/new/${encodeURIComponent(entry.url)}`);
              }}
              class="!p-2"
              color="light"
            >
              <i class="bx bx-folder-plus"></i>
            </Button></TableBodyCell
          >
          <TableBodyCell {tdClass}
            >{entry.name}{#if entry.publisher}
              &nbsp; <i class="bx bx-book"></i>{/if}</TableBodyCell
          >
          <TableBodyCell {tdClass}>{entry.url}</TableBodyCell>
        </tr>
        {#if entry.expand}
          {#each entry.availableFeeds as feed}
            <tr class="bg-slate-200 dark:bg-gray-700">
              <TableBodyCell {tdClass}></TableBodyCell>
              <TableBodyCell {tdClass}></TableBodyCell>
              <TableBodyCell colspan={2} {tdClass}>{feed}</TableBodyCell>
              <TableBodyCell {tdClass}></TableBodyCell>
            </tr>
          {/each}
        {/if}
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
      <div class="flex w-96 flex-col gap-2">
        <div>
          <Label>Name</Label>
          <Input bind:value={aggregator.name} on:input={checkName} color={nameColor}></Input>
        </div>
        <div>
          <Label>URL</Label>
          <Input bind:value={aggregator.url} on:input={checkUrl} color={urlColor}></Input>
        </div>
        <Button
          type="submit"
          class="mt-2 w-fit"
          color="light"
          disabled={validUrl === false ||
            validName === false ||
            aggregator.name === "" ||
            aggregator.url === ""}
        >
          <i class="bx bx-check me-2"></i>
          <span>Save aggregator</span>
        </Button>
      </div>
    </form>
  {/if}
  <ErrorMessage error={aggregatorSaveError}></ErrorMessage>

  <br />
</div>
