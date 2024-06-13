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
  import Diff from "$lib/Diff/DiffPage.svelte";
  import { wrap } from "svelte-spa-router/wrap";
  import Configuration from "$lib/Configuration/Overview.svelte";
  import Advisories from "$lib/Advisories/Overview.svelte";
  import Advisory from "$lib/Advisories/Advisory.svelte";
  import NotFound from "$lib/NotFound.svelte";
  import { appStore } from "$lib/store";
  import { push } from "svelte-spa-router";
  import Messages from "$lib/Messages/Messages.svelte";
  import Login from "$lib/Login/Login.svelte";
  import { configuration } from "$lib/configuration";
  import { type User, UserManager } from "oidc-client-ts";
  import { jwtDecode } from "jwt-decode";
  import QueryDesigner from "$lib/Queries/QueryDesigner.svelte";
  import QueryOverview from "$lib/Queries/Overview.svelte";
  import Test from "$lib/Test.svelte";

  let userManager = new UserManager(configuration.getConfiguration());
  userManager.events.addUserSignedIn(function () {
    console.log("User loaded");
  });
  userManager.events.addAccessTokenExpiring(function () {
    console.log("token expiring");
  });
  userManager.events.addSilentRenewError(function (e) {
    console.log("silent renew error", e.message);
    appStore.setIsUserLoggedIn(false);
    appStore.setSessionExpiredMessage(e.message);
    appStore.setSessionExpired(true);
    userManager.removeUser();
    push("/login");
  });
  userManager.getUser().then(async (user: User | null) => {
    if (!user) {
      userManager
        .signinRedirectCallback()
        .then(function (user) {
          appStore.setIsUserLoggedIn(true);
          appStore.setSessionExpired(false);
          appStore.setTokenParsed(jwtDecode(user.access_token));
          push("/");
          checkRoles();
        })
        .catch(function () {
          push("/login");
        });
    } else {
      appStore.setIsUserLoggedIn(true);
      appStore.setSessionExpired(false);
      appStore.setTokenParsed(jwtDecode(user.access_token));
      checkRoles();
    }
    appStore.setUserManager(userManager);
  });

  const checkRoles = () => {
    let hasRole =
      appStore.isAdmin() ||
      appStore.isEditor() ||
      appStore.isAuditor() ||
      appStore.isReviewer() ||
      appStore.isImporter();
    if (!hasRole) {
      appStore.setSessionExpired(true);
      appStore.setSessionExpiredMessage("User has no role");
      push("/login");
    }
  };

  const loginRequired = {
    loginRequired: true
  };

  const loginCondition = () => {
    if (!appStore.getUserManager()) return false;
    return appStore.getIsUserLoggedIn();
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
    "/queries/": wrap({
      component: QueryOverview,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/queries/new": wrap({
      component: QueryDesigner,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/queries/:id": wrap({
      component: QueryDesigner,
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
    "/test": wrap({
      component: Test,
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

<div class="flex bg-primary-700">
  <div>
    <SideNav></SideNav>
  </div>
  <main class="max-h-screen w-full overflow-auto bg-white p-6">
    {#if $appStore.app.userManager}
      <Router {routes} on:conditionsFailed={conditionsFailed} />
    {/if}
  </main>
  <Messages></Messages>
</div>
