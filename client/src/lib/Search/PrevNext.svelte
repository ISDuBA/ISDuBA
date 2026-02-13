<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { SidebarItem } from "flowbite-svelte";
  import { activeClass, sidebarItemClass, sidebarItemLinkClass } from "$lib/sidenav";
  import { appStore } from "$lib/store.svelte";
  import { request } from "$lib/request";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import { untrack } from "svelte";

  const indentedSidebarItemClass = `${sidebarItemLinkClass} !py-1 ps-10`;

  let loading = $state(false);
  let oldIndex: number | undefined = $state(undefined);

  let searchResults = $derived(appStore.state.app.search.results);
  let index: number = $derived.by(() => {
    const i = appStore.state.app.search.results?.findIndex((r) => {
      const result = $state.snapshot(r);
      const data = result.data;
      const params = $state.snapshot(appStore.state.app.routerParams);
      if (
        params &&
        result.id === Number(params.id) &&
        data[0].publisher === params.publisherNamespace &&
        data[0].tracking_id === params.trackingID
      ) {
        return true;
      }
      return false;
    });
    return i !== undefined ? i : -1;
  });
  let openedDocument = $derived(index !== -1 && searchResults ? searchResults[index] : undefined);

  // Always show 5 slots (including the indicators for leading/following documents)
  let indices = $derived.by((): number[] => {
    const offset = appStore.state.app.search.offset;
    const count = appStore.state.app.search.count;
    if (offset === null || count === null) return [];
    const a: number[] = [];
    if (count === 5 && index === 4) {
      a.push(index - 4);
    }
    if (index === count || offset + index > count - 1) {
      a.push(index - 3);
    }
    if (index < 3 || (count > 2 && index >= count - 1) || offset + index > count - 2) {
      a.push(index - 2);
    }
    if (count > 1 && index > 0) {
      a.push(index - 1);
    }
    a.push(index);
    if (count > 1 && index < count) {
      a.push(index + 1);
    }
    if (count > 2 && index + offset < 2) {
      a.push(index + 2);
    }
    if (count < 6 || index === 0) {
      a.push(index + 3);
    }
    if (count === 5 && index === 0) {
      a.push(index + 4);
    }
    return a;
  });

  const leading = $derived.by(() => {
    if (searchResults === null) return -1;
    const offset = appStore.state.app.search.offset;
    const count = appStore.state.app.search.count;
    if (offset !== null && count !== null && count > 5 && (offset > 0 || index > 1)) {
      return (
        offset +
        searchResults.length -
        (searchResults.length - index) -
        indices.filter((a) => a < index).length
      );
    }
    return 0;
  });

  const following = $derived.by(() => {
    const count = appStore.state.app.search.count;
    if (count === null) return 0;
    return count - leading - indices.length;
  });

  const loadResults = async () => {
    let url = untrack(() => appStore.state.app.search.requestURL);
    const offset = untrack(() => appStore.state.app.search.offset);
    if (!url || offset === null) return;
    loading = true;
    url = url?.concat(`&offset=${offset}&limit=10`);
    const response = await request(url, "GET");
    loading = false;
    if (response.ok) {
      let count, documents;
      if (appStore.state.app.search.type === SEARCHTYPES.EVENT) {
        count = response.content.count;
        documents = response.content.events;
      } else {
        ({ count, documents } = JSON.parse(response.content));
      }
      appStore.setSearchResults($state.snapshot(documents));
      appStore.setSearchResultCount($state.snapshot(count));
    }
  };

  $effect(() => {
    const offset = untrack(() => appStore.state.app.search.offset);
    const results = untrack(() => appStore.state.app.search.results);
    const count = untrack(() => appStore.state.app.search.count);
    if (!results || offset === null || count === null) return;
    if (oldIndex === undefined || oldIndex === -1) {
      oldIndex = index;
      return;
    } else if (
      index !== -1 &&
      oldIndex !== -1 &&
      ((index !== oldIndex && Math.abs(index - oldIndex) < 2) ||
        index === results.length - 1 ||
        (index === 0 && offset > 0))
    ) {
      if (index > results.length - 3 && index + offset < count - 4) {
        appStore.setSearchOffset(Math.min(Math.max(offset + 1, 0), count - 5));
      } else if (index < 2 && offset > 0) {
        appStore.setSearchOffset(Math.min(Math.max(offset - 1, 0), count - 5));
      }
      loadResults();
      oldIndex = index;
    }
  });
</script>

{#snippet leadingFollowingIndicator(text: string)}
  <li>
    <div class="ps-10 text-gray-200">
      <div class="ms-3">{text}</div>
    </div>
  </li>
{/snippet}

{#if searchResults && openedDocument}
  {#if leading > 0}
    {@render leadingFollowingIndicator(`${loading ? "" : leading} ...`)}
  {/if}
  {#each searchResults as result, i}
    {#if indices.includes(i)}
      {@const doc = result}
      {#if loading}
        <div
          class={`${indentedSidebarItemClass} ${i === index ? "bg-primary-200 dark:bg-gray-950" : ""} !text-gray-200 dark:!text-gray-400`}
        >
          <span class="ms-3">{doc.data[0].tracking_id}</span>
        </div>
      {:else}
        <SidebarItem
          class={i === index ? `${activeClass} px-0 py-0` : sidebarItemClass}
          aClass={`${indentedSidebarItemClass} ${i === index ? "!text-primary-900 dark:!text-white" : ""}`}
          label={doc.data[0].tracking_id}
          href={`#/advisories/${doc.data[0].publisher}/${doc.data[0].tracking_id}/documents/${doc.id}`}
        >
          {#snippet icon()}{/snippet}
        </SidebarItem>
      {/if}
    {/if}
  {/each}
  {#if following > 0}
    {@render leadingFollowingIndicator(`... ${loading ? "" : following}`)}
  {/if}
{/if}
