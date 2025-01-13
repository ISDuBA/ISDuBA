// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { sveltekit } from "@sveltejs/kit/vite";
import { coverageConfigDefaults, defineConfig } from "vitest/config";
import istanbul from "vite-plugin-istanbul";

export default defineConfig({
  server: {
    proxy: {
      "/api/": {
        target: "http://localhost:8081/",
        changeOrigin: true
      }
    }
  },
  plugins: [
    sveltekit(),
    istanbul({
      include: "src/*",
      exclude: ["node_modules", "test/"],
      extension: [".ts", ".svelte"],
      requireEnv: true
    })
  ],
  test: {
    include: ["src/**/*.{test,spec}.{js,ts}"],
    coverage: {
      provider: "v8",
      reporter: ["text", "json-summary", "json", "html"],
      exclude: [
        "**/build/**",
        "**/.svelte-kit/**",
        "**/*.config.*",
        ...coverageConfigDefaults.exclude
      ],
      thresholds: {
        lines: 60,
        branches: 60,
        functions: 60,
        statements: 60
      }
    }
  }
});
