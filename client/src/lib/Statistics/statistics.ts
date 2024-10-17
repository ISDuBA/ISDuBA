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

type StatisticFilter = {
  downloadFailed?: boolean;
  filenameFailed?: boolean;
  schemaFailed?: boolean;
  remoteFailed?: boolean;
  checksumFailed?: boolean;
  signatureFailed?: boolean;
  duplicateFailed?: boolean;
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

const fetchStatistic = async (
  from: Date,
  to: Date,
  step: string,
  filter?: StatisticFilter,
  id?: number,
  feed: boolean = false
): Promise<Result<Statistic, ErrorDetails>> => {
  let path = "/api/stats/imports";
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
    `${path}?from=${toLocaleISOString(from)}&to=${toLocaleISOString(to)}&step=${step}` +
      filterQuery,
    "GET"
  );
  if (resp.ok) {
    if (resp.content) {
      return {
        ok: true,
        value: resp.content
      };
    }
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load statistic`, resp)
  };
};

export { fetchStatistic, setToEndOfDay, toLocaleISOString };
export type { Statistic, StatisticEntry, StatisticFilter };
