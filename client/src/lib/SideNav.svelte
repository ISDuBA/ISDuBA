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
  const sidebarItemClass = "px-0 py-0";
  const sidebarItemLinkClass =
    "px-6 py-4 rounded-none! hover:text-primary-700 dark:hover:text-white";

  let transitionParams = {
    x: -320,
    duration: 200,
    easing: sineIn
  };
  let breakPoint: number = 1280;
  let width: number = $state(0);
  let drawerOpen: boolean = $state(true);
  $effect(() => {
    drawerOpen = width >= breakPoint;
  });

  const toggleDrawer = () => {
    drawerOpen = !drawerOpen;
  };
</script>

<svelte:window bind:innerWidth={width} />
{#if appStore.state.app.userManager && (appStore.state.app.isUserLoggedIn || appStore.state.app.sessionExpired)}
  <div class="flex">
    <Drawer
      {transitionParams}
      bind:open={drawerOpen}
      outsideclose={false}
      width="default"
      modal={false}
      class="bg-primary-700 static w-fit min-w-fit border-none! p-0 dark:bg-gray-900"
      id="sidebar"
    >
      <Sidebar
        class="bg-primary-700 sidebar relative w-full dark:bg-gray-900"
        classes={{
          div: "max-w-60 relative h-screen bg-primary-700 dark:bg-gray-900 px-0",
          nonactive: nonActiveClass,
          active: activeClass
        }}
        {activeUrl}
      >
        <SidebarWrapper
          class="bg-primary-700 sidebar-wrapper w-full overflow-y-auto rounded bg-gray-50 px-0 py-4 dark:bg-gray-900"
        >
          <SidebarGroup>
            <SidebarBrand
              classes={{
                span: "self-center text-4xl font-normal whitespace-nowrap text-white me-4"
              }}
              site={{ img: "favicon.svg", name: "ISDuBA", href: "/#/" }}
            ></SidebarBrand>
          </SidebarGroup>
          <SidebarGroup class="bg-primary-700 w-full space-y-0 dark:bg-gray-900">
            <!-- Entries which are available after login should go here-->
            <SidebarItem
              class={sidebarItemClass}
              aClass={sidebarItemLinkClass}
              label="Dashboard"
              href="/#/"
            >
              {#snippet icon()}
                <i class="bx bxs-dashboard"></i>
              {/snippet}
            </SidebarItem>
            <SidebarItem
              class={sidebarItemClass}
              aClass={sidebarItemLinkClass}
              label="Search"
              href="/#/search"
            >
              {#snippet icon()}
                <i class="bx bx-spreadsheet"></i>
              {/snippet}
            </SidebarItem>
            {#if appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager() || appStore.isImporter()}
              <SidebarItem
                class={sidebarItemClass}
                aClass={sidebarItemLinkClass}
                label="Sources"
                href="/#/sources"
              >
                {#snippet icon()}
                  <i class="bx bx-git-repo-forked"></i>
                {/snippet}
              </SidebarItem>
            {/if}
            {#if appStore.isAuditor() || appStore.isEditor() || appStore.isSourceManager()}
              <SidebarItem
                class={sidebarItemClass}
                aClass={sidebarItemLinkClass}
                label="Aggregators"
                href="/#/sources/aggregators"
              >
                {#snippet icon()}
                  <i class="bx bx-sitemap"></i>
                {/snippet}
              </SidebarItem>
            {/if}
            <SidebarItem
              class={sidebarItemClass}
              aClass={sidebarItemLinkClass}
              label="Statistics"
              href="/#/statistics"
            >
              {#snippet icon()}
                <i class="bx bx-bar-chart-square"></i>
              {/snippet}
            </SidebarItem>
            {#if !appStore.state.app.sessionExpired}
              <SidebarItem
                class={sidebarItemClass}
                aClass={sidebarItemLinkClass}
                label={truncate(appStore.state.app.tokenParsed?.preferred_username ?? "", 15)}
                href="/#/login"
              >
                {#snippet icon()}
                  <i class="bx bx-user"></i>
                {/snippet}
              </SidebarItem>
            {/if}
          </SidebarGroup>
        </SidebarWrapper>
      </Sidebar>
    </Drawer>
    <div class="h-screen w-12 bg-white p-1 dark:bg-gray-800">
      <button
        onclick={toggleDrawer}
        aria-label={drawerOpen ? "open navigation" : "close navigation"}
      >
        <i title={drawerOpen ? "open navigation" : "close navigation"} class="bx bx-menu text-2xl"
        ></i>
      </button>
    </div>
  </div>
{/if}
