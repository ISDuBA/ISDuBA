<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import ReplaceOperation from "./ReplaceOperation.svelte";

  export let content: any;
  export let isSideBySideViewActivated: boolean = true;
  export let operation: string;
  export let depth = 0;

  $: containerStyle = `padding-left: ${depth > 1 ? 6 * depth : 0}pt`;
</script>

<div style={containerStyle}>
  {#if Array.isArray(content) && !content[0]["m"]}
    {#each content as val, index}
      <div class="mb-4 flex">
        {index + 1}.&ensp;<svelte:self content={val} depth={depth + 1} {operation}></svelte:self>
      </div>
    {/each}
  {:else if Array.isArray(content) && operation === "replace"}
    <ReplaceOperation {content} {isSideBySideViewActivated}></ReplaceOperation>
  {:else if typeof content === "object"}
    {#each Object.keys(content) as key}
      <div>
        {key}:&ensp;
        {#if typeof content[key] === "string"}
          {content[key]}
        {:else}
          <svelte:self content={content[key]} depth={depth + 1} {operation}></svelte:self>
        {/if}
      </div>
    {/each}
  {:else}
    <span>{content}</span>
  {/if}
</div>
