// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { getAccessToken } from "$lib/request";
import type { WorkflowState } from "$lib/workflow";

type StateChange = {
  publisher: string;
  trackingID: string;
  state: WorkflowState;
};

async function updateMultipleStates(newStates: StateChange[]) {
  const token = await getAccessToken();
  return fetch(`/api/status`, {
    headers: {
      Authorization: `Bearer ${token}`
    },
    method: "PUT",
    body: JSON.stringify(newStates)
  });
}

export { updateMultipleStates };
