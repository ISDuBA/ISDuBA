<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import ProductIdentificationHelper from "./ProductIdentificationHelper.svelte";
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import { appStore } from "$lib/store.svelte";
  import { tick, untrack } from "svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";

  interface Props {
    product: any;
  }
  let { product }: Props = $props();
  async function updateUI() {
    await tick();
    document.getElementById(`${product?.product_id}`)?.scrollIntoView({ behavior: "smooth" });
  }

  let selectedProduct = $derived(appStore.state.webview.ui.selectedProduct);
  let productID = $derived(product?.product_id);
  let highlight = $derived(selectedProduct === productID);

  $effect(() => {
    untrack(() => selectedProduct);
    if (selectedProduct === productID) {
      updateUI();
    }
  });
</script>

<div id={product?.product_id}>
  <Collapsible
    header={product.name}
    level={4}
    {highlight}
    open={appStore.state.webview.ui.selectedProduct === product.product_id}
    onClose={() => {
      if (appStore.state.webview.ui.selectedProduct === product.product_id) {
        appStore.resetSelectedProduct();
      }
    }}
  >
    <KeyValue keys={["Name", "Product ID"]} values={[product.name, product.product_id]} />
    {#if product.product_identification_helper}
      <ProductIdentificationHelper helper={product.product_identification_helper} />
    {/if}
  </Collapsible>
</div>
