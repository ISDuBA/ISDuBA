<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

# Errorhandling

## Requests

Every request is done centralized in `client/src/utils.ts`.

### Requests per component

| Component           | Usecase                | Request(s)              |
| ------------------- | ---------------------- | ----------------------- |
| Login               | login                  | `/view`<br>`/about`<br> |
| Home                | newsflash              | _todo_                  |
| Advisories Overview | load documents         | `/documents`            |
| Advisory            | load one document      | `/documents`            |
|                     | load advisory versions | `/documents`            |
|                     | load SSVC              | `/documents`            |
|                     | load events            | `/events`               |
|                     | load comments          | `/commens`              |
|                     | create comment         | `/comments`             |
|                     | update state           | `/status`               |
|                     | load advisory state    | `/documents`            |
| Documents           | load documents         | `/documents`            |
| Diffpage            | diff documents         | _todo_                  |
| JSONDiff            | diff view              | `/diff`                 |


## Error categories

The messages are defined in `error.ts`.

### Not Found

When a (sub-) page could not be found we display `NotFound.svelte`.

### Network errors

Regarding Network errors we inform the user in the following categories:

  - `400`: `The request sent from the client could not be understood. Please contact an administrator.`
  - `401`: `You are unauthorized. Please re-login.`
  - `402`: `You are not allowed to do this. Please contact an administrator.`
  - `500`: `An error occured on the server. Please contact an administrator.`
  - `600`: `A network error occured. Try again later. If the error persists:Please contact an administrator.`

### Special case
  When client is unable to parse the JSON returned from the server:

  `The response from the server is not parsable. Please contact an administrator.`

### General Error
Every other error case is answered with
`An error occured. Please contact an administrator.`

