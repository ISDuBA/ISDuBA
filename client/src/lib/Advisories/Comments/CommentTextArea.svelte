<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Label, Textarea } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { canSetStateReview, isRoleIncluded } from "$lib/permissions";
  import { ARCHIVED, EDITOR, REVIEW, REVIEWER } from "$lib/workflow";
  import { appStore } from "$lib/store.svelte";
  import type { ErrorDetails } from "$lib/Errors/error";

  interface Props {
    value: string;
    buttonText: string;
    cancelable?: boolean;
    errorMessage: ErrorDetails | null;
    workflowState?: string;
    old?: string;
    onFocus?: (event: any) => void;
    onBlur?: (event: any) => void;
    onInput?: (event: any) => void;
    saveForReview?: () => void;
    saveForAssessing?: () => void;
    saveComment: () => void;
    cancel?: () => void;
  }

  let {
    value = $bindable(),
    buttonText,
    cancelable = false,
    errorMessage,
    workflowState = "",
    old = "",
    onFocus = undefined,
    onBlur = undefined,
    onInput = undefined,
    saveForReview = undefined,
    saveForAssessing = undefined,
    saveComment,
    cancel = undefined
  }: Props = $props();

  let count = $derived(value.length);
</script>

<Textarea
  bind:value
  on:focus={(event) => {
    if (onFocus) {
      onFocus(event);
    }
  }}
  on:blur={(event) => {
    if (onBlur) {
      onBlur(event);
    }
  }}
  on:input={(event) => {
    if (onInput) {
      onInput(event);
    }
  }}
  class="mb-2"
  id="comment-textarea"
>
  <div slot="footer" class="flex flex-col">
    <div class="flex items-start justify-between">
      <div>
        {#if !cancelable}
          {#if canSetStateReview(workflowState) && workflowState !== ARCHIVED && saveForReview}
            <Button size="xs" on:click={() => saveForReview()} disabled={count === 0} color="light">
              {#if count > 0}
                <span>Release for review with comment</span>
              {:else}
                <span>Release for review</span>
              {/if}
            </Button>
          {/if}
          {#if (workflowState === REVIEW && isRoleIncluded( appStore.getRoles(), [REVIEWER] )) || ((workflowState === REVIEW || workflowState === ARCHIVED) && isRoleIncluded( appStore.getRoles(), [EDITOR] ))}
            <Button
              size="xs"
              on:click={() => {
                if (saveForAssessing) saveForAssessing();
              }}
              disabled={count === 0}
              color="light"
            >
              <span>Send back to assessing</span>
            </Button>
          {/if}
          {#if workflowState === ARCHIVED && isRoleIncluded( appStore.getRoles(), [EDITOR] ) && saveForReview}
            <Button size="xs" on:click={() => saveForReview()} disabled={count === 0} color="light">
              <span>Send back to review</span>
            </Button>
          {/if}
        {/if}
        {#if !((workflowState === REVIEW || workflowState === ARCHIVED) && isRoleIncluded( appStore.getRoles(), [EDITOR] ))}
          <Button
            color="green"
            size="xs"
            on:click={() => saveComment()}
            disabled={count > 10000 || count === 0 || value === old}
          >
            <span>{buttonText}</span>
          </Button>
        {/if}
        {#if cancelable && cancel}
          <Button size="xs" on:click={() => cancel()} outline color="red">
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
