// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

const activeClass =
  "flex items-center p-2 text-base font-normal text-primary-900 bg-primary-200 dark:bg-gray-950 dark:text-white hover:bg-primary-100 dark:hover:bg-black";
const nonActiveClass =
  "flex items-center p-2 text-base font-normal text-white dark:text-white hover:bg-primary-100 dark:hover:bg-black hover:text-primary-900";

const sidebarItemClass = "px-0 py-0";
const sidebarItemLinkClass = "px-6 py-4 rounded-none! hover:text-primary-700 dark:hover:text-white";

export { activeClass, nonActiveClass, sidebarItemClass, sidebarItemLinkClass };
