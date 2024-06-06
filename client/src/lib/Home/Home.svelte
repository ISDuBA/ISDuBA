<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import Queries from "./Queries.svelte";
  import NewDocuments from "./NewDocuments.svelte";
</script>

<svelte:head>
  <title>Home</title>
</svelte:head>

{#if $appStore.app.isUserLoggedIn}
  <SectionHeader title="Overview"></SectionHeader>
  <hr class="mb-6" />
  <Queries></Queries>
  <div class="flex">
    <NewDocuments></NewDocuments>
  </div>
{/if}
