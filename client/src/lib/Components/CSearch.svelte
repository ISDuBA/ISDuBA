<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Input } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let containerClass: string | undefined = undefined;
  export let searchTerm: string;

  const dispatch = createEventDispatcher();

  const dispatchSearchEvent = () => {
    dispatch("search");
  };

  const clearInput = () => {
    searchTerm = "";
    dispatchSearchEvent();
  };
</script>

<div class={containerClass ?? "relative flex"}>
  <div class="relative w-96">
    <Input
      defaultClass="w-full disabled:cursor-not-allowed disabled:opacity-50 rtl:text-right !rounded-e-none"
      size="md"
      placeholder="Enter a search term"
      bind:value={searchTerm}
      on:keyup={(e) => {
        if (e.key === "Enter") dispatchSearchEvent();
      }}
    >
      <button
        on:click={clearInput}
        slot="right"
        class="group flex h-[26pt] w-[26pt] items-center justify-center rounded-md hover:bg-gray-200 dark:hover:bg-gray-500"
      >
        <i class="bx bx-x dark:group-hover:text-gray-800"></i>
      </button>
    </Input>
  </div>
  <Button size="xs" class="rounded-s-none" on:click={dispatchSearchEvent}>Search</Button>
</div>
