/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
 */

/* eslint-disable */
/**
* This file was automatically generated with the repo at https://heptapod.host/intevation/csaf-types.


/**
 * Contains the URL of the provider-metadata.json JSON schema which the document promises to be valid for.
 */
export type JSONSchema = "https://docs.oasis-open.org/csaf/csaf/v2.1/schema/provider.json";
/**
 * Contains the URL for this document.
 */
export type CanonicalURL = string;
/**
 * Contains a list of used distribution mechanisms.
 *
 * @minItems 1
 */
export type ListOfDistribution = [Distribution, ...Distribution[]];
/**
 * Provides the TLP label for the directory.
 */
export type TLPLabel = "AMBER" | "AMBER+STRICT" | "CLEAR" | "GREEN" | "RED";
/**
 * Contains the base url for the directory-based distribution.
 */
export type DirectoryURL = string;
/**
 * Contains a list of URLs which contain ROLIE category documents.
 *
 * @minItems 1
 */
export type ListOfROLIECategoryDocumentURLs = [ROLIECategoryDocumentURL, ...ROLIECategoryDocumentURL[]];
/**
 * Contains a URL of a ROLIE category document.
 */
export type ROLIECategoryDocumentURL = string;
/**
 * Contains a list of information about ROLIE feeds.
 *
 * @minItems 1
 */
export type ListOfROLIEFeeds = [ROLIEFeed, ...ROLIEFeed[]];
/**
 * Holds the date and time when the feed was last updated.
 */
export type LastUpdated = string;
/**
 * Contains a summary of the feed.
 */
export type SummaryOfTheFeed = string;
/**
 * Provides the TLP label for the feed.
 */
export type TLPLabel1 = "AMBER" | "AMBER+STRICT" | "CLEAR" | "GREEN" | "RED";
/**
 * Contains the URL of the feed.
 */
export type URLOfTheFeed = string;
/**
 * Contains a list of URLs which contain ROLIE service documents.
 *
 * @minItems 1
 */
export type ListOfROLIEServiceDocumentURLs = [ROLIEServiceDocumentURL, ...ROLIEServiceDocumentURL[]];
/**
 * Contains a URL of a ROLIE service document.
 */
export type ROLIEServiceDocumentURL = string;
/**
 * Holds the date and time when the document was last updated.
 */
export type LastUpdated1 = string;
/**
 * Decides whether this file should be linked in the list of a CSAF aggregator.
 */
export type ListOnCSAFAggregators = boolean;
/**
 * Holds the date and time from when the distributions within this CSAF provider are in a maintained state.
 */
export type MaintainedFrom = string;
/**
 * Holds the date and time until when the distributions within this CSAF provider are in a maintained state.
 */
export type MaintainedUntil = string;
/**
 * Gives the version of the CSAF provider metadata specification which the document was generated for.
 */
export type CSAFProviderMetadataVersion = "2.1";
/**
 * Decides whether the CSAF documents can be mirrored and provided by a CSAF aggregator.
 */
export type MirrorOnCSAFAggregators = boolean;
/**
 * Contains the fingerprint of the OpenPGP key.
 */
export type FingerprintOfTheKey = string;
/**
 * Contains the URL where the key can be retrieved.
 */
export type URLOfTheKey = string;
/**
 * Contains a list of OpenPGP keys used to sign CSAF documents.
 */
export type ListOfPublicOpenPGPKeys = PGPKeys[];
/**
 * Provides information about the category of publisher releasing the document.
 */
export type CategoryOfPublisher =
  "coordinator" | "discoverer" | "multiplier" | "other" | "translator" | "user" | "vendor";
/**
 * Contains details regarding ways to reach the party, e.g. through web sites, phone numbers, and postal mail addresses.
 */
export type ContactDetails = string;
/**
 * Contains the email address that can be used to reach the party.
 */
export type Email = string;
/**
 * Contains a URL pointing to a public OpenPGP key valid for the email of party provided in the sibling property `email`.
 */
export type PublicOpenPGPKeyURL = string;
/**
 * Contains a URL that can be used to reach the party.
 */
export type ContactURL = string;
/**
 * Provides information about the authority of the issuing party to release the document, in particular, the party's constituency and responsibilities or other obligations.
 */
export type IssuingAuthority = string;
/**
 * Contains the name of the issuing party.
 */
export type NameOfPublisher = string;
/**
 * Contains a URL which is under control of the issuing party and can be used as a globally unique identifier for that issuing party.
 */
export type NamespaceOfPublisher = string;
/**
 * Contains the role of the issuing party according to section 7 in the CSAF standard.
 */
export type RoleOfTheIssuingParty = "csaf_publisher" | "csaf_provider" | "csaf_trusted_provider";

/**
 * Representation of metadata information of a CSAF provider as a JSON document.
 */
export interface CSAFProviderMetadata {
  $schema: JSONSchema;
  canonical_url: CanonicalURL;
  distributions?: ListOfDistribution;
  last_updated: LastUpdated1;
  list_on_CSAF_aggregators: ListOnCSAFAggregators;
  maintained_from?: MaintainedFrom;
  maintained_until?: MaintainedUntil;
  metadata_version: CSAFProviderMetadataVersion;
  mirror_on_CSAF_aggregators: MirrorOnCSAFAggregators;
  public_openpgp_keys?: ListOfPublicOpenPGPKeys;
  publisher: Publisher;
  role: RoleOfTheIssuingParty;
}
/**
 * Contains the information of a used distribution mechanism.
 */
export interface Distribution {
  directory?: Directory;
  rolie?: ROLIE;
}
/**
 * Contains all information for directory-based distribution.
 */
export interface Directory {
  tlp_label: TLPLabel;
  url: DirectoryURL;
}
/**
 * Contains all information for ROLIE distribution.
 */
export interface ROLIE {
  categories?: ListOfROLIECategoryDocumentURLs;
  feeds: ListOfROLIEFeeds;
  services?: ListOfROLIEServiceDocumentURLs;
}
/**
 * Contains information about the ROLIE feed.
 */
export interface ROLIEFeed {
  last_updated: LastUpdated;
  summary?: SummaryOfTheFeed;
  tlp_label: TLPLabel1;
  url: URLOfTheFeed;
}
/**
 * Contains all information about an OpenPGP key used to sign CSAF documents.
 */
export interface PGPKeys {
  fingerprint: FingerprintOfTheKey;
  url: URLOfTheKey;
}
/**
 * Provides information about the publisher of the CSAF documents in this repository.
 */
export interface Publisher {
  category: CategoryOfPublisher;
  contact?: Contact;
  issuing_authority?: IssuingAuthority;
  name: NameOfPublisher;
  namespace: NamespaceOfPublisher;
}
/**
 * Contains information on how to contact the publisher.
 */
export interface Contact {
  details?: ContactDetails;
  email?: Email;
  public_openpgp_key_url?: PublicOpenPGPKeyURL;
  url?: ContactURL;
}
