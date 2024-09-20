// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

export const getRelativeTime = (date: Date) => {
  const now = Date.now();
  const passedTime = now - date.getTime();
  if (passedTime < 60000) {
    return "<1 min ago";
  } else if (passedTime < 3600000) {
    return `${Math.floor(passedTime / 60000)} min ago`;
  } else if (passedTime < 86400000) {
    return `${Math.floor(passedTime / 3600000)} hours ago`;
  } else {
    return `${Math.floor(passedTime / 86400000)} days ago`;
  }
};
