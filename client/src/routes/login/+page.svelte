<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { goto } from "$app/navigation";
  import { configuration } from "$lib/configuration";
  import Keycloak from "keycloak-js";
  import SideNav from "$lib/SideNav.svelte";

  onMount(async () => {
    appStore.setKeycloak(new Keycloak(configuration.getConfiguration()));
    await $appStore.app.keycloak
      .init({
        checkLoginIframe: false
        //onLoad: 'login-required'
      })
      .then(async (response: any) => {
        const profile = await $appStore.app.keycloak.loadUserProfile();
        console.log("Retrieved user profile:", profile);
        appStore.setLoginState(true);
        appStore.setUserProfile({
          firstName: profile.firstName,
          lastName: profile.lastName
        });
        goto("/");
      })
      .catch((error: any) => {
        console.log("error", error);
      });
  });
</script>

<SideNav></SideNav>
