<script lang="ts">
  import { Timeline, TimelineItem } from "flowbite-svelte";
  export let advisoryVersions: any;
  export let publisherNamespace: string;
  export let trackingID: string;
</script>

<div>
  <div class="flex">
    <h2 class="mr-3">Advisory Versions</h2>
    <Timeline class="flex flex-col-reverse">
      {#each advisoryVersions as version}
        <TimelineItem
          ><a href={`/advisories/${publisherNamespace}/${trackingID}/documents/${version.id}`}
            >{version.version}</a
          ></TimelineItem
        >
      {/each}
    </Timeline>
  </div>
</div>
