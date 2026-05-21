// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import type { PageLoad } from "./$types";
import { type User, UserManager } from "oidc-client-ts";
import { push } from "svelte-spa-router";
import { jwtDecode } from "jwt-decode";
import { configuration } from "$lib/configuration";
import { getErrorDetails } from "$lib/Errors/error";
import type { HttpResponse } from "$lib/types";
import { appStore } from "$lib/store.svelte";

const loadConfig = () => {
  return new Promise((resolve) => {
    fetch("api/client-config").then((response: any) => {
      if (response.ok) {
        response.json().then((content: any) => {
          appStore.setConfig(content);
          resolve(response);
        });
      } else {
        resolve(response);
      }
    });
  });
};

export const load: PageLoad = async () => {
  const response: any = await loadConfig();
  const userManager = new UserManager(configuration.getConfiguration());
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
          if (appStore.state.app.redirect) {
            push(appStore.state.app.redirect);
          } else {
            push("/");
          }
          if (!appStore.hasRole()) {
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
      if (!appStore.hasRole()) {
        appStore.setSessionExpired(true);
        appStore.setSessionExpiredMessage("User has no role");
        push("/login");
      }
    }
    appStore.setUserManager(userManager);
  });
  if (!response.ok) {
    const errorResponse: HttpResponse = response;
    errorResponse.error = response.status.toString();
    const loadConfigError = getErrorDetails(`Couldn't load Config.`, response);
    return { loadConfigError };
  }
  return;
};
