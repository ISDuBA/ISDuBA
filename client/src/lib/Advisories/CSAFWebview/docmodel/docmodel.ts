// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//

import {
  CSAFDocProps,
  EMPTY,
  Status,
  TLP,
  type AggregateSeverity,
  type DocModel,
  type Note,
  type Publisher,
  type Reference,
  type RevisionHistoryEntry
} from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
import {
  extractProducts,
  generateProductVulnerabilities
} from "../productvulnerabilities/productvulnerabilities";

/**
 * checkDocumentPresent checks whether the "document" property is present.
 * @param csafDoc
 * @returns true/false
 */
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
 * @param csafDoc
 * @returns true / false
 */
const checkVulnerabilities = (csafDoc: any): boolean => {
  return csafDoc[CSAFDocProps.VULNERABILITIES];
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

/**
 * getTitle retrieves title information.
 * @param csafDoc
 * @returns title | ""
 */
const getTitle = (csafDoc: any): string => {
  if (!checkDocumentPresent(csafDoc)) return EMPTY;
  return csafDoc.document[CSAFDocProps.TITLE] || EMPTY;
};

/**
 * getLanguage retrieves language information.
 * @param csafDoc
 * @returns language | ""
 */
const getLanguage = (csafDoc: any): string => {
  if (!checkDocumentPresent(csafDoc)) return EMPTY;
  return csafDoc.document[CSAFDocProps.LANG] || EMPTY;
};

/**
 * getCSAFVersion retrieves version information.
 * @param csafDoc
 * @returns version | ""
 */
const getCSAFVersion = (csafDoc: any): string => {
  if (!checkDocumentPresent(csafDoc)) return EMPTY;
  return csafDoc.document[CSAFDocProps.CSAFVERSION] || EMPTY;
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

/**
 * getTlp retrieves TLP information
 * @param csafDoc
 * @returns TLP | ""
 */
const getTlp = (csafDoc: any): TLP => {
  if (!checkTLPPresent(csafDoc)) return { label: "" };
  let label = "TLP.ERROR;";
  switch (csafDoc.document.distribution.tlp[CSAFDocProps.LABEL]) {
    case TLP.AMBER:
      label = TLP.AMBER;
      break;
    case TLP.GREEN:
      label = TLP.GREEN;
      break;
    case TLP.WHITE:
      label = TLP.WHITE;
      break;
    case TLP.RED:
      label = TLP.RED;
      break;
    default:
      label = TLP.ERROR;
      break;
  }
  return { label: label, url: csafDoc.document.distribution.tlp.url };
};

/**
 * getStatus retrieves the status of the document.
 * @param csafDoc
 * @returns status | ""
 */
const getStatus = (csafDoc: any): string => {
  if (!checkTrackingPresent(csafDoc)) return EMPTY;
  switch (csafDoc.document.tracking[CSAFDocProps.STATUS]) {
    case Status.draft:
      return Status.draft;
    case Status.final:
      return Status.final;
    case Status.interim:
      return Status.interim;
    default:
      return Status.ERROR;
  }
};

/**
 * getPublished retrieves the pubilshed info.
 * @param csafDoc
 * @returns info | ""
 */
const getPublished = (csafDoc: any): string => {
  if (!checkTrackingPresent(csafDoc)) return EMPTY;
  return csafDoc.document.tracking[CSAFDocProps.INITIALRELEASEDATE] || EMPTY;
};

/**
 * getLastUpdate retrieves the last update info.
 * @param csafDoc
 * @returns info | ""
 */
const getLastUpdate = (csafDoc: any): string => {
  if (!checkTrackingPresent(csafDoc)) return EMPTY;
  return csafDoc.document.tracking[CSAFDocProps.CURRENTRELEASEDATE] || EMPTY;
};

/**
 * getCategory retrieves the category info.
 * @param csafDoc
 * @returns info | ""
 */
const getCategory = (csafDoc: any): string => {
  if (!checkDocumentPresent(csafDoc)) return EMPTY;
  return csafDoc.document[CSAFDocProps.CATEGORY] || EMPTY;
};

/**
 * getPublisher retrieves publisher info.
 * @param csafDoc
 * @returns publisher info
 */
const getPublisher = (csafDoc: any): Publisher => {
  if (!checkPublisher(csafDoc)) {
    return {
      category: "",
      name: "",
      namespace: ""
    };
  }
  const publisher = csafDoc.document[CSAFDocProps.PUBLISHER];
  return {
    category: publisher[CSAFDocProps.PUBLISHER_CATEGORY],
    name: publisher[CSAFDocProps.PUBLISHER_NAME],
    namespace: publisher[CSAFDocProps.PUBLISHER_NAMESPACE],
    contact_details: publisher[CSAFDocProps.CONTACT_DETAILS],
    issuing_authority: publisher[CSAFDocProps.ISSUING_AUTHORITY]
  };
};

/**
 * getTrackingVersion retrieves tracking version.
 * @param csafDoc
 * @returns version | ""
 */
const getTrackingVersion = (csafDoc: any): string => {
  if (!checkTrackingPresent(csafDoc)) return EMPTY;
  return csafDoc.document.tracking[CSAFDocProps.TRACKINGVERSION] || EMPTY;
};

/**
 * getVulnerabilities retrieves the vulnerabilites section.
 * @param csafDoc
 * @returns vulnerabilities | []
 */
const getVulnerabilities = (csafDoc: any) => {
  if (!checkVulnerabilities(csafDoc)) return [];
  return csafDoc.vulnerabilities;
};

/**
 * getRevisionHistory retrieves revision history sorted by date.
 * @param csafDoc
 * @returns history | []
 */
const getRevisionHistory = (csafDoc: any): RevisionHistoryEntry[] => {
  if (!checkRevisionHistoryPresent(csafDoc)) return [];
  const result: RevisionHistoryEntry[] = csafDoc.document.tracking[CSAFDocProps.REVISIONHISTORY];
  result.sort((entry1: RevisionHistoryEntry, entry2: RevisionHistoryEntry) => {
    if (entry1.date < entry2.date) return 1;
    if (entry1.date > entry2.date) return -1;
    return 0;
  });
  return result;
};

/**
 * getProductTree retrieves the product tree.
 * @param csafDoc
 * @returns tree | []
 */
const getProductTree = (csafDoc: any) => {
  if (!checkproducTree(csafDoc)) return [];
  return csafDoc[CSAFDocProps.PRODUCTTREE];
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
 * getAcknowledgements retrieves ACKs.
 * @param csafDoc
 * @returns acks | []
 */
const getAcknowledgements = (csafDoc: any) => {
  if (!checkDocumentPresent(csafDoc)) return [];
  return csafDoc.document[CSAFDocProps.ACKNOWLEDGEMENTS];
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

/**
 * getReferences retrieves references.
 * @param csafDoc
 * @returns references | []
 */
const getReferences = (csafDoc: any): Reference[] => {
  if (!checkDocumentPresent(csafDoc)) return [];
  return csafDoc.document[CSAFDocProps.REFERENCES] || [];
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

/**
 * getGenerator retrieves generator info.
 * @param csafDoc
 * @returns generator info || null
 */
const getGenerator = (csafDoc: any) => {
  if (!checkTrackingPresent(csafDoc)) return null;
  return csafDoc.document.tracking[CSAFDocProps.GENERATOR] || null;
};

/**
 * getAliases retrieves aliases.
 * @param csafDoc
 * @returns aliases | null
 */
const getAliases = (csafDoc: any) => {
  if (!checkTrackingPresent(csafDoc)) return null;
  return csafDoc.document.tracking[CSAFDocProps.ALIASES] || null;
};

/**
 * convertToDocModel converts a CSAF document to a basic view model.
 * @param csafDoc
 * @returns DocModel
 */
const convertToDocModel = (csafDoc: any): DocModel => {
  const docModel: DocModel = {
    aggregateSeverity: getAggregateSeverity(csafDoc),
    acknowledgements: getAcknowledgements(csafDoc),
    aliases: getAliases(csafDoc),
    category: getCategory(csafDoc),
    csafVersion: getCSAFVersion(csafDoc),
    generator: getGenerator(csafDoc),
    id: getId(csafDoc),
    isDistributionPresent: checkDistributionPresent(csafDoc),
    isDocPresent: checkDocumentPresent(csafDoc),
    isProductTreePresent: checkproducTree(csafDoc),
    isPublisherPresent: checkPublisher(csafDoc),
    isRevisionHistoryPresent: checkRevisionHistoryPresent(csafDoc),
    isTLPPresent: checkTLPPresent(csafDoc),
    isTrackingPresent: checkTrackingPresent(csafDoc),
    isVulnerabilitiesPresent: checkVulnerabilities(csafDoc),
    lang: getLanguage(csafDoc),
    lastUpdate: getLastUpdate(csafDoc),
    notes: getNotes(csafDoc),
    productsByID: {},
    productTree: getProductTree(csafDoc),
    productVulnerabilities: [],
    published: getPublished(csafDoc),
    publisher: getPublisher(csafDoc),
    references: getReferences(csafDoc),
    revisionHistory: getRevisionHistory(csafDoc),
    status: getStatus(csafDoc),
    sourceLang: getSourceLang(csafDoc),
    title: getTitle(csafDoc),
    tlp: getTlp(csafDoc),
    trackingVersion: getTrackingVersion(csafDoc),
    vulnerabilities: getVulnerabilities(csafDoc)
  };
  const products = extractProducts(csafDoc);
  const productLookup = products.reduce((o: any, n: any) => {
    o[n.product_id] = n.name;
    return o;
  }, {});
  docModel.productsByID = productLookup;
  docModel.productVulnerabilities = generateProductVulnerabilities(
    csafDoc,
    products,
    productLookup
  );
  return docModel;
};

export { convertToDocModel };
