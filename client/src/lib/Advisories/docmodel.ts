// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//

import type {
  CSAFDocumentv2_0,
  CSAFVersion as CSAFVersion2_0,
  AggregateSeverity as AggregateSeverity2_0,
  DocumentLanguage,
  Title,
  TrafficLightProtocolTLP as TrafficLightProtocolTLP2_0,
  DocumentStatus,
  DocumentCategory,
  Note as Note2_0,
  Vulnerability as Vulnerability2_0,
  Vulnerabilities as Vulnerabilities2_0,
  Score,
  DocumentGenerator as DocumentGenerator2_0,
  ProductTree as ProductTree2_0,
  RevisionHistory as RevisionHistory2_0,
  DocumentReferences as DocumentReferences2_0
} from "./types/csaf-2.0";
import type {
  CSAFDocumentv2_1,
  CSAFVersion as CSAFVersion2_1,
  AggregateSeverity as AggregateSeverity2_1,
  TrafficLightProtocolTLP as TrafficLightProtocolTLP2_1,
  Vulnerability as Vulnerability2_1,
  Vulnerabilities as Vulnerabilities2_1,
  Metric,
  Note as Note2_1,
  CVSSv2,
  CVSSv3,
  CVSSv4,
  DocumentGenerator as DocumentGenerator2_1,
  ProductTree as ProductTree2_1,
  RevisionHistory as RevisionHistory2_1,
  DocumentReferences as DocumentReferences2_1,
  Version,
  SharingGroup,
  LicenseExpression,
  Contact,
  Involvement
} from "$lib/Advisories/types/csaf-2.1";

const EMPTY = "";

const isV2_1 = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1): boolean => {
  if (
    Object.keys(document).includes("document") &&
    (document as CSAFDocumentv2_1).document.csaf_version === "2.1"
  ) {
    return true;
  }
  return false;
};

const getTitle = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): Title | "" => {
  return document?.document.title || EMPTY;
};

const getLang = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): DocumentLanguage | "" => {
  return document?.document.lang || EMPTY;
};

const getCSAFVersion = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): CSAFVersion2_0 | CSAFVersion2_1 | "" => {
  return document?.document.csaf_version || EMPTY;
};

const getLicenseExpression = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): LicenseExpression | "" => {
  if (document?.document.csaf_version === "2.1") {
    return document.document.license_expression || EMPTY;
  }
  return EMPTY;
};

const getDistributionText = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): string => {
  return document?.document.distribution?.text || EMPTY;
};

const getDistributionSharingGroup = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): SharingGroup | undefined => {
  if (document?.document.csaf_version === "2.1") {
    return document.document.distribution?.sharing_group;
  }
  return undefined;
};

const getID = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): string => {
  return document?.document.tracking.id || EMPTY;
};

const getTLP = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): TrafficLightProtocolTLP2_0 | TrafficLightProtocolTLP2_1 | undefined => {
  return document?.document.distribution?.tlp;
};

const getStatus = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): DocumentStatus | "" => {
  return document?.document.tracking.status ?? EMPTY;
};

const getTrackingVersion = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): Version | string => {
  return document?.document.tracking.version ?? EMPTY;
};

const getInitialReleaseDate = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1) => {
  return document.document.tracking.initial_release_date;
};

const getCurrentReleaseDate = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1) => {
  return document.document.tracking.current_release_date;
};

const getCategory = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): DocumentCategory => {
  return document?.document.category || EMPTY;
};

const getPublisher = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null) => {
  return document?.document.publisher;
};

const getPublisherContact = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): Contact | undefined => {
  if (document?.document.csaf_version === "2.1") {
    return document?.document.publisher.contact;
  }
  return undefined;
};

const getVulnerabilities = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1
): Vulnerabilities2_0 | Vulnerabilities2_1 | [] => {
  return document.vulnerabilities ?? [];
};

const getCVSSOfVulnerability = (
  vulnerability: Vulnerability2_0 | Vulnerability2_1,
  cvssVersion: 2 | 3 | 4,
  csafVersion: "2.0" | "2.1"
): Array<CVSSv2 | CVSSv3 | CVSSv4> => {
  const listOfCVSS: Array<CVSSv2 | CVSSv3 | CVSSv4> = [];
  if (csafVersion === "2.0" && (vulnerability as Vulnerability2_0).scores) {
    (vulnerability as Vulnerability2_0).scores?.forEach((score: Score) => {
      if (cvssVersion === 2 && score.cvss_v2) listOfCVSS.push(score.cvss_v2);
      if (cvssVersion === 3 && score.cvss_v3) listOfCVSS.push(score.cvss_v3);
    });
  } else if (csafVersion === "2.1" && (vulnerability as Vulnerability2_1).metrics) {
    (vulnerability as Vulnerability2_1).metrics?.forEach((metric: Metric) => {
      if (cvssVersion === 2 && metric.content.cvss_v2) listOfCVSS.push(metric.content.cvss_v2);
      if (cvssVersion === 3 && metric.content.cvss_v3) listOfCVSS.push(metric.content.cvss_v3);
      if (cvssVersion === 4 && metric.content.cvss_v4) listOfCVSS.push(metric.content.cvss_v4);
    });
  }
  return listOfCVSS;
};

