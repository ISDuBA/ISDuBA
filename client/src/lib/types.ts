// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

/**
 * Contains a list of acknowledgment elements associated with the whole document.
 */
export type DocumentAcknowledgments = [Acknowledgment, ...Acknowledgment[]];
/**
 * Contains the names of contributors being recognized.
 */
export type ListOfAcknowledgedNames = [NameOfTheContributor, ...NameOfTheContributor[]];
/**
 * Contains the name of a single contributor being recognized.
 */
export type NameOfTheContributor = string;
/**
 * Contains the name of a contributing organization being recognized.
 */
export type ContributingOrganization = string;
/**
 * SHOULD represent any contextual details the document producers wish to make known about the acknowledgment or acknowledged parties.
 */
export type SummaryOfTheAcknowledgment = string;
/**
 * Specifies a list of URLs or location of the reference to be acknowledged.
 */
export type ListOfURLs = [URLOfAcknowledgment, ...URLOfAcknowledgment[]];
/**
 * Contains the URL or location of the reference to be acknowledged.
 */
export type URLOfAcknowledgment = string;
/**
 * Points to the namespace so referenced.
 */
export type NamespaceOfAggregateSeverity = string;
/**
 * Provides a severity which is independent of - and in addition to - any other standard metric for determining the impact or severity of a given vulnerability (such as CVSS).
 */
export type TextOfAggregateSeverity = string;
/**
 * Defines a short canonical name, chosen by the document producer, which will inform the end user as to the category of document.
 */
export type DocumentCategory = string;
/**
 * Gives the version of the CSAF specification which the document was generated for.
 */
export type CSAFVersion = "2.0";
/**
 * Provides a textual description of additional constraints.
 */
export type TextualDescription = string;
/**
 * Provides the TLP label of the document.
 */
export type LabelOfTLP = "AMBER" | "GREEN" | "RED" | "WHITE";
/**
 * Provides a URL where to find the textual description of the TLP version which is used in this document. Default is the URL to the definition by FIRST.
 */
export type URLOfTLPVersion = string;
/**
 * Identifies the language used by this document, corresponding to IETF BCP 47 / RFC 5646.
 */
export type DocumentLanguage = string;
/**
 * Holds notes associated with the whole document.
 */
export type DocumentNotes = [Note, ...Note[]];
/**
 * Indicates who is intended to read it.
 */
export type AudienceOfNote = string;
/**
 * Contains the information of what kind of note this is.
 */
export type NoteCategory =
  | "description"
  | "details"
  | "faq"
  | "general"
  | "legal_disclaimer"
  | "other"
  | "summary";
/**
 * Holds the content of the note. Content varies depending on type.
 */
export type NoteContent = string;
/**
 * Provides a concise description of what is contained in the text of the note.
 */
export type TitleOfNote = string;
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
export type ContactDetails = string;
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
 * Holds a list of references associated with the whole document.
 */
export type DocumentReferences = [Reference, ...Reference[]];
/**
 * Indicates whether the reference points to the same document or vulnerability in focus (depending on scope) or to an external resource.
 */
export type CategoryOfReference = "external" | "self";
/**
 * Indicates what this reference refers to.
 */
export type SummaryOfTheReference = string;
/**
 * Provides the URL for the reference.
 */
export type URLOfReference = string;
/**
 * If this copy of the document is a translation then the value of this property describes from which language this document was translated.
 */
export type SourceLanguage = string;
/**
 * This SHOULD be a canonical name for the document, and sufficiently unique to distinguish it from similar documents.
 */
export type TitleOfThisDocument = string;
/**
 * Contains a list of alternate names for the same document.
 */
export type Aliases = [AlternateName, ...AlternateName[]];
/**
 * Specifies a non-empty string that represents a distinct optional alternative ID used to refer to the document.
 */
export type AlternateName = string;
/**
 * The date when the current revision of this document was released
 */
export type CurrentReleaseDate = string;
/**
 * This SHOULD be the current date that the document was generated. Because documents are often generated internally by a document producer and exist for a nonzero amount of time before being released, this field MAY be different from the Initial Release Date and Current Release Date.
 */
export type DateOfDocumentGeneration = string;
/**
 * Represents the name of the engine that generated the CSAF document.
 */
