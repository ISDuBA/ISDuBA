<script lang="ts">
  import ReplaceOperation from "./ReplaceOperation.svelte";

  export let content: any;
  export let operation: string;
  export let depth = 0;

  $: containerStyle = `padding-left: ${depth > 1 ? 6 * depth : 0}pt`;
</script>

<div style={containerStyle}>
  {#if Array.isArray(content)}
    {#each content as val, index}
      <div class="mb-4 flex">
        {index + 1}.&ensp;<svelte:self content={val} depth={depth + 1}></svelte:self>
      </div>
    {/each}
  {:else if typeof content === "object"}
    {#each Object.keys(content) as key}
      <div>
        {key}:&ensp;
        {#if typeof content[key] === "string"}
          {content[key]}
        {:else}
          <svelte:self content={content[key]} depth={depth + 1}></svelte:self>
        {/if}
      </div>
    {/each}
  {:else if operation === "replace"}
    <ReplaceOperation {content}></ReplaceOperation>
  {:else}
    <span>{content}</span>
  {/if}
</div>
