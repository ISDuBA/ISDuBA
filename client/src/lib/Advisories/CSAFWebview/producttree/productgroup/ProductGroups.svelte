<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import ProductGroup from "./ProductGroup.svelte";
  export let productGroups: any;
</script>

{#if productGroups}
  {#each productGroups as productGroup}
    <ProductGroup {productGroup} />
  {/each}
{/if}
