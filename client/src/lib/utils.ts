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

export const request = async (
  path: string,
  requestMethod: string,
  formData?: FormData
): any | undefined => {
  await appStore.getKeycloak().updateToken(5);
  const response = await fetch(path, {
    headers: {
      Authorization: `Bearer ${await getAccessToken()}`
    },
    method: requestMethod,
    body: formData
  });
  if (response.ok) {
    return response;
  } else {
    try {
      const error = await response.json();
      appStore.displayErrorMessage(`${error.error ?? error.message}`);
    } catch {
      appStore.displayErrorMessage(`${response.status} - ${response.statusText}`);
    }
    return undefined;
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
