// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2025 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("Queries can be configured", async ({ page }) => {
  await page.goto("/#/queries");
  await page.getByRole("link", { name: "New query", exact: false }).first().click();
  const queryName = `Query ${Math.random()}`;
  await page.getByLabel("Name:").fill(queryName);
  await page.getByLabel(`Set visibility of column publisher`).check();
  await page.getByLabel(`Set visibility of column title`).check();

  // When invalid query criteria are entered the application should not navigate to the
  // query overview. Instead it should display an error message.
  await page.getByLabel("Query criteria:").fill("abc");
  await page.getByRole("button", { name: "Save", exact: false }).click();
  await page.getByText("Failed to save query");

  await page.getByLabel("Query criteria:").fill("");
  await page.getByRole("button", { name: "Save", exact: false }).click();
  const newQueryButton = await page.getByRole("link", { name: "New query", exact: false }).first();
  await expect(newQueryButton).toBeVisible();
  const table = page.getByRole("table").first();
  await expect(table).toContainText(queryName);
  await page.getByTitle(`delete ${queryName}`, { exact: false }).click();
  await page.getByRole("button", { name: "Yes", exact: false }).click();
  await expect(table).not.toContainText(queryName);
});
