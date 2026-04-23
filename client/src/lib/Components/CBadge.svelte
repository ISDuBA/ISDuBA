<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { twMerge } from "tailwind-merge";
  import { CloseButton } from "flowbite-svelte";
  import type { CloseButtonProps } from "flowbite-svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";
  import { fade, type TransitionConfig } from "svelte/transition";
  import type { Snippet } from "svelte";

  type TransitionFunc = (node: HTMLElement, params: any) => TransitionConfig;

  type Props = {
    color?: "dark" | "red" | "yellow" | "green" | "indigo" | "purple" | "pink" | "blue" | "primary";
    large?: boolean;
    dismissable?: boolean;
    showHoverEffect?: boolean;
    transition?: TransitionFunc;
    params?: object;
    onClose?: () => void;
    closeButtonSlot?: Snippet;
    children: Snippet;
  } & HTMLButtonAttributes;

  let {
    color = "primary",
    large = false,
    dismissable = false,
    showHoverEffect = false,
    transition = fade,
    params = {},
    onClose = undefined,
    closeButtonSlot = undefined,
    children,
    ...restProps
  }: Props = $props();

  let badgeStatus: boolean = $state(true);

  const colors = {
    primary: "bg-primary-100 text-primary-800 dark:bg-primary-900 dark:text-primary-300",
    dark: "bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-300",
    blue: "bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300",
    red: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
    green: "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
    yellow: "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
    indigo: "bg-indigo-100 text-indigo-800 dark:bg-indigo-900 dark:text-indigo-300",
    purple: "bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-300",
    pink: "bg-pink-100 text-pink-800 dark:bg-pink-900 dark:text-pink-300",
    none: ""
  };

  const hoverColors = {
    primary: "hover:bg-primary-200",
    dark: "hover:bg-gray-200",
    blue: "hover:bg-blue-200",
    red: "hover:bg-red-200",
    green: "hover:bg-green-200",
    yellow: "hover:bg-yellow-200",
    indigo: "hover:bg-indigo-200",
    purple: "hover:bg-purple-200",
    pink: "hover:bg-pink-200",
    none: ""
  };

  const close = () => {
    badgeStatus = false;
  };

  const baseClass: string =
    "font-medium inline-flex items-center justify-center px-2.5 py-0.5 rounded";

  let badgeClass: string = $derived(
    twMerge(
      baseClass,
      large ? "text-sm" : "text-xs",
      colors[color],
      showHoverEffect ? hoverColors[color] : "",
      `${restProps.class}`
    )
  );
</script>

{#if badgeStatus}
  <div transition:transition={params} class={badgeClass}>
    {@render children()}
    {#if dismissable}
      {#if closeButtonSlot}
        {@render closeButtonSlot()}
      {:else}
        <CloseButton
          class="ms-1.5 -me-1.5"
          color={color as CloseButtonProps["color"]}
          size={large ? "sm" : "xs"}
          ariaLabel="Remove badge"
          onclick={() => {
            close();
            if (onClose) {
              onClose();
            }
          }}
        />
      {/if}
    {/if}
  </div>
{/if}
