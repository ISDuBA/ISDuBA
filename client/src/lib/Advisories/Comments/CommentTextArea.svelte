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
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { canSetStateReview } from "$lib/permissions";

  export let value: string;
  export let buttonText: string;
  export let cancelable = false;
  export let errorMessage: string;
  export let state: string = "";
  export let old = "";
  $: count = value.length;

  const dispatch = createEventDispatcher();
</script>

<Textarea bind:value on:input={() => dispatch("input")} class="mb-2" id="comment-textarea">
  <div slot="footer" class="flex flex-col">
    <div class="flex items-start justify-between">
      <div>
        {#if !cancelable}
          <Button
            size="xs"
            on:click={() => dispatch("saveForReview")}
            disabled={!canSetStateReview(state) && count === 0}
            color="light"
          >
            {#if count > 0}
              <span>Release for review with comment</span>
            {:else}
              <span>Release for review</span>
            {/if}
          </Button>
        {/if}
        <Button
          size="xs"
          on:click={() => dispatch("saveComment")}
          disabled={count > 10000 || count === 0 || value === old}
        >
          <span>{buttonText}</span>
        </Button>
        {#if cancelable}
          <Button size="xs" on:click={() => dispatch("cancel")} outline color="red">
            <span>Cancel</span>
          </Button>
        {/if}
      </div>
      {#if count > 9500}
        <Label class={count < 10000 ? "text-gray-600" : "font-bold text-red-600"}>
          {`${count}/10000`}
        </Label>
      {/if}
    </div>
    <ErrorMessage message={errorMessage}></ErrorMessage>
  </div>
</Textarea>
