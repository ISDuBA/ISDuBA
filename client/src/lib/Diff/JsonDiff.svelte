<script lang="ts">
  import { Label } from "flowbite-svelte";
  import DiffEntry from "./DiffEntry.svelte";
  import type { JsonDiffResult } from "./JsonDiff";

  export let diff: any;
  $: groupedResults = [
    {
      op: "add",
      changes: diff.result.filter((result: JsonDiffResult) => result.op === "add")
    },
    {
      op: "replace",
      changes: diff.result.filter((result: JsonDiffResult) => result.op === "replace")
    },
    {
      op: "remove",
      changes: diff.result.filter((result: JsonDiffResult) => result.op === "remove")
    }
  ];

  const getBodyClass = (operation: string) => {
    let bodyClass = "mb-4 p-2";
    if (operation === "add") {
      return `${bodyClass} bg-green-100`;
    } else if (operation === "remove") {
      return `${bodyClass} bg-red-100`;
    } else {
      return `${bodyClass} bg-gray-100`;
    }
  };
</script>

<div>
  <Label class="text-lg">Compare Version {diff.docA} with Version {diff.docB}</Label>
  <span>{diff.result.length} changes</span>
  {#each groupedResults as result}
    {#if result.changes.length > 0}
      <div>
        {#if result.op === "add"}
          <div class="bg-green-800 py-1 ps-2 text-white">
            <i class="bx bx-plus text-lg"></i>
            <span>Added ({result.changes.length})</span>
          </div>
        {:else if result.op === "remove"}
          <div class="bg-red-800 py-1 ps-2 text-white">
            <i class="bx bx-x text-lg"></i>
            <span>Removed ({result.changes.length})</span>
          </div>
        {:else}
          <div class="bg-gray-800 py-1 ps-2 text-white">
            <i class="bx bx-pencil text-lg"></i>
            <span>Edited ({result.changes.length})</span>
          </div>
        {/if}
      </div>
      {#each result.changes as change}
        <div class={getBodyClass(change.op)}>
          <div class="mb-1">
            <b>
              <code>
                {change.path}
              </code>
            </b>
          </div>
          {#if change.value}
            <DiffEntry content={change.value} operation={change.op}></DiffEntry>
          {/if}
        </div>
      {/each}
    {/if}
  {/each}
</div>
