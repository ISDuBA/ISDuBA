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
  import {
    type AggregatorMetadata,
    type CSAFProviderEntry,
    type CSAFPublisherEntry,
    type Custom,
    type FeedSubscription,
    type Subscription
  } from "$lib/aggregatorTypes";
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";

  type FeedInfo = {
    id?: number;
    sourceID?: number;
    url: string;
  };

  type SourceInfo = {
    id?: number;
    name: string;
    feedsAvailable: number;
    feedsSubscribed: number;
    feeds: FeedInfo[];
    expand: boolean;
  };

  type AggregatorEntry = {
    name: string;
    role: string;
    url: string;
    availableSources: SourceInfo[];
    expand: boolean;
  };

  let loadingAggregators: boolean = false;
  let aggregators: Aggregator[] = [];
  let aggregatorData = new Map<number, AggregatorEntry[]>();

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

  const getSubsribedFeeds = (feeds: FeedSubscription[], sourceID: number): FeedInfo[] =>
    feeds.map(
      (f) =>
        <FeedInfo>{
          id: f.id,
          url: f.url,
          sourceID: sourceID
        }
    );

  const getFeeds = (
    sourceID: number | undefined,
    feeds: FeedSubscription[],
    availableFeeds?: string[]
  ) => {
    let subscribedFeeds = sourceID !== undefined ? getSubsribedFeeds(feeds, sourceID) : [];
    let unsubscribedFeeds =
      availableFeeds?.map(
        (feedURL) =>
          <FeedInfo>{
            url: feedURL
          }
      ) ?? [];
    unsubscribedFeeds = unsubscribedFeeds.filter(
      (f) => !subscribedFeeds.map((i) => i.url).includes(f.url)
    );
    return [...unsubscribedFeeds, ...subscribedFeeds];
  };

  const getSources = (entry: Subscription): SourceInfo[] =>
    entry.subscriptions?.map(
      (s) =>
        <SourceInfo>{
          id: s.id,
          name: s.name,
          expand: false,
          feedsAvailable: entry.available?.length ?? 0,
          feedsSubscribed: s.subscripted?.length ?? 0,
          feeds: getFeeds(s.id, s.subscripted ?? [], entry.available)
        }
    ) ?? [
      <SourceInfo>{
        name: "Not configured",
        feedsAvailable: entry.available?.length ?? 0,
        feedsSubscribed: entry.available?.length ?? 0,
        feeds: getFeeds(undefined, [], entry.available)
      }
    ];

  const findSubscription = (url: string, custom: Custom) =>
    custom.subscriptions.find((i) => i.url === url);

  const getAvailableSources = (url: string, custom: Custom) => {
    const subscription = findSubscription(url, custom);
    if (subscription) {
      return getSources(subscription);
    } else {
      return [];
    }
  };

  const parseAggregatorData = (data: AggregatorMetadata): AggregatorEntry[] => {
    const extractEntry = (i: CSAFProviderEntry | CSAFPublisherEntry) =>
      <AggregatorEntry>{
        name: i.metadata.publisher.name,
        url: i.metadata.url,
        availableSources: getAvailableSources(i.metadata.url, data.custom),
        role: i.metadata.role?.replace("csaf_", "").replace("_", " ")
      };

    const csafProviders = data.aggregator.csaf_providers.map(extractEntry);
    const csafPublisher = data.aggregator.csaf_publishers?.map(extractEntry) ?? [];

    return [...csafProviders, ...csafPublisher];
  };

  const resetAttention = async (aggregator: Aggregator) => {
    let resetResult = await resetAggregatorAttention(aggregator);
    if (resetResult.ok) {
      aggregator.attention = false;
    } else {
      aggregatorError = resetResult.error;
    }
  };

  const toggleAggregatorView = async (aggregator: Aggregator) => {
    if (aggregatorData.get(aggregator.id ?? -1)) {
      aggregatorData.delete(aggregator.id ?? -1);
      aggregatorData = aggregatorData;
      saveAggregatorExpand();
      return;
    }
    loadingAggregators = true;
    const resp = await fetchAggregatorData(aggregator.url);
    loadingAggregators = false;
    if (resp.ok) {
      aggregatorData.set(aggregator.id ?? -1, parseAggregatorData(resp.value));

      aggregatorData = aggregatorData;
      saveAggregatorExpand();
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

  const saveAggregatorExpand = () => {
    let idList = [...aggregatorData.keys()];
    sessionStorage.setItem("openAggregator", JSON.stringify(idList));
  };

  const restoreAggregatorExpand = () => {
    let idList = JSON.parse(sessionStorage.getItem("openAggregator") ?? "[]");
    if (idList) {
      for (let id of idList) {
        let aggregator = aggregators.find((a) => a.id === id);
        if (aggregator) {
          toggleAggregatorView(aggregator);
        }
      }
    }
  };

  onMount(async () => {
    await getAggregators();
    restoreAggregatorExpand();
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
        attribute: "delete",
        class: smallColumnClass
      },
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
      {
        label: "Name",
        attribute: "name"
      },
      {
        label: "Role",
        attribute: "role"
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

        <TableBodyCell tdClass={`${tdClass} ${smallColumnClass}`}>
          {#if list.length === 0}
            <i class="bx bx-plus"></i>
          {:else}
            <i class="bx bx-minus"></i>
          {/if}
        </TableBodyCell>
        {#if aggregator.attention}
          <TableBodyCell tdClass={`${tdClass} ${extraSmallColumnClass}`}>
            <button on:click={async () => resetAttention(aggregator)}>
              <i class="bx bx-info-square text-lg"></i>
            </button>
          </TableBodyCell>
        {:else}
          <TableBodyCell tdClass={`${tdClass} ${extraSmallColumnClass}`}></TableBodyCell>
        {/if}
        <TableBodyCell {tdClass}>{aggregator.name}</TableBodyCell>
        <TableBodyCell {tdClass}></TableBodyCell>
        <TableBodyCell {tdClass}>{aggregator.url}</TableBodyCell>
      </tr>
      {#each list as entry}
        <tr class="bg-slate-100 dark:bg-gray-700" on:click={() => (entry.expand = !entry.expand)}>
          <TableBodyCell {tdClass}>
            <Button
              on:click={async () => {
                await push(`/sources/new/${encodeURIComponent(entry.url)}`);
              }}
              class="!p-2"
              color="light"
            >
              <i class="bx bx-folder-plus"></i>
            </Button>
          </TableBodyCell>

          <TableBodyCell {tdClass}>
            {#if entry.expand}
              <i class="bx bx-minus"></i>
            {:else}
              <i class="bx bx-plus"></i>
            {/if}
          </TableBodyCell>

          <TableBodyCell {tdClass}></TableBodyCell>
          <TableBodyCell {tdClass}>{entry.name}</TableBodyCell>
          <TableBodyCell {tdClass}>{entry.role}</TableBodyCell>
          <TableBodyCell {tdClass}>{entry.url}</TableBodyCell>
        </tr>
        {#if entry.expand}
          {#each entry.availableSources as source}
            <tr
              class="bg-slate-300 dark:bg-gray-600"
              on:click={() => (source.expand = !source.expand)}
            >
              <TableBodyCell {tdClass}
                >{#if source.id !== undefined}<Button
                    on:click={async () => {
                      await push(`/sources/${source.id}`);
                    }}
                    class="!p-2"
                    color="light"
                  >
                    <i class="bx bx-folder-open"></i>
                  </Button>{/if}
              </TableBodyCell>
              <TableBodyCell {tdClass}>
                {#if source.expand}
                  <i class="bx bx-minus"></i>
                {:else}
                  <i class="bx bx-plus"></i>
                {/if}
              </TableBodyCell>
              <TableBodyCell {tdClass}></TableBodyCell>
              <TableBodyCell colspan={3} {tdClass}
                >{`${source.name} (${source.feedsSubscribed}/${source.feedsAvailable})`}</TableBodyCell
              >
            </tr>
            {#if source.expand}
              {#each source.feeds as feed}
                {@const feedClass =
                  tdClass + (feed.id === undefined ? " bg-amber-300 dark:bg-amber-600" : "")}
                <tr class="bg-slate-400 dark:bg-gray-500">
                  <TableBodyCell tdClass={feedClass}>
                    {#if feed.id !== undefined}
                      <Button
                        on:click={async () => {
                          sessionStorage.setItem("feedBlinkID", String(feed.id));
                          await push(`/sources/${feed.sourceID}`);
                        }}
                        class="!p-2"
                        color="light"
                      >
                        <i class="bx bx-folder-open"></i>
                      </Button>
                    {/if}
                  </TableBodyCell>
                  <TableBodyCell tdClass={feedClass}></TableBodyCell>
                  <TableBodyCell colspan={4} tdClass={feedClass}>{feed.url}</TableBodyCell>
                </tr>
              {/each}
            {/if}
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
