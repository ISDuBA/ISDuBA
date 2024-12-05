// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { expect, test } from "playwright-test-coverage";

test("Diff toolbox is working", async ({ page }) => {
  await page.goto("/#/search");
  await page.getByRole("button", { name: "Diff" }).click();
  await expect(
    page.getByText("Select a document or upload local ones.", { exact: true }).first()
  ).toBeVisible();
});
