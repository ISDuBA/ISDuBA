// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
import { getAccessToken, request } from "$lib/request";
import { appStore } from "$lib/store.svelte";
import type { WorkflowState } from "$lib/workflow";
import { push } from "$routes/router.svelte";
import {
  extractProducts,
  generateProductVulnerabilities
} from "./CSAFWebview/productvulnerabilities/productvulnerabilities";
import { getRevisionHistory } from "./docmodel";
import type { Revision as Revision2_0 } from "./types/csaf-2.0";
import type { Revision as Revision2_1 } from "./types/csaf-2.1";
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

interface AdvisoryVersionsResult {
  advisoryVersions?: AdvisoryVersion[];
  error?: ErrorDetails;
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

const loadAdvisoryVersions = async (
  encodedTrackingID: string,
  encodedPublisherNamespace: string
): Promise<AdvisoryVersionsResult | undefined> => {
  const response = await request(
    `/api/documents?&columns=id version tracking_id tracking_status&query=$tracking_id ${encodedTrackingID} = $publisher "${encodedPublisherNamespace}" = and`,
    "GET"
  );
  if (response.ok) {
    const result = await response.content;
    if (result.documents) {
      const advisoryVersions: AdvisoryVersion[] = result.documents.map((doc: any) => {
        return {
          id: doc.id,
          version: doc.version,
          tracking_id: doc.tracking_id,
          tracking_status: doc.tracking_status as TrackingStatus
        };
      });

      // Define the order of tracking statuses
      const statusOrder: Record<TrackingStatus, number> = {
        draft: 3,
        interim: 2,
        final: 1
      };

      // Sort the advisoryVersions array
      advisoryVersions.sort((a, b) => {
        // If versions are different, maintain original sort (or any default sort)
        if (a.version !== b.version) {
          return 0; // Keep original order for different versions
        }

        // If versions are the same, sort by tracking_status
        return statusOrder[a.tracking_status] - statusOrder[b.tracking_status];
      });

      return { advisoryVersions };
    }
  } else if (response.error) {
    return { error: getErrorDetails(`Could not load versions.`, response) };
  }
};

interface SearchMatch {
  path: string;
  positions: number[][];
}

const fetchSearchHits = async (id: number): Promise<SearchMatch[] | ErrorDetails> => {
  const query = `"${appStore.state.app.search.term}" search _clientSearch as`;
  const response = await request(
    `/api/documents/texts/${id}?query=${encodeURIComponent(query?.trim() ?? "")}`,
    "GET"
  );
  if (!response.ok) {
    return getErrorDetails("Could not load search matches.", response);
  }
  return response.content;
};

// Check if the loaded document (which is loaded via the internal ID) complies
// with the parameters in the URL.
const isResultConsistent = (params: any, document: any) => {
  return (
    params.trackingID &&
    params.publisherNamespace &&
    document.tracking.id === params.trackingID &&
    document.publisher.name === params.publisherNamespace
  );
};

type AdvisorySearchState = {
  scroll: boolean;
  searchMatches: SearchMatch[];
  matchIndex: number;
};

const advisorySearchState: AdvisorySearchState = $state({
  scroll: true,
  searchMatches: [],
  matchIndex: -1
});

const derivedHit = $derived(
  advisorySearchState.matchIndex !== -1
    ? advisorySearchState.searchMatches[advisorySearchState.matchIndex]
    : undefined
);

const getAdvisorySearchHit = () => {
  return derivedHit;
};

const doc = $derived(appStore.state.webview.doc);

const products = $derived(extractProducts(doc));

const productsByID = $derived.by(() => {
  return products.reduce((o: any, n: any) => {
    o[n.product_id] = n.name;
    return o;
  }, {});
});

const productVulnerabilities = $derived(
  generateProductVulnerabilities(doc, products, productsByID)
);

const getProductsByID = () => {
  return productsByID;
};

const getProductVulneravbilities = () => {
  return productVulnerabilities ?? [];
};

const sortedRevisionHistory = $derived(
  getRevisionHistory(doc)?.sort(
    (entry1: Revision2_0 | Revision2_1, entry2: Revision2_0 | Revision2_1) => {
      if (entry1.date < entry2.date) return 1;
      if (entry1.date > entry2.date) return -1;
      return 0;
    }
  ) ?? []
);

const getSortedRevisionHistory = () => sortedRevisionHistory;

const getAdvisoryLink = (doc: any) =>
  `/advisories/${doc.publisher}/${doc.tracking_id}/documents/${doc.id}`;
const getAdvisoryAnchorLink = (doc: any) => "#" + getAdvisoryLink(doc);

export {
  advisorySearchState,
  updateMultipleStates,
  loadAdvisoryVersions,
  fetchSearchHits,
  getAdvisoryAnchorLink,
  getAdvisorySearchHit,
  isResultConsistent,
  getProductsByID,
  getProductVulneravbilities,
  getSortedRevisionHistory
};
export type { AdvisoryVersion, SearchMatch };
