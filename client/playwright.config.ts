// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { type PlaywrightTestConfig, devices } from "@playwright/test";

const config: PlaywrightTestConfig = {
  projects: [
    // Setup project

    { name: "setup", testMatch: /setup\.ts/, teardown: "cleanup" },
    { name: "cleanup", testMatch: /teardown\.ts/ },
    {
      name: "chromium",
      use: {
        ...devices["Desktop Chrome"],
        // Use prepared auth state.
        storageState: "playwright/.auth/user.json"
      },
      dependencies: ["setup"]
    },

    {
      name: "firefox",
      use: {
        ...devices["Desktop Firefox"],
        // Use prepared auth state.
        storageState: "playwright/.auth/user.json"
      },
      dependencies: ["setup"]
    }
  ],
  webServer: {
    command: "npm run build && npm run preview",
    port: 4173
  },
  testDir: "tests",
  testMatch: /(.+\.)?(test|spec)\.[jt]s/
};

export default config;
