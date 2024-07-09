<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Badge, Input, Label, type InputType } from "flowbite-svelte";
  import { createEventDispatcher, getContext, onMount, tick } from "svelte";

  export let id: string;
  export let initialEntries: string[];
  export let label = "";
  export let placeholder = "";
  export let required = false;
  export let type: InputType = "text";
  export let pattern: string | undefined = undefined;
  export let minEntries: number | undefined = undefined;

  const dispatch = createEventDispatcher();
  let element: HTMLInputElement;
  let errorMessage = "-";
  let isValid = true;
  let active = false;
  let input = "";
  let entries: any[] = [];
  $: hasCorrectLength = minEntries ? entries.length >= minEntries : true;
  let borderInvalid = "border-[2px] dark:border-red-600 border-red-500 p-[4px]";
  let borderActive = "border-[2px] dark:border-gray-600 border-primary-500 p-[4px]";
  let borderInactive = "border-[1px] dark:border-gray-600 border-gray-300 p-[5px]";
  $: containerBorder = `border rounded-lg ${(isValid && hasCorrectLength) || !didBlur ? (active ? borderActive : borderInactive) : borderInvalid}`;
  let containerText =
    "dark:text-white dark:placeholder-gray-400 rtl:text-right text-sm text-gray-900";
  let containerCustom = "flex gap-x-2 gap-y-1 flex-wrap items-center";
  $: containerClass = `${containerBorder} ${containerText} ${containerCustom} focus:ring-primary-500 dark:focus:ring-primary-500 w-full bg-gray-50 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-gray-700`;
  $: if (initialEntries.length > 0) {
    entries = initialEntries;
  }
  let didBlur = false;

  onMount(() => {
    // Cannot add type as prop as it leads to this error:
    // 'type' attribute cannot be dynamic if input uses two-way binding
    element.type = type;
  });

  const validation: any = getContext("validation");

  const onKeypress = (event: any) => {
    validate();
    if ([" ", ",", "Enter"].includes(event.key)) {
      event.preventDefault();
      addEntry();
      if (event.key === "Enter") dispatch("submit");
    } else if (event.key === "Backspace" && input === "" && entries.length > 0) {
      removeEntry(entries.length - 1);
    }
  };

  const addEntry = () => {
    if (input.trim().length > 0 && isValid) {
      entries = entries.concat([input]);
      input = "";
      dispatch("edited", entries);
    }
  };

  const removeEntry = (index: number) => {
    didBlur = true;
    entries = entries.toSpliced(index, 1);
    dispatch("edited", entries);
    validateNumberOfEntries();
  };

  const validateNumberOfEntries = () => {
    if (minEntries !== undefined && entries.length <= minEntries) {
      errorMessage = `Add at least ${minEntries} entry${minEntries !== 1 ? "s" : ""}`;
      hasCorrectLength = false;
    }
  };

  const validate = async () => {
    if (!didBlur) return;
    await tick();
    isValid = true;
    if (!element.validity.valid && element.validationMessage) {
      errorMessage = element.validationMessage;
      isValid = false;
    } else if (minEntries && entries.length < minEntries) {
      validateNumberOfEntries();
    }
    if (isValid && validation) {
      validation.removeInvalidField(id);
    } else {
      validation.addInvalidField(id);
    }
  };
</script>

<div>
  <Label for={label} class="mb-1 block">
    <span>{label}</span>
    {#if required}
      <span class="text-red-500">*</span>
    {/if}
  </Label>
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

    <Input let:props>
      <input
        on:blur={() => {
          didBlur = true;
          active = false;
          addEntry();
        }}
        on:focus={() => (active = true)}
        on:keydown={onKeypress}
        bind:value={input}
        bind:this={element}
        id={label}
        {placeholder}
        {pattern}
        {...props}
        class="w-[unset] flex-grow border-0 p-1 focus:border-0 focus:outline-none focus:ring-0"
      />
    </Input>
  </div>
  <span
    class:text-xs={true}
    class:text-red-500={true}
    class:invisible={(isValid && hasCorrectLength) || !didBlur}>{errorMessage}</span
  >
</div>
