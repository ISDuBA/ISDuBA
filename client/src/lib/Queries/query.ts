// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

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

export {
  generateQueryString,
  COLUMNS,
  ORDERDIRECTIONS,
  SEARCHTYPES,
  SEARCHPAGECOLUMNS,
  type Column,
  type Search
};
