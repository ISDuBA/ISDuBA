<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    type Source,
    type Feed,
    fetchPMD,
    fetchSource,
    saveSource,
    fetchFeeds,
    calculateMissingFeeds,
    parseFeeds,
    saveFeeds
  } from "$lib/Sources/source";
  import { Button, Spinner, Modal, List, DescriptionList } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails, getErrorDetails } from "$lib/Errors/error";
  import type { CSAFProviderMetadata } from "$lib/provider";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { onDestroy, onMount } from "svelte";
  import SourceForm from "./SourceForm.svelte";
  import { request } from "$lib/request";
  import FeedView from "./FeedView.svelte";
  import { push } from "svelte-spa-router";
  import { DAY_MS } from "$lib/time";
  import SourceBasicStats from "./SourceBasicStats.svelte";
  import ImportStats from "$lib/Statistics/ImportStats.svelte";
  export let params: any = null;

  let sourceEdited: boolean = false;

  let modalOpen: boolean = false;
  let modalMessage = "";
  let modalTitle = "";
  let modalCallback: any;

  let saveSourceError: ErrorDetails | null;
  let loadSourceError: ErrorDetails | null;
  let loadFeedError: ErrorDetails | null;
  let saveFeedError: ErrorDetails | null;
  let loadPmdError: ErrorDetails | null;
  let feedError: ErrorDetails | null;
  let pmd: CSAFProviderMetadata;
  let feeds: Feed[] = [];

  let loadingFeeds: boolean = false;
  let loadingSource: boolean = false;
  let loadingPMD: boolean = false;

  let formClass = "max-w-[800pt]";

  let sourceForm: any;
  let updateSourceForm: any;
  let fillAgeDataFromSource: (source: Source) => void;

  let loadSource: boolean = true;

  let source: Source = {
    name: "",
    url: "",
    active: false,
    rate: 1,
    slots: 2,
    strict_mode: true,
    headers: [""],
    ignore_patterns: [""]
  };

  let oldSource = structuredClone(source);

  const dtClass: string = "ml-1 mt-1 text-gray-500 md:text-sm dark:text-gray-400";
  const ddClass: string = "break-words font-semibold ml-2 mb-1";

  let updateStats = setInterval(async () => {
    if (!source.id || source.id === 0) {
      return;
    }
    let result = await fetchSource(source.id, true);
    if (result.ok) {
      source.stats = result.value.stats;
    }
    let feedResult = await fetchFeeds(source.id ?? 0, true);
    if (feedResult.ok) {
      for (let feed of feedResult.value) {
        const find = feeds.find((i) => i.id === feed.id);
        if (find) {
          find.stats = feed.stats;
        }
      }
    }
    feeds = feeds;
  }, 30 * 1000);

  const loadSourceInfo = async (id: number) => {
    loadingSource = true;
    let result = await fetchSource(Number(id), true);
    if (result.ok) {
      source = result.value;
      loadSource = true;
      if (fillAgeDataFromSource) {
        fillAgeDataFromSource(source);
      }
      await updateSourceForm();
      oldSource = structuredClone(source);
      sourceEdited = false;
    } else {
      loadSourceError = result.error;
    }
    loadingSource = false;
  };

  const loadPMD = async () => {
    loadingPMD = true;
    let result = await fetchPMD(source.url);
    if (result.ok) {
      pmd = result.value;
    } else {
      loadPmdError = result.error;
    }
    loadingPMD = false;
  };

  const loadFeeds = async () => {
    if (source.id === undefined) {
      return;
    }
    loadingFeeds = true;
    let result = await fetchFeeds(source.id, true);
    if (result.ok) {
      feeds = result.value;
      feeds.map((f) => {
        f.enable = true;
      });
    } else {
      loadFeedError = result.error;
    }
    loadingFeeds = false;
  };

  const updateSource = async () => {
    await updateSourceForm();
    let result = await saveSource(source);
    if (!result.ok) {
      saveSourceError = result.error;
      return;
    }
    saveSourceError = null;
    await loadSourceInfo(source.id ?? 0);
  };

  const deleteSource = async () => {
    const resp = await request(`/api/sources/${source.id}`, "DELETE");
    if (resp.error) {
      saveSourceError = getErrorDetails(`Could not delete source`, resp);
    } else {
      push(`/sources`);
    }
  };

  const isDuplicateFeedLabel = (feed: Feed): boolean => {
    let found = feeds.filter((f) => f.id !== feed.id && f.label === feed.label);
    return found.length === 0 ? false : true;
  };

  const updateFeed = async (feed: Feed) => {
    if (isDuplicateFeedLabel(feed) || feed.label.length === 0) {
      return;
    }
    let result = await saveFeeds(source, [feed]);
    if (result.ok) {
      let id = result.value[0];
      if (id) {
        feed.id = id;
      }
    } else {
      saveFeedError = result.error;
    }
  };

  const sourceEqual = (a: Source, b: Source) => {
    let tmpA = structuredClone(a);
    let tmpB = structuredClone(b);

    tmpA.stats = undefined;
    tmpB.stats = undefined;

    if (!tmpA.headers) {
      tmpA.headers = [];
    }
    if (!tmpB.headers) {
      tmpB.headers = [];
    }

    if (tmpA.insecure === undefined) {
      tmpA.insecure = false;
    }
    if (tmpB.insecure === undefined) {
      tmpB.insecure = false;
    }

    if (tmpA.signature_check === undefined) {
      tmpA.signature_check = false;
    }
    if (tmpB.signature_check === undefined) {
      tmpB.signature_check = false;
    }

    if (tmpA.strict_mode === undefined) {
      tmpA.strict_mode = false;
    }
    if (tmpB.strict_mode === undefined) {
      tmpB.strict_mode = false;
    }
    return JSON.stringify(tmpA) === JSON.stringify(tmpB);
  };

  const inputChange = async () => {
    await updateSourceForm();
    if (sourceEqual(oldSource, source)) {
      sourceEdited = false;
    } else {
      sourceEdited = true;
    }
  };

  const clickFeed = async (feed: Feed) => {
    if (!feed.id) {
      return;
    }
    push(`/sources/logs/${feed.id}`);
  };

  onMount(async () => {
    updateSourceForm = sourceForm.updateSource;
    fillAgeDataFromSource = sourceForm.fillAgeDataFromSource;
    let id = params?.id;
    if (id) {
      let sourceID = Number(id);
      if (sourceID === 0) {
        source.id = sourceID;
      } else {
        await loadSourceInfo(sourceID);
        await loadPMD();
      }
      await loadFeeds();
      if (id !== 0) {
        let missingFeeds = calculateMissingFeeds(parseFeeds(pmd, feeds), feeds);
        missingFeeds.map((f) => {
          f.enable = false;
        });
        feeds.push(...missingFeeds);
      }
      feeds = feeds;
      fillAgeDataFromSource(source);
    }
  });

  onDestroy(() => {
    clearInterval(updateStats);
  });
