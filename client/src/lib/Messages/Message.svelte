<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Toast } from "flowbite-svelte";
  import { MESSAGE } from "./messagetypes";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { blur } from "svelte/transition";

  export let error: any = null;
  export let plain = false;
  export let dismissable = true;
  export let autoclose = true;
  $: divClass = plain
    ? "w-full max-w-xs p-0 text-gray-500 bg-white dark:text-gray-400 dark:bg-gray-800 gap-3"
    : "w-full max-w-xs p-4 text-gray-500 bg-white shadow dark:text-gray-400 dark:bg-gray-800 gap-3";

  let open = true;

  const coloryByType = (type: string) => {
    if (type === MESSAGE.ERROR) return "red";
    if (type === MESSAGE.WARNING) return "yellow";
    if (type === MESSAGE.SUCCESS) return "green";
    return "blue";
  };

  onMount(async () => {
    if (!autoclose) return;
    setTimeout(() => {
      open = false;
      appStore.removeError(error.id);
    }, 8000);
  });
</script>

{#if error}
  <Toast
    color={coloryByType(error.type)}
    bind:open
    transition={blur}
    on:close={() => {
      appStore.removeError(error.id);
    }}
    {divClass}
    {dismissable}
  >
    <svelte:fragment slot="icon">
      <i
        class:bx={true}
        class:bxs-message-rounded-x={error.type === MESSAGE.ERROR}
        class:bxs-message-rounded-error={error.type === MESSAGE.WARNING}
        class:bxs-message-rounded-check={error.type === MESSAGE.SUCCESS}
        class:bxs-message-rounded={error.type === MESSAGE.INFO}
      ></i>
    </svelte:fragment>
    <span class="mb-1 text-sm font-semibold text-gray-900 dark:text-white">{error.type}</span>
    <div class="mb-2 text-sm font-normal">{error.message}</div>
  </Toast>
{/if}
