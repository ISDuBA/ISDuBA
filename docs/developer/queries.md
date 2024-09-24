<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

## Using queries for dashboard configuration

**proposal**

We save in the `stored_queries` table if a query is a basis for
a dashboard display.

Optionally an entry has a role associated.
If there is no role, it is aimed for all roles.

It is indicated if a query is global and thus only editable for admins.


A table for ignored information allows the user to deselect
queries from considerations.
This is good for dashboard and for the list of queries in the search page.

Alternatives considered (and not implemented):
 * Doing the indication via a naming convention.
   (This would have been less explicit.)
 * Using some sort of `replaces` column in the query table.
   (This does not work well if a user wants to see global queries
    in the search list at the same time as the personal
    modified cloned queries.)

### Default dashboard queries

To support users to fulfill their role in a good way we store some queries
in `stored_queries` as part of the migrations. Since we don't know the names of
the accounts the admins will be using we add `system-default` as `definer` of the
these queries.

If some admin would decide to save these queries as non-global it could happen
that the queries would be lost because in this case they would belong to
`system-default`. If there would be no such user as admin this meant that no
one could ever set the queries to global again. Therefore, we prevent
setting the default queries to non-global.

### client does the selection

The client will select the first two eligible queries
and display the first two before displaying a brief stats section.

The first is to be displayed left and aims to show which new documents
have been incoming.

The second it to be displayed right and aims to show which changes
were done within our application to data that was already in.

Then the client displays the other eligible queries below.


#### calculating eligible queries

Take all personal queries that for the dashboard.

Behind them take all global queries for the dashboard that
are configured for the leading role and are not deselected
via the ignored table.


#### editing user interface

A special button ala "Clone the global dashboard queries for me"
will copy the two global queries, mark the copies as
for the dashboard, puts them into the first position
and marks the two global query as ignored.

This button avoids the problem that if only the second
global query is cloned, it would be the first to be displayed
under the calculation rule above.



