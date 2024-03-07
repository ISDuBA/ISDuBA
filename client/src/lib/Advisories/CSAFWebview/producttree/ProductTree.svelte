<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import Branch from "./branch/Branch.svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import ProductGroups from "./productgroup/ProductGroups.svelte";
  import ProductNames from "./product/ProductNames.svelte";
  import Relationships from "./relationship/Relationships.svelte";
</script>

{#if $appStore.webview.doc?.productTree.branches}
  {#each $appStore.webview.doc?.productTree.branches as branch}
    <Collapsible
      header="Branches"
      open={$appStore.webview.ui.isProductTreeVisible ||
        ($appStore.webview.ui.isProductTreeOpen &&
          !(
            $appStore.webview.doc?.productTree.relationships &&
            $appStore.webview.doc?.productTree.product_groups &&
            $appStore.webview.doc?.productTree.full_product_names
          ))}
    >
      <Branch {branch} />
    </Collapsible>
  {/each}
{/if}

{#if $appStore.webview.doc?.productTree.relationships}
  <Collapsible
    header="Relationships"
    open={$appStore.webview.ui.isProductTreeVisible ||
      ($appStore.webview.ui.isProductTreeOpen &&
        !(
          $appStore.webview.doc?.productTree.branches &&
          $appStore.webview.doc?.productTree.product_groups &&
          $appStore.webview.doc?.productTree.full_product_names
        ))}
  >
    <Relationships relationships={$appStore.webview.doc?.productTree.relationships} />
  </Collapsible>
{/if}

{#if $appStore.webview.doc?.productTree.product_groups}
  <Collapsible
    header="Product groups"
    open={$appStore.webview.ui.isProductTreeVisible ||
      ($appStore.webview.ui.isProductTreeOpen &&
        !(
          $appStore.webview.doc?.productTree.branches &&
          $appStore.webview.doc?.productTree.relationships &&
          $appStore.webview.doc?.productTree.full_product_names
        ))}
  >
    <ProductGroups productGroups={$appStore.webview.doc?.productTree.product_groups} />
  </Collapsible>
{/if}

{#if $appStore.webview.doc?.productTree.full_product_names}
  <Collapsible
    header="Full Product Names"
    open={$appStore.webview.ui.isProductTreeVisible ||
      ($appStore.webview.ui.isProductTreeOpen &&
        !(
          $appStore.webview.doc?.productTree.branches &&
          $appStore.webview.doc?.productTree.relationships &&
          $appStore.webview.doc?.productTree.product_groups
        ))}
  >
    <ProductNames productNames={$appStore.webview.doc?.productTree.full_product_names} />
  </Collapsible>
{/if}
