<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { Check, Copy } from "@boxicons/svelte";
  import type { CopyState } from "./types";

  interface Props {
    copyState?: CopyState;
    errorMessage: string;
    title: string;
    value: string;
  }
  let { copyState = $bindable(undefined), errorMessage, title, value }: Props = $props();

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(value);
      copyState = "success";
      setTimeout(() => {
        copyState = undefined;
      }, 2000);
    } catch (error) {
      console.error(error);
      copyState = "failure";
    }
  };
</script>

<div class="relative">
  <button onclick={copyToClipboard} class="cursor-pointer" disabled={!value} {title}>
    <Copy />
  </button>
  {#if copyState}
    <div
      class="tooltip absolute -top-[80%] left-[calc(100%+4px)] z-10 mt-1 rounded border-1 border-gray-400 bg-white p-1 text-xs text-gray-800 dark:bg-gray-800 dark:text-gray-200"
    >
      {#if copyState === "success"}
        <div class="flex items-center gap-1">
          <Check />
          <span>Copied</span>
        </div>
      {:else}
        {errorMessage}
      {/if}
    </div>
  {/if}
</div>
