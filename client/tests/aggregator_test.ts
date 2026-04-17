// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("Aggregator is working", async ({ page }) => {
  await page.goto("/#/sources/aggregators");
  await expect(page.getByRole("heading", { name: "Aggregators" })).toBeVisible();
  await page.getByRole("button", { name: "New aggregator" }).click();
  const aggName = `Aggregator ${Math.random()}`;
  await page.getByLabel("Name").fill(aggName);
  const aggURL = "https://wid.cert-bund.de/.well-known/csaf-aggregator/aggregator.json";
  await page.getByLabel("URL").fill(aggURL);
  await page.getByRole("button", { name: "Save aggregator" }).click();
  await expect(
    page.getByText(
      "These are the currently available providers. Please review their feeds and adjust the sources if needed."
    )
  ).toBeVisible();
  await page.getByText("Open-Xchange GmbH").click();
  await expect(page.getByText("As new source")).toBeVisible();
  const title = `Remove aggregator ${aggName}`;
  await page.getByTitle(title).click();
});
