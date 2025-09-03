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
    updateAggregator,
    fetchAggregatorData,
    fetchAggregators,
    deleteAggregator,
    type Aggregator,
    resetAggregatorAttention,
    dtClass,
    ddClass
  } from "$lib/Sources/source";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import {
    Accordion,
    Badge,
    DescriptionList,
    Input,
    List,
    Spinner,
    Label,
    Button,
    Toggle
  } from "flowbite-svelte";
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
  import { appStore } from "$lib/store.svelte";
  import { onMount } from "svelte";
  import { push } from "svelte-spa-router";
  import CAccordionItem from "$lib/Components/CAccordionItem.svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import { scale } from "svelte/transition";
  import FeedBulletPoint from "./FeedBulletPoint.svelte";
  import type { AggregatorEntry, AggregatorRole, FeedInfo, SourceInfo } from "./aggregator";
  import SourceContent from "./SourceContent.svelte";

  const textFlushOpen = "text-black dark:text-white";
  const accordionItemDefaultClass = `flex items-center gap-x-4 ${textFlushOpen} font-semibold w-full`;
  let loadingAggregators: boolean = false;
  let aggregators: Aggregator[] = [];
  let aggregatorData = new Map<number, AggregatorEntry[]>();
  let aggregatorMetaData = new Map<number, AggregatorMetadata>();

  let aggregatorError: ErrorDetails | null;
  let aggregatorSaveError: ErrorDetails | null;
  let aggregatorEditError: ErrorDetails | null;

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
  let validEditedName: boolean | null = null;
  let editedNameColor: "red" | "green" | "base" = "base";
  $: if (validEditedName !== undefined) {
    if (validEditedName === null) {
      editedNameColor = "base";
    } else if (validEditedName) {
      editedNameColor = "green";
    } else {
      editedNameColor = "red";
    }
  }
  let validEditedUrl: boolean | null = null;
  let editedUrlColor: "red" | "green" | "base" = "base";
  $: if (validEditedUrl !== undefined) {
    if (validEditedUrl === null) {
      editedUrlColor = "base";
    } else if (validUrl) {
      editedUrlColor = "green";
    } else {
      editedUrlColor = "red";
    }
  }

  let editedName: string = "";
  let editedUrl: string = "";

  let aggregator: Aggregator = {
    name: "",
    url: ""
  };

  let blinkId: number | undefined = undefined;
  let openAggregator: boolean[] = [];
  let showCreateForm = false;
  let aggregatorToEdit: number | undefined = undefined;
  let formClass = "max-w-[800pt]";

  const toggleCreateForm = () => {
    showCreateForm = !showCreateForm;
  };

  const toggleEditForm = (id: number) => {
    if (aggregatorToEdit) {
      aggregatorToEdit = undefined;
    } else {
      aggregatorToEdit = id;
    }
  };

  const checkUrl = (edit = false) => {
    const url = edit ? editedUrl : aggregator.url;
    if (url === "") {
      if (edit) {
        validEditedName = null;
      } else {
        validUrl = null;
      }
      return;
    }
    if (url.startsWith("https://") && url.endsWith("aggregator.json")) {
      if (edit) {
        validEditedUrl = null;
      } else {
        validUrl = null;
      }
      return;
    }
    if (edit) {
      validEditedUrl = false;
    } else {
      validUrl = false;
    }
  };

  const checkName = (id: number | undefined = undefined, edit = false) => {
    const name = edit ? editedName : aggregator.name;
    if (aggregators.find((i) => i.name === name && (i.id !== id || !id))) {
      if (edit) {
        validEditedName = false;
      } else {
        validName = false;
      }
      return;
    }
    if (edit) {
      validEditedName = null;
    } else {
      validName = null;
    }
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
          feedsSubscribed: getSubsribedValidFeedCount(entry.available ?? [], s.subscripted ?? []),
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

  const getRoleAbbreviation = (role: string | undefined): AggregatorRole => {
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
    return {
      label: role ?? "",
      abbreviation: role?.substring(0, 1) ?? ""
    };
  };

  const getAvailableFeedCount = (url: string, custom: Custom) => {
    const subscription = findSubscription(url, custom);
    if (subscription) {
      return subscription.available?.length ?? 0;
    }
    return 0;
  };

  const getSubsribedValidFeedCount = (available: string[], subscribedFeeds: FeedSubscription[]) => {
    return subscribedFeeds.filter((f) => available.includes(f.url)).length;
  };

  const getSubsribedFeedCount = (url: string, custom: Custom) => {
    const subscription = findSubscription(url, custom);
    if (subscription) {
      let available = subscription?.available ?? [];
      let subscribed = subscription.subscriptions ?? [];
      let subscribedFeeds = subscribed.flatMap((s) => s.subscripted ?? []).map((f) => f.url);
      let uniqueSubscribedFeeds = [...new Set(subscribedFeeds)];

      return uniqueSubscribedFeeds.filter((s) => available.includes(s)).length;
    }
    return 0;
  };

  const parseAggregatorData = (data: AggregatorMetadata): AggregatorEntry[] => {
    const extractEntry = (i: CSAFProviderEntry | CSAFPublisherEntry) =>
      <AggregatorEntry>{
        name: i.metadata.publisher.name,
        url: i.metadata.url,
        feedsAvailable: getAvailableFeedCount(i.metadata.url, data.custom),
        feedsSubscribed: getSubsribedFeedCount(i.metadata.url, data.custom),
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
      aggregators = aggregators;
    } else {
      aggregatorError = resetResult.error;
    }
  };

  const toggleActive = async (aggregator: Aggregator) => {
    aggregator.active = !aggregator.active;
    const result = await updateAggregator(aggregator);
    if (result.ok) {
      aggregators = aggregators;
    } else {
      aggregatorError = result.error;
    }
  };

  const toggleAggregatorView = async (aggregator: Aggregator) => {
    await navigator.locks.request("toggleAgg", async () => {
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
        aggregatorMetaData.set(aggregator.id, resp.value);
        aggregatorMetaData = aggregatorMetaData;
        saveAggregatorExpand();
      } else {
        aggregatorError = resp.error;
      }
    });
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
      showCreateForm = false;
      aggregator.name = "";
      aggregator.url = "";
      sessionStorage.setItem(
        "openAggregator",
        JSON.stringify([...aggregatorData.keys(), result.value])
      );
      await getAggregators();
      await restoreAggregatorExpand();
      await new Promise((res) => setTimeout(res, 500));
      document.getElementById(`aggregator-${result.value}`)?.scrollIntoView({ behavior: "smooth" });
      blinkId = result.value;
      await new Promise((res) => setTimeout(res, 5000));
      blinkId = undefined;
    }
  };

  const editAggregator = async (aggregator: Aggregator) => {
    let result = await updateAggregator(aggregator);
    if (!result.ok) {
      aggregatorEditError = result.error;
    } else {
      aggregatorToEdit = undefined;
      await getAggregators();
      await restoreAggregatorExpand();
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
  <title>Sources - Aggregators</title>
</svelte:head>

<div class="pb-10">
  <SectionHeader title="Aggregators"></SectionHeader>
  {#if appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager()}
    <Accordion flush multiple class="my-4">
      {#each aggregators as aggregator, index (index)}
        {@const list = aggregatorData.get(aggregator.id ?? -1) ?? []}
        {@const metadata = aggregatorMetaData.get(aggregator.id ?? -1)}
        <CAccordionItem
          id={`aggregator-${aggregator.id}`}
          paddingFlush="pt-0 py-2"
          defaultClass={`${accordionItemDefaultClass} ${aggregator.id === blinkId ? "blink" : ""}`}
          bind:open={openAggregator[index]}
          {textFlushOpen}
          toggleCallback={async () => {
            await toggleAggregatorView(aggregator);
          }}
        >
          <span slot="arrowup"></span>
          <span slot="arrowdown"> </span>
          <div slot="header" class="flex flex-col items-start gap-2">
            <div class="flex flex-wrap items-center gap-2">
              {#if list.length > 0}
                <i class="bx bx-chevron-up text-xl"></i>
              {:else}
                <i class="bx bx-chevron-down text-xl"></i>
              {/if}
              <span>{aggregator.name}</span>
              {#if aggregator.attention}
                <Badge class="h-fit">Sources changed</Badge>
              {/if}
              <div>
                {#if appStore.isSourceManager()}
                  <Button
                    on:click={async (event) => {
                      event.stopPropagation();
                      event.preventDefault();
                      if (aggregator.id !== undefined) {
                        await removeAggregator(aggregator.id);
                      }
                    }}
                    class="!p-2"
                    color="light"
                  >
                    <i class="bx bx-trash text-red-600"></i>
                  </Button>
                  {#if aggregator.id !== undefined && aggregator.id !== aggregatorToEdit}
                    <Button
                      on:click={(event) => {
                        event.stopPropagation();
                        event.preventDefault();
                        editedName = aggregator.name;
                        editedUrl = aggregator.url;
                        if (aggregator.id !== undefined) {
                          toggleEditForm(aggregator.id);
                        }
                      }}
                      class="!p-2"
                      color="light"
                    >
                      <i class="bx bx-pencil"></i>
                    </Button>
                  {/if}
                {/if}
              </div>
              {#if aggregator.active !== undefined && appStore.isSourceManager()}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <!-- svelte-ignore a11y-no-static-element-interactions -->
                <!--
                Stopping the event propagation in the click event of the Toggle doesn't work when using
                the mouse because it also consists of a span and a label element. These elements also fire
                click events which aren't stopped.
                Still, the Toggle needs the event listener so it can be toggled via keyboard.
              -->
                <div
                  on:click={(event) => {
                    event.preventDefault();
                    event.stopPropagation();
                    toggleActive(aggregator);
                  }}
                  class="mx-2 cursor-pointer p-2"
                >
                  <Toggle
                    on:click={(event) => {
                      event.stopPropagation();
                      event.preventDefault();
                      toggleActive(aggregator);
                    }}
                    bind:checked={aggregator.active}
                  >
                    {#if aggregator.active === true}
                      <span>Active</span>
                    {:else}
                      <span>Not active</span>
                    {/if}
                  </Toggle>
                </div>
              {/if}
            </div>
            {#if aggregator.id !== undefined && aggregator.id === aggregatorToEdit}
              <div class="flex flex-wrap gap-4">
                <div class="flex flex-col items-center gap-1 md:flex-row">
                  <Label>Name</Label>
                  <Input
                    class="h-fit w-fit"
                    bind:value={editedName}
                    on:click={(event) => {
                      event.stopPropagation();
                      event.preventDefault();
                    }}
                    on:input={() => {
                      if (aggregatorToEdit) {
                        checkName(aggregatorToEdit, true);
                      }
                    }}
                    color={editedNameColor}
                  ></Input>
                </div>
                <div class="flex flex-col items-center gap-1 md:flex-row">
                  <Label>URL</Label>
                  <Input
                    class="h-fit w-fit"
                    bind:value={editedUrl}
                    on:click={(event) => {
                      event.stopPropagation();
                      event.preventDefault();
                    }}
                    on:input={() => {
                      checkUrl(true);
                    }}
                    color={editedUrlColor}
                  ></Input>
                </div>
                <div class="mt-2 mb-2 flex flex-wrap gap-2">
                  <Button
                    class="w-fit"
                    on:click={(event) => {
                      event.stopPropagation();
                      event.preventDefault();
                      if (aggregator.id !== undefined) {
                        toggleEditForm(aggregator.id);
                      }
                    }}
                    color="light"><i class="bx bx-x"></i></Button
                  >
                  <Button
                    on:click={() => {
                      editAggregator({
                        id: aggregatorToEdit,
                        name: editedName,
                        url: editedUrl,
                        attention: aggregator.attention
                      });
                    }}
                    class="w-fit"
                    color="green"
                    disabled={validEditedUrl === false ||
                      validEditedName === false ||
                      editedName === "" ||
                      editedUrl === ""}
                  >
                    <i class="bx bx-check me-2"></i>
                    <span>Save</span>
                  </Button>
                </div>
                <ErrorMessage error={aggregatorEditError}></ErrorMessage>
              </div>
            {/if}
          </div>
          {#if list.length !== 0}
            <div
              class="mb-2 flex flex-col justify-between rounded-md border border-solid border-gray-300 px-4 py-2 break-all dark:border-gray-500"
            >
              <List tag="dl" class="w-full divide-y divide-gray-200 text-sm dark:divide-gray-600">
                <div>
                  <DescriptionList tag="dt" {dtClass}>URL</DescriptionList>
                  <DescriptionList tag="dd" {ddClass}>{aggregator.url}</DescriptionList>
                </div>
                {#if metadata?.aggregator}
                  {@const data = metadata.aggregator.aggregator}
                  <div>
                    <DescriptionList tag="dt" {dtClass}>Category</DescriptionList>
                    <DescriptionList tag="dd" {ddClass}>{data.category}</DescriptionList>
                  </div>
                  <div>
                    <DescriptionList tag="dt" {dtClass}>Last updated</DescriptionList>
                    <DescriptionList tag="dd" {ddClass}
                      >{metadata.aggregator.last_updated}</DescriptionList
                    >
                  </div>
                  <div>
                    <DescriptionList tag="dt" {dtClass}>Namespace</DescriptionList>
                    <DescriptionList tag="dd" {ddClass}>{data.namespace}</DescriptionList>
                  </div>
                  <div>
                    <DescriptionList tag="dt" {dtClass}>Contact details</DescriptionList>
                    <DescriptionList tag="dd" {ddClass}>{data.contact_details}</DescriptionList>
                  </div>
                  <div>
                    <DescriptionList tag="dt" {dtClass}>Issuing authority</DescriptionList>
                    <DescriptionList tag="dd" {ddClass}>{data.issuing_authority}</DescriptionList>
                  </div>
                {/if}
              </List>
            </div>
            {#if aggregator.attention && appStore.isSourceManager()}
              <Badge class="mb-2 h-fit p-1" dismissable>
                <p>
                  These are the currently available providers. Please review their feeds and adjust
                  the sources if needed.
                </p>
                <Button
                  slot="close-button"
                  let:close
                  color="light"
                  class="border-primary-700/55 text-primary-700 ms-1 min-h-[26px] min-w-[26px] rounded border bg-transparent p-0 hover:bg-white/50 dark:bg-transparent dark:hover:bg-white/20"
                  on:click={async (event) => {
                    event.stopPropagation();
                    event.preventDefault();
                    resetAttention(aggregator);
                    close();
                  }}
                >
                  <i class="bx bx-check"></i>
                </Button>
              </Badge>
            {/if}
            <div class="ps-4">
              {#each list as entry}
                <Collapsible header="" showBorder={false}>
                  <div slot="header" class="mb-2 flex items-center gap-2">
                    <div
                      class="flex flex-wrap items-center gap-x-3 gap-y-1 text-sm text-black dark:text-white"
                    >
                      <span>{entry.name}</span>
                      <span class="flex w-fit gap-1">
                        {#each new Array(entry.feedsSubscribed) as _a}
                          <FeedBulletPoint filled></FeedBulletPoint>
                        {/each}
                        {#each new Array(entry.feedsAvailable - entry.feedsSubscribed) as _a}
                          <FeedBulletPoint></FeedBulletPoint>
                        {/each}
                      </span>
                    </div>
                  </div>
                  <div class="mb-3 flex flex-col gap-3">
                    <List
                      tag="dl"
                      class="w-full divide-y divide-gray-200 text-sm dark:divide-gray-600"
                    >
                      <div>
                        <DescriptionList tag="dt" {dtClass}>URL</DescriptionList>
                        <DescriptionList tag="dd" {ddClass}>{entry.url}</DescriptionList>
                      </div>
                      <div>
                        <DescriptionList tag="dt" {dtClass}>Role</DescriptionList>
                        <DescriptionList tag="dd" {ddClass}>{entry.role.label}</DescriptionList>
                      </div>
                    </List>
                    {#each entry.availableSources as source}
                      {#if source.id === undefined || !appStore.isSourceManager()}
                        <div class="p-2">
                          <SourceContent {entry} {source}></SourceContent>
                        </div>
                      {:else}
                        <button
                          on:click={async () => {
                            await push(`/sources/${source.id}`);
                          }}
                          class={entry.feedsSubscribed === 0
                            ? ""
                            : "rounded-md border border-solid border-gray-300 p-2 hover:bg-gray-200 dark:hover:bg-gray-700"}
                        >
                          <SourceContent {entry} {source}></SourceContent>
                        </button>
                      {/if}
                    {/each}
                    {#if entry.feedsSubscribed > 0 && appStore.isSourceManager()}
                      <Button
                        href={`/#/sources/new/${encodeURIComponent(entry.url)}`}
                        class="mb-2 w-fit"
                        color="light"
                        size="xs"
                      >
                        <i class="bx bx-plus"></i>
                        <span>Again as another source</span>
                      </Button>
                    {/if}
                  </div>
                </Collapsible>
              {/each}
            </div>
          {/if}
        </CAccordionItem>
      {/each}
    </Accordion>
    <div class:invisible={!loadingAggregators} class={loadingAggregators ? "loadingFadeIn" : ""}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <ErrorMessage error={aggregatorError}></ErrorMessage>
    {#if appStore.isSourceManager()}
      <div class="min-h-64">
        {#if !showCreateForm}
          <Button class="mt-3 mb-2 w-fit" on:click={toggleCreateForm}
            ><i class="bx bx-plus me-2"></i>New aggregator</Button
          >
        {/if}
        {#if showCreateForm}
          <form transition:scale on:submit={submitAggregator} class={formClass}>
            <div class="flex w-96 flex-col gap-2">
              <div>
                <Label>Name</Label>
                <Input
                  bind:value={aggregator.name}
                  on:input={() => {
                    checkName();
                  }}
                  color={nameColor}
                ></Input>
              </div>
              <div>
                <Label>URL</Label>
                <Input
                  bind:value={aggregator.url}
                  on:input={() => {
                    checkUrl();
                  }}
                  color={urlColor}
                ></Input>
              </div>
              <div class="mt-2 mb-2 flex gap-2">
                <Button class="w-fit" on:click={toggleCreateForm} color="light"
                  ><i class="bx bx-x"></i></Button
                >
                <Button
                  type="submit"
                  class="w-fit"
                  color="green"
                  disabled={validUrl === false ||
                    validName === false ||
                    aggregator.name === "" ||
                    aggregator.url === ""}
                >
                  <i class="bx bx-check me-2"></i>
                  <span>Save aggregator</span>
                </Button>
              </div>
              <ErrorMessage error={aggregatorSaveError}></ErrorMessage>
            </div>
          </form>
        {/if}
      </div>
    {/if}
  {/if}
</div>
