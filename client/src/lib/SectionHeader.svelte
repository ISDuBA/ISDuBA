<script>
  import { Heading } from "flowbite-svelte";

  export let title = "";
</script>

<Heading tag="h2" customSize="text-lg" class="mb-3 text-lg text-gray-900 dark:text-white"
  >{title}</Heading
>
