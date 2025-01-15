// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { test as testBase } from "@playwright/test";
import { addCoverageReport } from "monocart-reporter";
import type { Page } from "@playwright/test";

// fixtures
const test = testBase.extend<{
  autoTestFixture: string;
}>({
  autoTestFixture: [
    async ({ page }: { page: Page }, use: any) => {
      // NOTE: it depends on your project name
      const isChromium = test.info().project.name === "chromium";

      // coverage API is chromium only
      if (isChromium) {
        await Promise.all([
          page.coverage.startJSCoverage({
            resetOnNavigation: false
          })
        ]);
      }

      await use("autoTestFixture");

      if (isChromium) {
        const [jsCoverage] = await Promise.all([page.coverage.stopJSCoverage()]);
        const coverageList = [...jsCoverage];
        await addCoverageReport(coverageList, test.info());
      }
    },
    {
      scope: "test",
      auto: true
    }
  ]
});

export { test };
