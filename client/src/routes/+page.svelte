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
  import Dashboard from "$lib/Dashboard/Dashboard.svelte";
  import Statistics from "$lib/Statistics/Overview.svelte";
  import Sources from "$lib/Sources/Overview.svelte";
  import Diff from "$lib/Diff/Diff.svelte";
  import { wrap } from "svelte-spa-router/wrap";
  import Configuration from "$lib/Configuration/Overview.svelte";
  import Search from "$lib/Search/Overview.svelte";
  import Advisory from "$lib/Advisories/Advisory.svelte";
  import NotFound from "$lib/NotFound.svelte";
  import { appStore } from "$lib/store.svelte";
  import { push } from "svelte-spa-router";
  import Messages from "$lib/Messages/Messages.svelte";
  import Login from "$lib/Login/Login.svelte";
  import { configuration } from "$lib/configuration";
  import { type User, UserManager } from "oidc-client-ts";
  import { jwtDecode } from "jwt-decode";
  import QueryDesigner from "$lib/Queries/QueryDesigner.svelte";
  import QueryOverview from "$lib/Queries/Overview.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import SourceEditor from "$lib/Sources/SourceEditor.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import type { HttpResponse } from "$lib/types";
  import DocumentUpload from "$lib/Sources/DocumentUpload.svelte";
  import SourceCreator from "$lib/Sources/SourceCreator.svelte";
  import AggregatorViewer from "$lib/Sources/Aggregators/AggregatorViewer.svelte";
  import FeedLogViewer from "$lib/Sources/FeedLogViewer.svelte";
  import { onMount } from "svelte";
  import RelatedDocuments from "$lib/Advisories/RelatedDocuments.svelte";

  let loadConfigError: ErrorDetails | null = $state(null);

  const loadConfig = () => {
    return new Promise((resolve) => {
      fetch("api/client-config").then((response: any) => {
        if (response.ok) {
          response.json().then((content: any) => {
            appStore.setConfig(content);
            resolve(response);
          });
        } else {
          let errorRespose: HttpResponse = response;
          errorRespose.error = response.status.toString();
          loadConfigError = getErrorDetails(`Couldn't load Config.`, response);
          resolve(response);
        }
      });
    });
  };

  let onConfigLoad: any;

  let inactivityTime = () => {
    let time: ReturnType<typeof setTimeout>;
    const timeoutTime = () => Number(appStore.getIdleTimeout() / (1000 * 1000));
    const resetTimer = () => {
      clearTimeout(time);
      localStorage.setItem("lastActivity", Date.now().valueOf().toString());
      time = setTimeout(logout, timeoutTime());
    };

    document.onmousemove = resetTimer;
    document.onkeydown = resetTimer;
    onConfigLoad = resetTimer;

    const logout = async () => {
      let lastActivity = localStorage.getItem("lastActivity");
      if (lastActivity) {
        const elapsedTime = Date.now() - parseInt(lastActivity, 10);
        if (elapsedTime < timeoutTime()) {
          return;
        }
      }
      appStore.setSessionExpired(true);
      appStore.setSessionExpiredMessage("Idle logout");
      sessionStorage.clear();
      await appStore.getUserManager()?.signoutRedirect();
    };
  };

  inactivityTime();

  loadConfig().then(() => {
    if (onConfigLoad) {
      onConfigLoad();
    }
    let userManager = new UserManager(configuration.getConfiguration());
    const sessionExpired = (e: any) => {
      appStore.setIsUserLoggedIn(false);
      if (e) appStore.setSessionExpiredMessage(e.message);
      appStore.setSessionExpired(true);
      userManager.removeUser();
      push("/login");
    };
    userManager.events.addSilentRenewError(sessionExpired);
    userManager.events.addAccessTokenExpired(sessionExpired);
    userManager.getUser().then(async (user: User | null) => {
      if (!user) {
        userManager
          .signinRedirectCallback()
          .then(function (user: any) {
            appStore.setIsUserLoggedIn(true);
            appStore.setSessionExpired(false);
            appStore.setTokenParsed(jwtDecode(user.access_token));
            push("/");
            const hasAnyRole = checkUserForRoles();
            if (!hasAnyRole) {
              appStore.setSessionExpired(true);
              appStore.setSessionExpiredMessage("User has no role");
              push("/login");
            }
          })
          .catch(function () {
            push("/login");
          });
      } else {
        appStore.setIsUserLoggedIn(true);
        appStore.setSessionExpired(false);
        appStore.setTokenParsed(jwtDecode(user.access_token));
        const hasAnyRole = checkUserForRoles();
        if (!hasAnyRole) {
          appStore.setSessionExpired(true);
          appStore.setSessionExpiredMessage("User has no role");
          push("/login");
        }
      }
      appStore.setUserManager(userManager);
    });
  });

  const checkUserForRoles = () => {
    let hasRole =
      appStore.isAdmin() ||
      appStore.isEditor() ||
      appStore.isAuditor() ||
      appStore.isReviewer() ||
      appStore.isSourceManager() ||
      appStore.isImporter();
    return hasRole;
  };

  const loginRequired = {
    loginRequired: true
  };

  const loginCondition = () => {
    if (!appStore.getUserManager()) return false;
    if (!checkUserForRoles()) return false;
    return appStore.getIsUserLoggedIn();
  };

  const routes = {
    "/": wrap({
      component: Dashboard as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/login": wrap({
      component: Login as any
    }),
    "/advisories/:publisherNamespace/:trackingID/documents/:id/:position?": wrap({
      component: Advisory as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/documents/:id/:position?": wrap({
      component: Advisory as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/advisories/:publisherNamespace/:trackingID/documents/:id/related/documents/:cve?": wrap({
      component: RelatedDocuments as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/documents/:id/related/documents/:cve?": wrap({
      component: RelatedDocuments as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/search": wrap({
      component: Search as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/configuration": wrap({
      component: Configuration as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/queries/": wrap({
      component: QueryOverview as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/queries/new": wrap({
      component: QueryDesigner as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/queries/:id": wrap({
      component: QueryDesigner as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/diff": wrap({
      component: Diff as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/statistics": wrap({
      component: Statistics as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources": wrap({
      component: Sources as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources/new": wrap({
      component: SourceCreator as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources/new/:domain": wrap({
      component: SourceCreator as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources/logs/:id": wrap({
      component: FeedLogViewer as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources/upload": wrap({
      component: DocumentUpload as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources/aggregators": wrap({
      component: AggregatorViewer as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "/sources/aggregators/:id": wrap({
      component: AggregatorViewer as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),

    "/sources/:id": wrap({
      component: SourceEditor as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "*": NotFound as any
  };

  const conditionsFailed = (event: any) => {
    if (event.detail.userData.loginRequired) {
      push("/login");
    }
  };

  onMount(() => {
    appStore.updateDarkMode();
    const darkModeObserver = new MutationObserver((mutations) => {
      mutations.forEach((mutation) => {
        if (mutation.attributeName === "class") {
          appStore.updateDarkMode();
        }
      });
    });
    darkModeObserver.observe(document.documentElement, {
      attributes: true
    });
  });
</script>

<!--
  To ensure current darkmode setting is always processed,
  not only when the DarkMode button is on screen.
-->
<svelte:head>
  <script>
    if ("THEME_PREFERENCE_KEY" in localStorage) {
      // explicit preference - overrides author's choice
      localStorage.getItem("THEME_PREFERENCE_KEY") === "dark"
        ? window.document.documentElement.classList.add("dark")
        : window.document.documentElement.classList.remove("dark");
    } else {
      // browser preference - does not overrides
      if (window.matchMedia("(prefers-color-scheme: dark)").matches)
        window.document.documentElement.classList.add("dark");
    }
  </script>
</svelte:head>

<div class="bg-primary-700 flex h-screen dark:bg-gray-800 dark:text-white">
  <div>
    <SideNav></SideNav>
  </div>
  <main
    class="flex max-h-screen w-full flex-col overflow-auto bg-white px-2 py-6 lg:px-6 dark:bg-gray-800"
  >
    {#if appStore.state.app.userManager}
      <Router {routes} on:conditionsFailed={conditionsFailed} />
    {/if}
    <ErrorMessage error={loadConfigError}></ErrorMessage>
  </main>
  <Messages></Messages>
</div>
