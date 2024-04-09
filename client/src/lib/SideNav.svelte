<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Heading, Sidebar, SidebarWrapper, SidebarGroup, SidebarItem } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import { page } from "$app/stores";

  async function logout() {
    $appStore.app.keycloak.logout();
  }

  function login() {
    $appStore.app.keycloak.login();
  }
  $: activeUrl = "/" + $page.url.hash;
  let activeClass =
    "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-primary-700 dark:text-white hover:bg-primary-100 dark:hover:bg-gray-700";
  let nonActiveClass =
    "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 hover:text-primary-900";
</script>

<Sidebar class="bg-primary-700 h-screen p-2" {activeUrl} {activeClass} {nonActiveClass}>
  <SidebarWrapper class="bg-primary-700">
    <Heading class="mb-6 text-white">ISDuBA</Heading>
    <SidebarGroup class="bg-primary-700">
      {#if $appStore.app.keycloak.authenticated}
        <!-- Entries which are available after login should go here-->
        <SidebarItem label="Home" href="/#/">
          <svelte:fragment slot="icon">
            <i class="bx bxs-dashboard"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Advisories" href="/#/advisories">
          <svelte:fragment slot="icon">
            <i class="bx bx-spreadsheet"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Compare" href="/#/diff">
          <svelte:fragment slot="icon">
            <i class="bx bx-transfer"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Documents" href="/#/documents">
          <svelte:fragment slot="icon">
            <i class="bx bx-spreadsheet"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Sources">
          <svelte:fragment slot="icon">
            <i class="bx bx-git-repo-forked"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Statistics">
          <svelte:fragment slot="icon">
            <i class="bx bx-bar-chart-square"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem label="Configuration">
          <svelte:fragment slot="icon">
            <i class="bx bx-cog"></i>
          </svelte:fragment>
        </SidebarItem>
        <SidebarItem
          label="Profile"
          target="_blank"
          href="http://localhost:8080/realms/isduba/account/#/"
        >
          <svelte:fragment slot="icon">
            <i class="bx bx-user"></i>
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
        <!-- Entries which should be available only if not logged in should go here-->
        <SidebarItem on:click={login} label="Login">
          <svelte:fragment slot="icon">
            <i class="bx bx-log-in"></i>
          </svelte:fragment>
        </SidebarItem>
      {/if}
      <SidebarItem label="About" href="/#/about"></SidebarItem>
    </SidebarGroup>
  </SidebarWrapper>
</Sidebar>
