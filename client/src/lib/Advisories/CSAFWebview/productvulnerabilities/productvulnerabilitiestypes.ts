// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2023 Intevation GmbH <https://intevation.de>

export const ProductStatus = {
  FIXED: "FIXED",
  UNDER_INVESTIGATION: "UNDER_INVESTIGATION",
  KNOWN_AFFECTED: "KNOWN_AFFECTED",
  NOT_AFFECTED: "NOT_AFFECTED",
  RECOMMENDED: "RECOMMENDED"
} as const;

export const ProductStatusSymbol = {
  FIXED: "F",
  UNDER_INVESTIGATION: "U",
  KNOWN_AFFECTED: "K",
  NOT_AFFECTED: "N",
  RECOMMENDED: "R"
} as const;

type StringObject = {
  [key: string]: string;
};

export type FullProductName = {
  name: string;
  product_id: string;
};

export type Product = { product_id: string; name: string };

export type Vulnerability = {
  cve?: string;
  known_affected?: StringObject;
  fixed?: StringObject;
  under_investigation?: StringObject;
  known_not_affected?: StringObject;
  recommended?: StringObject;
};

export type VulnerabilitesExtractionResult = {
  vulnerabilities: Vulnerability[];
  relevantProducts: StringObject;
};
