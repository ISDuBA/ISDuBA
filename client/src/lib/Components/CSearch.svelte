<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Input } from "flowbite-svelte";

  interface Props {
    buttonText?: string;
    containerClass?: string | undefined;
    onClear?: () => void;
    onKeyup?: (event: any) => void;
    placeholder?: string;
    searchTerm: string;
    search?: () => void;
  }

  let {
    buttonText = "Search",
    containerClass = undefined,
    onClear = undefined,
    onKeyup = undefined,
    placeholder = "Enter a search term",
    searchTerm = $bindable(),
    search
  }: Props = $props();

  const dispatchSearchEvent = () => {
    if (search) search();
  };

  const clearInput = () => {
    searchTerm = "";
    dispatchSearchEvent();
    if (onClear) onClear();
  };
</script>

<div class={containerClass ?? "relative flex w-full md:w-fit"}>
  <div class="relative w-full md:w-96">
    <Input
      class="w-full min-w-60 !rounded-e-none disabled:cursor-not-allowed disabled:opacity-50 rtl:text-right"
      size="md"
      {placeholder}
      bind:value={searchTerm}
      onkeyup={(e) => {
        if (onKeyup) onKeyup(e);
        if (e.key === "Enter") dispatchSearchEvent();
      }}
    >
      {#snippet right()}
        <button
          onclick={clearInput}
          aria-label="Clear search"
          class="group flex h-[26pt] w-[26pt] items-center justify-center rounded-md hover:bg-gray-200 dark:hover:bg-gray-500"
        >
          <i class="bx bx-x dark:group-hover:text-gray-800"></i>
        </button>
      {/snippet}
    </Input>
  </div>
  <Button size="sm" class="rounded-s-none" onclick={dispatchSearchEvent}>{buttonText}</Button>
</div>
