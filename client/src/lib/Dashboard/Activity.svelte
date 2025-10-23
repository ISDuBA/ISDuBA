<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { Snippet } from "svelte";
  interface Props {
    children: Snippet;
    topLeftSlot?: Snippet;
    topRightSlot?: Snippet;
    bottomLeftSlot?: Snippet;
    bottomRightSlot?: Snippet;
    bottomBottomSlot?: Snippet;
    onClicked: () => void;
  }

  let {
    onClicked,
    children,
    topLeftSlot = undefined,
    topRightSlot = undefined,
    bottomLeftSlot = undefined,
    bottomRightSlot = undefined,
    bottomBottomSlot = undefined
  }: Props = $props();
</script>

<button
  onclick={() => onClicked()}
  class="flex flex-col justify-between rounded-md border border-solid border-gray-300 px-4 py-2 break-all shadow-md hover:bg-gray-100 dark:hover:bg-gray-800"
>
  <div class="mb-2 flex w-full flex-row justify-between gap-2 text-left text-xs">
    <div>
      {#if topLeftSlot}
        {@render topLeftSlot()}
      {/if}
    </div>
    <div>
      {#if topRightSlot}
        {@render topRightSlot()}
      {/if}
    </div>
  </div>
  <hr class="mb-2" />
  <div class="mb-3 flex-grow text-left">
    {@render children()}
  </div>
  <div class="flex w-full flex-row justify-between gap-2 text-left text-xs">
    <div>
      {#if bottomLeftSlot}
        {@render bottomLeftSlot()}
      {/if}
    </div>
    <div>
      {#if bottomRightSlot}
        {@render bottomRightSlot()}
      {/if}
    </div>
  </div>
  <div class="w-full text-left">
    {#if bottomBottomSlot}
      {@render bottomBottomSlot()}
    {/if}
  </div>
</button>