export type EngineName = string;
/**
 * Contains the version of the engine that generated the CSAF document.
 */
export type EngineVersion = string;
/**
 * The ID is a simple label that provides for a wide range of numbering values, types, and schemes. Its value SHOULD be assigned and maintained by the original document issuing authority.
 */
export type UniqueIdentifierForTheDocument = string;
/**
 * The date when this document was first published.
 */
export type InitialReleaseDate = string;
/**
 * Holds one revision item for each version of the CSAF document, including the initial one.
 */
export type RevisionHistory = [Revision, ...Revision[]];
/**
 * The date of the revision entry
 */
export type DateOfTheRevision = string;
/**
 * Contains the version string used in an existing document with the same content.
 */
export type LegacyVersionOfTheRevision = string;
/**
 * Specifies a version string to denote clearly the evolution of the content of the document. Format must be either integer or semantic versioning.
 */
export type Version = string;
/**
 * Holds a single non-empty string representing a short description of the changes.
 */
export type SummaryOfTheRevision = string;
/**
 * Defines the draft status of the document.
 */
export type DocumentStatus = "draft" | "final" | "interim";
/**
 * Contains branch elements as children of the current element.
 */
export type ListOfBranches = [Branch, ...Branch[]];
/**
 * Describes the characteristics of the labeled branch.
 */
export type CategoryOfTheBranch =
  | "architecture"
  | "host_name"
  | "language"
  | "legacy"
  | "patch_level"
  | "product_family"
  | "product_name"
  | "product_version"
  | "product_version_range"
  | "service_pack"
  | "specification"
  | "vendor";
/**
 * Contains the canonical descriptor or 'friendly name' of the branch.
 */
export type NameOfTheBranch = string;
/**
 * The value should be the product’s full canonical name, including version number and other attributes, as it would be used in a human-friendly document.
 */
export type TextualDescriptionOfTheProduct = string;
/**
 * Token required to identify a full_product_name so that it can be referred to from other parts in the document. There is no predefined or required format for the product_id as long as it uniquely identifies a product in the context of the current document.
 */
export type ReferenceTokenForProductInstance = string;
/**
 * The Common Platform Enumeration (CPE) attribute refers to a method for naming platforms external to this specification.
 */
export type CommonPlatformEnumerationRepresentation = string;
/**
 * Contains a list of cryptographic hashes usable to identify files.
 */
export type ListOfHashes = [CryptographicHashes, ...CryptographicHashes[]];
/**
 * Contains a list of cryptographic hashes for this file.
 */
export type ListOfFileHashes = [FileHash, ...FileHash[]];
/**
 * Contains the name of the cryptographic hash algorithm used to calculate the value.
 */
export type AlgorithmOfTheCryptographicHash = string;
/**
 * Contains the cryptographic hash value in hexadecimal representation.
 */
export type ValueOfTheCryptographicHash = string;
/**
 * Contains the name of the file which is identified by the hash values.
 */
export type Filename = string;
/**
 * Contains a list of full or abbreviated (partial) model numbers.
 */
export type ListOfModels = [ModelNumber, ...ModelNumber[]];
/**
 * Contains a full or abbreviated (partial) model number of the component to identify.
 */
export type ModelNumber = string;
/**
 * The package URL (purl) attribute refers to a method for reliably identifying and locating software packages external to this specification.
 */
export type PackageURLRepresentation = string;
/**
 * Contains a list of URLs where SBOMs for this product can be retrieved.
 */
export type ListOfSBOMURLs = [SBOMURL, ...SBOMURL[]];
/**
 * Contains a URL of one SBOM for this product.
 */
export type SBOMURL = string;
/**
 * Contains a list of full or abbreviated (partial) serial numbers.
 */
export type ListOfSerialNumbers = [SerialNumber, ...SerialNumber[]];
/**
 * Contains a full or abbreviated (partial) serial number of the component to identify.
 */
export type SerialNumber = string;
/**
 * Contains a list of full or abbreviated (partial) stock keeping units.
 */
export type ListOfStockKeepingUnits = [StockKeepingUnit, ...StockKeepingUnit[]];
/**
 * Contains a full or abbreviated (partial) stock keeping unit (SKU) which is used in the ordering process to identify the component.
 */
