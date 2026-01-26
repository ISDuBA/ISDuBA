// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("cve comparison is working", async ({ page }) => {
  await page.goto("/#/search");
  await page.getByPlaceholder("Enter a search term").fill("avendor");
  await page.getByRole("button", { name: "Search", exact: true }).click();

  await page.getByText("Avendor-advisory-0005").first().click({ force: true });

  await expect(page.getByText("Test CSAF document")).toBeVisible();

  await page
    .locator("button", {
      has: page.locator(".bx.bx-link-alt.bx-rotate-90")
    })
    .click();

  // Test visibility of common components
  await expect(page.getByText("Documents having the same CVEs as")).toBeVisible();
  await expect(page.getByText("Avendor-advisory-0004").first()).toBeVisible();
  await expect(page.getByText("3 (interim)").first()).toBeVisible();
  await expect(page.getByText("3 (final)").first()).toBeVisible();
  await expect(page.getByText("CVE-2020-1234")).toBeVisible();
  await expect(page.locator(".text-center.px-2.py-2.w-fit.min-w-0.bx").first()).toBeVisible();

  await page.getByText("Avendor-advisory-0004").nth(3).click();
  await expect(page.getByText("Test CSAF document")).toBeVisible();

  await page.getByRole("tab", { name: "Vulnerabilities" }).click();
  const scoresCollapsible = await page.getByText("Scores").first();
  await scoresCollapsible.scrollIntoViewIfNeeded({ timeout: 2000 });
  await scoresCollapsible.click({ force: true });

  await page
    .locator("button", {
      has: page.locator(".bx.bx-link-alt.bx-rotate-90")
    })
    .nth(2)
    .click();
  await expect(page.getByText("Documents having the same CVEs as")).toBeVisible();
});
