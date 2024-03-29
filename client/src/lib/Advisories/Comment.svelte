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

  export let comment: any;
  let updatedComment = "";
  let isEditing = false;

  onMount(() => {
    updatedComment = comment.message;
  });

  function toggleEditing() {
    isEditing = !isEditing;
  }

  function updateComment() {
    const formData = new FormData();
    formData.append("message", updatedComment);
    $appStore.app.keycloak.updateToken(5).then(async () => {
      fetch(`/api/comments/${comment.id}`, {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        },
        method: "PUT",
        body: formData
      }).then((response) => {
        if (response.ok) {
          comment.message = updatedComment;
          toggleEditing();
          appStore.displaySuccessMessage("Comment updated.");
        } else {
          appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
        }
      });
    });
  }
</script>

<TimelineItem
  classLi="mb-4 ms-4"
  date={`${new Date(comment.time).toLocaleDateString("en-US")} - ${new Date(comment.time).toLocaleTimeString("en-US")}`}
>
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
      on:saveComment={updateComment}
      cancelable={true}
      buttonText="Save"
      bind:value={updatedComment}
    ></CommentTextArea>
  {/if}
</TimelineItem>
