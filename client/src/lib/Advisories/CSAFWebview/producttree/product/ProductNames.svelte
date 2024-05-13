<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

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
