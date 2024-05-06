<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { Button, Heading, Card } from "flowbite-svelte";
  import { PUBLIC_KEYCLOAK_URL, PUBLIC_KEYCLOAK_REALM } from "$env/static/public";
  import { A, P, Li, List } from "flowbite-svelte";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";
  import { request } from "$lib/utils";

  let error: string;

  async function logout() {
    appStore.setSessionExpired(true);
    appStore.setSessionExpiredMessage("Logout");
    await $appStore.app.userManager?.signoutRedirect();
  }

  async function login() {
    await $appStore.app.userManager?.signinRedirect();
  }

  let profileUrl = PUBLIC_KEYCLOAK_URL + "/realms/isduba/account/#/";

  async function getVersion() {
    const response = await request("api/about", "GET");
    if (response.ok) {
      const backendInfo = response.content;
      return backendInfo.version;
    } else if (response.error) {
      error = response.error;
    }
  }

  async function getView() {
    const response = await request("api/view", "GET");
    if (response.ok) {
      return new Map<string, [string]>(Object.entries(response.content));
    } else if (response.error) {
      error = response.error;
    }
    return new Map<string, [string]>();
  }
</script>

<svelte:head>
  <title>Login</title>
</svelte:head>

<div class="mt-60 flex items-center justify-center">
  <div class="flex w-96 flex-col gap-4">
    <div><Heading class="">ISDuBA</Heading></div>
    <Card>
      <div class="flex flex-col gap-4">
        <div>
          <img
            alt="Keycloak Logo"
            style="height:2rem;"
            src={`${PUBLIC_KEYCLOAK_URL}/resources/zph0a/admin/keycloak.v2/logo.svg`}
          />
        </div>
        <P class="flex flex-col"
          ><span><b>Server URL: </b>{PUBLIC_KEYCLOAK_URL}</span><span
            ><b>Realm: </b>{PUBLIC_KEYCLOAK_REALM}</span
          ></P
        >
        {#if $appStore.app.userManager && !$appStore.app.isUserLoggedIn}
          {#if $appStore.app.sessionExpired}
            <div class="text-yellow-400">
              <i class="bx bx-message-alt-error"></i> Your session is expired: {$appStore.app
                .sessionExpiredMessage || "Please login"}
            </div>
          {/if}
          <Button on:click={login}>Login</Button>
          <P>
            <A href="https://github.com/ISDuBA/" class="underline hover:no-underline"
              >Visit the ISDuBA project on Github</A
            ></P
          >
        {/if}
        {#if $appStore.app.userManager && $appStore.app.isUserLoggedIn}
          <Button href={profileUrl} target="_blank">Profile</Button>
          <Button on:click={logout}>Logout</Button>
        {/if}
      </div>
    </Card>
    {#if $appStore.app.isUserLoggedIn}
      <P class="mt-3">
        {#await getVersion() then version}
          Versions:
          <List tag="ul" class="space-y-1" list="none">
            <Li liClass="ml-3">ISDuBA: {version}</Li>
          </List>
        {/await}
        View:
        <List tag="ul" class="space-y-1" list="none">
          {#await getView() then view}
            {#each view.entries() as [publisher, tlp]}
              <Li liClass="ml-3">{publisher}: {tlp}</Li>
            {/each}
          {/await}
        </List>
        <ErrorMessage message={error}></ErrorMessage>
      </P>
    {/if}
  </div>
</div>
