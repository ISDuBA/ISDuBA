// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test("Diff toolbox is working", async ({ page }) => {
  await page.goto("/#/search");
  await page.getByRole("button", { name: "Diff" }).click();
  await expect(
    page.getByText("Select a document or upload local ones.", { exact: true }).first()
  ).toBeVisible();
  const table = page.getByRole("table").first();
  await table.getByTitle("Add to comparison").first().click();
  const compareButton = page.getByRole("button", { name: "Compare", exact: false });
  expect(compareButton).toBeDisabled();
  await table.getByTitle("Add to comparison", { exact: false }).nth(1).click();
  expect(compareButton).toBeEnabled();
  await page.getByTitle("Remove from selection").click();
  expect(compareButton).toBeDisabled();
  const exampleDocumentURL =
    "https://raw.githubusercontent.com/oasis-tcs/csaf/refs/heads/master/csaf_2.0/examples/csaf/cisco-sa-20180328-smi2.json";
  const response = await fetch(exampleDocumentURL);
  const arrayBuffer = await response.arrayBuffer();
  await page.locator('input[type="file"]').setInputFiles({
    name: "example-document.txt",
    mimeType: "text/json",
    buffer: Buffer.from(arrayBuffer)
  });
  await expect(page.getByText("Temporary documents:")).toBeVisible();
  const tempDocTable = page.getByRole("table").nth(1);
  // When this test is run by two browsers and both upload temporary documents there
  // are two documents and we just choose the first to prevent an error.
  await tempDocTable.getByTitle("Add to comparison:", { exact: false }).first().click();
  await compareButton.click();
  await expect(page.getByText(/\d+ changes/)).toBeVisible();
});
