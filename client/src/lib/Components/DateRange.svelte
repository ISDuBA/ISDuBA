<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Input, Label } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let from: string | undefined;
  export let to: string | undefined;

  const dispatch = createEventDispatcher();
  const uuid = crypto.randomUUID();
  const fromId = `from-${uuid}`;
  const toId = `to-${uuid}`;

  const onChange = () => {
    dispatch("change");
  };
</script>

<div class="flex gap-4">
  <div class="flex items-center gap-1">
    <Label for={fromId}>
      <span>From:</span>
    </Label>
    <Input class="h-fit" let:props>
      <input on:change={onChange} id={fromId} type="date" {...props} bind:value={from} />
    </Input>
  </div>
  <div class="flex items-center gap-1">
    <Label for={toId}>
      <span>To:</span>
    </Label>
    <Input class="h-fit" let:props>
      <input on:change={onChange} id={toId} type="date" {...props} bind:value={to} />
    </Input>
  </div>
</div>
