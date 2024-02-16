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
</script>

<Sidebar class="bg-primary-700 h-screen p-2">
  <SidebarWrapper class="bg-primary-700">
    <Heading class="mb-6 text-white">ISDuBA</Heading>
    <SidebarGroup class="bg-primary-700">
      {#if $appStore.app.isUserLoggedIn}
        <SidebarItem label="Advisories" href="/advisories" active={activeUrl === `/advisories`}>
          <svelte:fragment slot="icon">
            <i class="bx bx-spreadsheet"></i>
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
        <SidebarItem on:click={login} label="Login">
          <svelte:fragment slot="icon">
            <i class="bx bx-log-in"></i>
          </svelte:fragment>
        </SidebarItem>
      {/if}
    </SidebarGroup>
  </SidebarWrapper>
</Sidebar>
