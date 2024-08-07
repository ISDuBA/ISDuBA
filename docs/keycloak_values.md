<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

This documents gives a brief overview about the impact of some configurable values within Keycloak on the application.

## Groups
Which advisories any given user can access is regulated by their keycloak groups. Group access rights are additive, meaning a user has all rights of every group they are part of.
On the graphical interface of keycloak, groups are defined and can be edited under ```<url_of_keycloak>:/admin/master/console/#/isduba/groups```
A group has always at least an identifying unique name. For the purpose of this application, they also have a ```TLP``` attribute:

Key: ```TLP``` (This is the mandatory key name. A user without access to a group with the attribute ```TLP``` will be treated as being only in the default group.)

Value: ```{"$PUBLISHER": $TLPS}```

where

 - $PUBLISHER: the publisher this group can access. The value ```*``` allows access to all publishers. The respective value in advisories can be found under document\publisher\name\value. 
 - $TLPS: An array containing any combination of ```"WHITE"```, ```"GREEN"```, ```"AMBER"``` and ```"RED"```. The array-elements grant access to advisories of their respective TLP-level.

The current default is:

```{"*": ["WHITE"]}```

which allows access to all TLP:WHITE advisories of all publishers. Anyone without a group will have these priviledges.

The creation of new groups can be done via the graphical keycloak interface or via the [createGroup script.](./scripts/keycloak/createGroup.sh)

Editing any existing group can be done via the graphical interface.

Adding users to a group can be done via the graphical interface both within the group's own tab or under the ```group``` tab of a user
or via the [script designed to add users to roles or groups.](./scripts/keycloak/assignUserToRoleAndGroup.sh)

## Session lengths

The session lengths can be configured
in Keycloak under ```<url_of_keycloak>:/admin/master/console/#/isduba/realm-settings/```

A user stays logged in while they possess an active access token.

The key value is ```SSO Session Max``` found in the ```Sessions``` tab. This will limit the maximum duration of an active session to it's value. 
