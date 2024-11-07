<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# ISDuBA: How to

## Login prompt

The first site a user will be connected to.
It shows which Keycloak server url and realm a user will be connected to,
a possibility to login, and a link to this repository.

After logging in, the user will be redirected to the Dashboard.

## The Sidebar

On the left side of the application, the user can utilize a sidebar to navigate between the
different tabs of the application that they can access.

## Dashboard

The dashboards two main features are the dashboard queries,
as well as a short overview over the system activities
via statistics of imported advisories.

### Dashboard Queries

The dashboard [queries](#queries) are sets of up to six cards on the dashboard that show advisories
or events that have been deemed personally relevant. Which advisories or events
are shown is customizable. The default is that nearly every user, depending on role,
is shown events or advisories determined by two sets of queries, 
one which involves advisories that have become relevant
and one which shows whether they have been recently
mentioned or otherwise been involved. Clicking on any card will
redirect to the relevant [advisory](#advisories), while clicking on the "more" button will
redirect to the [search tab](#search) and automatically search for all relevant advisories.

The admins can modify these
default queries, while users can modify their own dashboard independently
in the [queries section of the Search tab](#queries).

### Changed Sources 

>[!Note]
>Only for Source-Managers

Additionally an overview over [different sources that may require attention](#sources)
has been given. Similarly, clicking on the cards will lead to the relevant sources.

### Dashboard statistics

A statistic that shows how many advisories were downloaded in any given
30-minute intervall within the last three calender days including today.

## Search

The search tab displays all advisories, documents and events within the system
and allows to search for specific criteria via customizable queries as well as a search
field. The centerpiece is a table that, per default, shows all the advisories within the system.
A 3-button-menue on the top right allows to switch between showing advisories, documents and events.
Within the table, the advisories can be sorted via all depicted criteria by clicking on it. 
Clicking on an entry will lead to the relevant advisory within it's [document view](#document-view).

The checkboxes on the left side allow for multioperations where applicable and
the Diff-module on the bottom left allows editors and reviewers to compare
any listed advisory to each other or to a manually uploaded advisory. Clicking on the Diff-button
will open a small box and show a new plusminus-icon next to the advisories.
After selecting two documents, the compare-button leads to a github-styled
diff-overview of these documents.


### Search and Advanced Search

The search bar allows searching for keywords within advisories.
Every occurance of a searched keyword will be listed
and will lead to the corresponding advisories [document view.](#document-view)

Triggering the advanced search toggle will allow for a filter expression based search
that can more accurately narrow down what the user wants to find. A guide on how the filter 
expression works [can be found within the filter expressions documentation](https://github.com/ISDuBA/ISDuBA/blob/main/docs/filter_expr.md).

### Queries

>[!Note]
>The query overview is different for admins. More below.

Utilizing the query cogwheel at the top left, a user can enter the query designer.
There are four different kind of queries:

 * Personal Queries
  * These are only visible and effective for the user who created them

 * Global relevant dashboard queries
  * These are the preset dashboard queries if not filtered out by the admins

 * Global dashboard queries (not displayed)
  * These are admin-created dashboard queries that users can copy if they so desire

 * Global search queries
  * These are admin-created queries that can be used to filter results on the search page.


Here a user can view the existing queries and copy them so they can edit them to fit their own needs.
They can also create new queries via the "New Query"-Button. 

>[!Note]
>For admins, the Global queries from before are replaced by a singular global query table

Admins are also able to edit or rename these queries.

Every query has 4 meta-attributes:
 * A name which will be a unique descriptor of the query
 * A description which will be displayed on the dashboard or the search page
 * A dashboard toggle to determine whether the query should be displayed on the dashboard or the search page
 * A hide toggle where a user can hide the query for themselfes to potentially reduce clutter

All of those can be set within the query overview.

Clicking on a query will lead to the query editor where these can be viewed or edited.

### Query Editor

Within the query editor, the queries filter options as well as meta-attributes can be set.
Admins can additionally set a whether the query should be global, and for which role the query should be.

The query type can be set (whether it's searching for advisories, documents or events),
as well as which information should be displayed about the results in which order.

Lastly and most importantly, the results can be filtered via the query criteria, which also
[use the filter expressions the advanced search uses.](https://github.com/ISDuBA/ISDuBA/blob/main/docs/filter_expr.md)

As such, the advanced search can be thought of as a live query for the search page.

## Document View

This is the users main hub for evaluating documents and advisories.
All information about a document is listed here, and using the labeled buttons
the user can look at any imported document of the advisory.
Using the "Show Changes"-button, users can also directly compare differences
between document versions. Additionally,
some actions can set the "state" of an advisory. The following states exist:

 * New: An advisory that noone has yet interacted with.
 * Read: An advisory someone has spend at least some time on.
 * Assessing: An advisory that has actively been worked on
 * Review: An advisory that has been evaluated and should be reviewed by a reviewer
 * Archived: An advisory that has been archived
 * Delete: An advisory that is marked for deletion

To see who can set which states, [refer to this graphic](https://github.com/ISDuBA/ISDuBA/blob/main/docs/workflow.svg).

Setting the SSVC-score can be done via the "Please enter a SSVC"-tool, where
users can either set the SSVC directly or use a decision tree via the "Evaluate"-button.

Lastly, to ease communication, there exists a comment-field where users can
communicate with each other as well as a history to see what has happened with the advisory thus far.

## Sources

The sources tab is the central hub for Source-Managers to configure the automatic download of
advisories from different Providers. The table "CSAF Provider" lists all providers of the system
as well as whether they are active, currently downloading and what they downloaded within the last day.
There is one special entry named "manual_imports" which tracks the manual imports done by the
[importer](#manual-imports).

>[!Note]
>The following actions can be performed by Source-Managers only!

Clicking onto an entry will lead to a more detailed view of the provider. After the info-section 
the source-manager can configure metainformation and activate or deactivate a source entirely.
Here, important information like the maximum age of to-be-downloaded advisories or the download-rate
can be configured. After that, the "Feeds" section allows to see information about every
single feed listed by the provider and clicking on them or the logs button
will show a detailed log feed about that csaf-feed. By clicking onto the trash can icon, a feed
can be deleted. It can then be readded by clicking the plus button. 

>[!Note]
>Doing so will cause this feed to be considered new and the application will attempt to restart the entire downloading process regarding that feed.

Last, some statistics show the import and error history of the provider, similar to what is
displayed for the entire system in [the statistics tab](#statistics)

## Statistics

The statistics tab gives a quick overview over the imports and advisories within the system.
There are three charts:

 * Advisories and Documents shows the creation date of documents and advisories in the system.

>[!Note]
>Any point in the graph represents how many advisories have been *created* prior. It does not correlate to the time of import!

 * Imports and CVEs shows when and how many documents have been imported and how many import failures happened. 

 * Critical value of imports breaks down the imports into their criticality values, with each section corresponding to a range of critical values.

The time range of each graph can be adjusted below the graph itself.
The display can also be switched to a table via the icons on top of the graph, which shows
the numerical values against timestamps. 
## Profile

The profile page contains a quick overview over the users personal information, such as roles or groups and
allows the user to quickly logout or view their keycloak profile.
