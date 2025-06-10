<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Badge, Button, Spinner, TableBodyCell } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { push } from "svelte-spa-router";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails, getErrorDetails } from "$lib/Errors/error";
  import { tdClass } from "$lib/Table/defaults";
  import { request } from "$lib/request";
  import { onDestroy, onMount } from "svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { type Source, fetchSources } from "$lib/Sources/source";
  import { appStore } from "$lib/store";
  import { DAY_MS } from "$lib/time";
  import ImportStats from "$lib/Statistics/ImportStats.svelte";
  import SourceBasicStats from "./SourceBasicStats.svelte";

  const shortLoadInterval = 5;
  const longLoadMultiplier = 6;

  let messageError: ErrorDetails | null;
  let sourcesError: ErrorDetails | null;

  let loadingSources: boolean = false;

  let sources: Source[] = [];
  let statsComponents: { [idToRole: string]: SourceBasicStats } = {};

  async function getMessage() {
    const response = await request("api/sources/message", "GET");
    if (response.ok) {
      return response.content;
    } else {
      messageError = getErrorDetails(`Couldn't load default message`, response);
    }
    return new Map<string, [string]>();
  }

  let updateIteration = 0;
  let sourceUpdate = setInterval(async () => {
    updateIteration = (updateIteration + 1) % longLoadMultiplier;
    if (updateIteration == 0) {
      for (let comp of Object.values(statsComponents)) {
        comp.reload();
      }
    }
    if (appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager()) {
      getSources();
    }
  }, shortLoadInterval * 1000);

  const getSources = async () => {
    loadingSources = true;
    const result = await fetchSources(true, true);
    loadingSources = false;
    if (result.ok) {
      sources = result.value;
    } else {
      sources = [];
      sourcesError = result.error;
    }
  };

  onMount(async () => {
    if (appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager()) {
      await getSources();
    }
  });

  onDestroy(() => {
    clearInterval(sourceUpdate);
  });
</script>

<svelte:head>
  <title>Sources</title>
</svelte:head>

<div>
  <SectionHeader title="Sources"></SectionHeader>
  <ImportStats axes={[{ label: "Imports", types: ["imports"] }]} divContainerClass="mb-8" title=""
  ></ImportStats>
  {#if appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager()}
    <CustomTable
      title="CSAF Provider"
      headers={[
        {
          label: "",
          attribute: "icon"
        },
        {
          label: "Name",
          attribute: "name"
        },
        {
          label: "Domain/PMD",
          attribute: "url"
        },
        {
          label: "Active",
          attribute: "active"
        },
        {
          label: "Loading/Queued",
          attribute: "stats"
        },
        {
          label: "Healthy",
          attribute: "healthy"
        },
        {
          label: "Imported (last 24h)",
          attribute: "statsHistory",
          progressDuration: shortLoadInterval * longLoadMultiplier
        }
      ]}
    >
      {#each sources as source, index (index)}
        <tr
          on:click={() => {
            if (appStore.isSourceManager()) {
              push(`/sources/${source.id}`);
            }
          }}
          on:blur={() => {}}
          on:focus={() => {}}
          class={appStore.isSourceManager() ? "cursor-pointer" : ""}
        >
          <TableBodyCell tdClass="w-fit max-w-10">
            <i class="bx bx-git-repo-forked"></i>
          </TableBodyCell>
          <TableBodyCell tdClass={`${tdClass} break-words hyphens-auto`}>
            <div class="whitespace-break-spaces">
              <span>{source.name}</span>
              {#if source.attention}
                <Badge class="min-w-fit">Source changed</Badge>
              {/if}
            </div>
          </TableBodyCell>
          {#if source.id !== 0}
            <TableBodyCell tdClass={`${tdClass} break-all whitespace-normal`}
              >{source.url}</TableBodyCell
            >
            <TableBodyCell {tdClass}
              ><i class={"bx " + (source.active ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
            >
            <TableBodyCell {tdClass}
              >{source.stats?.downloading}/{source.stats?.waiting}</TableBodyCell
            >
            <TableBodyCell {tdClass}
              ><i class={"bx " + (source.healthy ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
            >
          {:else}
            <TableBodyCell></TableBodyCell>
            <TableBodyCell></TableBodyCell>
            <TableBodyCell></TableBodyCell>
            <TableBodyCell></TableBodyCell>
          {/if}
          <TableBodyCell>
            {#if source.id !== undefined}
              {@const yesterday = Date.now() - DAY_MS}
              <SourceBasicStats bind:this={statsComponents[`${source.id}full`]} sourceID={source.id}
              ></SourceBasicStats>
              (<SourceBasicStats
                bind:this={statsComponents[`${source.id}last`]}
                from={new Date(yesterday)}
                sourceID={source.id}
              ></SourceBasicStats>)
            {/if}
          </TableBodyCell>
        </tr>
      {/each}
      <div slot="bottom">
        <div
          class:invisible={!loadingSources}
          class={loadingSources ? "loadingFadeIn" : ""}
          class:mb-4={true}
        >
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
        {#if appStore.isSourceManager()}
          <Button href="/#/sources/new" class="mb-2" color="primary" size="xs">
            <i class="bx bx-plus"></i>
            <span>Add source</span>
          </Button>
        {/if}
        <ErrorMessage error={sourcesError}></ErrorMessage>
      </div>
    </CustomTable>
  {/if}
  {#await getMessage() then resp}
    {#if resp.message}
      {#if !appStore.isSourceManager() && !appStore.isImporter()}
        {resp.message}
      {/if}
    {/if}
  {/await}
  <ErrorMessage error={sourcesError}></ErrorMessage>
  <ErrorMessage error={messageError}></ErrorMessage>

  <br />
  {#if appStore.isImporter()}
    <Button href="/#/sources/upload" class="my-2" color="primary" size="xs">
      <i class="bx bx-upload"></i>
      <span>Upload documents</span>
    </Button>
  {/if}
</div>
