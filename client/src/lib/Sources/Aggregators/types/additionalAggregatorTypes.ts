// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import type { CSAFAggregator as CSAFAggregator2_0 } from "#lib/Sources/Aggregators/types/aggregator-2.0.d.ts";
import type { CSAFAggregator as CSAFAggregator2_1 } from "#lib/Sources/Aggregators/types/aggregator-2.1.d.ts";

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
  aggregator: CSAFAggregator2_0 | CSAFAggregator2_1;
  custom: Custom;
}
