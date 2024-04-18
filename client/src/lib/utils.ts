/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

import { appStore } from "./store";
import { type HttpResponse } from "./types";
import { configuration } from "$lib/configuration";
import { push } from "svelte-spa-router";

export const request = async (
  path: string,
  requestMethod: string,
  formData?: FormData
): Promise<HttpResponse> => {
  try {
    const token = await getAccessToken();
    const response = await fetch(path, {
      headers: {
        Authorization: `Bearer ${token}`
      },
      method: requestMethod,
      body: formData
    });
    const contentType = response.headers.get("content-type");
    const isJson = contentType?.includes("application/json");
    if (response.ok) {
      if (contentType && isJson) {
        const json = await response.json();
        return { content: json, ok: true };
      } else {
        const text = await response.text();
        return { content: text, ok: true };
      }
    } else {
      if (contentType && isJson) {
        const json = await response.json();
        return { error: `${json.error ?? json.message}`, ok: false };
      } else {
        return { error: `${response.status}: ${response.statusText}`, ok: false };
      }
    }
  } catch (error: any) {
    return { error: `${error.name}: ${error.message}`, ok: false };
  }
};

const getAccessToken = async () => {
  const config = configuration.getConfiguration();
  const keycloak = appStore.getKeycloak();
  try {
    await keycloak.updateToken(config.updateIntervall);
    const expiry = new Date(keycloak.idTokenParsed.exp * 1000);
    appStore.setExpiryTime(expiry.toLocaleTimeString());
  } catch (error) {
    appStore.setSessionExpired(true);
    push("/login");
  }

  return keycloak.token;
};
