<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { browser } from "$app/environment";

  import { push } from "svelte-spa-router";

  let lastVisited = "/";
  if (browser) {
    lastVisited = localStorage.getItem("lastVisited") || "/";
  }
  onMount(async () => {
    await $appStore.app.keycloak
      .init({
        onLoad: "check-sso",
        checkLoginIframe: false,
        responseMode: "query"
      })
      .then(async (response: any) => {
        const profile = await $appStore.app.keycloak.loadUserProfile();
        appStore.setLoginState(true);
        appStore.setUserProfile({
          firstName: profile.firstName,
          lastName: profile.lastName
        });
        push(lastVisited);
      })
      .catch((error: any) => {
        console.log("error", error);
      });
  });
</script>
