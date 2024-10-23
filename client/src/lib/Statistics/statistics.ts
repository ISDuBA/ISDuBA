// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { getErrorDetails } from "$lib/Errors/error";
import type { ErrorDetails } from "$lib/Errors/error";
import { request } from "$lib/request";
import type { Result } from "$lib/types";

type StatisticEntry = [Date, number | null];
type Statistic = StatisticEntry[];

type CVSSTextualRating = "None" | "Low" | "Medium" | "High";
type CritStatisticEntry = [CVSSTextualRating | number];
type CritStatistic = [Date, CritStatisticEntry[]];

type StatisticType = "imports" | "importFailures" | "importFailuresCombined" | "cve" | "critical";

type StatisticFilter = {
  downloadFailed?: boolean;
  filenameFailed?: boolean;
  schemaFailed?: boolean;
  remoteFailed?: boolean;
  checksumFailed?: boolean;
  signatureFailed?: boolean;
  duplicateFailed?: boolean;
  [key: string]: boolean | undefined;
};

type StatisticGroup = {
  critical?: Statistic;
  cves?: Statistic;
  imports?: Statistic;
  signatureFailed?: Statistic;
  checksumFailed?: Statistic;
  filenameFailed?: Statistic;
  schemaFailed?: Statistic;
  downloadFailed?: Statistic;
  remoteFailed?: Statistic;
  duplicateFailed?: Statistic;
  importFailuresCombined?: Statistic;
  [key: string]: Statistic | undefined;
};

const getCVSSTextualRating = (CVSS: number): CVSSTextualRating => {
  if (CVSS === 0) {
    return "None";
  } else if (CVSS <= 3.9) {
    return "Low";
  } else if (CVSS <= 6.9) {
    return "Medium";
  } else {
    return "High";
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

const fetchBasicStatistic = async (
  from: Date,
  to: Date,
  step: number,
  type: StatisticType,
  id?: number,
  isFeed: boolean = false
): Promise<Result<StatisticGroup, ErrorDetails>> => {
  const stats: StatisticGroup = {};
  const response = await fetchStatistic(new Date(from), new Date(to), step, type, {}, id, isFeed);
  if (response.ok) {
    if (type === "imports") stats.imports = response.value;
    else stats[type] = response.value;
  } else if (response.error) {
    return {
      ok: false,
      error: response.error
    };
  }
  return { ok: true, value: stats };
};

const fetchImportFailuresStatistic = async (
  from: Date,
  to: Date,
  step: number,
  id?: number,
  isFeed: boolean = false
): Promise<Result<StatisticGroup, ErrorDetails>> => {
  const importStats: StatisticGroup = {};
  const failureTypes = [
    "signatureFailed",
    "checksumFailed",
    "filenameFailed",
    "schemaFailed",
    "downloadFailed",
    "remoteFailed",
    "duplicateFailed"
  ];
  for (let i = 0; i < failureTypes.length; i++) {
    const type: string = failureTypes[i];
    const filter: StatisticFilter = {};
    filter[type] = true;
    const response = await fetchStatistic(
      new Date(from),
      new Date(to),
      step,
      "imports",
      filter,
      id,
      isFeed
    );
    if (response.ok) {
      importStats[type] = response.value;
    } else if (response.error) {
      return {
        ok: false,
        error: response.error
      };
    }
  }
  return { ok: true, value: importStats };
};

const fetchStatistic = async (
  from: Date,
  to: Date,
  step: number,
  type: StatisticType,
  filter?: StatisticFilter,
  id?: number,
  feed: boolean = false
): Promise<Result<any, ErrorDetails>> => {
  let path = `/api/stats/${type}`;
  if (id && !feed) {
    path += `/source/${id}`;
  }
  if (id && feed) {
    path += `/feed/${id}`;
  }
  let filterQuery = "";
  if (filter) {
    filterQuery += filter.downloadFailed ? `&download_failed=true` : "";
    filterQuery += filter.filenameFailed ? `&filename_failed=true` : "";
    filterQuery += filter.schemaFailed ? `&schmema_failed=true` : "";
    filterQuery += filter.remoteFailed ? `&remote_failed=true` : "";
    filterQuery += filter.checksumFailed ? `&checksum_failed=true` : "";
    filterQuery += filter.signatureFailed ? `&signature_failed=true` : "";
    filterQuery += filter.duplicateFailed ? `&duplicate_failed=true` : "";
  }

  const resp = await request(
    `${path}?from=${toLocaleISOString(from)}&to=${toLocaleISOString(to)}&step=${step}ms` +
      filterQuery,
    "GET"
  );
  if (resp.ok) {
    if (resp.content) {
      for (let i = 0; i < resp.content.length; i++) {
        const date = new Date(resp.content[i][0]);
        date.setMinutes(date.getMinutes() + date.getTimezoneOffset());
        resp.content[i][0] = date;
      }
      return {
        ok: true,
        value: fillGaps(from, to, step, resp.content)
      };
    }
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load statistic`, resp)
  };
};

// Fill gaps with null values so the user can see at which times nothing was imported.
const fillGaps = (from: Date, to: Date, stepsInMilliseconds: number, values: any) => {
  const newStats: any = [];
  for (let i = from.getTime(); i <= to.getTime(); i += stepsInMilliseconds) {
    const foundValue: any = values.find((v: any) => v[0].getTime() === i);
    if (foundValue) {
      newStats.push(foundValue);
    } else {
      newStats.push([new Date(i), null]);
    }
  }
  return newStats;
};

export {
  fetchImportFailuresStatistic,
  fetchStatistic,
  fetchBasicStatistic,
  getCVSSTextualRating,
  pad,
  padMilliseconds,
  setToEndOfDay,
  toLocaleISOString
};
export type {
  StatisticGroup,
  Statistic,
  StatisticEntry,
  StatisticFilter,
  StatisticType,
  CritStatistic,
  CritStatisticEntry,
  CVSSTextualRating
};
