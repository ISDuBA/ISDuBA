<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { CSAFDocumentv2_1, ProductTree } from "$lib/Advisories/types/csaf-2.1";
  import { appStore } from "$lib/store.svelte";
  import { untrack } from "svelte";
  import Collapsible from "../../Collapsible.svelte";
  import { productTreeCutoffs } from "../../efficiencyCutoffs";
  import ProductPath from "./ProductPath.svelte";

  interface Props {
    basePath: string;
  }
  let { basePath = "" }: Props = $props();

  const uid = $props.id();

  let openRelationships = $state(false);

  let selectedProduct = $derived(appStore.state.webview.ui.selectedProduct);
  let doc: CSAFDocumentv2_1 | null = $derived(
    appStore.state.webview.doc as CSAFDocumentv2_1 | null
  );
  let productTree: ProductTree | undefined = $derived(doc?.product_tree);
  let productPaths = $derived(productTree?.product_paths);

  $effect(() => {
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
    const len = productTree?.product_paths?.length;
    openRelationships = (len && len !== 0) || 0 <= productTreeCutoffs.relations;
  });
</script>

<Collapsible
  header="Product Paths"
  open={!!selectedProduct || openRelationships}
  path="/product_tree"
>
  {#each productPaths as productPath, i (`productpaths-${uid}-${i}`)}
    <ProductPath {basePath} {productPath} path={`/product_tree/productpaths[${i}]`} />
  {/each}
</Collapsible>
