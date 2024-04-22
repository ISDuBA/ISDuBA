<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

This documents gives a brief overview about the impact of some configurable values within Keycloak on the application.

## Session lengths

The session lengths can be configured
in Keycloak under ```<url_of_keycloak>:/admin/master/console/#/isduba/realm-settings/```

A user stays logged in while they possess an active access token.

The key values are ```SSO Session Idle``` and ```SSO Session Max``` found in the ```Sessions``` tab
and ```Access Token Lifespan``` found in the ```Tokens``` tab.

An access token expires ```Access Token Lifespan``` after it's last refresh and is no longer valid. Since shortly before it expires (~ 5 seconds), any action
taken within the application will refresh it (even if it already expired) unless it's been longer than ```SSO Session Idle```
since the last refresh. 

(This also means setting the ```Access Token Lifespan``` to a value more than 5 seconds longer than the ```SSO Session Idle``` will
prevent the token from being refreshed at all.)

However, any session is limited by the ```SSO Session Max```
and any access tokens expiry time can be set to the session's start plus ```SSO Session Max``` at most and cannot be extended past that point in time.
