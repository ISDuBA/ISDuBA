// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("Advisory view is working", async ({ page }) => {
  await page.goto("/#/search");
  await page.getByPlaceholder("Enter a search term").fill("avendor");
  await page.getByRole("button", { name: "Search", exact: true }).click();
  await page.getByText("Avendor-advisory-0004").first().click({ force: true });
  await expect(page.getByText("Test CSAF document")).toBeVisible();
});
