<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { Snippet } from "svelte";
  import { advisorySearchState } from "../advisory.svelte";

  interface Props {
    header: string;
    title?: string;
    open?: boolean;
    showBorder?: boolean;
    level?: number;
    highlight?: boolean;
    // path used inside advisory view when this component should open
    // automatically if it contains the current search match.
    path?: string;
    onOpen?: () => any;
    onClose?: () => any;
    children: Snippet;
    headerSlot?: Snippet;
    headerRightSlot?: Snippet;
  }
  let {
    header,
    title = undefined,
    open = $bindable(false),
    showBorder = true,
    level = 2,
    highlight = false,
    path = undefined,
    onOpen = () => {
      //default: Do nothing
    },
    onClose = () => {
      //default: Do nothing
    },
    children,
    headerSlot = undefined,
    headerRightSlot = undefined
  }: Props = $props();
  const uuid = crypto.randomUUID();

  let visibility = $state("none");
  $effect(() => {
    if (open) {
      visibility = "block";
    } else {
      visibility = "none";
    }
  });
  /**
   * toggle toggles visibility of content.
   */
  const toggle = () => {
    if (visibility === "block") {
      if (onClose) {
        onClose();
      }
      open = false;
    } else {
      if (onOpen) {
        onOpen();
      }
      open = true;
      setTimeout(() => {
        const element = document.getElementById(`${uuid}`);
        if (element) {
          const y = element!.getBoundingClientRect().top + window.scrollY - 150;
          window.scrollTo({ top: y, behavior: "smooth" });
        }
      }, 200);
    }
  };
  let icon = $derived.by(() => {
    if (visibility === "block") {
      return "bx-chevron-down";
    } else {
      return "bx-chevron-right";
    }
  });

  const getClass = (level: number) => {
    switch (level) {
      case 2:
        return "text-xl";
      case 3:
        return "text-lg";
      case 4:
        return "";
      case 5:
        return "text-xs";
      default:
        return "";
    }
  };

  $effect(() => {
    if (path) {
      if (advisorySearchState.matchIndex !== -1) {
        const hitPath = advisorySearchState.searchMatches[advisorySearchState.matchIndex]?.path;
        const shouldOpen = hitPath !== undefined && hitPath.startsWith(path);
        if (shouldOpen) open = true;
      }
    }
  });
</script>

<div class:collapsible={true} class:highlight-section={highlight}>
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div {title} id={header}>
    <div class="inline-flex items-center" onclick={toggle}>
      <i class="bx {getClass(level)} {icon}"></i>
      <div class={getClass(level)}>
        {#if headerSlot}
          {@render headerSlot()}
        {:else}
          <span class={getClass(level)}>{header}</span>
        {/if}
      </div>
      <span class="ms-2">
        {@render headerRightSlot?.()}
      </span>
    </div>
    {#if visibility === "block"}
      <div id={uuid} class={`ml-2 pl-2 ${showBorder ? "border-l-2 border-l-gray-200" : ""}`}>
        {@render children()}
      </div>
    {/if}
  </div>
</div>
