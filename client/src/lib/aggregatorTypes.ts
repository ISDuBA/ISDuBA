// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

/**
 * Provides information about the category of aggregator.
 */
export type CategoryOfAggregator = "aggregator" | "lister";
/**
 * Information on how to contact the aggregator, possibly including details such as web sites, email addresses, phone numbers, and postal mail addresses.
 */
export type ContactDetails = string;
/**
 * Provides information about the authority of the aggregator to release the list, in particular, the party's constituency and responsibilities or other obligations.
 */
export type IssuingAuthority = string;
/**
 * Contains the name of the aggregator.
 */
export type NameOfAggregator = string;
/**
 * Contains a URL which is under control of the aggregator and can be used as a globally unique identifier for that aggregator.
 */
export type NamespaceOfAggregator = string;
/**
 * Gives the version of the CSAF aggregator specification which the document was generated for.
 */
export type CSAFAggregatorVersion = "2.0";
/**
 * Contains the URL for this document.
 */
export type CanonicalURL = string;
/**
 * Contains a list with information from CSAF providers.
 *
 * @minItems 1
 */
export type ListOfCSAFProviders = [CSAFProviderEntry, ...CSAFProviderEntry[]];
/**
 * Holds the date and time when this entry was last updated.
 */
export type LastUpdated = string;
/**
 * Provides information about the category of publisher releasing the document.
 */
export type CategoryOfPublisher =
  | "coordinator"
  | "discoverer"
  | "other"
  | "translator"
  | "user"
  | "vendor";
/**
 * Information on how to contact the publisher, possibly including details such as web sites, email addresses, phone numbers, and postal mail addresses.
 */
export type ContactDetails1 = string;
/**
 * Provides information about the authority of the issuing party to release the document, in particular, the party's constituency and responsibilities or other obligations.
 */
export type IssuingAuthority1 = string;
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
 * Contains the URL of the provider-metadata.json for that entry.
 */
export type URLOfTheMetadata = string;
/**
 * Contains a list of URLs or mirrors for this CSAF provider.
 *
 * @minItems 1
 */
export type ListOfMirrors = [Mirror, ...Mirror[]];
/**
 * Contains the base URL of the mirror for this issuing party.
 */
export type Mirror = string;
/**
 * Contains a list with information from CSAF publishers.
 *
 * @minItems 1
 */
export type ListOfCSAFPublishers = [CSAFPublisherEntry, ...CSAFPublisherEntry[]];
/**
 * Contains a list of URLs or mirrors for this CSAF publisher.
 *
 * @minItems 1
 */
export type ListOfMirrors1 = [Mirror, ...Mirror[]];
/**
 * Contains information about how often the CSAF publisher is checked for new CSAF documents.
 */
export type UpdateInterval = string;
/**
 * Holds the date and time when the document was last updated.
 */
export type LastUpdated1 = string;

/**
 * Representation of information where to find CSAF providers as a JSON document.
 */
export interface CSAFAggregator {
  aggregator: Aggregator;
  aggregator_version: CSAFAggregatorVersion;
  canonical_url: CanonicalURL;
  csaf_providers: ListOfCSAFProviders;
  csaf_publishers?: ListOfCSAFPublishers;
  last_updated: LastUpdated1;
  [k: string]: unknown;
}
/**
 * Provides information about the aggregator.
 */
export interface Aggregator {
  category: CategoryOfAggregator;
  contact_details?: ContactDetails;
  issuing_authority?: IssuingAuthority;
  name: NameOfAggregator;
  namespace: NamespaceOfAggregator;
  [k: string]: unknown;
}
/**
 * Contains information from a CSAF provider.
 */
export interface CSAFProviderEntry {
  metadata: CSAFProviderMetadata;
  mirrors?: ListOfMirrors;
  [k: string]: unknown;
}
/**
 * Contains the metadata of a single CSAF provider.
 */
export interface CSAFProviderMetadata {
  last_updated: LastUpdated;
  publisher: Publisher;
  role?: RoleOfTheIssuingParty;
  url: URLOfTheMetadata;
  [k: string]: unknown;
}
/**
 * Provides information about the issuing party for this entry.
 */
export interface Publisher {
  category: CategoryOfPublisher;
  contact_details?: ContactDetails1;
  issuing_authority?: IssuingAuthority1;
  name: NameOfPublisher;
  namespace: NamespaceOfPublisher;
  [k: string]: unknown;
}
/**
 * Contains information from a CSAF publisher.
 */
export interface CSAFPublisherEntry {
  metadata: CSAFPublisherMetadata;
  mirrors?: ListOfMirrors1;
  update_interval: UpdateInterval;
  [k: string]: unknown;
}
/**
 * Contains the metadata of a single CSAF publisher extracted from one of its CSAF documents.
 */
export interface CSAFPublisherMetadata {
  last_updated: LastUpdated;
  publisher: Publisher;
  role?: RoleOfTheIssuingParty;
  url: URLOfTheMetadata;
  [k: string]: unknown;
}

/**
 * Contains information about the subscribed feed.
 */
export interface FeedSubscription {
  id: number;
  url: string;
}

/**
 * Contains information about all subscribed sources.
 */
export interface SourceSubscription {
  url: string;
  id: number;
  name: string;
  subscripted?: FeedSubscription[];
}

/**
 * Contains information about all subscriptions.
 */
export interface Subscription {
  url: string;
  available?: string[];
  subscriptions?: SourceSubscription[];
}

/**
 * Contains additional metadata about the aggregator.
 */
export interface Custom {
  id: number;
  name: string;
  attention?: boolean;
  subscriptions: Subscription[];
}

/**
 * Contains the json document of the Aggregator with additional metadata
 */
export interface AggregatorMetadata {
  aggregator: CSAFAggregator;
  custom: Custom;
}
