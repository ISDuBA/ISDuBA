// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2025 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";
import { docLinkRegex } from "./utils";

const queryName = `Query ${Math.random()}`;

test("Queries can be configured", async ({ page }) => {
  await page.goto("/#/queries");
  await page.getByRole("link", { name: "New query", exact: false }).first().click();

  const docLinkHref = await page.getByRole("link", { name: "Documentation" }).getAttribute("href");
  expect(docLinkHref).toBeDefined();
  if (docLinkHref) expect(docLinkRegex.test(docLinkHref)).toBeTruthy();
  await page.getByLabel("Name:").fill(queryName);
  await page.getByLabel("Dashboard").check();
  await page.getByLabel("Hide").check();
  await page.getByLabel(`Set visibility of column publisher`).check();
  await page.getByLabel(`Set visibility of column title`).check();

  // When invalid query criteria are entered the application should not navigate to the
  // query overview. Instead it should display an error message.
  await page.getByLabel("Query criteria:").fill("abc");
  await page.getByRole("button", { name: "Save", exact: false }).click();
  await expect(page.getByText("Failed to save query")).toBeVisible();

  await page.getByLabel("Query criteria:").fill("");
  await page.getByRole("button", { name: "Save", exact: false }).click();
  const newQueryButton = page.getByRole("link", { name: "New query", exact: false }).first();
  await expect(newQueryButton).toBeVisible();
});

test("Query attributes 'dashboard', 'hide', and 'default' can be changed", async ({ page }) => {
  await page.goto("/#/queries");
  const table = page.getByRole("table").first();
  await expect(table).toContainText(queryName);
  const dashboardCheckbox = page.getByLabel(`Show query ${queryName} on dashboard`);
  await expect(dashboardCheckbox).toBeChecked();

  const queriesUrl = new RegExp(`.*\\/api\\/queries\\/\\d+`);
  const queryPromise = page.waitForRequest(queriesUrl, { timeout: 5000 });
  await dashboardCheckbox.uncheck();
  await queryPromise;
  await expect(dashboardCheckbox).not.toBeChecked();
  const hideCheckbox = page.getByLabel(`Hide query ${queryName} everywhere`);
  await expect(hideCheckbox).toBeChecked();
  const ignoreUrl = new RegExp(`.*\\/api\\/queries\\/ignore`);
  const ignorePromise = page.waitForRequest(ignoreUrl, { timeout: 5000 });
  await hideCheckbox.uncheck();
  await ignorePromise;
  await expect(hideCheckbox).not.toBeChecked();
});

test("Queries can be deleted", async ({ page }) => {
  await page.goto("/#/queries");
  await page.getByTitle(`delete ${queryName}`, { exact: false }).click();
  await page.getByRole("button", { name: "Yes", exact: false }).click();
  const table = page.getByRole("table").first();
  await expect(table).not.toContainText(queryName);
});
