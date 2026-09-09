<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { CSAFDocumentv2_0, ProductTree } from "$lib/Advisories/types/csaf-2.0";
  import { appStore } from "$lib/store.svelte";
  import { untrack } from "svelte";
  import Collapsible from "../../Collapsible.svelte";
  import { productTreeCutoffs } from "../../efficiencyCutoffs";
  import Relation from "./Relation.svelte";

  interface Props {
    basePath: string;
  }
  let { basePath = "" }: Props = $props();

  const uid = $props.id();

  let openRelationships = $state(false);

  let selectedProduct = $derived(appStore.state.webview.ui.selectedProduct);
  let doc: CSAFDocumentv2_0 | null = $derived(
    appStore.state.webview.doc as CSAFDocumentv2_0 | null
  );
  let productTree: ProductTree | undefined = $derived(doc?.product_tree);
  let relationships = $derived(productTree?.relationships);

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
    const len = productTree?.relationships?.length;
    openRelationships = (len && len !== 0) || 0 <= productTreeCutoffs.relations;
  });
</script>

<Collapsible
  header="Relationships"
  open={!!selectedProduct || openRelationships}
  path="/product_tree"
>
  {#each relationships as relation, i (`relationships-${uid}-${i}`)}
    <Relation {basePath} {relation} path={`/product_tree/relationships[${i}]`} />
  {/each}
</Collapsible>
