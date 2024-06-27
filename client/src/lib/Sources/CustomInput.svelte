<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Input, Label, type InputType } from "flowbite-svelte";

  export let value: string | number | Date | undefined;
  export let label: string;
  export let id: string;
  export let required = false;
  export let placeholder = "";
  export let type: InputType = "text";
</script>

<div>
  <Label for={id} class="mb-2 block">
    <span>{label}</span>
    {#if required}
      <span class="text-red-500">*</span>
    {/if}
  </Label>
  <Input {id} bind:value {placeholder} {required} {type} />
</div>
