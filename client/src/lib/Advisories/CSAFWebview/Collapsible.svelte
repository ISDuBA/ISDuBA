<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  export let header: string;
  export let open = false;
  export let level = 2;
  export let class_ = "pl-4";
  export let highlight = false;
  const uuid = crypto.randomUUID();
  export let onOpen = () => {
    //default: Do nothing
  };
  export let onClose = () => {
    //default: Do nothing
  };
  let visibility = "none";
  $: if (open) {
    visibility = "block";
  } else {
    visibility = "none";
  }
  let icon = "bx-chevron-down";
  /**
   * toggle toggles visibility of content.
   */
  const toggle = () => {
    if (visibility === "block") {
      onClose();
      visibility = "none";
    } else {
      onOpen();
      setTimeout(() => {
        const element = document.getElementById(`${uuid}`);
        if (element) {
          const y = element!.getBoundingClientRect().top + window.scrollY - 150;
          window.scrollTo({ top: y, behavior: "smooth" });
        }
      }, 200);
      visibility = "block";
    }
  };
  $: if (visibility === "block") {
    icon = "bx-chevron-down";
  } else {
    icon = "bx-chevron-right";
  }

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
</script>

<div class:collapsible={true} class:highlight-section={highlight}>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div title={header} id={header} class={class_}>
    <div class="inline-flex items-center" on:click={toggle}>
      <i class="bx {getClass(level)} {icon}" />
      <slot class={getClass(level)} name="header"
        ><span class={getClass(level)}>{header}</span></slot
      >
    </div>
    {#if visibility === "block"}
      <div id={uuid} class={class_}>
        <slot />
      </div>
    {/if}
  </div>
</div>
