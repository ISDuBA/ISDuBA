<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, ButtonGroup, Label, P, TimelineItem } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import CommentTextArea from "./CommentTextArea.svelte";
  import { request } from "$lib/utils";

  export let comment: any;
  let updatedComment = comment.message;
  let isEditing = false;
  let updateCommentError: string;

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
  }
</script>

<TimelineItem classLi="mb-4 ms-4" date={`${new Date(comment.time).toISOString()}`}>
  {#if !isEditing}
    <P class="mb-2">
      {comment.message}
      <small>({comment.commentator})</small>
      {#if $appStore.app.keycloak.tokenParsed.preferred_username === comment.commentator}
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
