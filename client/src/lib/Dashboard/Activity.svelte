<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
</script>

<button
  on:click={() => dispatch("click")}
  class="flex flex-col justify-between rounded-md border border-solid border-gray-300 px-4 py-2 shadow-md hover:bg-gray-100"
>
  <div class="mb-2 flex w-full flex-row justify-between text-left text-xs">
    <div>
      <slot name="top-left"></slot>
    </div>
    <div>
      <slot name="top-right"></slot>
    </div>
  </div>
  <hr class="mb-2" />
  <div class="mb-3 flex-grow text-left">
    <slot></slot>
  </div>
  <div class="flex flex-row justify-between text-left text-xs">
    <div>
      <slot name="bottom-left"></slot>
    </div>
    <div>
      <slot name="bottom-right"></slot>
    </div>
  </div>
  <div class="w-full text-left">
    <slot name="bottom-bottom"></slot>
  </div>
</button>
