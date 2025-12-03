// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2025 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("Sources are working", async ({ page }) => {
  await page.goto("/#/sources");
  await page.getByRole("link", { name: "Add source", exact: false }).click();
  await page.getByLabel("Domain/PMD").fill("lorem  ipsum");
  const checkUrlButton = await page.getByRole("button", {
    name: "Search and load provider metadata",
    exact: false
  });
  expect(checkUrlButton).toBeDisabled();
  await page.getByLabel("Domain/PMD").fill("intevation.de");
  await checkUrlButton.click();
  const sourceName = `Source ${Math.random()}`;
  await page.getByLabel("Name").fill(sourceName);
  await page.getByText("Advanced options").click();
  await page.getByLabel("Strict mode").check();
  const unsubscribeButton = await page.getByRole("button", {
    name: "Unsubscribe all",
    exact: false
  });
  await unsubscribeButton.click();
  expect(unsubscribeButton).toBeDisabled();
  await page.getByLabel("Enable feed with label", { exact: false }).first().click();
  await page.getByRole("button", { name: "Save source", exact: false }).click();

  // Test if input fields for time range fire requests to /api/sources/feed.
  await page.getByLabel("View feed details").click();
  const expectedUrl = new RegExp(
    `.*\\/api\\/sources\\/feeds\\/\\d+\\/log\\?.*from=\\d+-\\d+-\\d+T\\d+:\\d+:\\d+.\\d+Z&to=\\d+-\\d+-\\d+T\\d+:\\d+:\\d+.\\d+Z.*`
  );
  let feedLogRequestCount = 0;
  page.on("request", (data) => {
    if (expectedUrl.test(data.url())) feedLogRequestCount++;
  });
  // This should cause only one request since we limit requests with debounds.
  let requestPromise = page.waitForRequest(expectedUrl, { timeout: 5000 });
  await page.getByTitle("Hours").first().fill("6");
  await page.getByTitle("Minutes").first().fill("10");
  await requestPromise;

  requestPromise = page.waitForRequest(expectedUrl, { timeout: 5000 });
  await page.getByTitle("Hours").nth(1).fill("20");
  await requestPromise;

  expect(feedLogRequestCount).toBe(2);

  await page.goto("/#/sources");
  await page.getByText(sourceName).click();
  await page.getByRole("button", { name: "Delete source", exact: false }).click();
  await page.getByRole("button", { name: "Yes, I'm sure" }).click();
});
