// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { test as setup } from "@playwright/test";
import { fileURLToPath } from "url";
import path from "path";

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const authFile = path.join(__dirname, "../playwright/.auth/user.json");

setup("authenticate and upload document", async ({ page }) => {
  // Perform authentication steps. Replace these actions with your own.
  await page.goto("/");
  await page.getByRole("button", { name: "Login" }).click();
  await page.getByLabel("Username or email").fill("test-user");
  await page.getByLabel("Password", { exact: true }).fill("test-user");
  await page.getByRole("button", { name: "Sign In" }).click();
  // Wait until the page receives the cookies.
  //
  // Sometimes login flow sets cookies in the process of several redirects.
  // Wait for the final URL to ensure that the cookies are actually set.
  await page.waitForURL("/#/");
  // Alternatively, you can wait until the page reaches a state where all cookies are set.

  // Upload test document
  await page.goto("/#/sources");
  await page.getByRole("button", { name: " Upload documents" }).click();
  await page
    .locator('input[type="file"]')
    .setInputFiles("../docs/example-advisories/avendor-advisory-0004.json");
  await page.getByRole("button", { name: "Upload" }).click();

  // End of authentication steps.

  await page.context().storageState({ path: authFile });
});
