<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Input, Label, type InputType } from "flowbite-svelte";
  import { getContext, onMount, tick } from "svelte";

  export let value: string | number | Date | undefined;
  export let label: string;
  export let id: string;
  export let required = false;
  export let placeholder = "";
  export let type: InputType = "text";
  export let min: number | undefined = undefined;
  export let step: number | undefined = undefined;
  export let minlength: number | undefined = undefined;
  export let rules: any[] | undefined = undefined;
  export let errorMessage = "-";

  let element: HTMLInputElement;
  let isValid = true;
  let color: "base" | "red" | "green" | undefined;
  $: color = isValid ? "base" : "red";

  const validation: any = getContext("validation");

  onMount(() => {
    // Cannot add type as prop as it leads to this error:
    // 'type' attribute cannot be dynamic if input uses two-way binding
    element.type = type;
  });

  const validate = async () => {
    await tick();
    isValid = true;
    if (rules) {
      for (let i = 0; i < rules.length; i++) {
        const result = rules[i](value);
        if (typeof result === "string") {
          errorMessage = result;
          isValid = false;
          break;
        }
      }
    } else if (!element.validity.valid && element.validationMessage) {
      errorMessage = element.validationMessage;
      isValid = false;
    }
    if (isValid && validation) {
      validation.removeInvalidField(id);
    } else {
      validation.addInvalidField(id);
    }
  };

  const onEdited = () => {
    validate();
  };
</script>

<div>
  <Label for={id} class="mb-1 block">
    <span>{label}</span>
    {#if required}
      <span class="text-red-500">*</span>
    {/if}
  </Label>
  <Input {color} let:props>
    <input
      on:input={onEdited}
      bind:this={element}
      bind:value
      {id}
      {placeholder}
      {required}
      {min}
      {step}
      {minlength}
      {...props}
      class:border={true}
    />
  </Input>
  <span class:text-xs={true} class:text-red-500={true} class:invisible={isValid}
    >{errorMessage}</span
  >
</div>
