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
  import { onMount } from "svelte";
  import SourceForm from "./SourceForm.svelte";
  import { request } from "$lib/request";
  import FeedView from "./FeedView.svelte";
  import { push } from "svelte-spa-router";
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

  const dtClass: string = "ml-1 mt-1 text-gray-500 md:text-sm dark:text-gray-400";
  const ddClass: string = "break-words font-semibold ml-2 mb-1";

  const loadSourceInfo = async (id: number) => {
    loadingSource = true;
    let result = await fetchSource(Number(id), true);
    if (result.ok) {
      source = result.value;
      if (fillAgeDataFromSource) {
        fillAgeDataFromSource(source);
      }
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
    if (!source.id) {
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

  const inputChange = () => {
    sourceEdited = true;
  };

  const clickFeed = async (feed: Feed) => {
    if (!feed.id) {
      return;
    }
    push(`/sources/logs/${feed.id}`);
  };

  onMount(async () => {
    let id = params?.id;
    if (id) {
      await loadSourceInfo(Number(id));
      await loadPMD();
      await loadFeeds();
      let missingFeeds = calculateMissingFeeds(parseFeeds(pmd, feeds), feeds);
      missingFeeds.map((f) => {
        f.enable = false;
      });
      feeds.push(...missingFeeds);
      feeds = feeds;

      updateSourceForm = sourceForm.updateSource;
      fillAgeDataFromSource = sourceForm.fillAgeDataFromSource;
      fillAgeDataFromSource(source);
    }
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
    </List>
    {#if source.stats}
      <h4 class="mt-3">Statistics</h4>
      <List tag="dl" class="w-full divide-y divide-gray-200 text-sm">
        <div>
          <DescriptionList tag="dt" {dtClass}>Total downloading</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{source.stats.downloading}</DescriptionList>
        </div>
        <div>
          <DescriptionList tag="dt" {dtClass}>Total waiting</DescriptionList>
          <DescriptionList tag="dd" {ddClass}>{source.stats.waiting}</DescriptionList>
        </div>
      </List>
    {/if}

    <div class:invisible={!loadingPMD} class={!loadingPMD ? "loadingFadeIn" : ""} class:mb-4={true}>
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
    <SourceForm bind:this={sourceForm} {inputChange} {source} {formClass} enableActive={true}
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
  </div>
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
<ErrorMessage error={feedError}></ErrorMessage>
<ErrorMessage error={saveFeedError}></ErrorMessage>

<ErrorMessage error={saveSourceError}></ErrorMessage>
<ErrorMessage error={loadSourceError}></ErrorMessage>
<ErrorMessage error={loadPmdError}></ErrorMessage>
<ErrorMessage error={loadFeedError}></ErrorMessage>
