/* eslint-disable */
/**
* This file was automatically generated with the repo at https://heptapod.host/intevation/csaf-types.


/**
 * Contains the URL of the Aggregator JSON schema which the document promises to be valid for.
 */
export type JSONSchema = "https://docs.oasis-open.org/csaf/csaf/v2.1/schema/aggregator.json";
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
export type CSAFAggregatorVersion = "2.1";
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
  "coordinator" | "discoverer" | "multiplier" | "other" | "translator" | "user" | "vendor";
/**
 * Contains details regarding ways to reach the party, e.g. through web sites, phone numbers, and postal mail addresses.
 */
export type ContactDetails1 = string;
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
 * Holds the date and time from when this CSAF aggregator file is in a maintained state.
 */
export type MaintainedFrom = string;
/**
 * Holds the date and time until when this CSAF aggregator file is in a maintained state.
 */
export type MaintainedUntil = string;

/**
 * Representation of information where to find CSAF providers as a JSON document.
 */
export interface CSAFAggregator {
  $schema: JSONSchema;
  aggregator: Aggregator;
  aggregator_version: CSAFAggregatorVersion;
  canonical_url: CanonicalURL;
  csaf_providers: ListOfCSAFProviders;
  csaf_publishers?: ListOfCSAFPublishers;
  last_updated: LastUpdated1;
  maintained_from?: MaintainedFrom;
  maintained_until?: MaintainedUntil;
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
}
/**
 * Contains information from a CSAF provider.
 */
export interface CSAFProviderEntry {
  metadata: CSAFProviderMetadata;
  mirrors?: ListOfMirrors;
}
/**
 * Contains the metadata of a single CSAF provider.
 */
export interface CSAFProviderMetadata {
  last_updated: LastUpdated;
  publisher: Publisher;
  role: RoleOfTheIssuingParty;
  url: URLOfTheMetadata;
}
/**
 * Provides information about the issuing party for this entry.
 */
export interface Publisher {
  category: CategoryOfPublisher;
  contact?: Contact;
  issuing_authority?: IssuingAuthority1;
  name: NameOfPublisher;
  namespace: NamespaceOfPublisher;
}
/**
 * Contains information on how to contact the publisher.
 */
export interface Contact {
  details?: ContactDetails1;
  email?: Email;
  public_openpgp_key_url?: PublicOpenPGPKeyURL;
  url?: ContactURL;
}
/**
 * Contains information from a CSAF publisher.
 */
export interface CSAFPublisherEntry {
  metadata: CSAFPublisherMetadata;
  mirrors: ListOfMirrors1;
  update_interval: UpdateInterval;
}
/**
 * Contains the metadata of a single CSAF publisher extracted from one of its CSAF documents.
 */
export interface CSAFPublisherMetadata {
  last_updated: LastUpdated;
  publisher: Publisher;
  role: RoleOfTheIssuingParty;
  url: URLOfTheMetadata;
}
