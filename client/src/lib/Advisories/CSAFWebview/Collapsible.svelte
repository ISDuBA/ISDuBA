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
  export let level = "2";
  export let class_ = "";
  export let highlight = false;
  const uuid = crypto.randomUUID();
  export let onOpen = () => {
    //default: Do notthing
  };
  export let onClose = () => {
    //default: Do notthing
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
        const y = element!.getBoundingClientRect().top + window.scrollY - 150;
        window.scrollTo({ top: y, behavior: "smooth" });
      }, 200);
      visibility = "block";
    }
  };
  $: if (visibility === "block") {
    icon = "bx-chevron-down";
  } else {
    icon = "bx-chevron-right";
  }
</script>

<div class:collapsible={true} class:highlight-section={highlight}>
  {#if level == "2"}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div title={header} id={header} on:click={toggle} class={class_}>
      <h2><i class="bx {icon}" />{header}</h2>
    </div>
  {/if}
  {#if level == "3"}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div title={header} id={header} on:click={toggle} class={class_}>
      <h3><i class="bx {icon}" />{header}</h3>
    </div>
  {/if}
  {#if level == "4"}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div title={header} id={header} on:click={toggle} class={class_}>
      <h4><i class="bx {icon}" />{header}</h4>
    </div>
  {/if}
  {#if level == "5"}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div title={header} id={header} on:click={toggle} class={class_}>
      <h5><i class="bx {icon}" />{header}</h5>
    </div>
  {/if}
  {#if visibility === "block"}
    <div id={uuid} class="collapsible-body">
      <slot />
    </div>
  {/if}
</div>
