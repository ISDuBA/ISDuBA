// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { expect, test } from "@playwright/test";

test("Login page has expected heading ISDuBA", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("heading", { name: "ISDuBA" })).toBeVisible();
});

test("Login page has expected title 'Login'", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/Login/);
});

test("Login page has link to github page", async ({ page }) => {
  await page.goto("/");
  await expect(
    page.getByRole("link", { name: "Visit the ISDuBA project on Github" })
  ).toBeVisible();
});

test("Login page has Login button", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("button", { name: "Login" })).toBeVisible();
});

test("Login page has expected field Server URL", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText("Server URL")).toBeVisible();
});

test("Login page has expected field Realm", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByText("Realm")).toBeVisible();
});
