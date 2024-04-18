<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Keycloak from "keycloak-js";
  import { configuration } from "$lib/configuration";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import { Button, Heading, Card } from "flowbite-svelte";
  import { PUBLIC_KEYCLOAK_URL, PUBLIC_KEYCLOAK_REALM } from "$env/static/public";
  import { A, P, Li, List } from "flowbite-svelte";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";
  import { request } from "$lib/utils";

  let version: string = "Retrieving Version from server";
  let error: string;

  if (!$appStore.app.keycloak) appStore.setKeycloak(new Keycloak(configuration.getConfiguration()));

  async function logout() {
    appStore.setSessionExpired(true);
    $appStore.app.keycloak.logout();
  }

  async function login() {
    try {
      await $appStore.app.keycloak.login();
      appStore.setSessionExpired(false);
    } catch {
      appStore.setSessionExpired(true);
    }
  }

  let profileUrl = PUBLIC_KEYCLOAK_URL + "/realms/isduba/account/#/";

  onMount(async () => {
    if (!$appStore.app.keycloak.authenticated) {
      await $appStore.app.keycloak
        .init({
          onLoad: "check-sso",
          checkLoginIframe: false,
          responseMode: "query"
        })
        .then(async () => {
          if ($appStore.app.keycloak.authenticated) {
            const profile = await $appStore.app.keycloak.loadUserProfile();
            appStore.setUserProfile({
              firstName: profile.firstName,
              lastName: profile.lastName
            });
            const expiry = new Date($appStore.app.keycloak.idTokenParsed.exp * 1000);
            appStore.setExpiryTime(expiry.toLocaleTimeString());
            push("/");
          }
        })
        .catch((error: any) => {
          console.log("error", error);
        });
    } else {
      const response = await request("api/about", "GET");
      if (response.ok) {
        const backendInfo = response.content;
        version = backendInfo.version;
      } else if (response.error) {
        error = response.error;
      }
    }
  });
</script>

<div class="mt-60 flex items-center justify-center">
  <div class="flex w-96 flex-col gap-4">
    <div><Heading class="">ISDuBA</Heading></div>
    <Card>
      <div class="flex flex-col gap-4">
        <div>
          <img
            style="height:2rem;"
            src={`${PUBLIC_KEYCLOAK_URL}/resources/zph0a/admin/keycloak.v2/logo.svg`}
          />
        </div>
        <P class="flex flex-col"
          ><span><b>Server URL: </b>{PUBLIC_KEYCLOAK_URL}</span><span
            ><b>Realm: </b>{PUBLIC_KEYCLOAK_REALM}</span
          ></P
        >
        {#if $appStore.app.keycloak && !$appStore.app.keycloak.authenticated}
          {#if $appStore.app.sessionExpired}
            <div class="text-yellow-400">
              <i class="bx bx-message-alt-error"></i> Your session is expired
            </div>
          {/if}
          <Button on:click={login}>Login</Button>
          <P>
            <A href="https://github.com/ISDuBA/" class="underline hover:no-underline"
              >Visit the ISDuBA project on Github</A
            ></P
          >
        {/if}
        {#if $appStore.app.keycloak && $appStore.app.keycloak.authenticated}
          <Button href={profileUrl}>Profile</Button>
          <Button on:click={logout}>Logout</Button>
          <small class="text-gray-400"
            >Your session ends due to inactivity {$appStore.app.expiryTime}</small
          >
        {/if}
      </div>
    </Card>
    <P class="mt-3">
      Versions:
      <List tag="ul" class="space-y-1" list="none">
        <Li liClass="ml-3">ISDuBA: {version}</Li>
      </List>
      <ErrorMessage message={error} plain={true}></ErrorMessage>
    </P>
  </div>
</div>
