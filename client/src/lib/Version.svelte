<script lang="ts">
  import { Timeline, TimelineItem } from "flowbite-svelte";
  export let advisoryVersions: any;
  export let publisherNamespace: string;
  export let trackingID: string;

  let staticClass =
    "border-primary-800 border-solid border absolute -start-3 flex h-6 w-6 items-center justify-center rounded-full ring-8 ring-white";
  let selectedClass = `bg-primary-800 text-white cursor-default ${staticClass}`;
  let unselectedClass = `bg-white hover:bg-primary-600 hover:border-primary-600 hover:text-white ${staticClass}`;
</script>

<div>
  <div class="flex">
    <h2 class="mr-3">Advisory Versions</h2>
    <Timeline class="flex flex-col-reverse border-gray-400" order="vertical">
      {#each advisoryVersions as version}
        <TimelineItem classLi="mb-14">
          <svelte:fragment slot="icon">
            <a
              class={false ? selectedClass : unselectedClass}
              href={`/advisories/${publisherNamespace}/${trackingID}/documents/${version.id}`}
              >{version.version}
            </a>
          </svelte:fragment>
        </TimelineItem>
      {/each}
    </Timeline>
  </div>
</div>
