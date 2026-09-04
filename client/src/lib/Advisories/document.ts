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
  if (
    Object.keys(document).includes("document") &&
    (document as CSAFDocumentv2_1).document.csaf_version === "2.1"
  ) {
    return true;
  }
  return false;
};

const getTrackingVersion = (document: DocModel | CSAFDocumentv2_1): Version | string => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.tracking.version;
  } else {
    return (document as DocModel).trackingVersion;
  }
};

const getGenerator = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.tracking.generator;
  } else {
    return (document as DocModel).generator;
  }
};

const getTLP = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.distribution.tlp;
  } else {
    return (document as DocModel).tlp;
  }
};

const getPublisher = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.publisher;
  } else {
    return (document as DocModel).publisher;
  }
};

const getTitle = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.title;
  } else {
    return (document as DocModel).title;
  }
};

const getCategory = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.category;
  } else {
    return (document as DocModel).category;
  }
};

const getReferences = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).document.references;
  } else {
    return (document as DocModel).references;
  }
};

const getProductTree = (document: DocModel | CSAFDocumentv2_1) => {
  if (isV2_1(document)) {
    return (document as CSAFDocumentv2_1).product_tree;
  } else {
    return (document as DocModel).productTree;
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

export {
  isV2_1,
  fetchDocumentSSVC,
  getTrackingVersion,
  getGenerator,
  getTLP,
  getPublisher,
  getTitle,
  getCategory,
  getReferences,
  getProductTree
};
