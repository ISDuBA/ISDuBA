<!--
 This file is Free Software under the Apache-2.0 license
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

Part of the "source management" system is
how to discover new feeds and manage the ones the users want to follow.

# Discovering new Feeds

In ISDuBA's a _source entry_ means either one domain or a URL
to a `provider-metadata.json` (PMD).

Each source entry can lead to several _feeds_
to be potentially followed for new CSAF documents.

New feeds either come out of two events:
 1. a source entry's PMD lists more feeds
 2. an aggregator has a new entry that leads to new feeds.


## Implementation concept

**proposal**

Just like the feeds, the source entries and aggregators will be
periodically checked. We keep the last result of that query.

Whenever there is a difference to last result,
we create a new type of event.
These events can be seen and searched for by the source managers -
also show on the dashboard for them.

In case a source entry gets new feeds, we mark the source entry
as "needs review". This can be seen in the list of source entries.


### Review a source entry

If a source entry for review is shown in detail the user is
presented with the changed information and can decide to
do a selection and then save or abort the review.

New feeds will be indicated with a default value of being subscribed.
They can be deselected.

Old deselected feeds are shown differently. And can be selected as always.

We also show the event that has triggered the review.
This mean the user can see even if a feed was deleted from the event message.


### Aggregators

A subpage shows a list of aggregators.
We expect a handful up to a low two digit number.

When an aggregator changes we operate on each changed feed:
All source entries will be checked if the feed has a configuration there.

If we do not find a configuration for this feed, we create a fresh
source entry, make it inactive and mark it for review.
Together with the event, the source manager will then see it in their
list of sources.
(If we do check the sources for new feeds before we check the aggregators,
there will not be new feeds from the aggregator for old source entries
they the new feeds will have been found from the source entry check already.)

If an aggregator is inspected, we also check all source entries
for each feed and lists those sources where the feed has a configuration.
The user can directly jump to the source entry's detailed view from there.

It is possible that a feed is listed in several sources.
Or in none, which means the source entry was deleted by the users.
The list then offers a button to create a new source entry again.

That sources entries do not link to the aggregator entries they came
from or where they are listed is okay.
In the rare case that a user wants to know they can list all aggregators.

Aggregators shall only bring extra information, but the handling of the feeds
shall be integrated completely in the systems of the source entries.
