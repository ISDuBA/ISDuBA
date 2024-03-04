<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

export enum WorkflowStates {
  NEW = "new",
  READ = "read",
  ASSESSING = "assessing",
  REVIEW = "review",
  ARCHIVE = "archive",
  DELETE = "delete"
}
