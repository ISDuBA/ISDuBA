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
  import { Accordion, Badge, Input, Spinner, Label, Button, TableBodyCell } from "flowbite-svelte";
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
  import CAccordionItem from "$lib/Components/CAccordionItem.svelte";

  type FeedInfo = {
    id?: number;
    sourceID?: number;
    url: string;
    highlight: boolean;
  };

  type SourceInfo = {
    id?: number;
    name: string;
    feedsAvailable: number;
    feedsSubscribed: number;
    feeds: FeedInfo[];
    expand: boolean;
  };

  type AggregatorRole = {
    label: string;
    abbreviation: string;
  };

  type AggregatorEntry = {
    name: string;
    role: AggregatorRole;
    url: string;
    availableSources: SourceInfo[];
    expand: boolean;
  };

  const accordionItemDefaultClass = "flex items-center gap-x-4 text-gray-700 font-semibold w-full";
  const textFlushOpen = "text-gray-500 dark:text-white";
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

  let openAggregator: boolean[] = [];

  let formClass = "max-w-[800pt]";
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
      openAggregator = new Array(result.value.length).fill(false);
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
          highlight: false,
          sourceID: sourceID
        }
    );

  const getFeeds = (
    sourceID: number | undefined,
    feeds: FeedSubscription[],
    availableFeeds?: string[]
  ) => {
    let unsubscribedFeeds =
      availableFeeds?.map(
        (feedURL) =>
          <FeedInfo>{
            url: feedURL,
            highlight: true
          }
      ) ?? [];
    let subscribedFeeds = sourceID !== undefined ? getSubsribedFeeds(feeds, sourceID) : [];

    // Highlight the case, when a feed is configured that is no longer available
    subscribedFeeds.forEach((f) => {
      if (!unsubscribedFeeds.map((i) => i.url).includes(f.url)) {
        f.highlight = true;
      }
    });

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

  const getRoleAbbreviation = (role: string | undefined): AggregatorRole | undefined => {
    if (role) {
      switch (role) {
        case "csaf_publisher":
          return {
            label: "Publisher",
            abbreviation: "M"
          };
        case "csaf_trusted_provider":
          return {
            label: "Trusted provider",
            abbreviation: "TP"
          };
        case "csaf_provider":
          return {
            label: "Provider",
            abbreviation: "P"
          };
      }
    }
    return undefined;
  };

  const parseAggregatorData = (data: AggregatorMetadata): AggregatorEntry[] => {
    const extractEntry = (i: CSAFProviderEntry | CSAFPublisherEntry) =>
      <AggregatorEntry>{
        name: i.metadata.publisher.name,
        url: i.metadata.url,
        availableSources: getAvailableSources(i.metadata.url, data.custom),
        role: getRoleAbbreviation(i.metadata.role)
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
    if (aggregator.id === undefined) {
      return;
    }
    if (aggregatorData.get(aggregator.id)) {
      aggregatorData.delete(aggregator.id);
      aggregatorData = aggregatorData;
      saveAggregatorExpand();
      return;
    }
    loadingAggregators = true;
    const resp = await fetchAggregatorData(aggregator.url);
    loadingAggregators = false;
    if (resp.ok) {
      aggregatorData.set(aggregator.id, parseAggregatorData(resp.value));

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

  const restoreAggregatorExpand = async () => {
    let idList = JSON.parse(sessionStorage.getItem("openAggregator") ?? "[]");
    if (idList) {
      for (let id of idList) {
        let index = aggregators.findIndex((a) => a.id === id);
        if (index !== -1) {
          openAggregator[index] = true;
          await toggleAggregatorView(aggregators[index]);
        }
      }
    }
  };

  onMount(async () => {
    await getAggregators();
    await restoreAggregatorExpand();
  });
</script>

<svelte:head>
  <title>Sources - Aggregator</title>
</svelte:head>

<div>
  <SectionHeader title="Aggregator"></SectionHeader>
  <Accordion flush multiple class="my-8">
    {#each aggregators as aggregator, index (index)}
      {@const list = aggregatorData.get(aggregator.id ?? -1) ?? []}
      <CAccordionItem
        paddingFlush="pt-0 pb-3"
        defaultClass={accordionItemDefaultClass}
        bind:open={openAggregator[index]}
        {textFlushOpen}
        toggleCallback={async () => {
          await toggleAggregatorView(aggregator);
        }}
      >
        <div slot="header" class="flex flex-col gap-2">
          <div class="flex gap-2">
            <span class="me-4">{aggregator.name}</span>
            {#if aggregator.attention}
              <Badge dismissable
                >Feeds changed
                <Button
                  slot="close-button"
                  let:close
                  color="light"
                  class="ms-1 min-h-[26px] min-w-[26px] rounded border-0 bg-transparent p-0 text-primary-700 hover:bg-white/50 dark:bg-transparent dark:hover:bg-white/20"
                  on:click={async (event) => {
                    event.stopPropagation();
                    event.preventDefault();
                    resetAttention(aggregator);
                    close();
                  }}
                >
                  <i class="bx bx-x"></i>
                </Button>
              </Badge>
            {/if}
            <Button
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
          </div>
          <div class="flex gap-4">
            <span class="text-sm text-gray-800 dark:text-gray-300">{aggregator.url}</span>
          </div>
        </div>
        {#if list.length !== 0}
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
            {#each list as entry}
              <tr
                class="bg-slate-100 dark:bg-gray-700"
                on:click={() => (entry.expand = !entry.expand)}
              >
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

                <TableBodyCell {tdClass}>{entry.name}</TableBodyCell>
                <TableBodyCell {tdClass}>
                  <div class="min-w-6" title={entry.role.label}>{entry.role.abbreviation}</div>
                </TableBodyCell>
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
                    <TableBodyCell colspan={3} {tdClass}
                      >{`${source.name} (${source.feedsSubscribed}/${source.feedsAvailable})`}</TableBodyCell
                    >
                  </tr>
                  {#if source.expand}
                    {#each source.feeds as feed}
                      {@const feedClass =
                        tdClass + (feed.highlight ? " bg-amber-300 dark:bg-amber-600" : "")}
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
                        <TableBodyCell colspan={4} tdClass={feedClass}>{feed.url}</TableBodyCell>
                      </tr>
                    {/each}
                  {/if}
                {/each}
              {/if}
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
        {/if}
      </CAccordionItem>
    {/each}
  </Accordion>
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
