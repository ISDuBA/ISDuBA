// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

const MINUTE_MS = 1000 * 60;
const HOUR_MS = MINUTE_MS * 60;
const DAY_MS = HOUR_MS * 24;
const WEEK_MS = DAY_MS * 7;
const MONTH_MS = DAY_MS * 30;
const YEAR_MS = MONTH_MS * 12;

const getRelativeTime = (date: Date) => {
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

const setToEndOfDay = (date: Date) => {
  date.setHours(23);
  date.setMinutes(59);
  date.setSeconds(59);
  date.setMilliseconds(999);
  return date;
};

const pad = (n: number) => (n < 10 ? "0" + n : n);
const padMilliseconds = (n: number) => (n >= 100 ? n.toString() : n > 10 ? "0" + n : "00" + n);

const toLocaleISOString = (d: Date) => {
  return (
    d.getFullYear() +
    "-" +
    pad(d.getMonth() + 1) +
    "-" +
    pad(d.getDate()) +
    "T" +
    pad(d.getHours()) +
    ":" +
    pad(d.getMinutes()) +
    ":" +
    pad(d.getSeconds()) +
    "." +
    padMilliseconds(d.getMilliseconds()) +
    "Z"
  );
};

export {
  MINUTE_MS,
  HOUR_MS,
  DAY_MS,
  WEEK_MS,
  MONTH_MS,
  YEAR_MS,
  getRelativeTime,
  pad,
  setToEndOfDay,
  toLocaleISOString
};
