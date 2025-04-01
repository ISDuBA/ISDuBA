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
![Screenshot 2025-04-01 at 12-22-49 Sources - Aggregators](https://github.com/user-attachments/assets/3e2148e5-e320-4d3e-88c9-d0ba546c445b)



Alternatively, they can utilize the Sources-Tab to add a source directly, specifying the domain or the location of a provider-metadata.json directly. 
![Screenshot 2025-03-26 at 13-18-23 Sources](https://github.com/user-attachments/assets/e80fd4ff-203f-4289-8ee7-3425b964ab1c)
It's important to note that ISDuBA only starts importing once a source has been officially activiated by checking the "active" checkbox (highlighted by being underlined in the following screenshot).
![Screenshot 2025-04-01 at 12-29-26 Sources - Edit source1](https://github.com/user-attachments/assets/41a15af0-6095-4d51-839e-62198b6f67d2)



## Finding advisories

When trying to find advisories, the search tab can be used to view all advisories that are accessible via the current permissions. Via the search bar, terms can be searched. An advanced search function, triggered by flipping the `Advanced` switch allows filtering more precisely [via filter expressions.](./filter_expr.md)
![Screenshot 2025-03-26 at 13-41-20 Search](https://github.com/user-attachments/assets/53d39c68-004f-42f2-8c0b-30a4554b9ac7)


Furthermore, the cogwheel on the upper left of the page allows creating stored queries which can be safed and used to quickly filter advisories. The Dashboard option will cause the stored query to appear on the dashboard. The Hide option will cause the story query to be hidden everywhere, potentially to make room for other queries. `Query criteria` also allow filtering [via filter expressions.](./filter_expr.md)


When the desired avisory is found, a detailed advisory view can be entered by clicking on it. Here, SSVC can be set and the advisory can be commented.

## Comparing Documents
![Screenshot 2025-03-26 at 13-42-22 Search](https://github.com/user-attachments/assets/20edc7bb-806a-4766-8219-f2112066d1e9)

On the bottom left of the search page, a `diff`-widget allows opening the diff-dialogue. While it's open, advisories on the search page can be selected by clicking the +-icon next to them. After having selected 2 advisories, their newest documents can be compared against each other in a [github-pr-diff](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-comparing-branches-in-pull-requests)-based view. Users can also upload local documents temporarily to compare them against documents within the system. How many and for how long [is managed via the server configuration](https://github.com/ISDuBA/ISDuBA/blob/main/docs/isdubad-config.md#-section-temp_storage-temporary-document-storage).


There is also a possibility to compare different documents of an advisory by clicking the diff-icon next to the version numbers within an advisory. Which advisories are compared to each other can be changed via clicking on the version numbers.
![Screenshot 2025-03-26 at 13-42-49 Compare](https://github.com/user-attachments/assets/1504bae4-414d-4e8a-bc60-b3c2a6db6b68)
