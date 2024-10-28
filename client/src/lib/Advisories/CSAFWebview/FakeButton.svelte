<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  export let active = false;
  const activeClass =
    "h-7 py-1 px-3 border-gray-300 border text-xs bg-gray-200 hover:bg-gray-100 rounded-lg max-w-max";
  const inactiveClass =
    "h-7 py-1 px-3 border-gray-300 border text-xs hover:bg-gray-100 rounded-lg max-w-max";
</script>

<div class={active ? activeClass : inactiveClass}>
  <slot></slot>
</div>
