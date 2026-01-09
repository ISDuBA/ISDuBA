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
  // The tests run with two browsers so there will be two comments. The random
  // value helps to distinguish the comments.
  const comment = `Lorem ipsum ${Math.random()}`;
  await page.getByLabel("New Comment").fill(comment);
  await page.getByRole("button", { name: "Send" }).click();
  await expect(page.getByText(comment)).toBeVisible();
  await page.getByRole("button", { name: "History" }).click();
  await expect(page.getByText(comment)).toBeVisible();

  await page.getByRole("tab", { name: "Vulnerabilities" }).click();
  const scoresCollapsible = await page.getByText("Scores").first();
  await scoresCollapsible.scrollIntoViewIfNeeded({ timeout: 2000 });
  await scoresCollapsible.click({ force: true });

  // Test SSVC calculator
  await page.getByTitle("Edit SSVC").click();
  await page.getByRole("button", { name: "Evaluate" }).click();
  await page.getByRole("button", { name: "active" }).click();
  await page.getByRole("button", { name: "yes" }).click();
  await page.getByRole("button", { name: "total" }).click();
  await page.getByRole("button", { name: "Custom" }).click();
  await page.getByLabel("Essential").click();
  await page.getByLabel("Irreversible").click();
  await page.getByRole("button", { name: "Save" }).click();
  await page.getByRole("button", { name: "Save" }).click();
  const ssvcBadge = page.getByTitle("SSVCv2/E:A/A:Y/T:T/P:E/B:I/M:H/D:C/2025-10-15T17:35:23Z/");
  expect(ssvcBadge).toBeDefined();
});
