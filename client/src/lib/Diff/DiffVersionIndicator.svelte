<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  export let color: "red" | "green" | "gray";
  export let icon: "minus" | "plus";
  export let isDisabled = false;
  export let permanent = false;
  $: containerClass = `h-4 flex items-center rounded-full border-solid border-${color}-700 ${isDisabled ? "" : "group-hover:border"}`;
  $: iconClass = `bx bx-${icon} flex h-4 items-center text-${color}-700 ${permanent ? "" : "opacity-0 group-hover:opacity-100"}`;
</script>

<div class={containerClass}>
  <i class={iconClass}></i>
</div>
