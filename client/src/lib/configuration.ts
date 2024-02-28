// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: MIT
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

const configuration = {
  getConfiguration: () => {
    return {
      url: "http://localhost:8080",
      realm: "isduba",
      clientId: "auth"
    };
  }
};

export { configuration };
