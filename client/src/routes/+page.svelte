<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

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
  import Login from "$lib/Login/Login.svelte";
  import Sources from "$lib/Sources/Overview.svelte";
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

  appStore.setKeycloak(new Keycloak(configuration.getConfiguration()));

  const loginRequired = {
    loginRequired: true
  };

  const loginCondition = () => {
    return $appStore.app.isUserLoggedIn;
  };

  const routes = {
    "/": wrap({
      component: Home,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/login": wrap({
      component: Login
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
      push("/login");
    }
  };
</script>

<div class="bg-primary-700 flex">
  <SideNav></SideNav>
  <main class="w-full bg-white pl-6 pt-6">
    <Router {routes} on:conditionsFailed={conditionsFailed} />
  </main>
</div>
