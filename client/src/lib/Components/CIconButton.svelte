<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { type ColorVariant } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let icon = "";
  export let color: ColorVariant = "dark";
  export let title = "";

  const dispatch = createEventDispatcher();
</script>

<button
  on:click|preventDefault={(event) => {
    event.stopPropagation();
    dispatch("click", event);
  }}
  {title}
  class="p-1"
>
  <i class={`bx bx-${icon} text-${color}-600 text-lg ${$$props.class}`}></i>
</button>
