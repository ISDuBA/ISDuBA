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

/**
 * Splits a string so it can be easily rendered inside an each loop.
 * @param text String including search matches
 * @param positions Start index and length of a match inside the text. Coming from the backend.
 * @returns Strings, alternating non-match and match (always starts with non-match)
 */
const splitMatches = (text: string, positions: number[][]): string[] => {
  let lastPos = 0;
  const splits: string[] = [];
  for (let i = 0; i < positions.length; i++) {
    const pos = positions[i];
    const term = text.substring(pos[0], pos[0] + pos[1]);
    // Don't use the term to split the text although it would be easier because the method could find
    // other occurrences that were not considered by the backend.
    splits.push(text.slice(lastPos, pos[0]), term);
    lastPos = pos[0] + pos[1];
    if (i === positions.length - 1) {
      splits.push(text.slice(pos[0] + pos[1]));
    }
  }
  return splits;
};

export { truncate, areArraysEqual, addSlashes, splitMatches };
