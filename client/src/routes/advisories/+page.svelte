<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      const response = await fetch("/api/documents", {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      });
      console.log(response);
    }
  });
</script>

<RouteGuard>
  <h1 class="text-lg">Advisories</h1>
  {$appStore.app.keycloak.token}
</RouteGuard>
