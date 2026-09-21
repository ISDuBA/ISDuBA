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

test.beforeEach(async ({ page }) => {
  await page.goto("/#/login");
});

test("User page has header 'ISDuBA'", async ({ page }) => {
  const appName = page.getByText("ISDuBA", { exact: true });
  await expect(appName.first()).toBeVisible();
  await expect(appName.nth(1)).toBeVisible();
});

test("User page has Logout button", async ({ page }) => {
  await expect(page.getByRole("button", { name: "Logout" })).toBeVisible();
});

test("User page has Profile button", async ({ page }) => {
  await expect(page.getByRole("link", { name: "Profile" })).toBeVisible();
});

test("User page has expected field Server URL", async ({ page }) => {
  await expect(page.getByText("Server URL")).toBeVisible();
});

test("User page has expected field Realm", async ({ page }) => {
  await expect(page.getByText("Realm")).toBeVisible();
});

test("User page has link to Git repo and API", async ({ page }) => {
  await expect(page.getByRole("link", { name: "Github" })).toBeVisible();
  await expect(page.getByRole("link", { name: "API" })).toBeVisible();
});

test("User page has button to switch between light and dark mode", async ({ page }) => {
  const themeButton = page.getByRole("button", { name: /.* mode/ });
  await expect(themeButton).toBeVisible();
  await themeButton.click();
  await themeButton.click();
});

test("User page has footer text", async ({ page }) => {
  await expect(page.getByRole("link", { name: "Example Corp" })).toBeVisible();
});

test("User page shows all TLPs the test-user is allowed to see", async ({ page }) => {
  await expect(page.getByText("WHITE", { exact: true })).toBeVisible();
  await expect(page.getByText("GREEN", { exact: true })).toBeVisible();
  await expect(page.getByText("AMBER", { exact: true })).toBeVisible();
  await expect(page.getByText("RED", { exact: true })).toBeVisible();
});

test("User page shows all roles of test-user", async ({ page }) => {
  await expect(page.getByText("Admin", { exact: true })).toBeVisible();
  await expect(page.getByText("Importer", { exact: true })).toBeVisible();
  await expect(page.getByText("Editor", { exact: true })).toBeVisible();
  await expect(page.getByText("Source-Manager", { exact: true })).toBeVisible();
});
