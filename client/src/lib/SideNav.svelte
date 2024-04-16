<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    Drawer,
    Heading,
    Sidebar,
    SidebarWrapper,
    SidebarGroup,
    SidebarItem
  } from "flowbite-svelte";
  import { sineIn } from "svelte/easing";
  import { appStore } from "$lib/store";
  import { page } from "$app/stores";
  import "boxicons";
  import { PUBLIC_KEYCLOAK_URL } from "$env/static/public";
  async function logout() {
    $appStore.app.keycloak.logout();
  }

  let notactivated =
    "flex items-center p-2 text-base font-normal text-gray-400 dark:text-gray-400 hover:bg-primary-100 hover:text-primary-900";

  function login() {
    $appStore.app.keycloak.login();
  }
  $: activeUrl = "/" + $page.url.hash;
  let activeClass =
    "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-primary-700 dark:text-white hover:bg-primary-100 dark:hover:bg-gray-700";
  let nonActiveClass =
    "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 hover:text-primary-900";
  let profileUrl = PUBLIC_KEYCLOAK_URL + "/realms/isduba/account/#/";

  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn
  };
  let breakPoint: number = 1280;
  let width: number;
  let drawerHidden: boolean = false;
  $: if (width >= breakPoint) {
    drawerHidden = false;
  } else {
    drawerHidden = true;
  }
  const toggleSide = () => {
    if (width < breakPoint) {
      drawerHidden = !drawerHidden;
    }
  };
  const toggleDrawer = () => {
    drawerHidden = !drawerHidden;
  };
  $: activeUrl = $page.url.pathname;
</script>

<svelte:window bind:innerWidth={width} />

<div class="flex">
  <Drawer
    transitionType="fly"
    {transitionParams}
    bind:hidden={drawerHidden}
    activateClickOutside={false}
    width="w-45"
    backdrop={false}
    class="static h-screen bg-primary-700 p-2"
    id="sidebar"
  >
    <Sidebar {activeUrl} {activeClass} {nonActiveClass}>
      <SidebarWrapper class="bg-primary-700">
        <Heading class="mb-6 text-white">ISDuBA</Heading>
        <SidebarGroup class="bg-primary-700">
          {#if $appStore.app.keycloak.authenticated}
            <!-- Entries which are available after login should go here-->
            <SidebarItem label="Home" href="/#/" on:click={toggleSide}>
              <svelte:fragment slot="icon">
                <i class="bx bxs-dashboard"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Advisories" href="/#/advisories" on:click={toggleSide}>
              <svelte:fragment slot="icon">
                <i class="bx bx-spreadsheet"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Compare" href="/#/diff" on:click={toggleSide}>
              <svelte:fragment slot="icon">
                <i class="bx bx-transfer"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Documents" href="/#/documents" on:click={toggleSide}>
              <svelte:fragment slot="icon">
                <i class="bx bx-spreadsheet"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Sources" nonActiveClass={notactivated}>
              <svelte:fragment slot="icon">
                <i class="bx bx-git-repo-forked"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Statistics" nonActiveClass={notactivated}>
              <svelte:fragment slot="icon">
                <i class="bx bx-bar-chart-square"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Configuration" nonActiveClass={notactivated}>
              <svelte:fragment slot="icon">
                <i class="bx bx-cog"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem label="Profile" target="_blank" href={profileUrl}>
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
          <SidebarItem label="About" href="/#/about" on:click={toggleSide}></SidebarItem>
        </SidebarGroup>
      </SidebarWrapper>
    </Sidebar>
  </Drawer>
  <div class="h-screen bg-white p-2">
    <button on:click={toggleDrawer}>
      <box-icon name="menu" color="black" size="lg">Toggle menu</box-icon>
    </button>
  </div>
</div>
