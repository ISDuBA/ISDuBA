<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import { Button, Heading, Card } from "flowbite-svelte";
  import { A, P, Li, List } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { request } from "$lib/request";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import DarkMode from "flowbite-svelte/DarkMode.svelte";

  let viewError: ErrorDetails | null = $state(null);
  let versionError: ErrorDetails | null = $state(null);

  async function logout() {
    appStore.setSessionExpired(true);
    appStore.setSessionExpiredMessage("Logout");
    sessionStorage.clear();
    await appStore.state.app.userManager?.signoutRedirect();
  }

  async function login() {
    try {
      await appStore.state.app.userManager?.signinRedirect();
    } catch (e: any) {
      viewError = getErrorDetails(`Could not load login information: ` + e.message);
    }
  }

  let profileUrl = `${appStore.getKeycloakURL()}/realms/${appStore.getKeycloakRealm()}/account/?referrer=${appStore.getKeycloakClientID()}&referrer_uri=${encodeURIComponent(window.location.href)}`;

  async function getVersion() {
    const response = await request("api/about", "GET");
    if (response.ok) {
      const backendInfo = response.content;
      return backendInfo.version;
    } else {
      versionError = getErrorDetails(`Couldn't load version.`, response);
    }
  }

  async function getView() {
    const response = await request("api/view", "GET");
    if (response.ok) {
      return new Map<string, [string]>(Object.entries(response.content));
    } else {
      viewError = getErrorDetails(`Couldn't determine your role.`, response);
    }
    return new Map<string, [string]>();
  }

  const getTLPClass = (label: string) => {
    if (label === "WHITE") {
      return "tlpclear";
    } else if (label === "RED") {
      return "tlpred";
    } else if (label === "AMBER") {
      return "tlpamber";
    } else if (label === "GREEN") {
      return "tlpgreen";
    }
  };
</script>

<svelte:head>
  <title>Login</title>
</svelte:head>

<div class="flex h-screen items-center justify-center">
  <div class="flex w-96 flex-col gap-4">
    <div class="inline-flex flex-row justify-between">
      <Heading class="mb-2 flex items-center gap-4">
        <img class="h-10" src="favicon.svg" alt="Icon of ISDuBA" aria-hidden="true" />
        <span>ISDuBA</span>
      </Heading>
      <DarkMode />
    </div>
    <Card class="p-4 sm:p-6 md:p-8">
      <div class="flex flex-col gap-4">
        <P class="flex flex-col"
          ><span><b>Server URL:</b>&nbsp;{appStore.getKeycloakURL()}</span><span
            ><b>Realm:</b>&nbsp;{appStore.getKeycloakRealm()}</span
          ></P
        >
        {#if appStore.state.app.userManager && !appStore.state.app.isUserLoggedIn}
          {#if appStore.state.app.sessionExpired}
            <div class="text-orange-700 dark:text-yellow-400">
              <i class="bx bx-message-alt-error"></i> Your session is expired: {appStore.state.app
                .sessionExpiredMessage || "Please login"}
            </div>
          {/if}
          <Button onclick={login}><i class="bx bx-link-external mr-1"></i> Login</Button>
        {/if}
        {#if appStore.state.app.userManager && appStore.state.app.isUserLoggedIn}
          <Button href={profileUrl}><i class="bx bx-link-external mr-1"></i> Profile</Button>
          <Button onclick={logout}><i class="bx bx-link-external mr-1"></i> Logout</Button>
        {/if}
        <div class="flex flex-row gap-4">
          <div class="flex flex-grow flex-col">
            <A href="/swagger/index.html" class="text-left no-underline hover:underline">API</A>
          </div>
          <div class="flex flex-grow flex-col">
            <span class="text-right dark:text-white">
              ISDuBA is Free Software →<A
                href="https://github.com/ISDuBA/"
                class="no-underline hover:underline">Github</A
              >
            </span>
          </div>
        </div>
      </div>
    </Card>
    {#if appStore.state.app.isUserLoggedIn && !appStore.state.app.sessionExpired}
      <div class="mt-4 flex w-full flex-row gap-4">
        <div class="flex flex-grow flex-col">
          <span class="text-xl">User:</span>
          <span class="ml-3">{appStore.state.app.tokenParsed?.preferred_username}</span>
        </div>
        {#if !viewError}
          <div class="flex flex-grow flex-col">
            <span class="text-xl">View: </span>
            <List tag="ul" class="list-none space-y-1">
              {#await getView() then view}
                {#each view.entries() as [publisher, tlps]}
                  <Li class="ml-3"
                    >{publisher === "*" ? "all" : publisher}:
                    {#each tlps as tlp}
                      <div
                        class={getTLPClass(tlp)}
                        style="width: fit-content; display: inline; margin-right: 0.25em;"
                      >
                        {tlp}
                      </div>
                    {/each}
                  </Li>
                {/each}
              {/await}
            </List>
          </div>
          <div class="flex flex-col">
            <span class="text-xl">Roles:</span>
            <List tag="ul" class="list-none space-y-1">
              {#if appStore.isAdmin()}
                <Li class="ml-3">Admin</Li>
              {/if}
              {#if appStore.isReviewer()}
                <Li class="ml-3">Reviewer</Li>
              {/if}
              {#if appStore.isAuditor()}
                <Li class="ml-3">Auditor</Li>
              {/if}
              {#if appStore.isImporter()}
                <Li class="ml-3">Importer</Li>
              {/if}
              {#if appStore.isEditor()}
                <Li class="ml-3">Editor</Li>
              {/if}
              {#if appStore.isSourceManager()}
                <Li class="ml-3">Source-Manager</Li>
              {/if}
            </List>
          </div>
        {/if}
      </div>
      <P>
        {#await getVersion() then version}
          {#if !versionError}
            <span class="text-m">Version: {version}</span>
          {/if}
        {/await}
      </P>
    {/if}
    <ErrorMessage error={viewError}></ErrorMessage>
    <ErrorMessage error={versionError}></ErrorMessage>
  </div>
</div>

<style>
  .tlpclear {
    background: #000;
    color: #fff;
  }
  .tlpred {
    background: #000;
    color: #ff2b2b;
  }
  .tlpamber {
    background: #000;
    color: #ffc000;
  }
  .tlpgreen {
    background: #000;
    color: #33ff00;
  }
</style>
