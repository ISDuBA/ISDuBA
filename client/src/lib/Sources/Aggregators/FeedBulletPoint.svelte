<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  export let filled = false;

  const background = "bg-gray-600 dark:bg-gray-600";
  const borderColor = "gray";
  $: border = `border border-${borderColor}-600`;
</script>

<span class={`${filled ? background : ""} ${border} inline-block h-3 w-3 rounded-[3pt]`}></span>
