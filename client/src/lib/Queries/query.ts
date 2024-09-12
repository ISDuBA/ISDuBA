// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { request } from "$lib/request";
import type { Role } from "$lib/workflow";

const COLUMNS = {
  ADVISORY: [
    "critical",
    "id",
    "tracking_id",
    "version",
    "publisher",
    "current_release_date",
    "initial_release_date",
    "title",
    "tlp",
    "cvss_v2_score",
    "cvss_v3_score",
    "ssvc",
    "four_cves",
    "state",
    "comments",
    "recent",
    "versions"
  ],
  DOCUMENT: [
    "critical",
    "id",
    "tracking_id",
    "version",
    "publisher",
    "current_release_date",
    "initial_release_date",
    "title",
    "tlp",
    "cvss_v2_score",
    "cvss_v3_score",
    "four_cves",
    "comments"
  ],
  EVENT: [
    "critical",
    "id",
    "tracking_id",
    "version",
    "publisher",
    "current_release_date",
    "initial_release_date",
    "title",
    "tlp",
    "cvss_v2_score",
    "cvss_v3_score",
    "ssvc",
    "four_cves",
    "comments",
    "event",
    "event_state",
    "time",
    "actor",
    "comments_id"
  ]
};

enum ORDERDIRECTIONS {
  ASC = "asc",
  DESC = "desc"
}

enum SEARCHTYPES {
  ADVISORY = "advisories",
  DOCUMENT = "documents",
  EVENT = "events"
}

const SEARCHPAGECOLUMNS = {
  ADVISORY: [
    "critical",
    "cvss_v3_score",
    "cvss_v2_score",
    "ssvc",
    "state",
    "four_cves",
    "publisher",
    "title",
    "tracking_id",
    "initial_release_date",
    "current_release_date",
    "version",
    "comments",
    "recent",
    "versions"
  ],
  DOCUMENT: [
    "critical",
    "cvss_v3_score",
    "cvss_v2_score",
    "ssvc",
    "four_cves",
    "publisher",
    "title",
    "tracking_id",
    "initial_release_date",
    "current_release_date",
    "version",
    "comments"
  ]
};
interface Column {
  name: string;
  visible: boolean;
}

interface Search {
  searchType: SEARCHTYPES;
  columns: Column[];
  orderBy: [string, ORDERDIRECTIONS][];
  name: string;
  query: string;
  description: string;
  global: boolean;
  dashboard: boolean;
  role: Role | undefined;
}

type Query = {
  [key: string]: boolean | string | string[] | number | undefined;
  advisories: boolean;
  columns: string[];
  definer: string;
  global: boolean;
  id: number;
  name: string;
  kind: SEARCHTYPES;
  num: number;
  orders: string[] | undefined;
  query: string;
  description: string | undefined;
  dashboard: boolean;
  role: Role;
};

const generateQueryString = (currentSearch: Search) => {
  const chosenColumns = currentSearch.columns.filter((c: any) => {
    return c.visible === true;
  });
  const columns = /search msg as/.test(currentSearch.query)
    ? [{ name: "msg" }, ...chosenColumns]
    : chosenColumns;
  const columnsParam = `&columns=${columns.map((col: any) => col.name).join(" ")}`;
  const query = currentSearch.query ? `&query=${currentSearch.query}` : "";
  const advisoriesParam =
    currentSearch.searchType !== SEARCHTYPES.EVENT
      ? `advisories=${currentSearch.searchType === SEARCHTYPES.ADVISORY}`
      : "";
  const queryURL = `/api/${currentSearch.searchType === SEARCHTYPES.EVENT ? "events" : "documents"}?count=1&${advisoriesParam}${columnsParam}${query}`;
  return encodeURI(queryURL);
};

const createStoredQuery = (query: Query) => {
  return saveStoredQuery(query, "POST");
};

const updateStoredQuery = (query: Query) => {
  return saveStoredQuery(query, "PUT");
};

const saveStoredQuery = (query: Query, method: string) => {
  const formData = new FormData();
  if (method === "PUT") {
    formData.append("num", `${query.num}`);
  }
  formData.append("kind", query.kind);
  formData.append("name", query.name);
  formData.append("global", `${query.global}`);
  formData.append("dashboard", `${query.dashboard}`);
  if (query.role) {
    formData.append("role", `${query.role}`);
  } else {
    formData.append("role", "");
  }
  if (query.description && query.description.length > 0) {
    formData.append("description", query.description);
  }
  if (query.query.length > 0) {
    formData.append("query", query.query);
  }
  formData.append("columns", query.columns.join(" "));
  if (query.orders) {
    formData.append("orders", query.orders.join(" "));
  }
  const path = method === "PUT" ? `/${query.id}` : "";
  return request(`/api/queries${path}`, method, formData);
};

/**
 * Takes the list of existing queries, looks for already given clones and returns a proper name.
 * Expamples:
 *
 * For non existing clones
 *
 * Monat -> Monat (1)
 * Monat (1) -> Monat (1) (1)
 *
 * Say there is already a clone
 *
 * Monat and Monat (1) -> Monat (2)
 * Monat (1) and Monat (1) (1) -> Monat (1) (2)
 * Monat (1) (2) and Monat (1) (1) -> Monat (1) (3)
 *
 * And so on.
 *
 * @param queries list of queries
 * @param name name of the query
 */
const proposeName = (queries: Query[], name: string) => {
  const clones = queries
    .filter((r: any) => {
      const re = new RegExp(name.replaceAll("(", "\\(").replaceAll(")", "\\)") + " \\(\\d+\\)");
      return re.test(r.name);
    })
    .map((r: any) => {
      return r.name;
    })
    .sort((a: string, b: string) => a.localeCompare(b, "en", { numeric: true }));
  if (clones.length === 0) return `${name} (1)`;
  const highestIndex = parseInt(clones[clones.length - 1].split(name + " (")[1]);
  return `${name} (${highestIndex + 1})`;
};

export {
  generateQueryString,
  createStoredQuery,
  proposeName,
  updateStoredQuery,
  COLUMNS,
  ORDERDIRECTIONS,
  SEARCHTYPES,
  SEARCHPAGECOLUMNS
};
export type { Column, Query, Search };
