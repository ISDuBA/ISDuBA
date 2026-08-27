<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import type { PageProps } from "./$types";
  import Router from "svelte-spa-router";
  import "../app.css";
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
  import { push, replace } from "svelte-spa-router";
  import Messages from "$lib/Messages/Messages.svelte";
  import Login from "$lib/Login/Login.svelte";
  import QueryDesigner from "$lib/Queries/QueryDesigner.svelte";
  import QueryOverview from "$lib/Queries/Overview.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import SourceEditor from "$lib/Sources/SourceEditor.svelte";
  import DocumentUpload from "$lib/Sources/DocumentUpload.svelte";
  import SourceCreator from "$lib/Sources/SourceCreator.svelte";
  import AggregatorViewer from "$lib/Sources/Aggregators/AggregatorViewer.svelte";
  import FeedLogViewer from "$lib/Sources/FeedLogViewer.svelte";
  import { onMount } from "svelte";
  import RelatedDocuments from "$lib/Advisories/RelatedDocuments.svelte";
  import { routerState } from "./router.svelte";
  import FilterHelp from "$lib/Search/FilterHelp.svelte";

  let { data }: PageProps = $props();

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

  const loginRequired = {
    loginRequired: true
  };

  const loginCondition = async () => {
    if (!appStore.getUserManager()) return false;
    if (!appStore.hasRole()) return false;
    return appStore.getIsUserLoggedIn();
  };

  const routeLoaded = (detail: any) => {
    appStore.setRouterParams(detail.params);
    if (routerState.didPush === false) {
      appStore.setSearchResults(null);
      appStore.setSearchResultCount(null);
    }
    routerState.didPush = false;
  };

  const onRouteLoading = (detail: any) => {
    // Disable eslint warning which recommends to use SvelteURLSearchParams because we don't need reactivity in this place.
    // eslint-disable-next-line svelte/prefer-svelte-reactivity
    const searchParams = new URLSearchParams(detail.querystring);
    // If the URL contains the following params it is an URL created by Keycloak
    if (
      searchParams.get("state") &&
      searchParams.get("session_state") &&
      searchParams.get("iss") &&
      searchParams.get("code")
    ) {
      // These params have nothing to do with our client itself and might confuse users so we remove them
      searchParams.delete("state");
      searchParams.delete("session_state");
      searchParams.delete("iss");
      searchParams.delete("code");
      const sanitizedRoute = `${detail.location}${searchParams.toString()}`;
      replace(sanitizedRoute);
    }
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
    "/filter_help": wrap({
      component: FilterHelp as any,
      userData: loginRequired,
      conditions: [loginCondition]
    }),
    "*": NotFound as any
  };

  const conditionsFailed = (detail: any) => {
    if (detail.userData.loginRequired) {
      const location = detail.location;
      let redirectParam: string | undefined;
      if (location) {
        const redirectURL = `${window.location.origin}/#${location}`;
        appStore.setRedirect(redirectURL);
        redirectParam = `?redirect=${redirectURL}`;
      }
      push(`/login${redirectParam ?? ""}`);
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
      <Router
        {routes}
        onConditionsFailed={conditionsFailed}
        onRouteLoaded={routeLoaded}
        {onRouteLoading}
      />
    {/if}
    <ErrorMessage error={data.loadConfigError ?? null}></ErrorMessage>
  </main>
  <Messages></Messages>
</div>
