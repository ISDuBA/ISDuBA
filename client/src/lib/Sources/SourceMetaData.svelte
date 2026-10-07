<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { DescriptionList, List } from "flowbite-svelte";
  import type { Snippet } from "svelte";
  import { ddClass, dtClass } from "./source";

  interface Props {
    extension?: Snippet;
    pmd: any;
    source: any;
  }
  let { extension, pmd = null, source }: Props = $props();
</script>

<List tag="dl" class="mb-4 divide-y divide-gray-200 text-sm 2xl:w-max">
  <div>
    <DescriptionList tag="dt" class={dtClass}>Domain/PMD</DescriptionList>
    <DescriptionList tag="dd" class={ddClass}>{source.url}</DescriptionList>
  </div>
  {#if pmd}
    <div>
      <DescriptionList tag="dt" class={dtClass}>Canonical URL</DescriptionList>
      <DescriptionList tag="dd" class={ddClass}>{pmd.canonical_url}</DescriptionList>
    </div>
    <div>
      <DescriptionList tag="dt" class={dtClass}>Publisher Name</DescriptionList>
      <DescriptionList tag="dd" class={ddClass}>{pmd.publisher.name}</DescriptionList>
    </div>
    <div>
      <DescriptionList tag="dt" class={dtClass}>Publisher Contact</DescriptionList>
      <DescriptionList tag="dd" class={ddClass}>{pmd.publisher.contact_details}</DescriptionList>
    </div>
    <div>
      {#if pmd.publisher.issuing_authority}
        <DescriptionList tag="dt" class={dtClass}>Issuing Authority</DescriptionList>
        <DescriptionList tag="dd" class={ddClass}>{pmd.publisher.issuing_authority}</DescriptionList
        >
      {/if}
    </div>
    {#if extension}
      {@render extension()}
    {/if}
  {/if}
</List>
