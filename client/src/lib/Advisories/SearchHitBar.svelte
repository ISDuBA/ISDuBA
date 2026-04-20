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

  let index = $derived(advisorySearchState.hitIndex);
  let hits = $derived(advisorySearchState.searchHits);
</script>

<div class="flex items-center gap-2">
  <Button
    onclick={() => {
      advisorySearchState.scroll = true;
      advisorySearchState.hitIndex--;
    }}
    disabled={index === 0}
    class="h-7 w-7 p-1"
    color="light"
    title="Previous"
  >
    <i class="bx bx-chevron-left"></i>
  </Button>
  <small>{index + 1}/{hits?.length} hits</small>
  <Button
    onclick={() => {
      advisorySearchState.scroll = true;
      advisorySearchState.hitIndex++;
    }}
    disabled={index === hits.length - 1}
    class="h-7 w-7 p-1"
    color="light"
    title="Next"
  >
    <i class="bx bx-chevron-right"></i>
  </Button>
</div>