</script>

<svelte:head>
  <title>Sources - Edit source</title>
</svelte:head>

<Modal size="xs" title={modalTitle} bind:open={modalOpen} autoclose outsideclose>
  <div class="text-center">
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {modalMessage}
    </h3>
    <Button
      on:click={() => {
        modalCallback();
      }}
      color="red"
      class="me-2">Yes, I'm sure</Button
    >
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

<SectionHeader title={source.name}></SectionHeader>
<div class="mb-3 grid w-full grid-cols-1 justify-stretch gap-10 lg:grid-cols-2">
  <div class="w-full">
    <List tag="dl" class="w-full divide-y divide-gray-200 text-sm">
      <div>
        <DescriptionList tag="dt" {dtClass}>Domain/PMD</DescriptionList>
        <DescriptionList tag="dd" {ddClass}>{source.url}</DescriptionList>
      </div>
      {#if pmd}
        <div>
          <DescriptionList tag="dt" {dtClass}>Canonical URL</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{pmd.canonical_url}</DescriptionList>
        </div>
        <div>
          <DescriptionList tag="dt" {dtClass}>Publisher Name</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{pmd.publisher.name}</DescriptionList>
        </div>
        <div>
          <DescriptionList tag="dt" {dtClass}>Publisher Contact</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{pmd.publisher.contact_details}</DescriptionList>
        </div>
        <div>
          {#if pmd.publisher.issuing_authority}
            <DescriptionList tag="dt" {dtClass}>Issuing Authority</DescriptionList>
            <DescriptionList tag="dd" {ddClass}>{pmd.publisher.issuing_authority}</DescriptionList>
          {/if}
        </div>
      {/if}
      <div>
        <DescriptionList tag="dt" {dtClass}>Status</DescriptionList>
        {#if source.status}
          {#each source.status as s}
            <DescriptionList tag="dd" {ddClass}>{s}</DescriptionList>
          {/each}
        {:else}
          <DescriptionList tag="dd" {ddClass}>OK</DescriptionList>
        {/if}
      </div>
    </List>
    {#if source.stats}
      <h4 class="mt-3">Status</h4>
      <List tag="dl" class="flex w-full flex-wrap text-sm">
        <div>
          <DescriptionList tag="dt" {dtClass}>Loading</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{source.stats.downloading}</DescriptionList>
        </div>
        <div class="pl-4">
          <DescriptionList tag="dt" {dtClass}>Queued</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{source.stats.waiting}</DescriptionList>
        </div>
        <div class="pl-4">
          <DescriptionList tag="dt" {dtClass}>Imported (last 24h)</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>
            {#if source.id}
              {@const yesterday = Date.now() - DAY_MS}
              <SourceBasicStats sourceID={source.id}></SourceBasicStats>
              (<SourceBasicStats from={new Date(yesterday)} sourceID={source.id}
              ></SourceBasicStats>)
            {/if}
          </DescriptionList>
        </div>
      </List>
    {/if}
    <ErrorMessage error={loadSourceError}></ErrorMessage>
    <ErrorMessage error={loadPmdError}></ErrorMessage>

    <div class:invisible={!loadingPMD} class={!loadingPMD ? "loadingFadeIn" : ""} class:mb-4={true}>
      Loading PMD ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  </div>

  {#if source.id !== 0}
    <div class="w-full flex-auto">
      <div
        class:invisible={!loadingSource}
        class={!loadingSource ? "loadingFadeIn" : ""}
        class:mb-4={true}
      >
        Loading source configuration ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
      <SourceForm
        bind:this={sourceForm}
        bind:parseSource={loadSource}
        {inputChange}
        {source}
        {formClass}
        enableActive={true}
      ></SourceForm>
      <Button disabled={!sourceEdited} on:click={updateSource} color="light">
        <i class="bx bxs-save me-2"></i>
        <span>Save source</span>
      </Button>
      <Button
        on:click={(event) => {
          event.stopPropagation();
          modalCallback = () => {
            deleteSource();
          };
          modalMessage = "Are you sure you want to delete this source?";
          modalTitle = `Source ${source.name}`;
          modalOpen = true;
        }}
        title={`Delete source "${source.name}"`}
        color="light"
      >
        <i class="bx bx-trash me-2 text-red-500"></i>
        <span>Delete source</span>
      </Button>
      <ErrorMessage error={saveSourceError}></ErrorMessage>
    </div>
  {/if}
</div>

<FeedView {feeds} {clickFeed} {updateFeed} edit={true}></FeedView>
<div
  class:invisible={!loadingFeeds && !loadingPMD}
  class={!loadingFeeds && !loadingPMD ? "loadingFadeIn" : ""}
  class:mb-4={true}
>
  Loading ...
  <Spinner color="gray" size="4"></Spinner>
</div>
<ErrorMessage error={loadFeedError}></ErrorMessage>
<ErrorMessage error={feedError}></ErrorMessage>
<ErrorMessage error={saveFeedError}></ErrorMessage>

{#if source.id}
  <ImportStats
    axes={[{ label: "", types: ["imports"] }]}
    height="200pt"
    initialFrom={new Date(Date.now() - DAY_MS)}
    showModeToggle
    showRangeSelection
    source={{ id: source.id, isFeed: false }}
    title="Imports"
  ></ImportStats>
  <ImportStats
    axes={[{ label: "", types: ["importFailures"] }]}
    height="200pt"
    initialFrom={new Date(Date.now() - DAY_MS)}
    isStacked
    showLegend
    showModeToggle
    showRangeSelection
    source={{ id: source.id, isFeed: false }}
    title="Import errors"
  ></ImportStats>
{/if}
