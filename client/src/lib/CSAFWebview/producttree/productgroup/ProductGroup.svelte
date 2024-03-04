<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { ProductGroup } from "$lib/types";
  import ValueField from "$lib/CSAFWebview/ValueField.svelte";
  import ValueList from "$lib/CSAFWebview/ValueList.svelte";
  export let productGroup: ProductGroup;
</script>

<ValueList label="Group IDs" values={productGroup.group_ids} />
<ValueList label="Product IDs" values={productGroup.product_ids} />

{#if productGroup.summary}
  <ValueField label="Summary" value={productGroup.summary} />
{/if}
