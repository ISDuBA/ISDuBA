<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { browser } from "$app/environment";
  export let loginRequired: boolean = true;
  export let roles: any = [];

  $: activeUrl = $page.url.pathname;

  onMount(() => {
    if (browser) {
      localStorage.setItem("lastVisited", activeUrl);
    }
    if (loginRequired && $appStore.app.isUserLoggedIn === false) {
      goto("/login");
    }
  });
</script>

{#if loginRequired && $appStore.app.isUserLoggedIn}
  <slot />
{/if}
