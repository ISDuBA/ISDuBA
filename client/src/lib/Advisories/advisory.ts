// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { getAccessToken } from "$lib/request";
import { appStore } from "$lib/store";
import type { WorkflowState } from "$lib/workflow";
import { push } from "svelte-spa-router";

type StateChange = {
  publisher: string;
  trackingID: string;
  state: WorkflowState;
};

export type TrackingStatus = "draft" | "interim" | "final";

interface AdvisoryVersion {
  id: number;
  version: string;
  tracking_id: string;
  tracking_status: TrackingStatus;
}

// TODO: Refactor request from client/src/lib/request.ts so it can use JSON in body.
async function updateMultipleStates(newStates: StateChange[]) {
  try {
    const token = await getAccessToken();
    const response = await fetch(`/api/status`, {
      headers: {
        Authorization: `Bearer ${token}`
      },
      method: "PUT",
      body: JSON.stringify(newStates)
    });
    const contentType = response.headers.get("content-type");
    const isJson = contentType?.includes("application/json");
    let json;
    if (contentType && isJson) {
      try {
        json = await response.json();
      } catch (_) {
        return {
          error: "783", // Used by Shopify to indicate that the request includes a JSON syntax error. See https://shopify.dev/docs/api/usage/response-codes
          content: `${json.error}`,
          ok: false
        };
      }
    }
    const content = contentType && isJson ? json : await response.text();
    if (response.ok) {
      return { content: content, ok: true };
    }
    if (response.status == 401) {
      appStore.setSessionExpired(true);
      appStore.setSessionExpiredMessage("User unauthorized");
      await push("/login");
    }
    if (contentType && isJson) {
      return { error: `${response.status}`, content: json.error, ok: false };
    }
    return { error: `${response.status}`, content: content, ok: false };
  } catch (error: any) {
    if (/fetch/.test(error)) {
      return {
        error: "600",
        content: error,
        ok: false
      };
    }
    return {
      error: `${error.name}: ${error.message}`,
      ok: false
    };
  }
}

export const getAdvisoryLink = (item: any) =>
  `/advisories/${item.publisher}/${item.tracking_id}/documents/${item.id}`;
export const getAdvisoryAnchorLink = (item: any) => "#" + getAdvisoryLink(item);

export { updateMultipleStates };
export type { AdvisoryVersion };
