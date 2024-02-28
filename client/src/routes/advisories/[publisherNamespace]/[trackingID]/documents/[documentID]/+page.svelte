<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { page } from "$app/stores";
  import { Button, Drawer, Label, Textarea, Timeline } from "flowbite-svelte";
 SPDX-License-Identifier: Apache-2.0lte/easing";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Comment from "$lib/Comment.svelte";
  import Version from "$lib/Version.svelte";

  let document = {};
  let hideComments = false;
  let comment: string = "";
  let comments: any = [];
  let advisoryVersions: string[] = [];

  let transitionParams = {
    x: 320,
    duration: 120,
    easing: sineIn
  };

  const loadAdvisoryVersions = async () => {
    const response = await fetch(
      `/api/documents?&columns=id version&query=$tracking_id ${$page.params.trackingID} = $publisher "${$page.params.publisherNamespace}" = and`,
      {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      }
    );
    if (response.ok) {
      const result = await response.json();
      advisoryVersions = result.documents.map((doc: any) => {
        return { id: doc.id, version: doc.version };
      });
    } else {
      // Do errorhandling
    }
  };

  const loadDocument = async () => {
    const response = await fetch(`/api/documents/${$page.params.documentID}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      }
    });
    if (response.ok) {
      ({ document } = await response.json());
    } else {
      // Do errorhandling
    }
  };

  function toggleComments() {
    hideComments = !hideComments;
  }
  function loadComments() {
    fetch(`/api/comments/${$page.params.documentID}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      }
    }).then((response) => {
      if (response.ok) {
        response.json().then((json) => {
          json.forEach((c: any) => {
            c.documentID = $page.params.documentID;
          });
          comments = json;
        });
      } else {
        // Do errorhandling
      }
    });
  }
  function createComment() {
    const formData = new FormData();
    formData.append("message", comment);
    fetch(`/api/comments/${$page.params.documentID}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      },
      method: "POST",
      body: formData
    }).then((response) => {
      if (response.ok) {
        comment = "";
        loadComments();
      } else {
        // Do errorhandling
      }
    });
  }

  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      loadDocument();
      loadComments();
      loadAdvisoryVersions();
    }
  });
</script>

<RouteGuard>
  <div class="flex">
    <div class="grow">
      <table>
        <tr>
          <td>PublisherNamespace:</td><td class="pl-3">{$page.params.publisherNamespace}</td>
        </tr>
        <tr>
          <td>TrackingId:</td><td class="pl-3">{$page.params.trackingID}</td>
        </tr>
        <tr>
          <td>DocumentID:</td><td class="pl-3">{$page.params.documentID}</td>
        </tr>
        {#if document}
          <tr>
            <td>Current release date:</td><td class="pl-3"
              >{document.tracking?.current_release_date}</td
            >
          </tr>
        {/if}
      </table>
    </div>
    <Version
      publisherNamespace={$page.params.publisherNamespace}
      trackingID={$page.params.trackingID}
      {advisoryVersions}
    ></Version>
    <Button
      on:click={toggleComments}
      outline={true}
      class="absolute right-2 top-2 z-10 !p-2"
      size="lg"
    >
      <i class={hideComments ? "bx bx-chevron-left" : "bx bx-chevron-right"}></i>
    </Button>
    <Drawer
      activateClickOutside={false}
      backdrop={false}
      class="relative flex flex-col"
      placement="right"
      width="w-1/3"
      hidden={hideComments}
      transitionType="in:slide"
      {transitionParams}
    >
      <div class="overflow-y-scroll pl-2">
        <Timeline class="flex flex-col-reverse">
          {#each comments as comment}
            <Comment {comment}></Comment>
          {/each}
        </Timeline>
      </div>
      <div>
        <Label class="mb-2" for="comment-textarea">Comment:</Label>
        <Textarea bind:value={comment} class="mb-2" id="comment-textarea"></Textarea>
        <Button on:click={createComment}>Send</Button>
      </div>
    </Drawer>
  </div>
</RouteGuard>
