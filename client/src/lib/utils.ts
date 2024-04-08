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

export const request = async (
  path: string,
  requestMethod: string,
  formData?: FormData
): Promise<HttpResponse> => {
  try {
    const response = await fetch(path, {
      headers: {
        Authorization: `Bearer ${await getAccessToken()}`
      },
      method: requestMethod,
      body: formData
    });
    const json = await response.json();
    if (response.ok) {
      return { content: json, ok: true };
    } else {
      return { error: `${json.error ?? json.message}`, ok: false };
    }
  } catch (error: any) {
    return { error: `${error.name}: ${error.message}`, ok: false };
  }
};

const getAccessToken = async () => {
  const keycloak = appStore.getKeycloak();
  try {
    await keycloak.updateToken(5);
  } catch (error) {
    await keycloak.login();
  }

  return keycloak.token;
};
