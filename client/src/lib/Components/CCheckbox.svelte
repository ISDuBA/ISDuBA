<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Checkbox, type CheckboxItem } from "flowbite-svelte";
  import type { FormColorType } from "./types";
  import type { Snippet } from "svelte";
  import type { HTMLButtonAttributes } from "svelte/elements";

  type Props = {
    ariaLabel?: string;
    checked?: boolean;
    choices?: CheckboxItem[];
    color?: FormColorType;
    custom?: boolean;
    disabled?: boolean;
    group?: string[];
    inline?: boolean;
    name?: string | undefined;
    value?: string | number;
    onChanged?: (event: any) => void;
    onClicked?: (event: any) => void;
    children?: Snippet;
  } & HTMLButtonAttributes;

  let {
    ariaLabel = undefined,
    checked = $bindable(false),
    choices = [],
    color = "primary",
    custom = false,
    disabled = false,
    group = [],
    inline = false,
    name = undefined,
    value = undefined,
    onChanged = undefined,
    onClicked = undefined,
    children,
    ...restProps
  }: Props = $props();
</script>

<Checkbox
  bind:checked
  aria-label={ariaLabel}
  class="min-h-[20px] min-w-[20px]"
  classes={{ div: `p-[6px]! py-[6px]! cursor-pointer ${restProps.class}` }}
  {choices}
  {color}
  {custom}
  {disabled}
  {group}
  {inline}
  {name}
  {value}
  onchange={(event) => {
    if (onChanged) {
      onChanged(event);
    }
  }}
  onclick={(event) => {
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
