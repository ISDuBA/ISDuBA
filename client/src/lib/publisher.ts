/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */
export const getPublisher = (publisher: string, width?: number) => {
  if (width && width > 1280) return publisher;
  switch (publisher) {
    case "Red Hat Product Security":
      return "RH";
    case "Siemens ProductCERT":
      return "SI";
    case "Bundesamt für Sicherheit in der Informationstechnik":
      return "BSI";
    case "SICK PSIRT":
      return "SCK";
    default:
      return publisher;
  }
};
