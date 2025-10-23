<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  interface Props {
    color: "red" | "green" | "gray";
    icon: "minus" | "plus";
    hoverIcon: "minus" | "plus" | undefined;
    isDisabled?: boolean;
    permanent?: boolean;
  }

  let { color, icon, hoverIcon, isDisabled = false, permanent = false }: Props = $props();
  let containerClass = $derived(
    `h-4 flex items-center rounded-full border-solid border-${color}-700 dark:border-${color}-300 ${isDisabled ? "" : "group-hover:border"}`
  );
  let basicIconClass = "bx flex h-4 items-center";
  let iconClass = $derived(
    `${basicIconClass} bx-${icon} text-${color}-700 dark:text-${color}-300 ${!permanent ? "!hidden" : ""} ${isDisabled ? "" : "group-hover:!hidden"}`
  );
  let hoverClass = $derived(
    `${basicIconClass} bx-${hoverIcon ?? icon} text-${color}-700 dark:text-${color}-300 !hidden ${isDisabled ? "" : "group-hover:!inline"}`
  );
</script>

<div class={containerClass}>
  <i class={iconClass}></i>
  <i class={hoverClass}></i>
</div>
