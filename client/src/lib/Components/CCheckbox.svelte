<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Checkbox, type CheckboxItem, type FormColorType } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let checked: boolean | undefined = undefined;
  export let choices: CheckboxItem[] = [];
  export let color: FormColorType = "primary";
  export let custom = false;
  export let disabled = false;
  export let group: string[] = [];
  export let groupInputClass = "";
  export let groupLabelClass = "";
  export let inline = false;
  export let name: string | undefined = undefined;
  export let spacing = "";
  export let value: string | number = "on";

  const dispatch = createEventDispatcher();
</script>

<Checkbox
  {checked}
  class={`min-h-[20px] min-w-[20px] cursor-pointer !p-[6px] !py-[6px] ${$$props.class}`}
  {choices}
  {color}
  {custom}
  {disabled}
  {group}
  {groupInputClass}
  {groupLabelClass}
  {inline}
  {name}
  {spacing}
  {value}
  on:change={(event) => {
    dispatch("change", event);
  }}
  on:click={(event) => {
    event.stopPropagation();
    dispatch("click", event);
  }}
></Checkbox>
