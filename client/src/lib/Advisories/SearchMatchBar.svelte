<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { Button } from "flowbite-svelte";
  import { advisorySearchState } from "./advisory.svelte";
  import { appStore } from "$lib/store.svelte";

  let index = $derived(advisorySearchState.matchIndex);
  let matches = $derived(advisorySearchState.searchMatches);

  const prev = () => {
    advisorySearchState.scroll = true;
    if (index <= 0) {
      advisorySearchState.matchIndex = advisorySearchState.searchMatches.length - 1;
    } else {
      advisorySearchState.matchIndex--;
    }
  };

  const next = () => {
    advisorySearchState.scroll = true;
    if (index >= advisorySearchState.searchMatches.length - 1) {
      advisorySearchState.matchIndex = 0;
    } else {
      advisorySearchState.matchIndex++;
    }
  };
</script>

<div class="sticky z-100 flex items-center gap-1 md:static xl:w-fit">
  <Button
    onclick={prev}
    disabled={matches.length < 2}
    class="h-7 w-7 p-1"
    color="light"
    title="Previous match"
  >
    <i class="bx bx-chevron-up"></i>
  </Button>
  <Button
    onclick={next}
    disabled={matches.length < 2}
    class="h-7 w-7 p-1"
    color="light"
    title="Next match"
  >
    <i class="bx bx-chevron-down"></i>
  </Button>
  {#if matches.length === 0}
    <small>No matches for "{appStore.state.app.search.term}"</small>
  {:else}
    <small>{index + 1}/{matches?.length} matches for "{appStore.state.app.search.term}"</small>
  {/if}
</div>
