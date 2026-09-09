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
  DocumentLanguage,
  Title,
  TrafficLightProtocolTLP as TrafficLightProtocolTLP2_0,
  DocumentStatus,
  DocumentCategory,
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
  TrafficLightProtocolTLP as TrafficLightProtocolTLP2_1,
  Vulnerability as Vulnerability2_1,
  Vulnerabilities as Vulnerabilities2_1,
  Metric,
  CVSSv2,
  CVSSv3,
  CVSSv4,
  DocumentGenerator as DocumentGenerator2_1,
  ProductTree as ProductTree2_1,
  RevisionHistory as RevisionHistory2_1,
  DocumentReferences as DocumentReferences2_1,
  Version
} from "$lib/Advisories/types/csaf-2.1";
import {
  CSAFDocProps,
  EMPTY,
  type AggregateSeverity,
  type Note
} from "$lib/Advisories/types/docmodeltypes";

const isV2_1 = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1): boolean => {
  if (
    Object.keys(document).includes("document") &&
    (document as CSAFDocumentv2_1).document.csaf_version === "2.1"
  ) {
    return true;
  }
  return false;
};

const checkDocumentPresent = (csafDoc: any): boolean => {
  return csafDoc[CSAFDocProps.DOCUMENT];
};

/**
 * checkTrackingPresent checks whether the "tracking" property is present.
 * @param csafDoc
 * @returns true/false
 */
const checkTrackingPresent = (csafDoc: any): boolean => {
  return checkDocumentPresent(csafDoc) && csafDoc.document[CSAFDocProps.TRACKING];
};

/**
 * checkDistributionPresent checks whether the "distribution" property is present.
 * @param csafDoc
 * @returns true/false
 */
const checkDistributionPresent = (csafDoc: any): boolean => {
  return checkDocumentPresent(csafDoc) && csafDoc.document[CSAFDocProps.DISTRIBUTION];
};

/**
 * checkTLPPresent checks whether the "TLP" property is present.
 * @param csafDoc
 * @returns true/false
 */
const checkTLPPresent = (csafDoc: any): boolean => {
  return (
    checkDistributionPresent(csafDoc) &&
    csafDoc.document.distribution[CSAFDocProps.TLP] &&
    csafDoc.document.distribution[CSAFDocProps.TLP][CSAFDocProps.LABEL]
  );
};

/**
 * checkPublisher checks whether the "Publisher" property is present.
 * @param csafDoc
 * @returns true / false
 */
const checkPublisher = (csafDoc: any): boolean => {
  return checkDocumentPresent(csafDoc) && csafDoc.document[CSAFDocProps.PUBLISHER];
};

/**
 * checkVulnerabilities checks whether the "vulnerabitlites" section is present.
 * @param document
 * @returns true / false
 */
const checkVulnerabilities = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1): boolean => {
  return document[CSAFDocProps.VULNERABILITIES] !== undefined;
};

/**
 * checkproducTree checks whether the "product tree" section is present.
 * @param csafDoc
 * @returns true / false
 */
const checkproducTree = (csafDoc: any): boolean => {
  return csafDoc[CSAFDocProps.PRODUCTTREE];
};

/**
 * checkRevisionHistoryPresent checks whether the "revision history" section is present.
 * @param csafDoc
 * @returns true / false
 */
const checkRevisionHistoryPresent = (csafDoc: any): boolean => {
  return checkTrackingPresent(csafDoc) && csafDoc.document.tracking[CSAFDocProps.REVISIONHISTORY];
};

const getTitle = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): Title | "" => {
  return document?.document.title || EMPTY;
};

const getLanguage = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1): DocumentLanguage | "" => {
  if (!checkDocumentPresent(document)) return EMPTY;
  return document.document.lang || EMPTY;
};

const getCSAFVersion = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): CSAFVersion2_0 | CSAFVersion2_1 | "" => {
  return document?.document.csaf_version || EMPTY;
};

const getDistributionText = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null): string => {
  return document?.document.distribution?.text || EMPTY;
};

/**
 * getId retrieves a document ID.
 * @param csafDoc
 * @returns id | ""
 */
const getId = (csafDoc: any): string => {
  if (!checkTrackingPresent(csafDoc)) return EMPTY;
  return csafDoc.document.tracking[CSAFDocProps.ID] || EMPTY;
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
  if (!checkDocumentPresent(document)) return EMPTY;
  return document?.document.category || EMPTY;
};

const getPublisher = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null) => {
  return document?.document.publisher;
};

const getVulnerabilities = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1
): Vulnerabilities2_0 | Vulnerabilities2_1 | [] => {
  if (!checkVulnerabilities(document) || document.vulnerabilities === undefined) return [];
  return document.vulnerabilities;
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

/**
 * getNotes retrieves notes.
 * @param csafDoc
 * @returns notes | []
 */
const getNotes = (csafDoc: any): Note[] => {
  if (!checkDocumentPresent(csafDoc)) return [];
  return csafDoc.document[CSAFDocProps.NOTES];
};

/**
 * getAcknowledgments retrieves ACKs.
 * @param csafDoc
 * @returns acks | []
 */
const getAcknowledgments = (csafDoc: any) => {
  if (!checkDocumentPresent(csafDoc)) return [];
  return csafDoc.document[CSAFDocProps.ACKNOWLEDGMENTS];
};

/**
 * getSourceLang retrieves the source language.
 * @param csafDoc
 * @returns lang | ""
 */
const getSourceLang = (csafDoc: any): string => {
  if (!checkDocumentPresent(csafDoc)) return EMPTY;
  return csafDoc.document[CSAFDocProps.SOURCELANG] || EMPTY;
};

const getReferences = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): DocumentReferences2_0 | DocumentReferences2_1 | [] => {
  return document?.document.references || [];
};

/**
 * getAggregateSeverity retrieves the aggregate severity info.
 * @param csafDoc
 * @returns info | null
 */
const getAggregateSeverity = (csafDoc: any): AggregateSeverity | null => {
  if (!checkDocumentPresent(csafDoc)) return null;
  return csafDoc.document[CSAFDocProps.AGGREGATE_SEVERITY] || null;
};

const getGenerator = (
  document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null
): DocumentGenerator2_0 | DocumentGenerator2_1 | null => {
  return document?.document.tracking.generator || null;
};

const getAliases = (document: CSAFDocumentv2_0 | CSAFDocumentv2_1 | null) => {
  return document?.document.tracking.aliases || null;
};

export {
  isV2_1,
  getAliases,
  getCategory,
  getCSAFVersion,
  getDistributionText,
  getInitialReleaseDate,
  getCurrentReleaseDate,
  getGenerator,
  getHighestScore,
  getProductTree,
  getPublisher,
  getReferences,
  getRevisionHistory,
  getStatus,
  getVulnerabilities,
  getTLP,
  getTitle,
  getTrackingVersion
};
