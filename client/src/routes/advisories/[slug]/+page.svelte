<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import type { PageData } from "./$types";
  import { Button, Drawer, Label, Textarea, Timeline, TimelineItem } from "flowbite-svelte";

  export let data: PageData;
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

  function toggleComments() {
    hideComments = !hideComments;
  }
</script>

<RouteGuard>
  <div class="flex">
    <div class="grow">
      <div>Slug ID: {data.id}</div>
    </div>
    <Drawer
      activateClickOutside={false}
      backdrop={false}
      class="relative overflow-visible transition-all duration-100 ease-linear motion-reduce:transition-none"
      width={hideComments ? "w-0" : "w-1/3"}
      placement="right"
      hidden={false}
    >
      <Button
        on:click={toggleComments}
        outline={true}
        class="absolute -left-8 top-8 z-10 !p-2"
        size="lg"
      >
        <i class={hideComments ? "bx bx-chevron-left" : "bx bx-chevron-right"}></i>
      </Button>
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