export type StockKeepingUnit = string;
/**
 * Contains a list of identifiers which are either vendor-specific or derived from a standard not yet supported.
 */
export type ListOfGenericURIs = [GenericURI, ...GenericURI[]];
/**
 * Refers to a URL which provides the name and knowledge about the specification used or is the namespace in which these values are valid.
 */
export type NamespaceOfTheGenericURI = string;
/**
 * Contains the identifier itself.
 */
export type URI = string;
/**
 * Contains a list of full product names.
 */
export type ListOfFullProductNames = [FullProductName, ...FullProductName[]];
/**
 * Contains a list of product groups.
 */
export type ListOfProductGroups = [ProductGroup, ...ProductGroup[]];
/**
 * Token required to identify a group of products so that it can be referred to from other parts in the document. There is no predefined or required format for the product_group_id as long as it uniquely identifies a group in the context of the current document.
 */
export type ReferenceTokenForProductGroupInstance = string;
/**
 * Lists the product_ids of those products which known as one group in the document.
 */
export type ListOfProductIDs = [
  ReferenceTokenForProductInstance,
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Gives a short, optional description of the group.
 */
export type SummaryOfTheProductGroup = string;
/**
 * Contains a list of relationships.
 */
export type ListOfRelationships = [Relationship, ...Relationship[]];
/**
 * Defines the category of relationship for the referenced component.
 */
export type RelationshipCategory =
  | "default_component_of"
  | "external_component_of"
  | "installed_on"
  | "installed_with"
  | "optional_component_of";
/**
 * Token required to identify a full_product_name so that it can be referred to from other parts in the document. There is no predefined or required format for the product_id as long as it uniquely identifies a product in the context of the current document.
 */
export type ReferenceTokenForProductInstance1 = string;
/**
 * Token required to identify a full_product_name so that it can be referred to from other parts in the document. There is no predefined or required format for the product_id as long as it uniquely identifies a product in the context of the current document.
 */
export type ReferenceTokenForProductInstance2 = string;
/**
 * Represents a list of all relevant vulnerability information items.
 */
export type Vulnerabilities = [Vulnerability, ...Vulnerability[]];
/**
 * Contains a list of acknowledgment elements associated with this vulnerability item.
 */
export type VulnerabilityAcknowledgments = [Acknowledgment, ...Acknowledgment[]];
/**
 * Holds the MITRE standard Common Vulnerabilities and Exposures (CVE) tracking number for the vulnerability.
 */
export type CVE = string;
/**
 * Holds the ID for the weakness associated.
 */
export type WeaknessID = string;
/**
 * Holds the full name of the weakness as given in the CWE specification.
 */
export type WeaknessName = string;
/**
 * Holds the date and time the vulnerability was originally discovered.
 */
export type DiscoveryDate = string;
/**
 * Contains a list of machine readable flags.
 */
export type ListOfFlags = [Flag, ...Flag[]];
/**
 * Contains the date when assessment was done or the flag was assigned.
 */
export type DateOfTheFlag = string;
/**
 * Specifies a list of product_group_ids to give context to the parent item.
 */
export type ListOfProductGroupIds = [
  ReferenceTokenForProductGroupInstance,
  ...ReferenceTokenForProductGroupInstance[]
];
/**
 * Specifies the machine readable label.
 */
export type LabelOfTheFlag =
  | "component_not_present"
  | "inline_mitigations_already_exist"
  | "vulnerable_code_cannot_be_controlled_by_adversary"
  | "vulnerable_code_not_in_execute_path"
  | "vulnerable_code_not_present";
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Represents a list of unique labels or tracking IDs for the vulnerability (if such information exists).
 */
export type ListOfIDs = [ID, ...ID[]];
/**
 * Indicates the name of the vulnerability tracking or numbering system.
 */
export type SystemName = string;
/**
 * Is unique label or tracking ID for the vulnerability (if such information exists).
 */
export type Text = string;
/**
 * Contains a list of involvements.
 */
export type ListOfInvolvements = [Involvement, ...Involvement[]];
/**
 * Holds the date and time of the involvement entry.
 */
export type DateOfInvolvement = string;
/**
 * Defines the category of the involved party.
 */
export type PartyCategory = "coordinator" | "discoverer" | "other" | "user" | "vendor";
/**
 * Defines contact status of the involved party.
 */
export type PartyStatus =
  | "completed"
  | "contact_attempted"
  | "disputed"
  | "in_progress"
  | "not_contacted"
  | "open";
/**
 * Contains additional context regarding what is going on.
 */
export type SummaryOfTheInvolvement = string;
/**
 * Holds notes associated with this vulnerability item.
 */
export type VulnerabilityNotes = [Note, ...Note[]];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds1 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds2 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds3 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds4 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds5 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds6 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds7 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Specifies a list of product_ids to give context to the parent item.
 */
export type ListOfProductIds8 = [
  ReferenceTokenForProductInstance,
  ...ReferenceTokenForProductInstance[]
];
/**
 * Holds a list of references associated with this vulnerability item.
 */
export type VulnerabilityReferences = [Reference, ...Reference[]];
/**
 * Holds the date and time the vulnerability was originally released into the wild.
 */
export type ReleaseDate = string;
/**
 * Contains a list of remediations.
 */
export type ListOfRemediations = [Remediation, ...Remediation[]];
/**
 * Specifies the category which this remediation belongs to.
 */
export type CategoryOfTheRemediation =
  | "mitigation"
  | "no_fix_planned"
  | "none_available"
  | "vendor_fix"
  | "workaround";
/**
 * Contains the date from which the remediation is available.
 */
export type DateOfTheRemediation = string;
/**
 * Contains a thorough human-readable discussion of the remediation.
 */
export type DetailsOfTheRemediation = string;
/**
 * Contains a list of entitlements.
 */
export type ListOfEntitlements = [EntitlementOfTheRemediation, ...EntitlementOfTheRemediation[]];
/**
 * Contains any possible vendor-defined constraints for obtaining fixed software or hardware that fully resolves the vulnerability.
 */
export type EntitlementOfTheRemediation = string;
/**
 * Specifies what category of restart is required by this remediation to become effective.
 */
export type CategoryOfRestart =
  | "connected"
  | "dependencies"
  | "machine"
  | "none"
  | "parent"
  | "service"
  | "system"
  | "vulnerable_component"
  | "zone";
/**
 * Provides additional information for the restart. This can include details on procedures, scope or impact.
 */
export type AdditionalRestartInformation = string;
/**
 * Contains the URL where to obtain the remediation.
 */
export type URLToTheRemediation = string;
/**
 * Contains score objects for the current vulnerability.
 */
export type ListOfScores = [Score, ...Score[]];
/**
 * Contains information about a vulnerability that can change with time.
 */
export type ListOfThreats = [Threat, ...Threat[]];
/**
 * Categorizes the threat according to the rules of the specification.
 */
export type CategoryOfTheThreat = "exploit_status" | "impact" | "target_set";
/**
 * Contains the date when the assessment was done or the threat appeared.
 */
export type DateOfTheThreat = string;
/**
 * Represents a thorough human-readable discussion of the threat.
 */
export type DetailsOfTheThreat = string;
/**
 * Gives the document producer the ability to apply a canonical name or title to the vulnerability.
 */
export type Title = string;

/**
 * Representation of security advisory information as a JSON document.
 */
export interface CommonSecurityAdvisoryFramework {
  document: DocumentLevelMetaData;
  product_tree?: ProductTree;
  vulnerabilities?: Vulnerabilities;
  [k: string]: unknown;
}
/**
 * Captures the meta-data about this document describing a particular set of security advisories.
 */
export interface DocumentLevelMetaData {
  acknowledgments?: DocumentAcknowledgments;
  aggregate_severity?: AggregateSeverity;
  category: DocumentCategory;
  csaf_version: CSAFVersion;
  distribution?: RulesForSharingDocument;
  lang?: DocumentLanguage;
  notes?: DocumentNotes;
  publisher: Publisher;
  references?: DocumentReferences;
  source_lang?: SourceLanguage;
  title: TitleOfThisDocument;
  tracking: Tracking;
  [k: string]: unknown;
}
/**
 * Acknowledges contributions by describing those that contributed.
 */
export interface Acknowledgment {
  names?: ListOfAcknowledgedNames;
  organization?: ContributingOrganization;
  summary?: SummaryOfTheAcknowledgment;
  urls?: ListOfURLs;
  [k: string]: unknown;
}
/**
 * Is a vehicle that is provided by the document producer to convey the urgency and criticality with which the one or more vulnerabilities reported should be addressed. It is a document-level metric and applied to the document as a whole — not any specific vulnerability. The range of values in this field is defined according to the document producer's policies and procedures.
 */
export interface AggregateSeverity {
  namespace?: NamespaceOfAggregateSeverity;
  text: TextOfAggregateSeverity;
  [k: string]: unknown;
}
/**
 * Describe any constraints on how this document might be shared.
 */
export interface RulesForSharingDocument {
  text?: TextualDescription;
  tlp?: TrafficLightProtocolTLP;
  [k: string]: unknown;
}
/**
 * Provides details about the TLP classification of the document.
 */
export interface TrafficLightProtocolTLP {
  label: LabelOfTLP;
  url?: URLOfTLPVersion;
  [k: string]: unknown;
}
/**
 * Is a place to put all manner of text blobs related to the current context.
 */
export interface Note {
  audience?: AudienceOfNote;
  category: NoteCategory;
  text: NoteContent;
  title?: TitleOfNote;
  [k: string]: unknown;
}
/**
 * Provides information about the publisher of the document.
 */
export interface Publisher {
  category: CategoryOfPublisher;
  contact_details?: ContactDetails;
  issuing_authority?: IssuingAuthority;
  name: NameOfPublisher;
  namespace: NamespaceOfPublisher;
  [k: string]: unknown;
}
/**
 * Holds any reference to conferences, papers, advisories, and other resources that are related and considered related to either a surrounding part of or the entire document and to be of value to the document consumer.
 */
export interface Reference {
  category?: CategoryOfReference;
  summary: SummaryOfTheReference;
  url: URLOfReference;
  [k: string]: unknown;
}
/**
 * Is a container designated to hold all management attributes necessary to track a CSAF document as a whole.
 */
export interface Tracking {
  aliases?: Aliases;
  current_release_date: CurrentReleaseDate;
  generator?: DocumentGenerator;
  id: UniqueIdentifierForTheDocument;
  initial_release_date: InitialReleaseDate;
  revision_history: RevisionHistory;
  status: DocumentStatus;
  version: Version;
  [k: string]: unknown;
}
/**
 * Is a container to hold all elements related to the generation of the document. These items will reference when the document was actually created, including the date it was generated and the entity that generated it.
 */
export interface DocumentGenerator {
  date?: DateOfDocumentGeneration;
  engine: EngineOfDocumentGeneration;
  [k: string]: unknown;
}
/**
 * Contains information about the engine that generated the CSAF document.
 */
export interface EngineOfDocumentGeneration {
  name: EngineName;
  version?: EngineVersion;
  [k: string]: unknown;
}
/**
 * Contains all the information elements required to track the evolution of a CSAF document.
 */
export interface Revision {
  date: DateOfTheRevision;
  legacy_version?: LegacyVersionOfTheRevision;
  number: Version;
  summary: SummaryOfTheRevision;
  [k: string]: unknown;
}
/**
 * Is a container for all fully qualified product names that can be referenced elsewhere in the document.
 */
export interface ProductTree {
  branches?: ListOfBranches;
  full_product_names?: ListOfFullProductNames;
  product_groups?: ListOfProductGroups;
  relationships?: ListOfRelationships;
  [k: string]: unknown;
}
/**
 * Is a part of the hierarchical structure of the product tree.
 */
export interface Branch {
  branches?: ListOfBranches;
  category: CategoryOfTheBranch;
  name: NameOfTheBranch;
  product?: FullProductName;
  [k: string]: unknown;
}
/**
 * Specifies information about the product and assigns the product_id.
 */
export interface FullProductName {
  name: TextualDescriptionOfTheProduct;
  product_id: ReferenceTokenForProductInstance;
  product_identification_helper?: HelperToIdentifyTheProduct;
  [k: string]: unknown;
}
/**
 * Provides at least one method which aids in identifying the product in an asset database.
 */
export interface HelperToIdentifyTheProduct {
  cpe?: CommonPlatformEnumerationRepresentation;
  hashes?: ListOfHashes;
  model_numbers?: ListOfModels;
  purl?: PackageURLRepresentation;
  sbom_urls?: ListOfSBOMURLs;
  serial_numbers?: ListOfSerialNumbers;
  skus?: ListOfStockKeepingUnits;
  x_generic_uris?: ListOfGenericURIs;
  [k: string]: unknown;
}
/**
 * Contains all information to identify a file based on its cryptographic hash values.
 */
export interface CryptographicHashes {
  file_hashes: ListOfFileHashes;
  filename: Filename;
  [k: string]: unknown;
}
/**
 * Contains one hash value and algorithm of the file to be identified.
 */
export interface FileHash {
  algorithm: AlgorithmOfTheCryptographicHash;
  value: ValueOfTheCryptographicHash;
  [k: string]: unknown;
}
/**
 * Provides a generic extension point for any identifier which is either vendor-specific or derived from a standard not yet supported.
 */
export interface GenericURI {
  namespace: NamespaceOfTheGenericURI;
  uri: URI;
  [k: string]: unknown;
}
/**
 * Defines a new logical group of products that can then be referred to in other parts of the document to address a group of products with a single identifier.
 */
export interface ProductGroup {
  group_id: ReferenceTokenForProductGroupInstance;
  product_ids: ListOfProductIDs;
  summary?: SummaryOfTheProductGroup;
  [k: string]: unknown;
}
/**
 * Establishes a link between two existing full_product_name_t elements, allowing the document producer to define a combination of two products that form a new full_product_name entry.
 */
export interface Relationship {
  category: RelationshipCategory;
  full_product_name: FullProductName;
  product_reference: ReferenceTokenForProductInstance1;
  relates_to_product_reference: ReferenceTokenForProductInstance2;
  [k: string]: unknown;
}
/**
 * Is a container for the aggregation of all fields that are related to a single vulnerability in the document.
 */
export interface Vulnerability {
  acknowledgments?: VulnerabilityAcknowledgments;
  cve?: CVE;
  cwe?: CWE;
  discovery_date?: DiscoveryDate;
  flags?: ListOfFlags;
  ids?: ListOfIDs;
  involvements?: ListOfInvolvements;
  notes?: VulnerabilityNotes;
  product_status?: ProductStatus;
  references?: VulnerabilityReferences;
  release_date?: ReleaseDate;
  remediations?: ListOfRemediations;
  scores?: ListOfScores;
  threats?: ListOfThreats;
  title?: Title;
  [k: string]: unknown;
}
/**
 * Holds the MITRE standard Common Weakness Enumeration (CWE) for the weakness associated.
 */
export interface CWE {
  id: WeaknessID;
  name: WeaknessName;
  [k: string]: unknown;
}
/**
 * Contains product specific information in regard to this vulnerability as a single machine readable flag.
 */
export interface Flag {
  date?: DateOfTheFlag;
  group_ids?: ListOfProductGroupIds;
  label: LabelOfTheFlag;
  product_ids?: ListOfProductIds;
  [k: string]: unknown;
}
/**
 * Contains a single unique label or tracking ID for the vulnerability.
 */
export interface ID {
  system_name: SystemName;
  text: Text;
  [k: string]: unknown;
}
/**
 * Is a container, that allows the document producers to comment on the level of involvement (or engagement) of themselves or third parties in the vulnerability identification, scoping, and remediation process.
 */
export interface Involvement {
  date?: DateOfInvolvement;
  party: PartyCategory;
  status: PartyStatus;
  summary?: SummaryOfTheInvolvement;
  [k: string]: unknown;
}
/**
 * Contains different lists of product_ids which provide details on the status of the referenced product related to the current vulnerability.
 */
export interface ProductStatus {
  first_affected?: ListOfProductIds1;
  first_fixed?: ListOfProductIds2;
  fixed?: ListOfProductIds3;
  known_affected?: ListOfProductIds4;
  known_not_affected?: ListOfProductIds5;
  last_affected?: ListOfProductIds6;
  recommended?: ListOfProductIds7;
  under_investigation?: ListOfProductIds8;
  [k: string]: unknown;
}
/**
 * Specifies details on how to handle (and presumably, fix) a vulnerability.
 */
export interface Remediation {
  category: CategoryOfTheRemediation;
  date?: DateOfTheRemediation;
  details: DetailsOfTheRemediation;
  entitlements?: ListOfEntitlements;
  group_ids?: ListOfProductGroupIds;
  product_ids?: ListOfProductIds;
  restart_required?: RestartRequiredByRemediation;
  url?: URLToTheRemediation;
  [k: string]: unknown;
}
/**
 * Provides information on category of restart is required by this remediation to become effective.
 */
export interface RestartRequiredByRemediation {
  category: CategoryOfRestart;
  details?: AdditionalRestartInformation;
  [k: string]: unknown;
}
/**
 * Specifies information about (at least one) score of the vulnerability and for which products the given value applies.
 */
export interface Score {
  cvss_v2?: JSONSchemaForCommonVulnerabilityScoringSystemVersion20;
  cvss_v3?:
    | JSONSchemaForCommonVulnerabilityScoringSystemVersion30
    | JSONSchemaForCommonVulnerabilityScoringSystemVersion31;
  products: ListOfProductIds;
  [k: string]: unknown;
}
export interface JSONSchemaForCommonVulnerabilityScoringSystemVersion20 {
  /**
   * CVSS Version
   */
  version: "2.0";
  vectorString: string;
  accessVector?: "NETWORK" | "ADJACENT_NETWORK" | "LOCAL";
  accessComplexity?: "HIGH" | "MEDIUM" | "LOW";
  authentication?: "MULTIPLE" | "SINGLE" | "NONE";
  confidentialityImpact?: "NONE" | "PARTIAL" | "COMPLETE";
  integrityImpact?: "NONE" | "PARTIAL" | "COMPLETE";
  availabilityImpact?: "NONE" | "PARTIAL" | "COMPLETE";
  baseScore: number;
  exploitability?: "UNPROVEN" | "PROOF_OF_CONCEPT" | "FUNCTIONAL" | "HIGH" | "NOT_DEFINED";
  remediationLevel?:
    | "OFFICIAL_FIX"
    | "TEMPORARY_FIX"
    | "WORKAROUND"
    | "UNAVAILABLE"
    | "NOT_DEFINED";
  reportConfidence?: "UNCONFIRMED" | "UNCORROBORATED" | "CONFIRMED" | "NOT_DEFINED";
  temporalScore?: number;
  collateralDamagePotential?:
    | "NONE"
    | "LOW"
    | "LOW_MEDIUM"
    | "MEDIUM_HIGH"
    | "HIGH"
    | "NOT_DEFINED";
  targetDistribution?: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  confidentialityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  integrityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  availabilityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  environmentalScore?: number;
  [k: string]: unknown;
}
export interface JSONSchemaForCommonVulnerabilityScoringSystemVersion30 {
  /**
   * CVSS Version
   */
  version: "3.0";
  vectorString: string;
  attackVector?: "NETWORK" | "ADJACENT_NETWORK" | "LOCAL" | "PHYSICAL";
  attackComplexity?: "HIGH" | "LOW";
  privilegesRequired?: "HIGH" | "LOW" | "NONE";
  userInteraction?: "NONE" | "REQUIRED";
  scope?: "UNCHANGED" | "CHANGED";
  confidentialityImpact?: "NONE" | "LOW" | "HIGH";
  integrityImpact?: "NONE" | "LOW" | "HIGH";
  availabilityImpact?: "NONE" | "LOW" | "HIGH";
  baseScore: number;
  baseSeverity: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  exploitCodeMaturity?: "UNPROVEN" | "PROOF_OF_CONCEPT" | "FUNCTIONAL" | "HIGH" | "NOT_DEFINED";
  remediationLevel?:
    | "OFFICIAL_FIX"
    | "TEMPORARY_FIX"
    | "WORKAROUND"
    | "UNAVAILABLE"
    | "NOT_DEFINED";
  reportConfidence?: "UNKNOWN" | "REASONABLE" | "CONFIRMED" | "NOT_DEFINED";
  temporalScore?: number;
  temporalSeverity?: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  confidentialityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  integrityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  availabilityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  modifiedAttackVector?: "NETWORK" | "ADJACENT_NETWORK" | "LOCAL" | "PHYSICAL" | "NOT_DEFINED";
  modifiedAttackComplexity?: "HIGH" | "LOW" | "NOT_DEFINED";
  modifiedPrivilegesRequired?: "HIGH" | "LOW" | "NONE" | "NOT_DEFINED";
  modifiedUserInteraction?: "NONE" | "REQUIRED" | "NOT_DEFINED";
  modifiedScope?: "UNCHANGED" | "CHANGED" | "NOT_DEFINED";
  modifiedConfidentialityImpact?: "NONE" | "LOW" | "HIGH" | "NOT_DEFINED";
  modifiedIntegrityImpact?: "NONE" | "LOW" | "HIGH" | "NOT_DEFINED";
  modifiedAvailabilityImpact?: "NONE" | "LOW" | "HIGH" | "NOT_DEFINED";
  environmentalScore?: number;
  environmentalSeverity?: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  [k: string]: unknown;
}
export interface JSONSchemaForCommonVulnerabilityScoringSystemVersion31 {
  /**
   * CVSS Version
   */
  version: "3.1";
  vectorString: string;
  attackVector?: "NETWORK" | "ADJACENT_NETWORK" | "LOCAL" | "PHYSICAL";
  attackComplexity?: "HIGH" | "LOW";
  privilegesRequired?: "HIGH" | "LOW" | "NONE";
  userInteraction?: "NONE" | "REQUIRED";
  scope?: "UNCHANGED" | "CHANGED";
  confidentialityImpact?: "NONE" | "LOW" | "HIGH";
  integrityImpact?: "NONE" | "LOW" | "HIGH";
  availabilityImpact?: "NONE" | "LOW" | "HIGH";
  baseScore: number;
  baseSeverity: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  exploitCodeMaturity?: "UNPROVEN" | "PROOF_OF_CONCEPT" | "FUNCTIONAL" | "HIGH" | "NOT_DEFINED";
  remediationLevel?:
    | "OFFICIAL_FIX"
    | "TEMPORARY_FIX"
    | "WORKAROUND"
    | "UNAVAILABLE"
    | "NOT_DEFINED";
  reportConfidence?: "UNKNOWN" | "REASONABLE" | "CONFIRMED" | "NOT_DEFINED";
  temporalScore?: number;
  temporalSeverity?: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  confidentialityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  integrityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  availabilityRequirement?: "LOW" | "MEDIUM" | "HIGH" | "NOT_DEFINED";
  modifiedAttackVector?: "NETWORK" | "ADJACENT_NETWORK" | "LOCAL" | "PHYSICAL" | "NOT_DEFINED";
  modifiedAttackComplexity?: "HIGH" | "LOW" | "NOT_DEFINED";
  modifiedPrivilegesRequired?: "HIGH" | "LOW" | "NONE" | "NOT_DEFINED";
  modifiedUserInteraction?: "NONE" | "REQUIRED" | "NOT_DEFINED";
  modifiedScope?: "UNCHANGED" | "CHANGED" | "NOT_DEFINED";
  modifiedConfidentialityImpact?: "NONE" | "LOW" | "HIGH" | "NOT_DEFINED";
  modifiedIntegrityImpact?: "NONE" | "LOW" | "HIGH" | "NOT_DEFINED";
  modifiedAvailabilityImpact?: "NONE" | "LOW" | "HIGH" | "NOT_DEFINED";
  environmentalScore?: number;
  environmentalSeverity?: "NONE" | "LOW" | "MEDIUM" | "HIGH" | "CRITICAL";
  [k: string]: unknown;
}
/**
 * Contains the vulnerability kinetic information. This information can change as the vulnerability ages and new information becomes available.
 */
export interface Threat {
  category: CategoryOfTheThreat;
  date?: DateOfTheThreat;
  details: DetailsOfTheThreat;
  group_ids?: ListOfProductGroupIds;
  product_ids?: ListOfProductIds;
  [k: string]: unknown;
}
