<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Checkbox, type CheckboxItem, type FormColorType } from "flowbite-svelte";
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";

  type Props = {
    checked?: boolean;
    choices?: CheckboxItem[];
    color?: FormColorType;
    custom?: boolean;
    disabled?: boolean;
    group?: string[];
    groupInputClass?: string;
    groupLabelClass?: string;
    inline?: boolean;
    name?: string | undefined;
    spacing?: string;
    value?: string | number;
    onChanged?: (event: any) => void;
    onClicked?: (event: any) => void;
    children?: Snippet;
  } & HTMLButtonAttributes;

  let {
    checked = $bindable(false),
    choices = [],
    color = "primary",
    custom = false,
    disabled = false,
    group = [],
    groupInputClass = "",
    groupLabelClass = "",
    inline = false,
    name = undefined,
    spacing = "",
    value = "on",
    onChanged = undefined,
    onClicked = undefined,
    children,
    ...restProps
  }: Props = $props();
</script>

<Checkbox
  bind:checked
  class={`min-h-[20px] min-w-[20px] cursor-pointer !p-[6px] !py-[6px] ${restProps.class}`}
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
    if (onChanged) {
      onChanged(event);
    }
  }}
  on:click={(event) => {
    event.stopPropagation();
    if (onClicked) {
      onClicked(event);
    }
  }}
>
  <span class="ps-2">
    {@render children?.()}
  </span>
</Checkbox>
