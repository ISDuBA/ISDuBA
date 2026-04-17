<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { advisorySearchState, type SearchHit } from "../advisory.svelte";

  interface Props {
    text: string | undefined;
    textPath: string;
  }
  let { textPath = "", text }: Props = $props();
  let uid = $props.id();
  let elementID = $derived(`SearchableText-${uid}`);

  let hit: SearchHit | undefined = $derived.by(() => {
    if (advisorySearchState.hitIndex !== -1) {
      return advisorySearchState.searchHits[advisorySearchState.hitIndex];
    }
    return undefined;
  });

  let path: string | undefined = $derived(hit?.path);
  let term = $derived(appStore.state.app.search.term);

  let splitted: string[] | undefined = $derived.by(() => {
    if (path != undefined && textPath == path && hit?.text && term) {
      return hit.text.split(new RegExp(RegExp.escape(term), "ig"));
    }
  });

  $effect(() => {
    if (path != undefined && textPath == path && advisorySearchState.scroll) {
      document.getElementById(elementID)?.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  });
</script>

<span class="flex inline flex-nowrap">
  {#if path != undefined && textPath == path && splitted}
    {#each splitted as s, index (`SearchableText-${uid}-${index}`)}
      <span>{s}</span>
      {#if index < splitted.length - 1}
        <span id={elementID} class="bg-yellow-200 dark:bg-yellow-800">{term}</span>
      {/if}
    {/each}
  {:else}
    {text}
  {/if}
</span>
