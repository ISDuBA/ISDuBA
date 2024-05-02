// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

const COLUMNS = {
  ADVISORY: [
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
    "state"
  ],
  DOCUMENT: [
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
    "four_cves"
  ]
};

const ORDERDIRECTIONS = {
  ASC: "asc",
  DESC: "desc"
};

const SEARCHTYPES = {
  ADVISORY: "advisories",
  DOCUMENT: "documents"
};

const newQuery = () => {
  const activeColumns: any[] = [];
  return {
    currentStep: 1,
    searchType: SEARCHTYPES.ADVISORY,
    chosenColumns: activeColumns,
    activeColumns: [...COLUMNS.ADVISORY],
    description: "New Query",
    query: "",
    global: false
  };
};

const generateQueryString = (currentSearch: any) => {
  const columns = /search msg as/.test(currentSearch.query)
    ? [{ name: "msg" }, ...currentSearch.chosenColumns]
    : currentSearch.chosenColumns;
  const columnsParam = `&columns=${columns.map((col: any) => col.name).join(" ")}`;
  const order =
    currentSearch.chosenColumns.length > 0
      ? `&order=${currentSearch.chosenColumns
          .map((col: any) => {
            return col.searchOrder === ORDERDIRECTIONS.ASC ? col.name : `-${col.name}`;
          })
          .join(" ")}`
      : "";
  const query = currentSearch.query ? `&query=${currentSearch.query}` : "";
  const queryURL = `/api/documents?count=1&advisories=${currentSearch.searchType === SEARCHTYPES.ADVISORY}${columnsParam}${order}${query}`;
  return encodeURI(queryURL);
};

export { generateQueryString, COLUMNS, ORDERDIRECTIONS, SEARCHTYPES, newQuery };
