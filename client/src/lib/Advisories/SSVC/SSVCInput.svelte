<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Label, Input } from "flowbite-svelte";
  import { vectorStart } from "./SSVCCalculator";
  import { onMount, tick } from "svelte";

  interface Props {
    autofocus: boolean;
    disabled: boolean;
    isValid: boolean;
    value: string;
    onKeyup?: (event: any) => void;
  }

  let {
    autofocus = false,
    disabled = false,
    /* eslint-disable-next-line no-useless-assignment */
    isValid = $bindable(),
    value = $bindable(),
    onKeyup = undefined
  }: Props = $props();

  const defaultInputClass = "border-gray-600 rounded-s-none h-10";

  let minLengthReached = $state(false);
  let endsWithSlash = $state(false);
  let containsValidDate = $state(false);
  let inputValue = $state("");

  onMount(() => {
    const newValue = value.replace(vectorStart, "");
    isValid = isVectorValid(newValue);
    inputValue = newValue;
  });

  /**
   * Checks if the vector has a valid format. It doesn't validate the semantic value.
   */
  const isVectorValid = (valueToValidate: string): boolean => {
    const parts = valueToValidate.split("/");
    const decisions = parts.filter((p) => p.length === 3 && p.charAt(1) === ":");
    minLengthReached = decisions.length >= 5;
    endsWithSlash = valueToValidate.charAt(inputValue.length - 1) === "/";
    const lastPart = parts[parts.length - 1];
    const secondToLastPart = parts[parts.length - 2];
    containsValidDate =
      (lastPart === "" && secondToLastPart !== undefined && isDateValid(secondToLastPart)) ||
      (lastPart !== undefined && isDateValid(lastPart));
    return minLengthReached && endsWithSlash && containsValidDate;
  };

  const isDateValid = (date: string) => {
    const dateRegex = /\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\dZ/;
    const result = date.match(dateRegex);
    return result?.length === 1;
  };

  /**
   * Increase the usability by adding characters automatically which always have to be entered.
   */
  const autoAddCharacters = () => {
    const letterRegex = /[a-zA-Z]/;
    const lowerCaseRegex = /[a-z]/;
    const last = inputValue.charAt(inputValue.length - 1);
    if (inputValue.length > 1 && last.match(letterRegex)) {
      if (last.match(lowerCaseRegex)) {
        inputValue = inputValue.slice(0, -1) + last.toUpperCase();
      }
      const secondToLast = inputValue.charAt(inputValue.length - 2);
      if (secondToLast === ":") {
        inputValue = `${inputValue}/`;
      } else if (secondToLast === "/") {
        inputValue = `${inputValue}:`;
      }
    } else if (inputValue.length === 1 && inputValue.match(letterRegex)) {
      if (last.match(lowerCaseRegex)) {
        inputValue = inputValue.slice(0, -1) + last.toUpperCase();
      }
      inputValue = `${inputValue}:`;
    }
  };

  const handleInput = async (event: any) => {
    await tick();
    const newInput = event.target.value;
    value = `${vectorStart}${inputValue}`;
    isValid = isVectorValid(newInput);
    if (onKeyup) {
      onKeyup(event);
    }
  };

  const handlePasteEvent = async (event: any) => {
    handleInput(event);
  };

  const handleKeyEvent = (event: any) => {
    if (event.key === "v" && event.ctrlKey === true) {
      // Already handled by listener of the paste event
      return;
    }
    if (!["Backspace", "Escape"].includes(event.key)) {
      autoAddCharacters();
    }
    handleInput(event);
  };
</script>

<div class="mb-3 flex w-full">
  <Label
    class="flex h-10 items-center rounded-s-md border border-r-0 border-gray-400 px-2 text-gray-500"
    >{vectorStart}</Label
  >
  <div class="flex w-full flex-col gap-y-2">
    <Input
      onkeyup={handleKeyEvent}
      onpaste={handlePasteEvent}
      bind:value={inputValue}
      {autofocus}
      {disabled}
      class={defaultInputClass}
      type="text"
    />
    <div class="flex flex-col gap-1 ps-2">
      <div
        class="flex items-baseline gap-1"
        class:text-gray-600={!minLengthReached}
        class:dark:text-gray-400={!minLengthReached}
        class:text-green-600={minLengthReached}
        class:dark:text-green-400={minLengthReached}
      >
        {#if minLengthReached}
          <i class="bx bx-check-circle"></i>
        {:else}
          <i class="bx bx-x-circle"></i>
        {/if}
        <span>At least 5 key pairs</span>
      </div>
      <div
        class="flex items-baseline gap-1"
        class:text-gray-600={!containsValidDate}
        class:dark:text-gray-400={!containsValidDate}
        class:text-green-600={containsValidDate}
        class:dark:text-green-400={containsValidDate}
      >
        {#if containsValidDate}
          <i class="bx bx-check-circle"></i>
        {:else}
          <i class="bx bx-x-circle"></i>
        {/if}
        <span>Contains valid date (yyyy-mm-ddThh:mm:ssZ) after last key pair</span>
      </div>
      <div
        class="flex items-baseline gap-1"
        class:text-gray-600={!endsWithSlash}
        class:dark:text-gray-400={!endsWithSlash}
        class:text-green-600={endsWithSlash}
        class:dark:text-green-400={endsWithSlash}
      >
        {#if endsWithSlash}
          <i class="bx bx-check-circle"></i>
        {:else}
          <i class="bx bx-x-circle"></i>
        {/if}
        <span>Ends with "/"</span>
      </div>
    </div>
  </div>
</div>
