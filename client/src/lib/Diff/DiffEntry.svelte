<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import DiffEntry from "./DiffEntry.svelte";
  import ReplaceOperation from "./ReplaceOperation.svelte";

  interface Props {
    content: any;
    isSideBySideViewActivated?: boolean;
    operation: string;
    depth?: number;
  }

  let { content, isSideBySideViewActivated = true, operation, depth = 0 }: Props = $props();

  const uid = $props.id();

  let containerStyle = $derived(`padding-left: ${depth > 1 ? 6 * depth : 0}pt`);
</script>

<div style={containerStyle}>
  {#if Array.isArray(content) && !content[0]["m"]}
    {#each content as val, index (`diffentry-1-${uid}-${index}`)}
      <div class="mb-4 flex">
        {index + 1}.&ensp;<DiffEntry content={val} depth={depth + 1} {operation}></DiffEntry>
      </div>
    {/each}
  {:else if Array.isArray(content) && operation === "replace"}
    <ReplaceOperation {content} {isSideBySideViewActivated}></ReplaceOperation>
  {:else if typeof content === "object"}
    {#each Object.keys(content) as key, index (`diffentry-2-${uid}-${index}`)}
      <div>
        {key}:&ensp;
        {#if typeof content[key] === "string"}
          {content[key]}
        {:else}
          <DiffEntry content={content[key]} depth={depth + 1} {operation}></DiffEntry>
        {/if}
      </div>
    {/each}
  {:else}
    <span>{content}</span>
  {/if}
</div>
