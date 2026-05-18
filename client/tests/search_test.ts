// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import { expect } from "@playwright/test";
import { test } from "./fixtures";
import { docLinkRegex } from "./utils";

test("Advisory view is working", async ({ page }) => {
  await page.goto("/#/search");
  // Doesn't work without "force"
  await page.getByLabel("Advanced").check({ force: true });
  const docLinkHref = await page.getByRole("link", { name: "Documentation" }).getAttribute("href");
  expect(docLinkHref).toBeDefined();
  if (docLinkHref) expect(docLinkRegex.test(docLinkHref)).toBeTruthy();
});
