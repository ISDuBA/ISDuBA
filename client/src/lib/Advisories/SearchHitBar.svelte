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

  let index = $derived(advisorySearchState.hitIndex);
  let hits = $derived(advisorySearchState.searchHits);

  const prev = () => {
    advisorySearchState.scroll = true;
    if (index <= 0) {
      advisorySearchState.hitIndex = advisorySearchState.searchHits.length - 1;
    } else {
      advisorySearchState.hitIndex--;
    }
  };

  const next = () => {
    advisorySearchState.scroll = true;
    if (index >= advisorySearchState.searchHits.length - 1) {
      advisorySearchState.hitIndex = 0;
    } else {
      advisorySearchState.hitIndex++;
    }
  };
</script>

<div class="flex items-center gap-1">
  <Button onclick={prev} class="h-7 w-7 p-1" color="light" title="Previous hit">
    <i class="bx bx-chevron-up"></i>
  </Button>
  <Button onclick={next} class="h-7 w-7 p-1" color="light" title="Next hit">
    <i class="bx bx-chevron-down"></i>
  </Button>
  <small>{index + 1}/{hits?.length} hits for "{appStore.state.app.search.term}"</small>
</div>
