<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { searchColumnName } from "$lib/Table/defaults";
  import { Button } from "flowbite-svelte";

  let routerParams = $derived(appStore.state.app.routerParams);
  let results = $derived(appStore.state.app.search.results);
  let index = $derived(appStore.state.app.search.indexInsideDoc);
  let doc: any = $derived.by(() => {
    if (!routerParams) return;
    return results?.find((result) => {
      const r = $state.snapshot(result);
      return (
        r.id === Number(routerParams.id) &&
        r.data[0].publisher === routerParams.publisherNamespace &&
        r.data[0].tracking_id === routerParams.trackingID
      );
    });
  });
  let hits = $derived(doc?.data.map((d: any) => d[searchColumnName]));
  let hit = $derived(index !== -1 ? hits?.[index] : undefined);
</script>

<div class="flex gap-2 items-center">
  <Button class="h-7 w-7 p-1" color="light" title="Previous">
    <i class="bx bx-chevron-left"></i>
  </Button>
  <small>{index + 1}/{hits?.length} hits</small>
  <Button class="h-7 w-7 p-1" color="light" title="Next">
    <i class="bx bx-chevron-right"></i>
  </Button>
</div>
