// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import {
  PUBLIC_KEYCLOAK_URL,
  PUBLIC_KEYCLOAK_REALM,
  PUBLIC_KEYCLOAK_CLIENTID
  // PUBLIC_UPDATE_INTERVALL
} from "$env/static/public";
import { type UserManagerSettings, WebStorageStateStore } from "oidc-client-ts";

const url = window.location.origin;

const configuration = {
  getConfiguration: (): UserManagerSettings => {
    return {
      authority: PUBLIC_KEYCLOAK_URL + "/realms/" + PUBLIC_KEYCLOAK_REALM,
      client_id: PUBLIC_KEYCLOAK_CLIENTID,
      redirect_uri: url + "/#/",
      post_logout_redirect_uri: url + "/#/login",
      response_type: "code",
      response_mode: "fragment",
      scope: "openid",

      automaticSilentRenew: true,

      filterProtocolClaims: true,
      userStore: new WebStorageStateStore({ store: localStorage })
    };
  }
};

export { configuration };
