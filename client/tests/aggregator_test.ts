// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

const aggName = `Aggregator ${Math.random()}`.replace(".", "");
const editedAggName = aggName + "-edited";

test("Aggregator is working", async ({ page }) => {
  await page.goto("/#/sources/aggregators");
  await expect(page.getByRole("heading", { name: "Aggregators" })).toBeVisible();
  await page.getByRole("button", { name: "New aggregator" }).click();
  const newAggregatorNameInput = page.getByLabel("Name");
  await newAggregatorNameInput.fill(aggName);
  const aggURL = "https://wid.cert-bund.de/.well-known/csaf-aggregator/aggregator.json";
  await page.getByLabel("URL").fill(aggURL);
  await page.getByRole("button", { name: "Save aggregator" }).click();
  await newAggregatorNameInput.waitFor({ state: "hidden" });
  // Accordion item of the new aggregator should be opened right away
  const hintText =
    "These are the currently available providers. Please review their feeds and adjust the sources if needed.";
  await expect(page.getByText(hintText)).toBeVisible();
  await page.getByText("Open-Xchange GmbH").click();
  await expect(page.getByText("As new source")).toBeVisible();
  const accordionHeader = page.getByRole("button", {
    name: new RegExp(`^${aggName}.*`),
    exact: true
  });
  await accordionHeader.scrollIntoViewIfNeeded();
  // Prevent that PW hits the toggle button
  await accordionHeader.click({ position: { x: 1, y: 1 }, scroll: "none" });
  await page.getByText(hintText).waitFor({ state: "hidden" });

  // Edit aggregator
  await page.getByTitle(`Edit aggregator ${aggName}`).click();
  expect(page.getByTitle(`Save ${aggName}`)).toBeVisible();
  await page.getByTitle(`Cancel editing ${aggName}`).click();
  expect(page.getByTitle(`Save ${aggName}`)).not.toBeVisible();
  await page.getByTitle(`Edit aggregator ${aggName}`).click();
  const nameLabel = page.getByLabel("Name");
  await nameLabel.clear();
  await nameLabel.pressSequentially(editedAggName);
  await page.getByTitle(`Save ${aggName}`).click();
  await nameLabel.waitFor({ state: "hidden" });

  // Delete aggregator
  const title = `Remove aggregator ${editedAggName}`;
  await page.getByTitle(title).click();
});
