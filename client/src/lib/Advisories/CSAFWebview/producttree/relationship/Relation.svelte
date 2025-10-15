<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { tick, untrack } from "svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import ProductIdentificationHelper from "../product/ProductIdentificationHelper.svelte";
  import type { Relationship } from "$lib/pmdTypes";
  import { A } from "flowbite-svelte";

  interface Props {
    basePath: string;
    relation: Relationship;
  }
  let { basePath, relation }: Props = $props();

  let blink = $state(false);
  async function updateUI() {
    await tick();
    document
      .getElementById(`${relation.full_product_name.product_id}`)
      ?.scrollIntoView({ behavior: "smooth" });
    blink = true;
    await new Promise((res) => setTimeout(res, 5000));
    blink = false;
  }
  let selectedProduct = $derived(appStore.state.webview.ui.selectedProduct);
  let productID = $derived(relation.full_product_name.product_id);
  let highlight = $derived(selectedProduct === productID);
  $effect(() => {
    untrack(() => selectedProduct);
    untrack(() => blink);
    if (selectedProduct === productID) {
      updateUI();
    }
  });
</script>

<Collapsible
  header={`${relation.full_product_name.product_id}`}
  level={4}
  open={relation.full_product_name.product_id === appStore.state.webview.ui.selectedProduct}
  {highlight}
  onClose={() => {
    if (appStore.state.webview.ui.selectedProduct === relation.full_product_name.product_id) {
      appStore.resetSelectedProduct();
    }
  }}
>
  <div id={relation.full_product_name.product_id} class={blink ? "blink" : ""}>
    <KeyValue
      keys={["Category", "Name", "Product ID"]}
      values={[
        relation.category,
        relation.full_product_name.name,
        relation.full_product_name.product_id
      ]}
    />
    {#if relation.full_product_name.product_identification_helper}
      <ProductIdentificationHelper
        helper={relation.full_product_name.product_identification_helper}
      />
    {/if}
    <table>
      <tbody>
        <tr>
          <td>Product reference</td>
          <td
            ><A
              class="text-primary-700 dark:text-primary-400"
              id={crypto.randomUUID()}
              href={`${basePath}product-${encodeURIComponent(relation.product_reference)}`}
              >{relation.product_reference}</A
            ></td
          >
        </tr>
        <tr>
          <td>Relates to</td>
          <td
            ><A
              class="text-primary-700 dark:text-primary-400"
              id={crypto.randomUUID()}
              href={`${basePath}product-${encodeURIComponent(relation.relates_to_product_reference)}`}
              >{relation.relates_to_product_reference}</A
            ></td
          >
        </tr>
      </tbody>
    </table>
  </div>
</Collapsible>
