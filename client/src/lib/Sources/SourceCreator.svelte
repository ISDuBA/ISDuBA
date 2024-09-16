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
    getSourceName,
    parseFeeds,
    saveSource,
    saveFeeds
  } from "$lib/Sources/source";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Input, Label, Button, Spinner, List, DescriptionList } from "flowbite-svelte";
  import SourceForm from "./SourceForm.svelte";
  import type { CSAFProviderMetadata } from "$lib/provider";
  import { push } from "svelte-spa-router";
  import FeedView from "./FeedView.svelte";
  import { onMount } from "svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";

  export let params: any = null;

  let errorMessage: ErrorDetails | null;

  let sourceForm: any;
  let updateSourceForm: any;

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

  let source: Source = {
    name: "",
    url: "",
    rate: undefined,
    slots: undefined,
    active: false,
    strict_mode: true,
    headers: [""],
    ignore_patterns: [""]
  };

  let formClass = "max-w-[800pt]";
  const dtClass: string = "ml-1 mt-1 text-gray-500 md:text-sm dark:text-gray-400";
  const ddClass: string = "break-words font-semibold ml-2 mb-1";
  let loadingPMD: boolean = false;

  let pmd: CSAFProviderMetadata | null = null;
  let pmdFeeds: Feed[] = [];

  const loadPMD = async () => {
    loadingPMD = true;
    let result = await fetchPMD(source.url);
    loadingPMD = false;
    if (result.ok) {
      if (!params?.domain) {
        push(`/sources/new/${encodeURIComponent(source.url)}`);
      } else {
        pmd = result.value;
        source.name = await getSourceName(pmd);
        pmdFeeds = parseFeeds(pmd, []);
      }
    } else {
      errorMessage = result.error;
    }
  };

  const checkUrl = async () => {};

  const saveAll = async () => {
    updateSourceForm();
    if (source.age === "") {
      source.age = undefined;
    }
    let result = await saveSource(source);
    if (!result.ok) {
      errorMessage = result.error;
      return;
    }
    let feedResult = await saveFeeds(source, pmdFeeds);
    if (!feedResult.ok) {
      errorMessage = feedResult.error;
    }
    push(`/sources/`);
  };

  onMount(async () => {
    let domain = params?.domain;
    if (domain) {
      source.url = domain;
      await loadPMD();
      updateSourceForm = sourceForm.updateSource;
    }
  });
</script>

<svelte:head>
  <title>Sources - Add source</title>
</svelte:head>

<SectionHeader title="Add new CSAF trusted provider"></SectionHeader>
{#if params?.domain}
  <List tag="dl" class="divide-y divide-gray-200 text-sm 2xl:w-max">
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

  <SourceForm bind:this={sourceForm} {formClass} {source}></SourceForm>
  <FeedView feeds={pmdFeeds}></FeedView>

  <Button on:click={saveAll} color="light">
    <i class="bx bxs-save me-2"></i>
    <span>Save source</span>
  </Button>
{:else}
  <form on:submit={loadPMD} class={formClass}>
    <Label>Domain/PMD</Label>
    <Input bind:value={source.url} on:input={checkUrl} color={urlColor}></Input>
    <br />
    <div class:hidden={!loadingPMD} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button type="submit" color="light">
      <i class="bx bx-check me-2"></i>
      <span>Search and load provider metadata</span>
    </Button>
  </form>
{/if}

<br />
<ErrorMessage error={errorMessage}></ErrorMessage>
