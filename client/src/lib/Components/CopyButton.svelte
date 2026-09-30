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
  import { Button } from "flowbite-svelte";

  interface Props {
    copyState?: CopyState;
    errorMessage: string;
    paddingClass?: string;
    showBorder?: boolean;
    successMessage?: string;
    title: string;
    tooltipPlacement?: "topright" | "bottomright";
    value: string;
  }
  let {
    copyState = $bindable(undefined),
    errorMessage,
    paddingClass = "py-1",
    showBorder = false,
    successMessage = "Copied",
    title,
    tooltipPlacement = "topright",
    value
  }: Props = $props();

  let buttonClass = $derived.by(() => {
    let btnClass = "cursor-pointer h-7 " + paddingClass;
    if (!showBorder) {
      btnClass += " border-0";
    }
    return btnClass;
  });

  let tooltipPlacementClass = $derived.by(() => {
    switch (tooltipPlacement) {
      case "bottomright":
        return "-bottom-[80%]";
      case "topright":
        return "-top-[80%]";
    }
  });

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
  <Button
    onclick={copyToClipboard}
    class={buttonClass}
    color="light"
    disabled={!value}
    size="xs"
    {title}
  >
    <Copy />
  </Button>
  {#if copyState}
    <div
      class={`tooltip absolute ${tooltipPlacementClass} left-[calc(100%+4px)] z-10 mt-1 rounded border border-gray-400 bg-white p-1 text-xs text-gray-800 dark:bg-gray-800 dark:text-gray-200`}
    >
      {#if copyState === "success"}
        <div class="flex items-center gap-1">
          <Check />
          <span>{successMessage}</span>
        </div>
      {:else}
        {errorMessage}
      {/if}
    </div>
  {/if}
</div>
