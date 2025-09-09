<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import Product from "$lib/Advisories/CSAFWebview/producttree/product/Product.svelte";
  import type { Branch } from "$lib/pmdTypes";
  import { Badge } from "flowbite-svelte";
  import Self from "./Branch.svelte";

  interface Props {
    branch: Branch;
    open: boolean;
    openSubBranches: boolean;
  }
  let { branch, open, openSubBranches = false }: Props = $props();
</script>

<div class="pl-3">
  <Collapsible {open} header={branch.category + ": " + branch.name}>
    {#snippet headerSlot()}
      <div class="py-2">
        <Badge rounded large color="dark">{branch.category}</Badge>
        {branch.name}
      </div>
    {/snippet}
    {#if branch.branches}
      {#each branch.branches as b}
        <Self branch={b} open={openSubBranches} {openSubBranches} />
      {/each}
    {/if}
    {#if branch.product}
      <Product product={branch.product} />
    {/if}
  </Collapsible>
</div>
