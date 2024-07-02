<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import JsonDiff from "./JsonDiff.svelte";
  import { appStore } from "$lib/store";
  $: $appStore.app.diff.docA, compare();
  $: $appStore.app.diff.docB, compare();

  let diffDocuments: any;
  let title: string;

  const compare = () => {
    if ($appStore.app.diff.docA && $appStore.app.diff.docB) {
      diffDocuments = {
        docA: $appStore.app.diff.docB,
        docB: $appStore.app.diff.docA
      };
      title = `Changes from ${diffDocuments.docB.tracking_id} (Version ${diffDocuments.docB.version}) to ${diffDocuments.docB.tracking_id} (Version ${diffDocuments.docA.version})`;
    }
  };

  onMount(async () => {
    compare();
  });
</script>

<svelte:head>
  <title>Compare</title>
</svelte:head>

<SectionHeader title="Comparison"></SectionHeader>
<hr class="mb-6" />
{#if diffDocuments}
  <JsonDiff {diffDocuments} {title}></JsonDiff>
{/if}
