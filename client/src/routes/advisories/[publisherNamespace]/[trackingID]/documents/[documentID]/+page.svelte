<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { page } from "$app/stores";
  import { Button, Drawer, Label, Textarea, Timeline, TimelineItem } from "flowbite-svelte";
  import { sineIn } from "svelte/easing";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";

  let document = {};
  let hideComments = false;

  let comments = [
    {
      author: "Beate Bearbeiterin",
      comment:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras in mi neque. Nam et pretium purus, vel condimentum magna. Vestibulum gravida, felis non efficitur imperdiet, arcu orci commodo ligula, vel pharetra tortor mi at felis. Proin eleifend dolor vitae lacinia luctus. Praesent sed justo quis eros convallis lacinia. Sed pharetra sollicitudin dui. Nam molestie convallis venenatis. Phasellus luctus felis at magna venenatis pellentesque. Integer mattis odio ac sapien pulvinar finibus. Pellentesque sit amet enim vitae ligula rutrum laoreet. Phasellus a placerat erat. Aliquam tortor eros, dignissim quis vulputate et, interdum vel tellus.",
      date: "2024-02-19"
    },
    {
      author: "Rene Reviewer",
      comment:
        "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Cras in mi neque. Nam et pretium purus, vel condimentum magna. Vestibulum gravida, felis non efficitur imperdiet, arcu orci commodo ligula, vel pharetra tortor mi at felis. Proin eleifend dolor vitae lacinia luctus. Praesent sed justo quis eros convallis lacinia. Sed pharetra sollicitudin dui. Nam molestie convallis venenatis. Phasellus luctus felis at magna venenatis pellentesque. Integer mattis odio ac sapien pulvinar finibus. Pellentesque sit amet enim vitae ligula rutrum laoreet. Phasellus a placerat erat. Aliquam tortor eros, dignissim quis vulputate et, interdum vel tellus.",
      date: "2024-02-20"
    }
  ];

  let transitionParams = {
    x: 320,
    duration: 120,
    easing: sineIn
  };

  function toggleComments() {
    hideComments = !hideComments;
  }
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      const response = await fetch(`/api/documents/${$page.params.documentID}`, {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      });
      if (response.ok) {
        ({ document } = await response.json());
        console.log(document);
      } else {
        // Do errorhandling
      }
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
      class="relative overflow-visible"
      placement="right"
      width="w-1/3"
      hidden={hideComments}
      transitionType="in:slide"
      {transitionParams}
    >
      <Timeline>
        {#each comments as comment}
          <TimelineItem date={comment.date} title={comment.author}>
            <span>Version 1.0</span>
            <p>
              {comment.comment}
            </p>
          </TimelineItem>
        {/each}
      </Timeline>
      <Label class="mb-2" for="comment-textarea">Kommentar</Label>
      <Textarea class="mb-2" id="comment-textarea"></Textarea>
      <Button>Send</Button>
    </Drawer>
  </div>
</RouteGuard>
