<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { configuration } from "$lib/configuration";
  import Keycloak from "keycloak-js";
  import { Heading, Sidebar, SidebarWrapper, SidebarGroup, SidebarItem } from "flowbite-svelte";

  let keycloak: any;
  let token: string;
  let firstName: string;
  let lastName: string;

  onMount(async () => {
    keycloak = new Keycloak(configuration.getConfiguration());
    const authenticated = await keycloak
      .init({
        checkLoginIframe: false
        //onLoad: 'login-required'
      })
      .then(async (response: any) => {
        token = keycloak.token;
        const profile = await keycloak.loadUserProfile();
        console.log("Retrieved user profile:", profile);
        firstName = profile.firstName;
        lastName = profile.lastName;
      })
      .catch((error: any) => {
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
        <SidebarItem on:click={logout} label="Logout ({firstName} {lastName})">
          <svelte:fragment slot="icon">
            <i class="bx bx-power-off"></i>
          </svelte:fragment>
        </SidebarItem>
      {:else}
        <SidebarItem on:click={login} label="Login">
          <svelte:fragment slot="icon">
            <i class="bx bx-log-in"></i>
          </svelte:fragment>
        </SidebarItem>
      {/if}
    </SidebarGroup>
  </SidebarWrapper>
</Sidebar>
