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
  const loadPath: string = window.location.hash.replace("#", "");
  const response: any = await loadConfig();
  if (!response.ok) {
    const errorResponse: HttpResponse = response;
    errorResponse.error = response.status.toString();
    const loadConfigError = getErrorDetails(`Couldn't load Config.`, response);
    return { loadConfigError };
  }
  const userManager = new UserManager(configuration.getConfiguration());
  appStore.setUserManager(userManager);
  const sessionExpired = (e?: any) => {
    appStore.setIsUserLoggedIn(false);
    if (e) appStore.setSessionExpiredMessage(e.message);
    appStore.setSessionExpired(true);
    userManager.removeUser();
    if (appStore.state.app.redirect && window.location.hash.includes("/login")) {
      push(appStore.state.app.redirect);
    } else {
      push("/login");
    }
  };
  const user: User | null = await userManager.getUser();
  let isExpired = false;

  if (!user) {
    try {
      const user: any = await userManager.signinRedirectCallback();
      appStore.setIsUserLoggedIn(true);
      appStore.setSessionExpired(false);
      appStore.setTokenParsed(jwtDecode(user.access_token));
      push("/");
      if (!appStore.hasRole()) {
        appStore.setSessionExpired(true);
        appStore.setSessionExpiredMessage("User has no role");
        push("/login");
      }
    } catch (e) {
      console.error(e);
      push("/login");
    }
  } else {
    appStore.setIsUserLoggedIn(true);
    appStore.setSessionExpired(false);
    appStore.setTokenParsed(jwtDecode(user.access_token));
    if (!appStore.hasRole()) {
      appStore.setSessionExpired(true);
      appStore.setSessionExpiredMessage("User has no role");
      push("/login");
    } else {
      // We have to save the current URL at this place or we don't know where to redirect the user
      // since if sessionExpired is called by the userManager the URL is already changed.
      // It's also important that we only want to redirect the user if they if they initially opened
      // a page and their session is already expired.
      const token = jwtDecode(user.access_token);
      if (token.exp) {
        const expDate = new Date(token.exp * 1000);
        if (expDate.getTime() < Date.now()) {
          if (loadPath) {
            appStore.setRedirect(loadPath);
          }
          sessionExpired();
          isExpired = true;
        }
      }
    }
  }
  if (!isExpired) {
    userManager.events.addSilentRenewError(sessionExpired);
    userManager.events.addAccessTokenExpired(sessionExpired);
  }

  return;
};
