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
    saveFeeds,
    resetSourceAttention,
    dtClass,
    ddClass
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
  import CBadge from "$lib/Components/CBadge.svelte";

  interface Props {
    params: any;
  }
  let { params = null }: Props = $props();

  const uid = $props.id();

  const shortLoadInterval = 5;
  const longLoadMultiplier = 6;

  let sourceEdited: boolean = $state(false);

  let modalOpen: boolean = $state(false);
  let modalMessage = $state("");
  let modalTitle = $state("");
  let modalCallback: any = $state(null);

  let saveSourceError: ErrorDetails | null = $state(null);
  let loadSourceError: ErrorDetails | null = $state(null);
  let loadFeedError: ErrorDetails | null = $state(null);
  let saveFeedError: ErrorDetails | null = $state(null);
  let loadPmdError: ErrorDetails | null = $state(null);
  let feedError: ErrorDetails | null = $state(null);
  let pmd: CSAFProviderMetadata | null = $state(null);
  let feeds: Feed[] = $state([]);

  let loadingFeeds: boolean = $state(false);
  let loadingSource: boolean = $state(false);
  let loadingPMD: boolean = $state(false);

  let formClass = "max-w-[800pt]";

  let sourceForm: any = $state(null);
  let updateSourceForm: any = $state(null);
  let fillAgeDataFromSource: (source: Source) => void;

  let loadSource: boolean = $state(true);

  const initialSource = {
    name: "",
    url: "",
    active: false,
    rate: 1,
    slots: 2,
    strict_mode: true,
    headers: [""],
    ignore_patterns: [""],
    attention: false,
    client_cert_passphrase: ""
  };
  let source: Source = $state(structuredClone(initialSource));

  let oldSource: Source = $state(structuredClone(initialSource));

  let sourceStatFull: SourceBasicStats | null = $state(null);
  let sourceStatLast: SourceBasicStats | null = $state(null);

  let updateIteration = 0;
  let updateStats = setInterval(async () => {
    updateIteration = (updateIteration + 1) % 6;
    if (updateIteration == 0 && sourceStatFull && sourceStatLast) {
      sourceStatFull.reload();
      sourceStatLast.reload();
    }
    if (!source.id || source.id === 0) {
      return;
    }
    let result = await fetchSource(source.id, true);
    if (result.ok) {
      source.stats = result.value.stats;
    }
    let feedResult = await fetchFeeds(source.id ?? 0, true, true);
    if (feedResult.ok) {
      for (let feed of feedResult.value) {
        const find = feeds.find((i) => i.id === feed.id);
        if (find) {
          find.stats = feed.stats;
        }
      }
    }
    feeds = feeds;
  }, shortLoadInterval * 1000);

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
      oldSource = structuredClone($state.snapshot(source));
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
      await markAsDone();
    } else {
      saveFeedError = result.error;
    }
  };

  const markAsDone = async () => {
    let result = await resetSourceAttention(source);
    if (!result.ok) {
      saveSourceError = result.error;
    }

    await loadSourceInfo(source.id ?? 0);
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

    if (tmpA.secure === undefined) {
      tmpA.secure = false;
    }
    if (tmpB.secure === undefined) {
      tmpB.secure = false;
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
    if (sourceEqual($state.snapshot(oldSource), $state.snapshot(source))) {
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
        source.id = 0;
      } else {
        await loadSourceInfo(sourceID);
        await loadPMD();
      }
      await loadFeeds();
      if (source.id !== 0 && pmd) {
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
      onclick={() => {
        modalCallback();
      }}
      color="red"
      class="me-2">Yes, I'm sure</Button
    >
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

<SectionHeader title={source.name}></SectionHeader>
{#if source.id !== 0}
  <div class="mb-3 grid w-full grid-cols-1 justify-stretch gap-10 lg:grid-cols-2">
    <div class="w-full">
      <List tag="dl" class="w-full divide-y divide-gray-200 text-sm">
        <div>
          <DescriptionList tag="dt" class={dtClass}>Domain/PMD</DescriptionList>
          <DescriptionList tag="dd" class={ddClass}>{source.url}</DescriptionList>
        </div>
        {#if pmd}
          <div>
            <DescriptionList tag="dt" class={dtClass}>Canonical URL</DescriptionList>
            <DescriptionList tag="dd" class={ddClass}>{pmd.canonical_url}</DescriptionList>
          </div>
          <div>
            <DescriptionList tag="dt" class={dtClass}>Publisher Name</DescriptionList>
            <DescriptionList tag="dd" class={ddClass}>{pmd.publisher.name}</DescriptionList>
          </div>
          <div>
            <DescriptionList tag="dt" class={dtClass}>Publisher Contact</DescriptionList>
            <DescriptionList tag="dd" class={ddClass}
              >{pmd.publisher.contact_details}</DescriptionList
            >
          </div>
          <div>
            {#if pmd.publisher.issuing_authority}
              <DescriptionList tag="dt" class={dtClass}>Issuing Authority</DescriptionList>
              <DescriptionList tag="dd" class={ddClass}
                >{pmd.publisher.issuing_authority}</DescriptionList
              >
            {/if}
          </div>
        {/if}
        <div>
          <DescriptionList tag="dt" class={dtClass}>Status</DescriptionList>
          {#if source.status}
            {#each source.status as s, i (`sourceeditor-${uid}-${i}`)}
              <DescriptionList tag="dd" class={ddClass}>{s}</DescriptionList>
            {/each}
          {:else}
            <DescriptionList tag="dd" class={ddClass}>OK</DescriptionList>
          {/if}
        </div>
      </List>
      {#if source.stats}
        <h4 class="mt-3">Status</h4>
        <div class="grid w-full grid-cols-[max-content_max-content_max-content] gap-x-4 text-sm">
          <DescriptionList tag="dt" class={dtClass}>Loading</DescriptionList>
          <DescriptionList tag="dt" class={dtClass + " mr-1"}>Queued</DescriptionList>
          <div>
            <DescriptionList tag="dt" class={dtClass + " mr-1"}>Imported (last 24h)</DescriptionList
            >
            <div class="mt-1 mb-1 h-1 min-h-1">
              <div class="progressmeter">
                <span class="w-full"
                  ><span
                    style="animation-duration: {shortLoadInterval * longLoadMultiplier}s"
                    class="infiniteprogress bg-primary-500"
                  ></span></span
                >
              </div>
            </div>
          </div>
          <DescriptionList tag="dd" class={ddClass}>{source.stats.downloading}</DescriptionList>
          <DescriptionList tag="dd" class={ddClass}>{source.stats.waiting}</DescriptionList>
          <DescriptionList tag="dd" class={ddClass}>
            {#if source.id}
              {@const yesterday = Date.now() - DAY_MS}
              <SourceBasicStats bind:this={sourceStatFull} sourceID={source.id}></SourceBasicStats>
              (<SourceBasicStats
                bind:this={sourceStatLast}
                from={new Date(yesterday)}
                sourceID={source.id}
              ></SourceBasicStats>)
            {/if}
          </DescriptionList>
        </div>
      {/if}
      <ErrorMessage error={loadSourceError}></ErrorMessage>
      <ErrorMessage error={loadPmdError}></ErrorMessage>

      <div
        class:invisible={!loadingPMD}
        class={!loadingPMD ? "loadingFadeIn" : ""}
        class:mb-4={true}
      >
        Loading PMD ...
        <Spinner color="gray" size="4"></Spinner>
      </div>
    </div>

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
        {oldSource}
        {formClass}
        enableActive={true}
      ></SourceForm>
      <Button disabled={!sourceEdited} onclick={updateSource} color="green">
        <i class="bx bxs-save me-2"></i>
        <span>Save source</span>
      </Button>
      <Button
        onclick={(event: any) => {
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
  </div>
{/if}

<FeedView {feeds} placeholderFeed={source.id === 0} {clickFeed} {updateFeed} edit={true}>
  {#snippet feedViewTopSlot()}
    <div>
      {#if source.attention}
        <CBadge class="mb-2 h-fit p-1" dismissable>
          <p>
            These are the currently available feeds. Please review them and adjust the subscriptions
            if needed.
          </p>
          {#snippet closeButtonSlot()}
            <Button
              color="light"
              class="border-primary-700/55 text-primary-700 ms-1 min-h-[26px] min-w-[26px] rounded border bg-transparent p-0 hover:bg-white/50 dark:bg-transparent dark:hover:bg-white/20"
              onclick={async () => {
                markAsDone();
                close();
              }}
            >
              <i class="bx bx-check"></i>
            </Button>
          {/snippet}
        </CBadge>
      {/if}
    </div>
  {/snippet}
</FeedView>
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

{#if source.id !== undefined}
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
