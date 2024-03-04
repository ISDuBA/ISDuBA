// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2023 Intevation GmbH <https://intevation.de>

type hashEntry = {
  algorithm: string;
  value: string;
};

type xGenericURI = {
  namespace: string;
  uri: string;
};

export type { hashEntry, xGenericURI };
