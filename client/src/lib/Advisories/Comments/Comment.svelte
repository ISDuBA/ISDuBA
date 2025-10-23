<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell } from "flowbite-svelte";
  import { appStore } from "$lib/store.svelte";
  import CommentTextArea from "./CommentTextArea.svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { ARCHIVED, ASSESSING, NEW, READ, REVIEW } from "$lib/workflow";
  import { getReadableDateString } from "../CSAFWebview/helpers";
  import { untrack } from "svelte";

  interface Props {
    comment: any;
    fullHistory: boolean;
    workflowState: string;
    onCommentUpdated: () => void;
  }

  let { comment, fullHistory, workflowState = "", onCommentUpdated }: Props = $props();
  const intlFormat = new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "medium"
  });
  let updatedComment = $state(comment.message);
  let commentID = $state(comment.comment_id);
  $effect(() => {
    untrack(() => commentID);
    if (commentID !== comment.comment_id) {
      updatedComment = comment.message;
      commentID = comment.comment_id;
    }
  });
  let isEditing = $state(false);
  let updateCommentError: ErrorDetails | null = $state(null);
  let lastEdited = $derived.by(() => {
    if (comment.times) {
      let latest = [...comment.times].sort().reverse()[0];
      latest = latest.replace("T", " ").split(".")[0];
      return `(edited ${latest})`;
    }
  });

  let isCommentingAllowed = $derived.by(() => {
    if ([NEW, READ, ASSESSING, REVIEW, ARCHIVED].includes(workflowState)) {
      if (appStore.isReviewer() && [NEW, READ, ARCHIVED].includes(workflowState)) {
        return false;
      } else {
        return appStore.isEditor() || appStore.isReviewer();
      }
    } else {
      return false;
    }
  });

  const tdClass = "py-2 px-2";

  function toggleEditing() {
    isEditing = !isEditing;
  }

  async function updateComment() {
    updateCommentError = null;
    const formData = new FormData();
    formData.append("message", updatedComment);
    const response = await request(`/api/comments/post/${comment.comment_id}`, "PUT", formData);
    if (response.ok) {
      comment.message = updatedComment;
      toggleEditing();
    } else if (response.error) {
      updateCommentError = getErrorDetails(`Could not update comment.`, response);
    }
    onCommentUpdated();
  }
</script>

<TableBodyCell class={tdClass}>
  <div class="flex flex-col">
    <div class="flex flex-row items-baseline">
      <small class="w-40 text-xs text-slate-400" title={comment.time}
        >{getReadableDateString(comment.time, intlFormat)}</small
      >
      <small class="ml-1 flex-grow"
        >{fullHistory ? `Comment (${comment.actor})` : `${comment.actor}`}
      </small>
      <small class="ml-1 text-xs text-slate-400">on version: {comment.documentVersion}</small>
    </div>
    {#if !isEditing}
      <div class="mt-1 flex flex-row items-center">
        <div style="white-space: pre-wrap">{comment.message}</div>
        <div class="ml-auto">
          {#if appStore.state.app.tokenParsed?.preferred_username === comment.actor && isCommentingAllowed}
            <button class="h-7 !p-2" onclick={toggleEditing} aria-label="Edit comment">
              <i class="bx bx-edit text-lg"></i>
            </button>
          {/if}
        </div>
      </div>
    {:else}
      <CommentTextArea
        cancel={toggleEditing}
        onInput={() => (updateCommentError = null)}
        saveComment={updateComment}
        cancelable={true}
        buttonText="Save"
        errorMessage={updateCommentError}
        bind:value={updatedComment}
        old={comment.message}
      ></CommentTextArea>
    {/if}
    <div class="mt-1">
      <small class="text-gray-400">{lastEdited}</small>
    </div>
  </div>
</TableBodyCell>
