// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import type { Query, SEARCHTYPES } from "$lib/Queries/query";

interface SearchParameters {
  advanced?: boolean;
  currentPage?: number;
  limit?: number;
  orderBy?: string[];
  type?: SEARCHTYPES;
  detailed?: boolean;
  searchTerm?: string;
  queryID?: number;
}

type QueryState = {
  ignoredQueries: number[];
  queries: any[];
};

const queryState: QueryState = $state({
  ignoredQueries: [],
  queries: []
});

const defaultQueries = $derived.by(() => {
  return (
    queryState.queries
      ?.filter((q: Query) => q.default_query && !q.dashboard)
      .sort((a: any, b: any) => {
        if (a.global && b.default_query) return 1;
        if (a.default_query && b.global) return -1;
        return 0;
      }) ?? []
  );
});
const derivedDefaultQuery = $derived(defaultQueries.length > 0 ? defaultQueries[0] : undefined);
const defaultQuery = () => derivedDefaultQuery;

const {
  derivedVisibleQueries,
  derivedInvisibleQueries
}: { derivedVisibleQueries: Query[]; derivedInvisibleQueries: Query[] } = $derived.by(() => {
  const state = $state.snapshot(queryState);
  const queries = state.queries;
  const ignoredQueries = state.ignoredQueries;
  if (queries === null) return { derivedVisibleQueries: [], derivedInvisibleQueries: [] };
  const visible: Query[] = [];
  if (defaultQueries.length > 0) visible.push($state.snapshot(defaultQueries[0]));
  const invisible: Query[] = [];
  for (let i = 0; i < queries.length; i++) {
    const q = queries[i];
    const contained = visible.find((qu) => qu.id === q.id);
    if (contained) continue;
    if (!q.dashboard && !ignoredQueries.includes(q.id)) {
      visible.push(q);
    } else {
      invisible.push(q);
    }
  }
  return { derivedVisibleQueries: visible, derivedInvisibleQueries: invisible };
});
const invisibleQueries = () => derivedInvisibleQueries;

const derivedSortedQueries = $derived(
  derivedVisibleQueries.toSorted((a: any, b: any) => {
    if (a.global && !b.global) {
      return -1;
    } else if (!a.global && b.global) {
      return 1;
    }
    return 0;
  })
);
const sortedQueries = () => derivedSortedQueries;

export { queryState, defaultQuery, invisibleQueries, sortedQueries };
export type { SearchParameters };
