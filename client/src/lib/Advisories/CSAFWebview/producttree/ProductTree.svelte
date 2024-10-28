<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

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
  import { productTreeCutoffs } from "../efficiencyCutoffs";
  export let basePath = "";

  let openSubBranches = false;
  let openBranches = false;
  let openRelationships = false;
  $: selectedProduct = $appStore.webview.ui.selectedProduct;

  $: {
    let size = 0;
    for (let branch of $appStore.webview.doc?.productTree.branches ?? []) {
      if (branch.branches) {
        size = size + branch.branches.length;
      }
      if (size >= productTreeCutoffs.level2Upper) {
        break;
      }
    }
    openBranches = !!selectedProduct || size <= productTreeCutoffs.level2Upper;
    openSubBranches = !!selectedProduct || size <= productTreeCutoffs.level2Lower;
    openRelationships =
      $appStore.webview.doc?.productTree.relationships?.length ?? 0 <= productTreeCutoffs.relations;
  }
</script>

{#if $appStore.webview.doc?.productTree.branches}
  <Collapsible
    header="Branches"
    open={!!selectedProduct ||
      $appStore.webview.doc?.productTree.branches.length <= productTreeCutoffs.level1}
  >
    {#each $appStore.webview.doc?.productTree.branches as branch}
      <Branch {branch} {openSubBranches} open={openBranches} />
    {/each}
  </Collapsible>
{/if}

{#if $appStore.webview.doc?.productTree.relationships}
  <Collapsible header="Relationships" open={!selectedProduct && openRelationships}>
    <Relationships {basePath} relationships={$appStore.webview.doc?.productTree.relationships} />
  </Collapsible>
{/if}

{#if $appStore.webview.doc?.productTree.product_groups}
  <Collapsible header="Product groups" open>
    <ProductGroups
      productGroups={!selectedProduct && $appStore.webview.doc?.productTree.product_groups}
    />
  </Collapsible>
{/if}

{#if $appStore.webview.doc?.productTree.full_product_names}
  <Collapsible header="Full Product Names" open>
    <ProductNames
      productNames={!selectedProduct && $appStore.webview.doc?.productTree.full_product_names}
    />
  </Collapsible>
{/if}
