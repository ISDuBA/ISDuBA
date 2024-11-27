<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Input, Label } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let clearable = false;
  export let from: string | undefined;
  export let to: string | undefined;

  const dispatch = createEventDispatcher();
  const uuid = crypto.randomUUID();
  const fromId = `from-${uuid}`;
  const toId = `to-${uuid}`;
  const iconClass = "bx bx-x text-lg";
  const resetButtonClass = "rounded-s-none px-3";
  const defaultInputClass = "h-fit";
  const inputClass = `${defaultInputClass} ${clearable ? "rounded-e-none" : ""}`;
  let hideFrom = false;
  let hideTo = false;

  const onChange = () => {
    dispatch("change");
  };

  /*
  This method results in removal and addition of the date input from and to the DOM.
  Otherwise it is not possible to remove the value from the input field once a selection was made.
  */
  const clearFrom = () => {
    hideFrom = true;
    from = undefined;
    hideFrom = false;
    onChange();
  };

  const clearTo = () => {
    hideTo = true;
    to = undefined;
    hideTo = false;
    onChange();
  };
</script>

<div class="flex gap-4">
  <div class="flex items-center gap-1">
    <Label for={fromId}>
      <span>From:</span>
    </Label>
    <div class="flex">
      {#if !hideFrom}
        <Input class={inputClass} let:props>
          <input on:change={onChange} id={fromId} type="date" {...props} bind:value={from} />
        </Input>
      {/if}
      {#if clearable}
        <Button color="light" class={resetButtonClass} on:click={clearFrom}>
          <i class={iconClass}></i>
        </Button>
      {/if}
    </div>
  </div>
  <div class="flex items-center gap-1">
    <Label for={toId}>
      <span>To:</span>
    </Label>
    <div class="flex">
      {#if !hideTo}
        <Input class={inputClass} let:props>
          <input on:change={onChange} id={toId} type="date" {...props} bind:value={to} />
        </Input>
      {/if}
      {#if clearable}
        <Button color="light" class={resetButtonClass} on:click={clearTo}>
          <i class={iconClass}></i>
        </Button>
      {/if}
    </div>
  </div>
</div>
