// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

const tdClass = "whitespace-nowrap py-2 px-2";
const tablePadding = "px-2";
const title =
  "whitespace-normal overflow-hidden leading-6 h-12 py-2 px-2 title-column overflow-clip text-ellipsis w-full";
const publisher = "whitespace-nowrap w-48 max-w-48 overflow-clip text-ellipsis";

const searchColumnName = "_clientSearch";

type TableHeader = {
  label: string;
  attribute: string | undefined;
  class?: string;
  clickCallBack?: () => void;
  progressDuration?: number;
};

export { tdClass, tablePadding, title, publisher, searchColumnName, type TableHeader };
