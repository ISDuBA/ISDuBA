// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import * as fs from "fs";
import * as path from "path";
import * as crypto from "crypto";
import { test as baseTest } from "@playwright/test";

const istanbulCLIOutput = path.join(process.cwd(), ".nyc_output");

export function generateUUID(): string {
  return crypto.randomBytes(16).toString("hex");
}

export const test = baseTest.extend({
  // @ts-expect-error: Ignore missing typecasts
  context: async ({ context }, use) => {
    await context.addInitScript(() =>
      window.addEventListener("beforeunload", () =>
        (window as any).collectIstanbulCoverage(JSON.stringify((window as any).__coverage__))
      )
    );
    await fs.promises.mkdir(istanbulCLIOutput, { recursive: true });
    await context.exposeFunction("collectIstanbulCoverage", (coverageJSON: string) => {
      if (coverageJSON)
        fs.writeFileSync(
          path.join(istanbulCLIOutput, `playwright_coverage_${generateUUID()}.json`),
          coverageJSON
        );
    });
    await use(context);
    for (const page of context.pages()) {
      await page.evaluate(() =>
        (window as any).collectIstanbulCoverage(JSON.stringify((window as any).__coverage__))
      );
    }
  }
});

export const expect = test.expect;
