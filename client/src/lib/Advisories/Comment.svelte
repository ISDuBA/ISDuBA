<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  /* eslint-disable svelte/no-at-html-tags */
  import { Button, ButtonGroup, Label, P, TimelineItem } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import CommentTextArea from "./CommentTextArea.svelte";
  import { request } from "$lib/utils";
  import { marked } from "marked";
  import DOMPurify from "dompurify";
  import { createEventDispatcher } from "svelte";

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
    const response = await request(`/api/comments/${comment.id}`, "PUT", formData);
    if (response.ok) {
      comment.message = updatedComment;
      toggleEditing();
      appStore.displaySuccessMessage("Comment updated.");
    } else if (response.error) {
      updateCommentError = response.error;
    }
    dispatch("commentUpdate");
  }

  const parseMarkdown = (markdown: string) => {
    let html = marked.parse(markdown) as string;
    return DOMPurify.sanitize(html);
  };
</script>

<TimelineItem classLi="mb-4 ms-4" date={`${new Date(comment.time).toISOString()}`}>
  {#if !isEditing}
    <P class="mb-2">
      <div class="display-markdown">
        {@html parseMarkdown(comment.message)}
      </div>
      <small>({comment.commentator})</small>
      {#if $appStore.app.tokenParsed?.preferred_username === comment.commentator}
        <ButtonGroup>
          <Button class="!p-2" on:click={toggleEditing}>
            <i class="bx bx-edit text-lg"></i>
          </Button>
        </ButtonGroup>
      {/if}
    </P>
    <Label class="text-xs text-slate-400">Document-Version: {comment.documentVersion}</Label>
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
</TimelineItem>

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
