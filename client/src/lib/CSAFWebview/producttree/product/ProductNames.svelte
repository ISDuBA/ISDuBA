<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import ProductName from "./ProductName.svelte";

  export let productNames: any;
  productNames = productNames.sort((p1: any, p2: any) => {
    if (p1.name < p2.name) return -1;
    if (p1.name > p2.name) return 1;
    return 0;
  });
</script>

{#each productNames as product}
  <ProductName {product} />
{/each}
