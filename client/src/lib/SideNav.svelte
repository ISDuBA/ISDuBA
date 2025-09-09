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
  import { appStore } from "$lib/store.svelte";
  import { page } from "$app/stores";
  import { truncate } from "$lib/utils";

  let activeUrl = $derived("/" + $page.url.hash);

  let activeClass =
    "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-gray-950 dark:text-white hover:bg-primary-100 dark:hover:bg-black";
  let nonActiveClass =
    "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 dark:hover:bg-black hover:text-primary-900";

  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn
  };
  let breakPoint: number = 1280;
  let width: number = $state(0);
  let drawerHidden: boolean = $state(false);
  $effect(() => {
    drawerHidden = width < breakPoint;
  });

  const toggleDrawer = () => {
    drawerHidden = !drawerHidden;
  };
</script>

<svelte:window bind:innerWidth={width} />
{#if appStore.state.app.userManager && (appStore.state.app.isUserLoggedIn || appStore.state.app.sessionExpired)}
  <div class="flex">
    <Drawer
      transitionType="fly"
      {transitionParams}
      bind:hidden={drawerHidden}
      activateClickOutside={false}
      width="fit-content"
      backdrop={false}
      class="bg-primary-700 static min-w-fit p-0 dark:bg-gray-900"
      id="sidebar"
    >
      <Sidebar
        asideClass="w-fit max-w-60 relative h-full"
        class="bg-primary-700 dark:bg-gray-900"
        {activeUrl}
        {activeClass}
        {nonActiveClass}
      >
        <SidebarWrapper
          class="bg-primary-700 px-0 dark:bg-gray-900"
          divClass="overflow-y-auto py-4 px-3 bg-gray-50 rounded"
        >
          <SidebarGroup>
            <SidebarBrand
              spanClass="self-center text-4xl font-normal whitespace-nowrap text-white me-4"
              site={{ img: "favicon.svg", name: "ISDuBA", href: "/#/" }}
            ></SidebarBrand>
          </SidebarGroup>
          <SidebarGroup class="bg-primary-700 space-y-0 dark:bg-gray-900">
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
            {#if appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager() || appStore.isImporter()}
              <SidebarItem class="px-6 py-2.5" label="Sources" href="/#/sources">
                <svelte:fragment slot="icon">
                  <i class="bx bx-git-repo-forked"></i>
                </svelte:fragment>
              </SidebarItem>
            {/if}
            {#if appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager()}
              <SidebarItem class="px-6 py-2.5" label="Aggregators" href="/#/sources/aggregators">
                <svelte:fragment slot="icon">
                  <i class="bx bx-sitemap"></i>
                </svelte:fragment>
              </SidebarItem>
            {/if}
            <SidebarItem class="px-6 py-2.5" label="Statistics" href="/#/statistics">
              <svelte:fragment slot="icon">
                <i class="bx bx-bar-chart-square"></i>
              </svelte:fragment>
            </SidebarItem>
            {#if !appStore.state.app.sessionExpired}
              <SidebarItem
                class="px-6 py-2.5"
                label={truncate(appStore.state.app.tokenParsed?.preferred_username ?? "", 15)}
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
    <div class="h-screen w-16 bg-white p-2 dark:bg-gray-800">
      <button
        onclick={toggleDrawer}
        aria-label={drawerHidden ? "open navigation" : "close navigation"}
      >
        <i title={drawerHidden ? "open navigation" : "close navigation"} class="bx bx-menu text-2xl"
        ></i>
      </button>
    </div>
  </div>
{/if}
