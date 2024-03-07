// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//

export const CSAFDocProps = {
  ACKNOWLEDGEMENTS: "acknowledgements",
  AGGREGATE_SEVERITY: "aggregate_severity",
  ALIASES: "aliases",
  CATEGORY: "category",
  CONTACT_DETAILS: "contact_details",
  CSAFVERSION: "csaf_version",
  CURRENTRELEASEDATE: "current_release_date",
  DISTRIBUTION: "distribution",
  DOCUMENT: "document",
  GENERATOR: "generator",
  ID: "id",
  INITIALRELEASEDATE: "initial_release_date",
  ISSUING_AUTHORITY: "issuing_authority",
  LABEL: "label",
  LANG: "lang",
  NOTES: "notes",
  PRODUCTTREE: "product_tree",
  PUBLISHER_CATEGORY: "category",
  PUBLISHER_NAME: "name",
  PUBLISHER_NAMESPACE: "namespace",
  PUBLISHER: "publisher",
  REFERENCES: "references",
  REVISIONHISTORY: "revision_history",
  STATUS: "status",
  SOURCELANG: "sourcelang",
  TITLE: "title",
  TLP: "tlp",
  TRACKING: "tracking",
  TRACKINGVERSION: "version",
  VULNERABILITIES: "vulnerabilities"
} as const;

export const TLP = {
  AMBER: "AMBER",
  GREEN: "GREEN",
  RED: "RED",
  WHITE: "WHITE",
  ERROR: "Invalid TLP"
} as const;

export const EMPTY = "";

export type TLPKeys = (typeof TLP)[keyof typeof TLP];

export const Status = {
  draft: "draft",
  final: "final",
  interim: "interim",
  ERROR: "Invalid Status"
} as const;
export type StatusKeys = (typeof Status)[keyof typeof Status];

export const DocumentCategory = {
  CSAF_SECURITY_ADVISORY: "csaf_security_advisory",
  CSAF_BASE: "csaf_base",
  CSAF_VEX: "csaf_vex"
} as const;

export type Publisher = {
  category: string;
  name: string;
  namespace: string;
  contact_details?: string;
  issuing_authority?: string;
};

export type RevisionHistoryEntry = {
  date: string;
  legacyVersion?: string;
  number: number;
  summary: string;
};

export type DocModel = {
  acknowledgements: Acknowledgement[];
  aggregateSeverity: AggregateSeverity | null;
  aliases: string[];
  category: string;
  csafVersion: string;
  generator: any;
  id: string;
  isDistributionPresent: boolean;
  isDocPresent: boolean;
  isProductTreePresent: boolean;
  isPublisherPresent: boolean;
  isRevisionHistoryPresent: boolean;
  isTLPPresent: boolean;
  isTrackingPresent: boolean;
  isVulnerabilitiesPresent: boolean;
  lang: string;
  lastUpdate: string;
  notes: Note[];
  productsByID: any;
  productTree: any;
  productVulnerabilities: any;
  published: string;
  publisher: Publisher;
  references: Reference[];
  revisionHistory: RevisionHistoryEntry[];
  status: string;
  sourceLang: string;
  title: string;
  tlp: TLP;
  trackingVersion: string;
  vulnerabilities: any;
};

export type DocModelKey = keyof DocModel;

export type TLP = {
  label: string;
  url?: string;
};

export type Note = {
  category: string;
  text: string;
  audience?: string;
  title?: string;
};

export type Acknowledgement = {
  names: string[];
  organization: string;
  summary: string;
  urls: string[];
};

export type Reference = {
  url: string;
  summary: string;
  category: string;
};

export type AggregateSeverity = {
  namespace?: string;
  text: string;
};

export type Engine = {
  name: string;
  version?: string;
};

export type Generator = {
  engine: Engine;
  date?: string;
};
