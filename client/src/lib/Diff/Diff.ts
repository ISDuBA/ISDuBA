// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { request } from "$lib/request";

type DiffOperation = "add" | "remove" | "replace";

type JsonDiffResult = {
  op: DiffOperation;
  path: string;
  value?: string | object | object[];
};

type JsonDiffResultWrapper = {
  result: JsonDiffResult | JsonDiffResult[];
};

const fetchDiffEntry = (urlPath: string, operation: DiffOperation, jsonPath: string) => {
  const requestPath = encodeURI(`${urlPath}&item_op=${operation}&item_path=${jsonPath}`);
  return request(requestPath, "GET");
};

export { fetchDiffEntry };
export type { DiffOperation, JsonDiffResult, JsonDiffResultWrapper };
