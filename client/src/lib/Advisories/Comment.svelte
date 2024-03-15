<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, ButtonGroup, Label, P, Textarea, TimelineItem } from "flowbite-svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";

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
      } else {
        // Do errorhandling
      }
    });
  }
</script>

<TimelineItem
  date={`${new Date(comment.time).toLocaleDateString("en-US")} - ${new Date(comment.time).toLocaleTimeString("en-US")}`}
  title={comment.commentator}
>
  <Label class="text-xs text-slate-400">Document-Version: {comment.documentID}</Label>
  {#if !isEditing}
    <P class="mb-2">
      {comment.message}
    </P>
    {#if $appStore.app.keycloak.tokenParsed.preferred_username === comment.commentator}
      <ButtonGroup>
        <Button class="!p-2" on:click={toggleEditing}>
          <i class="bx bx-edit text-lg"></i>
        </Button>
      </ButtonGroup>
    {/if}
  {:else}
    <Textarea bind:value={updatedComment}></Textarea>
    <Button color="red" outline={true} class="!p-2" on:click={toggleEditing}>
      <i class="bx bx-x text-lg"></i>
    </Button>
    <Button color="green" class="!p-2" on:click={updateComment}>
      <i class="bx bx-check text-lg"></i>
    </Button>
  {/if}
</TimelineItem>
