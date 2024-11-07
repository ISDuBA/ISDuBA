<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  export let content: any[];
  export let isSideBySideViewActivated: boolean = true;
  $: sideBySideContent = isSideBySideViewActivated
    ? [
        content?.filter((element) => ["d", "o"].includes(element.m)),
        content?.filter((element) => ["i", "o"].includes(element.m))
      ]
    : [];

  const getSpanClass = (type: string) => {
    if (type === "i") return "bg-green-200";
    if (type === "d") return "bg-red-200";
  };
</script>

<div>
  {#if sideBySideContent.length > 0}
    <div class="flex justify-between gap-2">
      <div class="flex w-6/12 items-center gap-1">
        <i class="bx bx-minus"></i>
        <div class="h-fit w-fit bg-red-200 dark:bg-[#522630]">
          {#each sideBySideContent[0] as part}
            <span>{part.t}</span>
          {/each}
        </div>
      </div>
      <div class="flex w-6/12 items-center gap-1">
        <i class="bx bx-plus"></i>
        <div class="h-fit w-fit bg-green-200 dark:bg-[#173d3e]">
          {#each sideBySideContent[1] as part}
            <span>{part.t}</span>
          {/each}
        </div>
      </div>
    </div>
  {:else if content}
    {#each content as part}
      <span class={getSpanClass(part.m)}>{part.t}</span>
    {/each}
  {/if}
</div>
