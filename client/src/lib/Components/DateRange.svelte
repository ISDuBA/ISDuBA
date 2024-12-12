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
  import TimeInput from "./TimeInput.svelte";

  export let clearable = false;
  export let from: Date | undefined;
  export let to: Date | undefined;
  export let showTimeControls = false;

  const dispatch = createEventDispatcher();
  const uuid = crypto.randomUUID();
  const fromId = `from-${uuid}`;
  const toId = `to-${uuid}`;
  const iconClass = "bx bx-x text-lg";
  const resetButtonClass = "rounded-s-none px-3";
  const defaultInputClass = "h-fit";
  const inputClass = `${defaultInputClass} ${clearable || showTimeControls ? "rounded-e-none" : ""}`;
  let hideFrom = false;
  let hideTo = false;
  let fromString: string | undefined;
  let toString: string | undefined;
  let fromTimeInput: any;
  let toTimeInput: any;

  $: if (from) {
    fromString = from.toISOString().split("T")[0];
  }
  $: if (to) {
    toString = to.toISOString().split("T")[0];
  }

  const onChange = () => {
    dispatch("change");
  };

  /*
  This method results in removal and addition of the date input from and to the DOM.
  Otherwise it is not possible to remove the value from the input field once a selection was made.
  */
  const clearFrom = () => {
    hideFrom = true;
    fromString = undefined;
    hideFrom = false;
    from = undefined;
    fromTimeInput.clearInput();
    onChange();
  };

  const clearTo = () => {
    hideTo = true;
    toString = undefined;
    hideTo = false;
    to = undefined;
    toTimeInput.clearInput();
    onChange();
  };

  const onDateChanged = (event: any) => {
    if (event.target.id.startsWith("from")) {
      fromString = event.target.value;
    } else if (event.target.id.startsWith("to")) {
      toString = event.target.value;
    }
    if (fromString) {
      const fromDate = new Date(fromString);
      if (!from) {
        from = fromDate;
      } else {
        from.setDate(fromDate.getDate());
        from.setMonth(fromDate.getMonth());
        from.setFullYear(fromDate.getFullYear());
      }
    }
    if (toString) {
      const toDate = new Date(toString);
      if (!to) {
        to = toDate;
      } else {
        to.setDate(toDate.getDate());
        to.setMonth(toDate.getMonth());
        to.setFullYear(toDate.getFullYear());
      }
    }
    onChange();
  };

  const onFromTimeChanged = (event: any) => {
    if (event.detail) {
      const detail = event.detail;
      if (!from) from = new Date();
      if (detail.hours !== undefined) {
        from.setUTCHours(Number(detail.hours));
      }
      if (detail.minutes !== undefined) {
        from.setUTCMinutes(Number(detail.minutes));
      }
      from = from;
      onChange();
    }
  };

  const onToTimeChanged = (event: any) => {
    if (event.detail) {
      const detail = event.detail;
      if (!to) to = new Date();
      if (detail.hours !== undefined) {
        to.setUTCHours(Number(detail.hours));
      }
      if (detail.minutes !== undefined) {
        to.setUTCMinutes(Number(detail.minutes));
      }
      to = to;
      onChange();
    }
  };
</script>

<div class="flex flex-wrap gap-4">
  <div class="flex items-center gap-1">
    <Label for={fromId}>
      <span>From</span>
      {#if showTimeControls}
        <span>(UTC)</span>
      {/if}
      <span>:</span>
    </Label>
    <div class="flex">
      {#if !hideFrom}
        <Input class={inputClass} let:props>
          <input on:input={onDateChanged} id={fromId} type="date" {...props} value={fromString} />
        </Input>
      {/if}
      {#if showTimeControls}
        <TimeInput
          bind:this={fromTimeInput}
          on:timeChanged={onFromTimeChanged}
          roundEnd={!clearable}
        ></TimeInput>
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
      <span>To</span>
      {#if showTimeControls}
        <span>(UTC)</span>
      {/if}
      <span>:</span>
    </Label>
    <div class="flex">
      {#if !hideTo}
        <Input class={inputClass} let:props>
          <input on:input={onDateChanged} id={toId} type="date" {...props} value={toString} />
        </Input>
      {/if}
      {#if showTimeControls}
        <TimeInput bind:this={toTimeInput} on:timeChanged={onToTimeChanged} roundEnd={!clearable}
        ></TimeInput>
      {/if}
      {#if clearable}
        <Button color="light" class={resetButtonClass} on:click={clearTo}>
          <i class={iconClass}></i>
        </Button>
      {/if}
    </div>
  </div>
</div>
