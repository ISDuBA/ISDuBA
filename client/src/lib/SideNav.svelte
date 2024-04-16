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

  let notactivated =
    "flex items-center p-2 text-base font-normal text-gray-400 dark:text-gray-400 hover:bg-primary-100 hover:text-primary-900";

  $: activeUrl = "/" + $page.url.hash;
  let activeClass =
    "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-primary-700 dark:text-white hover:bg-primary-100 dark:hover:bg-gray-700";
  let nonActiveClass =
    "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 hover:text-primary-900";
</script>

{#if $appStore.app.keycloak && ($appStore.app.keycloak.authenticated || $appStore.app.sessionExpired)}
  <Sidebar class="bg-primary-700 h-screen p-2" {activeUrl} {activeClass} {nonActiveClass}>
    <SidebarWrapper class="bg-primary-700">
      <Heading class="mb-6 text-white">ISDuBA</Heading>
      <SidebarGroup class="bg-primary-700">
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
        {#if !$appStore.app.sessionExpired}
          <SidebarItem
            label={$appStore.app.keycloak.idTokenParsed.preferred_username}
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
{/if}
