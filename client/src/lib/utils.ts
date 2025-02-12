/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

const truncate = (str: string, n: number) => {
  return str.length > n ? str.slice(0, n - 1) + "…" : str;
};

const areArraysEqual = (
  a: (number | string)[],
  b: (number | string)[],
  sameOrder: boolean = false
) => {
  if (a === b) return true;
  if (a.length !== b.length || typeof a[0] !== typeof b[0]) return false;
  for (let i = 0; i < a.length; i++) {
    if (sameOrder && a[i] !== b[i]) return false;
    else if (!b.includes(a[i])) return false;
  }
  return true;
};

const addSlashes = (str: string) => {
  return (str + "").replace(/[\\"']/g, "\\$&").replace(/\u0000/g, "\\0");
};

export { truncate, areArraysEqual, addSlashes };
