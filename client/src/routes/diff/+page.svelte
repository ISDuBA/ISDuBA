<script lang="ts">
  import { onMount } from "svelte";
  import Diff from "$lib/Diff.svelte";

  let diff: string;
  onMount(async () => {
    fetch("advisory.diff").then((response) => {
      response.text().then((text) => {
        diff = text;
      });
    });
  });
</script>

<h1 class="text-lg">Comparison</h1>
<Diff {diff}></Diff>
