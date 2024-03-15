<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Reference from "./Reference.svelte";
  export let references: any;
</script>

{#if references}
  {#each references as reference, index}
    <Reference {reference} />
    {#if references.length > 1 && index < references.length - 1}
      <hr />
    {/if}
  {/each}
{/if}
