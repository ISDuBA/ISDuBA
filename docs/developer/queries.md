## Using queries for dashboard configuration

**proposal**

We save in the `stored_queries` table if a query is a basis for
a dashboard display.

Optionally an entry has a role associcated.
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


### client does the selection

The client will select the first two eligable queries
and display the first two before displaying a brief stats section.

The first is to be displayed left and aims to show which new documents
have been incoming.

The second it to be displayed right and aims to show which changes
were done within our application to data that was alread in.

Then the client displays the other eligable queries below.


#### calculating eligable queries

Take all personal queries that for the dashboard.

Behind them take all global queries for the dasboard that
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



