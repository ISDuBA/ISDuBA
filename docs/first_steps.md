<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

## Adding sources

Sources are what are primarily providing security advisories to ISDuBA. Managing which sources' advisories are important is the responsibility of the `source-manager` role.

An easy way to find new sources is via aggregators within the `aggregators` tab, which per default contains the preconfigures but inactive "BSI CSAF Lister" as well as an option to configure additional aggregators. Using these aggregators, the source-manager can create new sources by expanding an aggregator and then the desired source.

Alternatively, they can utilize the Sources-Tab to add a source directly, specifying the domain or the location of a provider-metadata.json directly. It's important to note that ISDuBA only starts importing once a source has been officially activiated by checking the "active" checkbox.

## Finding advisories

When trying to find advisories, the search tab can be used to view all advisories that are accessible via the current permissions. Via the search bar, terms can be searched. An advanced search function, triggered by flipping the "Advanced" switch allows filtering more precisely [via filter expressions.](./filter_expr.md)

Furthermore, the cogwheel on the upper left of the page allows creating stored queries which can be safed and used to quickly filter advisories. The Dashboard option will cause the stored query to appear on the dashboard. The Hide option will cause the story query to be hidden everywhere, potentially to make room for other queries. `Query criteria` also allow filtering [via filter expressions.](./filter_expr.md)

When the desired avisory is found, a detailed advisory view can be entered by clicking on it. Here, SSVC can be set and the advisory can be commented.

