<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script>
  import { Heading } from "flowbite-svelte";

  export let title = "";
</script>

<div class="mb-3 flex break-all text-lg text-gray-900 dark:text-white">
  <Heading tag="h2" class="me-4 w-fit" customSize="text-lg">{title} <slot /></Heading>
  <slot name="right"></slot>
</div>
