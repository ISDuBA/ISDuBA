<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import Router from "svelte-spa-router";
  import "../app.pcss";
  import "boxicons/css/boxicons.min.css";
  import SideNav from "$lib/SideNav.svelte";
  import Home from "$lib/Home/Home.svelte";
  import Statistics from "$lib/Statistics/Overview.svelte";
  import Sources from "$lib/Sources/Overview.svelte";
  import About from "$lib/About/About.svelte";
  import Diff from "$lib/Diff/DiffPage.svelte";
  import { wrap } from "svelte-spa-router/wrap";
  import Configuration from "$lib/Configuration/Overview.svelte";
  import Documents from "$lib/Documents/Overview.svelte";
  import Advisories from "$lib/Advisories/Overview.svelte";
  import Advisory from "$lib/Advisories/Advisory.svelte";
  import NotFound from "$lib/NotFound.svelte";
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import Keycloak from "keycloak-js";
  import { configuration } from "$lib/configuration";
  import { onMount } from "svelte";
  import Messages from "$lib/Messages/Messages.svelte";

  appStore.setKeycloak(new Keycloak(configuration.getConfiguration()));

  onMount(async () => {
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
        }
      })
      .catch((error: any) => {
        console.log("error", error);
      });
  });

  const loginRequired = {
    loginRequired: true
  };

  const loginCondition = async () => {
    if (!$appStore.app.keycloak.authenticated) return false;
    const keycloak = appStore.getKeycloak();
    try {
      await keycloak.updateToken(5);
      return true;
    } catch (error) {
      await keycloak.login();
      return false;
    }
  };

  const routes = {
    "/": wrap({
      component: Home,
      userData: loginRequired,
      conditions: []
    }),
    "/about": wrap({
      component: About
    }),
    "/advisories/:publisherNamespace/:trackingID/documents/:id": wrap({
      component: Advisory,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/advisories": wrap({
      component: Advisories,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/configuration": wrap({
      component: Configuration,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/documents": wrap({
      component: Documents,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/diff": wrap({
      component: Diff,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/statistics": wrap({
      component: Statistics,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources": wrap({
      component: Sources,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "*": NotFound
  };

  const conditionsFailed = (event: any) => {
    if (event.detail.userData.loginRequired) {
      push("/");
    }
  };
</script>

<div class="flex bg-primary-700">
  <div>
    <SideNav></SideNav>
    {#if $appStore.app.keycloak.authenticated}
      <div style="position:absolute; top:4.5em; left:1.5em; color:white">
        Session ends at {$appStore.app.expiryTime}
      </div>
    {/if}
  </div>
  <main class="max-h-screen w-full bg-white p-6">
    <Router {routes} on:conditionsFailed={conditionsFailed} />
  </main>
  <Messages></Messages>
</div>
