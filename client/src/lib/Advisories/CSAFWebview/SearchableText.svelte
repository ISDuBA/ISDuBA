<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { splitMatches } from "$lib/utils";
  import { advisorySearchState, type SearchMatch } from "../advisory.svelte";

  interface Props {
    text: string | undefined;
    textPath: string;
  }
  let { textPath = "", text }: Props = $props();
  let uid = $props.id();
  let elementID = $derived(`SearchableText-${uid}`);

  let match: SearchMatch | undefined = $derived.by(() => {
    if (advisorySearchState.matchIndex !== -1) {
      return advisorySearchState.searchMatches[advisorySearchState.matchIndex];
    }
    return undefined;
  });

  let path: string | undefined = $derived.by(() => {
    return match?.path;
  });

  let splitted: string[] | undefined = $derived.by(() => {
    if (path != undefined && textPath == path && match && text != undefined) {
      return splitMatches(text, match?.positions);
    }
    return [];
  });

  $effect(() => {
    // Scroll to a match when the advisory view was just opened and when the user switches to the previous/next match.
    if (path != undefined && textPath == path && advisorySearchState.scroll) {
      document
        .getElementById(elementID)
        ?.scrollIntoView({ behavior: "smooth", block: "center", inline: "center" });
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
