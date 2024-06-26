<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Badge, Input, Label } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let initialEntries: string[];
  export let label = "";
  export let placeholder = "";

  const dispatch = createEventDispatcher();
  let input = "";
  let entries: any[] = [];
  let containerBorder =
    "border focus:border-primary-500 dark:border-gray-600 dark:focus:border-primary-500 border-gray-300 rounded-lg";
  let containerText =
    "dark:text-white dark:placeholder-gray-400 rtl:text-right text-sm text-gray-900";
  let containerCustom = "flex gap-x-2 gap-y-1 p-1 flex-wrap";
  let containerClass = `${containerBorder} ${containerText} ${containerCustom} focus:ring-primary-500 dark:focus:ring-primary-500 w-full bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-gray-700`;
  $: if (initialEntries.length > 0) {
    entries = initialEntries;
  }

  const onKeypress = (event: any) => {
    if (["Enter", ","].includes(event.key) && input.length > 0) {
      entries = entries.concat([input]);
      input = "";
      dispatch("edited", entries);
    } else if (event.key === "Backspace" && input === "" && entries.length > 0) {
      entries = entries.toSpliced(entries.length - 1, 1);
    }
  };

  const removeEntry = (index: number) => {
    entries = entries.toSpliced(index, 1);
    dispatch("edited", entries);
  };
</script>

<div>
  <Label for={label} class="mb-2 block">{label}</Label>
  <div class={containerClass}>
    {#each entries as entry, index (index)}
      <Badge border dismissable>
        <button
          on:click={() => removeEntry(index)}
          slot="close-button"
          class="-my-2 -me-1 ms-2 text-lg"
        >
          <i class="bx bx-x text-primary-600"></i>
        </button>
        <span>{entry}</span>
      </Badge>
    {/each}

    <Input
      on:keydown={onKeypress}
      size="sm"
      class="border-0"
      id={label}
      bind:value={input}
      {placeholder}
    />
  </div>
</div>
