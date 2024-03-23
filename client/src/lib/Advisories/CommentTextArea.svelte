<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Label, Textarea } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let value: string;
  export let buttonText: string;
  export let cancelable = false;
  $: count = value.length;

  const dispatch = createEventDispatcher();
</script>

<Textarea bind:value class="mb-2" id="comment-textarea">
  <div slot="footer" class="flex items-start justify-between">
    <div>
      <Button on:click={() => dispatch("saveComment")} disabled={count > 10000 || count === 0}>
        <span>{buttonText}</span>
      </Button>
      {#if cancelable}
        <Button on:click={() => dispatch("cancel")} outline color="red">
          <span>Cancel</span>
        </Button>
      {/if}
    </div>
    <Label class={count < 10000 ? "text-gray-600" : "font-bold text-red-600"}
      >{`${count}/10000`}</Label
    >
  </div>
</Textarea>
