// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import type { ErrorDetails } from "$lib/Errors/error";
import { request } from "$lib/request";
import { getErrorDetails } from "$lib/Errors/error";
import type { DocModel } from "./types/docmodeltypes";
import type { CSAFDocumentv2_1, Version } from "./types/csaf-2.1";

const isV2_1 = (document: DocModel | CSAFDocumentv2_1): boolean => {
  return Object.keys(document).includes("$schema");
};

const getTrackingVersion = (document: DocModel | CSAFDocumentv2_1): Version | string => {
  if (isV2_1(document)) {
    return (document as DocModel).trackingVersion;
  } else {
    return (document as CSAFDocumentv2_1).document.tracking.version;
  }
};

const fetchDocumentSSVC = async (
  documentId: string | number,
  abortController?: AbortController
): Promise<string | ErrorDetails | undefined> => {
  const response = await request(
    `/api/ssvc/documents/${documentId}`,
    "GET",
    undefined,
    abortController
  );

  // Any error
  if (!response.ok) {
    if (response.error !== "AbortError") {
      return getErrorDetails("Could not load SSVC.", response);
    } else {
      return undefined;
    }
  }

  const result = await response.content;

  // got a non-empty result
  if (result && typeof result.ssvc === "string" && result.ssvc !== "") {
    return result.ssvc;
  }
  // no SSVC
  return undefined;
};

export { fetchDocumentSSVC, getTrackingVersion };
