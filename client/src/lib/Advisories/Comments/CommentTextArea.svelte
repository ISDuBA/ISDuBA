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
  import { canSetStateReview, isRoleIncluded } from "$lib/permissions";
  import { ARCHIVED, EDITOR, REVIEW, REVIEWER } from "$lib/workflow";
  import { appStore } from "$lib/store";
  import type { ErrorDetails } from "$lib/Errors/error";

  export let value: string;
  export let buttonText: string;
  export let cancelable = false;
  export let errorMessage: ErrorDetails | null;
  export let state: string = "";
  export let old = "";
  $: count = value.length;

  const dispatch = createEventDispatcher();
</script>

<Textarea
  bind:value
  on:focus={() => dispatch("focus")}
  on:blur={() => dispatch("blur")}
  on:input={() => dispatch("input")}
  class="mb-2"
  id="comment-textarea"
>
  <div slot="footer" class="flex flex-col">
    <div class="flex items-start justify-between">
      <div>
        {#if !cancelable}
          {#if canSetStateReview(state) && state !== ARCHIVED}
            <Button
              size="xs"
              on:click={() => dispatch("saveForReview")}
              disabled={count === 0}
              color="light"
            >
              {#if count > 0}
                <span>Release for review with comment</span>
              {:else}
                <span>Release for review</span>
              {/if}
            </Button>
          {/if}
          {#if (state === REVIEW && isRoleIncluded( appStore.getRoles(), [REVIEWER] )) || ((state === REVIEW || state === ARCHIVED) && isRoleIncluded( appStore.getRoles(), [EDITOR] ))}
            <Button
              size="xs"
              on:click={() => dispatch("saveForAssessing")}
              disabled={count === 0}
              color="light"
            >
              <span>Send back to assessing</span>
            </Button>
          {/if}
          {#if state === ARCHIVED && isRoleIncluded(appStore.getRoles(), [EDITOR])}
            <Button
              size="xs"
              on:click={() => dispatch("saveForReview")}
              disabled={count === 0}
              color="light"
            >
              <span>Send back to review</span>
            </Button>
          {/if}
        {/if}
        {#if !((state === REVIEW || state === ARCHIVED) && isRoleIncluded( appStore.getRoles(), [EDITOR] ))}
          <Button
            size="xs"
            on:click={() => dispatch("saveComment")}
            disabled={count > 10000 || count === 0 || value === old}
          >
            <span>{buttonText}</span>
          </Button>
        {/if}
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
    <ErrorMessage error={errorMessage}></ErrorMessage>
  </div>
</Textarea>
