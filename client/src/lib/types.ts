// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

export type HttpResponse = {
  content?: any;
  error?: string;
  ok: boolean;
};

export type Result<T, E> = { ok: true; value: T } | { ok: false; error: E };

export type CurrentSearchQuery = {
  query: string;
  searchTerm: string;
  offset: number;
  orderBy: string[];
};
