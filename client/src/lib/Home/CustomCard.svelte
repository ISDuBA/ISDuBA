<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
</script>

<div class="rounded-md border border-solid border-gray-300 px-4 py-2 shadow-md">
  <div class="mb-2 flex flex-row justify-between text-xs">
    <div>
      <slot name="top-left"></slot>
    </div>
    <div>
      <slot name="top-right"></slot>
    </div>
  </div>
  <hr class="mb-2" />
  <div class="mb-3">
    <slot></slot>
  </div>
  <div class="flex flex-row justify-between text-xs">
    <div>
      <slot name="bottom-left"></slot>
    </div>
    <div>
      <slot name="bottom-right"></slot>
    </div>
  </div>
</div>
