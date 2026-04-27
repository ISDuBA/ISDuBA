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
    if (advisorySearchState.matchIndex !== -1) {
      return advisorySearchState.searchMatches[advisorySearchState.matchIndex];
    }
    return undefined;
  });

  let path: string | undefined = $derived.by(() => {
    return hit?.path;
  });

  let splitted: string[] | undefined = $derived.by(() => {
    if (path != undefined && textPath == path && hit && text != undefined) {
      let t = text;
      const splits: string[] = [];
      for (let i = 0; i < hit.positions.length; i++) {
        const pos = hit.positions[i];
        const term = text.substring(pos[0], pos[0] + pos[1]);
        const splittedText = t.split(term);
        splits.push(splittedText[0], term);
        t = splittedText[1];
      }
      if (t) splits.push(t);
      return splits;
    }
    return [];
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
      {#if (index + 1) % 2 === 0}
        <span id={elementID} class="bg-yellow-200 dark:bg-yellow-800">{s}</span>
      {:else}
        <span>{s}</span>
      {/if}
    {/each}
  {:else}
    {text}
  {/if}
</span>
