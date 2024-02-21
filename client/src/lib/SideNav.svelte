<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Heading, Sidebar, SidebarWrapper, SidebarGroup, SidebarItem } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import { page } from "$app/stores";

  async function logout() {
    appStore.setLoginState(false);
    $appStore.app.keycloak.logout();
  }

  function login() {
    $appStore.app.keycloak.login();
  }
  $: activeUrl = $page.url.pathname;
  let activeClass =
    "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-primary-700 dark:text-white hover:bg-primary-100 dark:hover:bg-gray-700";
  let nonActiveClass =
    "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 hover:text-primary-900";
</script>

<Sidebar class="h-screen bg-primary-700 p-2" {activeUrl} {activeClass} {nonActiveClass}>
  <SidebarWrapper class="bg-primary-700">
    <Heading class="mb-6 text-white">ISDuBA</Heading>
    <SidebarGroup class="bg-primary-700">
      {#if $appStore.app.isUserLoggedIn}
        <!-- Entries which are available after login should go here-->
        <SidebarItem label="Home" href="/">
          <svelte:fragment slot="icon">
            <i class="bx bxs-dashboard"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Advisories" href="/advisories">
          <svelte:fragment slot="icon">
            <i class="bx bx-spreadsheet"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Single Advisory" href="/advisories/1">
          <svelte:fragment slot="icon">
            <i class="bx bx-spreadsheet"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Compare" href="/diff">
          <svelte:fragment slot="icon">
            <i class="bx bx-transfer"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem
          on:click={logout}
          label="Logout ({$appStore.app.userProfile.firstName} {$appStore.app.userProfile
            .lastName})"
        >
          <svelte:fragment slot="icon">
            <i class="bx bx-power-off"></i>
          </svelte:fragment>
        </SidebarItem>
      {:else}
        <!-- Entries which should be available independend from login status should go here-->
        <SidebarItem on:click={login} label="Login">
          <svelte:fragment slot="icon">
            <i class="bx bx-log-in"></i>
          </svelte:fragment>
        </SidebarItem>
      {/if}
    </SidebarGroup>
  </SidebarWrapper>
</Sidebar>
