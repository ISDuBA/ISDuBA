<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import Keycloak from "keycloak-js";
  import { Heading, Sidebar, SidebarWrapper, SidebarGroup, SidebarItem } from "flowbite-svelte";

  let keycloak;
  let token;
  let firstName;
  let lastName;

  onMount(async () => {
    keycloak = new Keycloak({
      url: "/auth/",
      realm: "isduba",
      clientId: "auth"
    });
    const authenticated = await keycloak
      .init({
        checkLoginIframe: false
        //onLoad: 'login-required'
      })
      .then(async (response) => {
        token = keycloak.token;
        const profile = await keycloak.loadUserProfile();
        console.log("Retrieved user profile:", profile);
        firstName = profile.firstName;
        lastName = profile.lastName;
      })
      .catch((error) => {
        console.log("error", error);
      });
  });

  async function logout() {
    keycloak.logout();
  }

  function login() {
    keycloak.login();
  }
</script>

<Sidebar class="bg-primary-700 h-screen p-2">
  <SidebarWrapper class="bg-primary-700">
    <Heading class="mb-6 text-white">ISDuBA</Heading>
    <SidebarGroup class="bg-primary-700">
      {#if token}
        <SidebarItem on:click={logout} label="Logout ({firstName} {lastName})"></SidebarItem>
      {:else}
        <SidebarItem on:click={login} label="Login"></SidebarItem>
      {/if}
    </SidebarGroup>
  </SidebarWrapper>
</Sidebar>
