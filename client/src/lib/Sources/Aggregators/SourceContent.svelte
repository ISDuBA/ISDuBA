<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { AggregatorEntry, FeedInfo, SourceInfo } from "./aggregator";
  import { Button } from "flowbite-svelte";
  import { tdClass } from "$lib/Table/defaults";
  import FeedBulletPoint from "./FeedBulletPoint.svelte";
  import { appStore } from "$lib/store";

  export let source: SourceInfo;
  export let entry: AggregatorEntry;

  const sortFeeds = (a: FeedInfo, b: FeedInfo) => {
    if (a.highlight && !b.highlight) {
      return 1;
    } else if (!a.highlight && b.highlight) {
      return -1;
    }
    return 0;
  };
</script>

<div class="mb-2 flex items-center gap-2 text-sm text-black dark:text-white">
  {#if source.id}
    <i class="bx bx-git-repo-forked text-lg"></i>
  {/if}
  {source.name}
  {#if entry.feedsSubscribed === 0 && appStore.isSourceManager()}
    <Button href={`/#/sources/new/${encodeURIComponent(entry.url)}`} color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>As new source</span>
    </Button>
  {/if}
</div>
<div class="flex flex-col items-start">
  {#each source.feeds.toSorted(sortFeeds) as feed}
    {@const feedClass = `text-sm ${tdClass} ${feed.highlight ? "text-amber-600" : "text-black dark:text-white"}`}
    <div class="mb-2 ms-8">
      <div>
        <FeedBulletPoint filled={!feed.highlight}></FeedBulletPoint>
        <span class={feedClass}>{feed.url}</span>
      </div>
    </div>
  {/each}
</div>
