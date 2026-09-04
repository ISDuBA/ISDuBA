<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { tick, untrack } from "svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import ProductIdentificationHelper from "../product/ProductIdentificationHelper.svelte";
  import type { ProductPath } from "$lib/Advisories/types/csaf-2.1";
  import { A } from "flowbite-svelte";
  import SearchableText from "../../SearchableText.svelte";
  import Link from "$lib/Components/Link.svelte";

  interface Props {
    basePath: string;
    path: string;
    productPath: ProductPath;
  }
  let { basePath, path, productPath }: Props = $props();

  let blink = $state(false);
  async function updateUI() {
    await tick();
    document
      .getElementById(`${productPath.full_product_name.product_id}`)
      ?.scrollIntoView({ behavior: "smooth" });
    blink = true;
    await new Promise((res) => setTimeout(res, 5000));
    blink = false;
  }
  let selectedProduct = $derived(appStore.state.webview.ui.selectedProduct);
  let productID = $derived(productPath.full_product_name.product_id);
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
  header={`${productPath.full_product_name.product_id}`}
  level={4}
  open={productPath.full_product_name.product_id === appStore.state.webview.ui.selectedProduct}
  {highlight}
  onClose={() => {
    if (appStore.state.webview.ui.selectedProduct === productPath.full_product_name.product_id) {
      appStore.resetSelectedProduct();
    }
  }}
  {path}
>
  <div id={productPath.full_product_name.product_id} class={blink ? "blink" : ""}>
    <KeyValue
      keys={["Name", "Product ID"]}
      values={[productPath.full_product_name.name, productPath.full_product_name.product_id]}
      paths={[`${path}/full_product_name/name`, `${path}/full_product_name/product_id`]}
    />
    {#if productPath.full_product_name.product_identification_helper}
      <ProductIdentificationHelper
        helper={productPath.full_product_name.product_identification_helper}
        path={`${path}/product_identification_helper`}
      />
    {/if}
    <table class="border-separate border-spacing-3">
      <tbody>
        <tr>
          <td>Beginning product reference</td>
          <td
            ><A
              class="text-primary-700 dark:text-primary-400"
              id={crypto.randomUUID()}
              href={`${basePath}product-${encodeURIComponent(productPath.beginning_product_reference)}`}
            >
              <SearchableText
                textPath={`${path}/beginning_product_reference`}
                text={productPath.beginning_product_reference}
              />
            </A>
          </td>
        </tr>
        <Collapsible header="Subpaths" open path={`${path}/subpaths`}>
          {#each productPath.subpaths as subpath, index (`productpath-${index}`)}
            <tr>
              <td>Category</td>
              <td>
                <SearchableText
                  textPath={`${path}/subpaths/${index}/category`}
                  text={subpath.category}
                />
              </td>
            </tr>
            <tr>
              <td>Next product reference</td>
              <td>
                <Link href={subpath.next_product_reference} class="underline">
                  <SearchableText
                    textPath={`${path}/subpaths/${index}/next_product_reference`}
                    text={subpath.next_product_reference}
                  />
                </Link>
              </td>
            </tr>
          {/each}
        </Collapsible>
      </tbody>
    </table>
  </div>
</Collapsible>
