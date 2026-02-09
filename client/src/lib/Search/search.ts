// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import type { SEARCHTYPES } from "$lib/Queries/query";

interface SearchParameters {
  currentPage?: number;
  limit?: number;
  orderBy?: string[];
  type?: SEARCHTYPES;
  detailed?: boolean;
  searchTerm?: string;
}

export type { SearchParameters };
