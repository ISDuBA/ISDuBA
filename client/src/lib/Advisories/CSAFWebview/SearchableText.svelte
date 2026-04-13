<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
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

  let [before, highlighted, after] = $derived.by(() => {
    if (path != undefined && textPath == path && hit?.text) {
      const firstSplit = hit.text.split("{-");
      const bef = firstSplit[0];
      let [high, aft] = firstSplit[1].split("-}");
      return [bef, high, aft];
    }
    return ["", "", ""];
  });

  $effect(() => {
    if (path != undefined && textPath == path && advisorySearchState.scroll) {
      document.getElementById(elementID)?.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  });
</script>

<span class="flex inline flex-nowrap">
  {#if path != undefined && textPath == path}
    <span>{before}</span>
    <span id={elementID} class="bg-yellow-200 dark:bg-yellow-800">{highlighted}</span>
    <span>{after}</span>
  {:else}
    {text}
  {/if}
</span>
