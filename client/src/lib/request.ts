/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

import { appStore } from "./store.svelte";
import { push } from "$routes/router.svelte";
import type { User } from "oidc-client-ts";
import type { HttpResponse } from "./types";
import { jwtDecode } from "jwt-decode";

const requestData = async (
  abortController: AbortController | undefined,
  path: string,
  token: any,
  requestMethod: string,
  formData?: FormData | string
) => {
  if (abortController) {
    return fetch(path, {
      headers: {
        Authorization: `Bearer ${token}`
      },
      method: requestMethod,
      body: formData,
      signal: abortController.signal
    });
  } else {
    return await fetch(path, {
      headers: {
        Authorization: `Bearer ${token}`
      },
      method: requestMethod,
      body: formData
    });
  }
};

export const request = async (
  path: string,
  requestMethod: string,
  formData?: FormData | string | undefined,
  abortController?: AbortController
): Promise<HttpResponse> => {
  const httpResponse: HttpResponse = {
    ok: true
  };
  try {
    const token = await getAccessToken();
    const response = await requestData(abortController, path, token, requestMethod, formData);
    httpResponse.ok = response.ok;
    httpResponse.status = response.status;
    const contentType = response.headers.get("content-type");
    const isJson = contentType?.includes("application/json");
    let json;
    if (isJson) {
      try {
        json = await response.json();
      } catch (e) {
        httpResponse.ok = false;
        if (e instanceof DOMException && e.name == "AbortError") {
          httpResponse.error = e.name;
        } else {
          // 783 used by Shopify to indicate that the request includes a JSON syntax error. See https://shopify.dev/docs/api/usage/response-codes
          httpResponse.error = "783";
        }
        httpResponse.content = `${json.error}`;
        return httpResponse;
      }
    }
    const content = isJson ? json : await response.text();
    if (response.ok) {
      httpResponse.content = content;
      return httpResponse;
    }
    if (response.status == 401) {
      appStore.setSessionExpired(true);
      appStore.setSessionExpiredMessage("User unauthorized");
      await push("/login");
    }
    httpResponse.error = `${response.status}`;
    if (isJson) {
      // Handle pmd proxy errors
      if (json.messages) {
        httpResponse.content = json.messages;
      } else {
        httpResponse.content = json.error;
      }
    }
    return httpResponse;
  } catch (error: any) {
    httpResponse.ok = false;
    if (error.name === "AbortError") {
      httpResponse.error = error.name;
    }
    if (/fetch/.test(error)) {
      httpResponse.error = "600";
      httpResponse.content = error;
    }
    httpResponse.error = `${error.name}: ${error.message}`;
    return httpResponse;
  }
};

export const getAccessToken = async () => {
  const userManager = appStore.getUserManager();
  if (!userManager) {
    await push("/login");
    return;
  }

  return userManager.getUser().then(async (user: User | null) => {
    if (user) {
      appStore.setTokenParsed(jwtDecode(user.access_token));
      return user.access_token;
    } else {
      await push("/login");
    }
  });
};
