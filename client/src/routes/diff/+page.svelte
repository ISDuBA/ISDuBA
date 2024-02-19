<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Diff from "$lib/Diff.svelte";

  let diff: string;
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      fetch("advisory.diff").then((response) => {
        response.text().then((text) => {
          diff = text;
        });
      });
    }
  });
</script>

<RouteGuard>
  <h1 class="text-lg">Comparison</h1>
  <Diff {diff}></Diff>
</RouteGuard>