/**
 * Retrieves the CVSS object with the highest base score,
 * prioritizing CVSS v3 over v2 across the entire document.
 * If any CVSS v3 score exists, all CVSS v2 scores are ignored.
 */
const getHighestScore = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): CVSSv2 | CVSSv3 | CVSSv4 | null => {
  if (!document?.vulnerabilities?.length) {
    return null;
  }

  let hasCvssV4 = false;
  let hasCvssV3 = false;
  let highestScoreObject: CVSSv2 | CVSSv3 | CVSSv4 | null = null;
  let highestBaseScore = -1;

  // First pass: Check if any CVSS v3 or v4 scores exist
  if (document.document.csaf_version === "2.1") {
    for (const vulnerability of document.vulnerabilities as Vulnerabilities2_1) {
      if (vulnerability.metrics?.some((score: Metric) => score.content.cvss_v4)) {
        hasCvssV4 = true;
        break;
      }
    }
  }
  for (const vulnerability of document.vulnerabilities) {
    if (document.document.csaf_version === "2.1") {
      if (
        (vulnerability as Vulnerability2_1).metrics?.some(
          (metric: Metric) => metric.content.cvss_v3
        )
      ) {
        hasCvssV3 = true;
        break;
      }
    } else if (document.document.csaf_version === "2.0") {
      if ((vulnerability as Vulnerability2_0).scores?.some((score: Score) => score.cvss_v3)) {
        hasCvssV3 = true;
        break;
      }
    }
  }

  // Second pass: Find the highest score based on the first pass's result
  for (const vulnerability of document.vulnerabilities) {
    if (
      (document.document.csaf_version === "2.0" && !(vulnerability as Vulnerability2_0).scores) ||
      (document.document.csaf_version === "2.1" && !(vulnerability as Vulnerability2_1).metrics)
    ) {
      continue;
    }

    let cvssList: Array<CVSSv2 | CVSSv3 | CVSSv4>;
    // Consider only the scores with the highest CVSS version
    if (hasCvssV4) {
      cvssList = getCVSSOfVulnerability(vulnerability, 4, document.document.csaf_version);
    } else if (hasCvssV3) {
      cvssList = getCVSSOfVulnerability(vulnerability, 3, document.document.csaf_version);
    } else {
      cvssList = getCVSSOfVulnerability(vulnerability, 2, document.document.csaf_version);
    }
    cvssList.forEach((cvss: CVSSv2 | CVSSv3 | CVSSv4) => {
      if (cvss.baseScore && cvss.baseScore > highestBaseScore) {
        highestBaseScore = cvss.baseScore;
        highestScoreObject = cvss;
      }
    });
  }

  return highestScoreObject;
};

const getRevisionHistory = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): RevisionHistory2_0 | RevisionHistory2_1 | null => {
  return document?.document.tracking.revision_history ?? null;
};

const getProductTree = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): ProductTree2_0 | ProductTree2_1 | undefined => {
  return document?.product_tree;
};

const getNotes = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): Note2_0[] | Note2_1[] | undefined => {
  return document?.document.notes;
};

const getAcknowledgments = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null) => {
  return document?.document.acknowledgments;
};

const getSourceLang = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): string => {
  return document?.document.source_lang || EMPTY;
};

const getReferences = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): DocumentReferences2_0 | DocumentReferences2_1 | [] => {
  return document?.document.references || [];
};

const getAggregateSeverity = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): AggregateSeverity2_0 | AggregateSeverity2_1 | null => {
  return document?.document.aggregate_severity || null;
};

const getGenerator = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): DocumentGenerator2_0 | DocumentGenerator2_1 | null => {
  return document?.document.tracking.generator || null;
};

const getAliases = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null) => {
  return document?.document.tracking.aliases || null;
};

const getInvolvement = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): Involvement | undefined => {
  if (document?.document.csaf_version === "2.1") {
    return document?.document.involvement;
  }
  return undefined;
};

export {
  isV2_1,
  getAcknowledgments,
  getAggregateSeverity,
  getAliases,
  getCategory,
  getCSAFVersion,
  getLicenseExpression,
  getDistributionSharingGroup,
  getDistributionText,
  getID,
  getInitialReleaseDate,
  getCurrentReleaseDate,
  getGenerator,
  getHighestScore,
  getInvolvement,
  getLang,
  getNotes,
  getProductTree,
  getPublisher,
  getPublisherContact,
  getReferences,
  getRevisionHistory,
  getSourceLang,
  getStatus,
  getVulnerabilities,
  getTLP,
  getTitle,
  getTrackingVersion
};
