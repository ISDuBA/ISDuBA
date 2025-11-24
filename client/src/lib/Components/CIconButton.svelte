<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { ColorVariant } from "./types";
  import type { HTMLButtonAttributes } from "svelte/elements";

  type Props = {
    ariaLabel?: string;
    color?: ColorVariant;
    icon?: string;
    disabled?: boolean;
    title?: string;
    onClicked?: (event?: any) => void;
  } & HTMLButtonAttributes;

  let {
    ariaLabel = undefined,
    color = "dark",
    disabled = false,
    icon = "",
    title = "",
    onClicked = () => {},
    ...restProps
  }: Props = $props();
</script>

<button
  onclick={(event) => {
    event.preventDefault();
    onClicked(event);
  }}
  {disabled}
  {title}
  aria-label={ariaLabel || `Icon button ${icon}`}
  id={restProps.id}
  class="p-1"
>
  <i
    aria-hidden="true"
    class={`bx bx-${icon} text-${color}-600 text-lg ${restProps.class} ${disabled ? "contrast-0 saturate-0" : ""}`}
  ></i>
</button>
