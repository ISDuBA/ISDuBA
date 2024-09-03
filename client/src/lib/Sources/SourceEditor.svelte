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
  import { Button, TableBodyCell, Spinner, Modal, Table, TableBodyRow } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails } from "$lib/Errors/error";
  import type { CSAFProviderMetadata } from "$lib/provider";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { onMount } from "svelte";
  import SourceForm from "./SourceForm.svelte";
  import FeedView from "./FeedView.svelte";
  export let params: any = null;

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

  const loadSourceInfo = async (id: number) => {
    loadingSource = true;
    let result = await fetchSource(Number(id), true);
    if (result.ok) {
      source = result.value;
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
  };

  const updateFeed = async (feed: Feed) => {
    let result = await saveFeeds(source, [feed]);
    if (result.ok) {
      feed.id = result.value[0];
    } else {
      saveFeedError = result.error;
    }
  };

  onMount(async () => {
    let id = params?.id;
    if (id) {
      await loadSourceInfo(Number(id));
      await loadPMD();
      await loadFeeds();
      let missingFeeds = calculateMissingFeeds(parseFeeds(pmd), feeds);
      missingFeeds.map((f) => {
        f.enable = false;
      });
      feeds.push(...missingFeeds);
      feeds = feeds;

      updateSourceForm = sourceForm.updateSource;
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
<div class="flex">
  <div class="flex-auto">
    <Table class="2xl:w-max" noborder>
      <TableBodyRow>
        <TableBodyCell>Domain/PMD</TableBodyCell>
        <TableBodyCell>{source.url}</TableBodyCell>
      </TableBodyRow>
      {#if pmd}
        <TableBodyRow>
          <TableBodyCell>Canonical URL</TableBodyCell>
          <TableBodyCell>{pmd.canonical_url}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>Publisher Name</TableBodyCell>
          <TableBodyCell>{pmd.publisher.name}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>Publisher Contact</TableBodyCell>
          <TableBodyCell>{pmd.publisher.contact_details}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>Issuing Authority</TableBodyCell>
          <TableBodyCell>{pmd.publisher.issuing_authority}</TableBodyCell>
        </TableBodyRow>
      {/if}
      {#if source.stats}
        <TableBodyRow>
          <TableBodyCell>Downloading</TableBodyCell>
          <TableBodyCell>{source.stats.downloading}</TableBodyCell>
        </TableBodyRow>
        <TableBodyRow>
          <TableBodyCell>Waiting</TableBodyCell>
          <TableBodyCell>{source.stats.waiting}</TableBodyCell>
        </TableBodyRow>
      {/if}
    </Table>
    <div class:hidden={!loadingPMD} class:mb-4={true}>
      Loading PMD ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  </div>

  <div class="flex-auto">
    <div class:hidden={!loadingSource} class:mb-4={true}>
      Loading source configuration ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <SourceForm bind:this={sourceForm} {source} {formClass} enableActive={true}></SourceForm>
    <Button on:click={updateSource} color="light">
      <i class="bx bxs-save me-2"></i>
      <span>Save source</span>
    </Button>
  </div>
</div>

<FeedView {feeds} {updateFeed} edit={true}></FeedView>
<div class:hidden={!loadingFeeds && !loadingPMD} class:mb-4={true}>
  Loading ...
  <Spinner color="gray" size="4"></Spinner>
</div>
<ErrorMessage error={feedError}></ErrorMessage>
<ErrorMessage error={saveFeedError}></ErrorMessage>

<ErrorMessage error={saveSourceError}></ErrorMessage>
<ErrorMessage error={loadSourceError}></ErrorMessage>
<ErrorMessage error={loadPmdError}></ErrorMessage>
<ErrorMessage error={loadFeedError}></ErrorMessage>
