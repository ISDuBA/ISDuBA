// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

type AggregatorRole = {
  label: string;
  abbreviation: string;
};

type AggregatorEntry = {
  name: string;
  role: AggregatorRole;
  url: string;
  feedsAvailable: number;
  feedsSubscribed: number;
  availableSources: SourceInfo[];
  expand: boolean;
};

type FeedInfo = {
  id?: number;
  sourceID?: number;
  url: string;
  highlight: boolean;
};

type SourceInfo = {
  id?: number;
  name: string;
  feedsAvailable: number;
  feedsSubscribed: number;
  feeds: FeedInfo[];
  expand: boolean;
};

export type { AggregatorEntry, AggregatorRole, FeedInfo, SourceInfo };
