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
  import type { Snippet } from "svelte";

  interface Props {
    edit?: boolean;
    feeds?: Feed[];
    placeholderFeed?: boolean;
    updateFeed?: (feed: Feed) => Promise<void>;
    clickFeed?: (feed: Feed) => Promise<void>;
    feedViewTopSlot?: Snippet;
  }

  let {
    edit = false,
    feeds = [],
    placeholderFeed = false,
    updateFeed = async (_feed: Feed) => {},
    clickFeed = async (_feed: Feed) => {},
    feedViewTopSlot = undefined
  }: Props = $props();

  let subscribedFeeds = $derived(feeds.filter((feed) => feed.enable));

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

  let tableHeaders = $derived.by(() => {
    if (edit !== undefined || placeholderFeed !== undefined) {
      if (placeholderFeed) {
        return headerPlaceholder;
      } else if (edit) {
        return headersEdit;
      } else {
        return headers;
      }
    } else {
      return [];
    }
  });

  let logLevels: { value: LogLevel; name: string }[] = $state([]);

  let loadConfigError: ErrorDetails | null = $state(null);

  let feedBlinkID = $state(-1);
  let isSubscribingAll = $state(false);
  let isUnSubscribingAll = $state(false);

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
    {#snippet topSlot()}
      {#if feedViewTopSlot}
        <div>
          {@render feedViewTopSlot()}
        </div>
      {/if}
    {/snippet}
    {#snippet headerRightSlot()}
      <div class="flex gap-2">
        {#if !placeholderFeed}
          <Button
            onclick={() => {
              changeAllSubscriptions();
            }}
            class="flex gap-2"
            color="light"
            disabled={subscribedFeeds.length === feeds.length ||
              isSubscribingAll ||
              isUnSubscribingAll}
            size="sm"
          >
            Subscribe all
            {#if isSubscribingAll}
              <Spinner color="gray" size="4"></Spinner>
            {/if}
          </Button>
          <Button
            onclick={() => {
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
        {/if}
      </div>
    {/snippet}
    {#snippet mainSlot()}
      {#each feeds as feed, index (index)}
        <tr class={feed.id === feedBlinkID ? "blink" : ""}>
          {#if placeholderFeed}
            <TableBodyCell class={tdClass}>{feed.label}</TableBodyCell>
            <TableBodyCell onclick={async () => await clickFeed(feed)} class={tdClass}>
              <a
                href={"javascript:void(0);"}
                onclick={async () => await clickFeed(feed)}
                aria-label="View feed archive"
              >
                <i class="bx bx-archive"> </i></a
              >
            </TableBodyCell>
          {:else}
            <TableBodyCell class={tdClass}>
              {#if feed.enable && !placeholderFeed}
                <CIconButton
                  onClicked={async () => {
                    feed.enable = false;
                    await updateFeed(feed);
                    feed.id = undefined;
                  }}
                  icon="trash"
                ></CIconButton>
              {:else}
                <CIconButton
                  onClicked={async () => {
                    feed.enable = true;
                    await updateFeed(feed);
                  }}
                  ariaLabel={`Enable feed with label ${feed.label}`}
                  icon="plus"
                ></CIconButton>
              {/if}
            </TableBodyCell>
            <TableBodyCell
              onclick={async () => await clickFeed(feed)}
              class={`${tdClass} break-all whitespace-normal`}
            >
              {#if edit && feed.enable}
                <a
                  href={"javascript:void(0);"}
                  onclick={async () => await clickFeed(feed)}
                  aria-label="View feed details">{feed.url}</a
                >
              {:else}
                <span class="text-amber-600">
                  {feed.url}
                </span>
              {/if}
            </TableBodyCell>
            <TableBodyCell class={tdClass}
              ><Select
                items={logLevels}
                bind:value={feed.log_level}
                onchange={async () => await updateFeed(feed)}
              /></TableBodyCell
            >
            {#if edit && !feed.enable}
              <TableBodyCell class={tdClass}>N/A</TableBodyCell>
            {:else}
              <TableBodyCell class={tdClass}
                ><Input bind:value={feed.label} oninput={async () => await updateFeed(feed)}
                ></Input></TableBodyCell
              >
            {/if}
            {#if edit}
              <TableBodyCell class={tdClass}
                >{(feed.stats?.downloading ?? 0) + "/" + (feed.stats?.waiting ?? 0)}</TableBodyCell
              >
              <TableBodyCell class={tdClass}
                ><i class={"bx " + (feed.healthy ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
              >
              {#if feed.enable}
                <TableBodyCell onclick={async () => await clickFeed(feed)} class={tdClass}>
                  <a
                    href={"javascript:void(0);"}
                    onclick={async () => await clickFeed(feed)}
                    aria-label="View feed archive"
                  >
                    <i class="bx bx-archive"> </i></a
                  >
                </TableBodyCell>
              {/if}
            {/if}
          {/if}
        </tr>
      {/each}
    {/snippet}
  </CustomTable>
{/if}
<ErrorMessage error={loadConfigError}></ErrorMessage>
