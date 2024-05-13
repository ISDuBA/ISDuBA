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

The key value is ```SSO Session Max``` found in the ```Sessions``` tab. This will limit the maximum duration of an active session to it's value. 
