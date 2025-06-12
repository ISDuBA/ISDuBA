<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { type Feed, getLogLevels, LogLevel } from "$lib/Sources/source";
  import { Select, Input, TableBodyCell, Button, Spinner } from "flowbite-svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { onMount } from "svelte";
  import { tdClass, type TableHeader } from "$lib/Table/defaults";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";
  import CIconButton from "$lib/Components/CIconButton.svelte";

  export let placeholderFeed: boolean = false;
  export let feeds: Feed[] = [];
  export let edit: boolean = false;

  export let updateFeed = async (_feed: Feed) => {};
  export let clickFeed = async (_feed: Feed) => {};
  export let showProgress = false;

  $: subscribedFeeds = feeds.filter((feed) => feed.enable);

  $: {
    let loadingHeader = headersEdit.find((header) => header.label == "Loading/Queued");
    if (loadingHeader) {
      if (showProgress) {
        // loadingHeader.progressDuration = shortLoadInterval;
      } else {
        loadingHeader.progressDuration = undefined;
      }
    }
  }

  // const shortLoadInterval = 5;

  let headers: TableHeader[] = [
    {
      label: "",
      attribute: "enable"
    },
    {
      label: "URL",
      attribute: "url"
    },
    {
      label: "Log level",
      attribute: "log_level",
      class: "min-w-32"
    },
    {
      label: "Label",
      attribute: "label",
      class: "min-w-32"
    }
  ];

  let headersEdit: TableHeader[] = [
    ...headers,
    { label: "Loading/Queued", attribute: "stats" },
    { label: "Healthy", attribute: "healthy" },
    { label: "Logs", attribute: "logs" }
  ];

  let headerPlaceholder = headersEdit.filter(
    (i) => i.attribute === "label" || i.attribute === "logs"
  );

  let tableHeaders = headers;
  $: if (edit !== undefined || placeholderFeed !== undefined) {
    if (placeholderFeed) {
      tableHeaders = headerPlaceholder;
    } else if (edit) {
      tableHeaders = headersEdit;
    } else {
      tableHeaders = headers;
    }
  }

  let logLevels: { value: LogLevel; name: string }[] = [];

  let loadConfigError: ErrorDetails | null;

  let feedBlinkID = -1;
  let isSubscribingAll = false;
  let isUnSubscribingAll = false;

  onMount(async () => {
    const resp = await getLogLevels(!edit);
    if (resp.ok) {
      logLevels = resp.value;
    } else {
      loadConfigError = resp.error;
    }
    feedBlinkID = Number(sessionStorage.getItem("feedBlinkID") ?? "-1");
    sessionStorage.removeItem("feedBlinkID");
  });

  const changeAllSubscriptions = async (subscribe = true) => {
    if (subscribe) {
      isSubscribingAll = true;
    } else {
      isUnSubscribingAll = true;
    }
    const feedsCopy = feeds;
    for (let i = 0; i < feedsCopy.length; i++) {
      const feed = feedsCopy[i];
      if (feed.enable !== subscribe) {
        feed.enable = subscribe;
        await updateFeed(feed);
        if (!subscribe) {
          feeds[i].id = undefined;
        }
      }
    }
    if (subscribe) {
      isSubscribingAll = false;
    } else {
      isUnSubscribingAll = false;
    }
    feeds = feeds;
  };
</script>

{#if logLevels}
  <CustomTable title="Feeds" headers={tableHeaders}>
    <div slot="top">
      <slot name="top"></slot>
    </div>
    <div slot="header-right" class="flex gap-2">
      <Button
        on:click={() => {
          changeAllSubscriptions();
        }}
        class="flex gap-2"
        color="light"
        disabled={subscribedFeeds.length === feeds.length || isSubscribingAll || isUnSubscribingAll}
        size="sm"
      >
        Subscribe all
        {#if isSubscribingAll}
          <Spinner color="gray" size="4"></Spinner>
        {/if}
      </Button>
      <Button
        on:click={() => {
          changeAllSubscriptions(false);
        }}
        class="flex gap-2"
        color="light"
        disabled={subscribedFeeds.length === 0 || isSubscribingAll || isUnSubscribingAll}
        size="sm"
        >Unsubscribe all
        {#if isUnSubscribingAll}
          <Spinner color="gray" size="4"></Spinner>
        {/if}
      </Button>
    </div>
    {#each feeds as feed, index (index)}
      <tr class={feed.id === feedBlinkID ? "blink" : ""}>
        {#if placeholderFeed}
          <TableBodyCell {tdClass}>{feed.label}</TableBodyCell>
          <TableBodyCell on:click={async () => await clickFeed(feed)} {tdClass}>
            <a href={"javascript:void(0);"} on:click={async () => await clickFeed(feed)}>
              <i class="bx bx-archive"> </i></a
            >
          </TableBodyCell>
        {:else}
          <TableBodyCell {tdClass}>
            {#if feed.enable && !placeholderFeed}
              <CIconButton
                on:click={async () => {
                  feed.enable = false;
                  await updateFeed(feed);
                  feed.id = undefined;
                }}
                icon="trash"
              ></CIconButton>
            {:else}
              <CIconButton
                on:click={async () => {
                  feed.enable = true;
                  await updateFeed(feed);
                }}
                icon="plus"
              ></CIconButton>
            {/if}
          </TableBodyCell>
          <TableBodyCell
            on:click={async () => await clickFeed(feed)}
            tdClass={`${tdClass} break-all whitespace-normal`}
          >
            {#if edit && feed.enable}
              <a href={"javascript:void(0);"} on:click={async () => await clickFeed(feed)}
                >{feed.url}</a
              >
            {:else}
              <span class="text-amber-600">
                {feed.url}
              </span>
            {/if}
          </TableBodyCell>
          <TableBodyCell {tdClass}
            ><Select
              items={logLevels}
              bind:value={feed.log_level}
              on:change={async () => await updateFeed(feed)}
            /></TableBodyCell
          >
          {#if edit && !feed.enable}
            <TableBodyCell {tdClass}>N/A</TableBodyCell>
          {:else}
            <TableBodyCell {tdClass}
              ><Input bind:value={feed.label} on:input={async () => await updateFeed(feed)}
              ></Input></TableBodyCell
            >
          {/if}
          {#if edit}
            <TableBodyCell {tdClass}
              >{(feed.stats?.downloading ?? 0) + "/" + (feed.stats?.waiting ?? 0)}</TableBodyCell
            >
            <TableBodyCell {tdClass}
              ><i class={"bx " + (feed.healthy ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
            >
            {#if feed.enable}
              <TableBodyCell on:click={async () => await clickFeed(feed)} {tdClass}>
                <a href={"javascript:void(0);"} on:click={async () => await clickFeed(feed)}>
                  <i class="bx bx-archive"> </i></a
                >
              </TableBodyCell>
            {/if}
          {/if}
        {/if}
      </tr>
    {/each}
  </CustomTable>
{/if}
<ErrorMessage error={loadConfigError}></ErrorMessage>
