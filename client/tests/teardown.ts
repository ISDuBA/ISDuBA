// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("Delete documents", async ({ page }) => {
  await page.goto("/");
  await page.getByRole("button", { name: "Login" }).click();
  await page.getByLabel("Username or email").fill("test-user");
  await page.getByLabel("Password", { exact: true }).fill("test-user");
  await page.getByRole("button", { name: "Sign In" }).click();
  await page.waitForURL("/#/");
  await page.goto("/#/search");
  await expect(page.getByText("advisories in total")).toBeVisible();
  await page.getByPlaceholder("Enter a search term").fill("avendor");
  await page.getByRole("button", { name: "Search", exact: true }).click();
  await page.getByTitle("delete Avendor-advisory-0004").click();
  await page.getByText("Yes").click();
  await page.getByTitle("delete Avendor-advisory-0005").click();
  await page.getByText("Yes").click();
  await expect(page.getByText("No results were found.")).toBeVisible();
});
