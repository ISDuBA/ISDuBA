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
  PUBLIC_KEYCLOAK_CLIENTID,
  PUBLIC_UPDATE_INTERVALL
} from "$env/static/public";

const configuration = {
  getConfiguration: () => {
    return {
      updateIntervall: PUBLIC_UPDATE_INTERVALL,
      url: PUBLIC_KEYCLOAK_URL,
      realm: PUBLIC_KEYCLOAK_REALM,
      clientId: PUBLIC_KEYCLOAK_CLIENTID
    };
  }
};

export { configuration };
