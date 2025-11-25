// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";

test.describe.configure({ mode: "parallel" });

test("User page has Logout button", async ({ page }) => {
  await page.goto("/#/login");
  await expect(page.getByRole("button", { name: "Logout" })).toBeVisible();
});

test("User page has expected field Server URL", async ({ page }) => {
  await page.goto("/#/login");
  await expect(page.getByText("Server URL")).toBeVisible();
});

test("User page has expected field Realm", async ({ page }) => {
  await page.goto("/#/login");
  await expect(page.getByText("Realm")).toBeVisible();
});

test("User page has link to Git repo and API", async ({ page }) => {
  await page.goto("/#/login");
  await expect(page.getByRole("link", { name: "Github" })).toBeVisible();
  await expect(page.getByRole("link", { name: "API" })).toBeVisible();
});
