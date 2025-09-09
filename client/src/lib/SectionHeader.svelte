<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Heading } from "flowbite-svelte";
  import type { Snippet } from "svelte";

  interface Props {
    title: string;
    mainSlot?: Snippet;
    rightSlot?: Snippet;
  }
  let { title = "", mainSlot = undefined, rightSlot = undefined }: Props = $props();
</script>

<div class="mb-3 flex text-lg break-all text-gray-900 dark:text-white">
  <Heading tag="h2" class="me-4 w-fit" customSize="text-lg"
    >{title}
    {#if mainSlot}
      {@render mainSlot()}
    {/if}
  </Heading>
  {#if rightSlot}
    {@render rightSlot()}
  {/if}
</div>
