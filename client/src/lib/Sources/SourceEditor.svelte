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
    deleteFeed,
    fetchFeeds,
    logLevels,
    calculateMissingFeeds
  } from "$lib/Sources/source";
  import {
    Input,
    Label,
    Button,
    TableBodyCell,
    Spinner,
    Modal,
    Select,
    Table,
    TableBodyRow
  } from "flowbite-svelte";
  import { tdClass } from "$lib/Table/defaults";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails } from "$lib/Errors/error";
  import type { CSAFProviderMetadata } from "$lib/provider";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { onMount } from "svelte";
  import SourceForm from "./SourceForm.svelte";
  import FeedView from "./FeedView.svelte";
  import { push } from "svelte-spa-router";
  export let params: any = null;

  let modalOpen: boolean = false;
  let modalMessage = "";
  let modalTitle = "";
  let modalCallback: any;

  let saveSourceError: ErrorDetails | null;
  let loadSourceError: ErrorDetails | null;
  let loadFeedError: ErrorDetails | null;
  let loadPmdError: ErrorDetails | null;
  let feedError: ErrorDetails | null;
  let pmd: CSAFProviderMetadata;
  let pmdFeeds: Feed[] = [];
  let missingFeeds: Feed[] = [];
  let feeds: Feed[] = [];

  let feedEdit: Feed | null;

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
    let result = await fetchSource(Number(id));
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
    let result = await fetchFeeds(source.id);
    if (result.ok) {
      feeds = result.value;
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

  onMount(async () => {
    let id = params?.id;
    if (id) {
      await loadSourceInfo(Number(id));
      await loadPMD();
      await loadFeeds();
      missingFeeds = calculateMissingFeeds(pmdFeeds, feeds);

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
<div class="flex">
  <div class="w-1/2 flex-auto">
    <CustomTable
      title="Feeds"
      headers={[
        {
          label: "Label",
          attribute: "label"
        },
        {
          label: "Domain/PMD",
          attribute: "url"
        },
        {
          label: "Rolie",
          attribute: "rolie"
        },
        {
          label: "Log level",
          attribute: "log_level"
        }
      ]}
    >
      {#each feeds as feed, index (index)}
        <tr
          on:click={() => {
            if (feed.id) {
              push(`/sources/feed/${feed.id}`);
            }
          }}
          class="cursor-pointer"
        >
          <TableBodyCell {tdClass}>{feed.label}</TableBodyCell>
          <TableBodyCell {tdClass}>{feed.url}</TableBodyCell>
          <TableBodyCell {tdClass}>{feed.rolie}</TableBodyCell>
          <TableBodyCell {tdClass}>{feed.log_level}</TableBodyCell>
          <td>
            <Button
              on:click={() => {
                feedEdit = feed;
              }}
              title={`Edit feed "${feed.label}"`}
              class="border-0 p-2"
              color="light"
            >
              <i class="bx bx-edit text-xl"></i>
            </Button>
          </td>
          <td>
            <Button
              on:click={(event) => {
                event.stopPropagation();
                modalCallback = async () => {
                  if (feed.id) {
                    await deleteFeed(feed.id);
                  }
                };
                modalMessage = "Are you sure you want to delete this feed?";
                modalTitle = `Feed ${feed.label}`;
                modalOpen = true;
              }}
              title={`Delete feed "${feed.label}"`}
              class="border-0 p-2"
              color="light"
            >
              <i class="bx bx-trash text-xl text-red-500"></i>
            </Button>
          </td>
        </tr>
      {/each}
      <div slot="bottom">
        <div class:hidden={!loadingFeeds} class:mb-4={true}>
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
        <ErrorMessage error={feedError}></ErrorMessage>
      </div>
    </CustomTable>
  </div>

  <div class="w-1/2 flex-auto">
    <CustomTable
      title="Missing feeds"
      headers={[
        {
          label: "Label",
          attribute: "label"
        },
        {
          label: "Domain/PMD",
          attribute: "url"
        },
        {
          label: "Rolie",
          attribute: "rolie"
        },
        {
          label: "Log level",
          attribute: "log_level"
        }
      ]}
    >
      {#each missingFeeds as feed, index (index)}
        <tr
          class="cursor-pointer"
          on:click={() => {
            feedEdit = feed;
          }}
        >
          <TableBodyCell {tdClass}>{feed.label}</TableBodyCell>
          <TableBodyCell {tdClass}>{feed.url}</TableBodyCell>
          <TableBodyCell {tdClass}>{feed.rolie}</TableBodyCell>
          <TableBodyCell {tdClass}>{feed.log_level}</TableBodyCell>
        </tr>
      {/each}
      <div slot="bottom">
        <div class:hidden={!loadingFeeds && !loadingPMD} class:mb-4={true}>
          Loading ...
          <Spinner color="gray" size="4"></Spinner>
        </div>
        <ErrorMessage error={feedError}></ErrorMessage>
      </div>
    </CustomTable>
  </div>
</div>
<!-- TODO: Move into feed viewer -->
{#if feedEdit}
  <SectionHeader title={feedEdit.enable ? "New feed" : "Edit feed"}>
    <div slot="right">
      <slot name="header-right"></slot>
    </div>
  </SectionHeader>
  <form
    on:submit={async () => {
      if (feedEdit) {
        feedEdit.enable = true;
        feedEdit = null;
      }
    }}
    class={formClass}
  >
    <Label>URL</Label>
    <Input readonly bind:value={feedEdit.url}></Input>
    <Label>Log level</Label>
    <Select items={logLevels} bind:value={feedEdit.log_level} />
    <Label>Label</Label>
    <Input bind:value={feedEdit.label}></Input>
    <br />
    <Button type="submit" color="light">
      <i class="bx bxs-save me-2"></i>
      <span>Save feed</span>
    </Button>
  </form>{/if}
<br />

<FeedView feeds={pmdFeeds}></FeedView>

<ErrorMessage error={saveSourceError}></ErrorMessage>
<ErrorMessage error={loadSourceError}></ErrorMessage>
<ErrorMessage error={loadPmdError}></ErrorMessage>
<ErrorMessage error={loadFeedError}></ErrorMessage>
