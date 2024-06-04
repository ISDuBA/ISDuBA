<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  /* eslint-disable svelte/no-at-html-tags */
  import { Label, TableBodyCell } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import CommentTextArea from "./CommentTextArea.svelte";
  import { request } from "$lib/utils";
  import { marked } from "marked";
  import DOMPurify from "dompurify";
  import { createEventDispatcher } from "svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import { tdClass } from "$lib/table/defaults";

  export let comment: any;
  let updatedComment = comment.message;
  let isEditing = false;
  let updateCommentError: string;

  const dispatch = createEventDispatcher();

  function toggleEditing() {
    isEditing = !isEditing;
  }

  async function updateComment() {
    updateCommentError = "";
    const formData = new FormData();
    formData.append("message", updatedComment);
    const response = await request(`/api/comments/${comment.comment_id}`, "PUT", formData);
    if (response.ok) {
      comment.message = updatedComment;
      toggleEditing();
    } else if (response.error) {
      updateCommentError = getErrorMessage(response.error);
    }
    dispatch("commentUpdate");
  }

  const parseMarkdown = (markdown: string) => {
    let html = marked.parse(markdown) as string;
    return DOMPurify.sanitize(html);
  };
</script>

<TableBodyCell {tdClass}>
  <div class="ml-1 flex flex-col">
    <small class="text-xs text-slate-400">{comment.time}</small>
    {#if !isEditing}
      <div class="flex flex-row items-center">
        <div class="display-markdown">
          {@html parseMarkdown(comment.message)}
        </div>
        <div>
          {#if $appStore.app.tokenParsed?.preferred_username === comment.actor}
            <button class="h-7 !p-2" on:click={toggleEditing}>
              <i class="bx bx-edit text-lg"></i>
            </button>
          {/if}
        </div>
      </div>
      <small>({comment.actor})</small>
      <Label class="text-xs text-gray-400">Document-Version: {comment.documentVersion}</Label>
    {:else}
      <CommentTextArea
        on:cancel={toggleEditing}
        on:input={() => (updateCommentError = "")}
        on:saveComment={updateComment}
        cancelable={true}
        buttonText="Save"
        errorMessage={updateCommentError}
        bind:value={updatedComment}
      ></CommentTextArea>
    {/if}
  </div>
</TableBodyCell>

<style>
  /* Reset styles inside markdown block */
  .display-markdown :global(a) {
    text-decoration: underline;
  }
  .display-markdown :global(ol) {
    display: block;
    list-style-type: decimal;
    margin-block-start: 1em;
    margin-block-end: 1em;
    padding-inline-start: 40px;
  }
  .display-markdown :global(ul) {
    display: block;
    list-style-type: disc;
    margin-block-start: 1em;
    margin-block-end: 1em;
    padding-inline-start: 40px;
  }
  .display-markdown :global(blockquote) {
    display: block;
    margin-block: 1em;
    margin-inline: 40px;
  }
  .display-markdown :global(table) {
    border: 1px solid;
  }
  .display-markdown :global(th) {
    border: 1px solid;
  }
  .display-markdown :global(td) {
    border: 1px solid;
  }
  .display-markdown :global(h1) {
    display: block;
    font-size: 2em;
    font-weight: bold;
    margin-block: 0.67em;
  }
  .display-markdown :global(h2) {
    display: block;
    font-size: 1.5em;
    font-weight: bold;
    margin-block: 0.83em;
  }
  .display-markdown :global(h3) {
    display: block;
    font-size: 1.17em;
    font-weight: bold;
    margin-block: 1em;
  }
  .display-markdown :global(h4) {
    display: block;
    font-size: 1em;
    font-weight: bold;
    margin-block: 1.33em;
  }

  .display-markdown :global(h5) {
    display: block;
    font-size: 0.83em;
    font-weight: bold;
    margin-block: 1.67em;
  }
  .display-markdown :global(h6) {
    display: block;
    font-size: 0.67em;
    font-weight: bold;
    margin-block: 2.33em;
  }
</style>
