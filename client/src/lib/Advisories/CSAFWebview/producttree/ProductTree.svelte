<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import Branch from "./branch/Branch.svelte";
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import ProductGroups from "./productgroup/ProductGroups.svelte";
  import ProductNames from "./product/ProductNames.svelte";
  import Relationships from "./relationship/Relationships.svelte";
  import { productTreeCutoffs } from "../efficiencyCutoffs";
  import { untrack } from "svelte";
  import type { ProductTree } from "$lib/Advisories/types/csaf-2.0";

  interface Props {
    basePath: string;
  }
  let { basePath }: Props = $props();

  const uid = $props.id();

  let openSubBranches = $state(false);
  let openBranches = $state(false);
  let openRelationships = $state(false);
  let selectedProduct = $derived(appStore.state.webview.ui.selectedProduct);
  let doc = $derived(appStore.state.webview.doc);
  let productTree = $derived(doc?.product_tree);
  let csafVersion = $derived(doc?.document.csaf_version);

  $effect(() => {
    untrack(() => openSubBranches);
    untrack(() => openBranches);
    untrack(() => openRelationships);
    let size = 0;
    for (let branch of productTree?.branches ?? []) {
      if (branch.branches) {
        size = size + branch.branches.length;
      }
      if (size >= productTreeCutoffs.level2Upper) {
        break;
      }
    }
    openBranches = !!selectedProduct || size <= productTreeCutoffs.level2Upper;
    openSubBranches = !!selectedProduct || size <= productTreeCutoffs.level2Lower;
  });
</script>

{#if doc && productTree}
  {#if productTree.branches}
    <Collapsible
      header="Branches"
      open={!!selectedProduct || productTree.branches.length <= productTreeCutoffs.level1}
      path="/product_tree"
    >
      {#each productTree.branches as branch, i (`producttree-${uid}-${i}`)}
        {#if csafVersion === "2.0"}
          <Branch
            {branch}
            {openSubBranches}
            open={openBranches}
            path={`/product_tree/branches[${i}]`}
          />
        {/if}
      {/each}
    </Collapsible>
  {/if}

  {#if csafVersion === "2.0" && (productTree as ProductTree | undefined)?.relationships}
    <Relationships {basePath} />
  {/if}

  {#if productTree.product_groups}
    <Collapsible header="Product groups" open path="/product_tree">
      <ProductGroups productGroups={productTree.product_groups} />
    </Collapsible>
  {/if}

  {#if productTree.full_product_names}
    <Collapsible header="Full Product Names" open path="/product_tree">
      <ProductNames productNames={productTree.full_product_names} />
    </Collapsible>
  {/if}
{/if}
