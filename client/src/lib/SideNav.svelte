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
    Sidebar,
    SidebarWrapper,
    SidebarGroup,
    SidebarItem,
    SidebarBrand
  } from "flowbite-svelte";
  import { sineIn } from "svelte/easing";
  import { appStore } from "$lib/store";
  import { page } from "$app/stores";

  let notactivated =
    "flex items-center p-2 text-base font-normal text-gray-400 dark:text-gray-400 hover:bg-primary-100 hover:text-primary-900";

  $: activeUrl = "/" + $page.url.hash;

  let activeClass =
    "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-primary-700 dark:text-white hover:bg-primary-100 dark:hover:bg-gray-700";
  let nonActiveClass =
    "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 hover:text-primary-900";

  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn
  };
  let breakPoint: number = 1280;
  let width: number;
  let drawerHidden: boolean = false;
  $: drawerHidden = width < breakPoint;

  const toggleDrawer = () => {
    drawerHidden = !drawerHidden;
  };
</script>

<svelte:window bind:innerWidth={width} />
{#if $appStore.app.userManager && ($appStore.app.isUserLoggedIn || $appStore.app.sessionExpired)}
  <div class="flex">
    <Drawer
      transitionType="fly"
      {transitionParams}
      bind:hidden={drawerHidden}
      activateClickOutside={false}
      width="fit-content"
      backdrop={false}
      class="static min-w-fit bg-primary-700 p-0"
      id="sidebar"
    >
      <Sidebar
        asideClass="w-fit max-w-60"
        class="bg-primary-700"
        {activeUrl}
        {activeClass}
        {nonActiveClass}
      >
        <SidebarWrapper class="bg-primary-700 px-0">
          <SidebarGroup>
            <SidebarBrand
              spanClass="self-center text-4xl font-normal whitespace-nowrap text-white me-4"
              site={{ img: "favicon.svg", name: "ISDuBA", href: "/#/" }}
            ></SidebarBrand>
          </SidebarGroup>
          <SidebarGroup class="space-y-0 bg-primary-700">
            <!-- Entries which are available after login should go here-->
            <SidebarItem class="px-6 py-2.5" label="Dashboard" href="/#/">
              <svelte:fragment slot="icon">
                <i class="bx bxs-dashboard"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem class="px-6 py-2.5" label="Search" href="/#/search">
              <svelte:fragment slot="icon">
                <i class="bx bx-spreadsheet"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem class="px-6 py-2.5" label="Sources" href="/#/sources">
              <svelte:fragment slot="icon">
                <i class="bx bx-git-repo-forked"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem
              class="px-6 py-2.5"
              label="Statistics"
              href="javascript: void(0)"
              nonActiveClass={notactivated}
            >
              <svelte:fragment slot="icon">
                <i class="bx bx-bar-chart-square"></i>
              </svelte:fragment>
            </SidebarItem>
            <SidebarItem
              class="px-6 py-2.5"
              label="Configuration"
              href="javascript: void(0)"
              nonActiveClass={notactivated}
            >
              <svelte:fragment slot="icon">
                <i class="bx bx-cog"></i>
              </svelte:fragment>
            </SidebarItem>
            {#if !$appStore.app.sessionExpired}
              <SidebarItem
                class="px-6 py-2.5"
                label={$appStore.app.tokenParsed?.preferred_username}
                href="/#/login"
              >
                <svelte:fragment slot="icon">
                  <i class="bx bx-user"></i>
                </svelte:fragment>
              </SidebarItem>
            {/if}
          </SidebarGroup>
        </SidebarWrapper>
      </Sidebar>
    </Drawer>
    <div class="h-screen bg-white p-2">
      <button on:click={toggleDrawer}>
        <i title={drawerHidden ? "open navigation" : "close navigation"} class="bx bx-menu text-2xl"
        ></i>
      </button>
    </div>
  </div>
{/if}
