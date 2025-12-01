// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { type PlaywrightTestConfig, devices } from "@playwright/experimental-ct-svelte";

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
        storageState: "playwright/.auth/user.json",
        screenshot: "only-on-failure",
        trace: "retain-on-failure"
      },
      dependencies: ["setup"],
      outputDir: "./playwright/output"
    },

    {
      name: "firefox",
      use: {
        ...devices["Desktop Firefox"],
        // Use prepared auth state.
        storageState: "playwright/.auth/user.json",
        screenshot: "only-on-failure",
        trace: "retain-on-failure"
      },
      dependencies: ["setup"],
      outputDir: "./playwright/output"
    }
  ],
  webServer: {
    command: "npm run build:test && npm run preview",
    port: 4173
  },
  testDir: "tests",
  testMatch: /(.+\.)?(test|spec)\.[jt]s/,
  reporter: [
    [
      "monocart-reporter",
      {
        name: "Monocart Reporter test report",
        logging: "info",
        outputFile: "./test-results/monocart-report/index.html",

        coverage: {
          entryFilter: (_entry: any) => true,

          // "sourceFilter" and "all" are modified configs from
          // https://github.com/cenfun/playwright-ct-svelte/blob/main/playwright-ct.config.ts
          sourceFilter: {
            "**/node_modules/**": false,
            "**/*.svg": false,
            "**/src/**": true
          },

          all: {
            dir: "./src",
            filter: {
              "**/*.js": true,
              "**/*.ts": true,
              "**/*.svelte": true
            }
          },

          reports: ["v8", "console-details"]
        }
      }
    ]
  ]
};

export default config;
