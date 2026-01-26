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

  const indentedSidebarItemClass = `${sidebarItemLinkClass} !py-1 ps-10`;

  let searchResults = $derived(appStore.state.app.search.results);
  let index = $derived(
    appStore.state.app.search.results?.findIndex((r) => {
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
    })
  );
  const openedDocument = $derived(
    index !== undefined && index !== -1 && searchResults ? searchResults[index] : undefined
  );
  let indices = $derived(
    index !== undefined && index !== -1
      ? [index - 2, index - 1, index, index + 1, index + 2]
      : undefined
  );
</script>

{#if searchResults && openedDocument}
  {#each searchResults as result, i}
    {#if indices?.includes(i)}
      {@const doc = result}
      <SidebarItem
        class={i === index ? `${activeClass} px-0 py-0` : sidebarItemClass}
        aClass={indentedSidebarItemClass}
        label={doc.data[0].tracking_id}
        href={`#/advisories/${doc.data[0].publisher}/${doc.data[0].tracking_id}/documents/${doc.id}`}
      >
        {#snippet icon()}{/snippet}
      </SidebarItem>
    {/if}
  {/each}
{/if}
