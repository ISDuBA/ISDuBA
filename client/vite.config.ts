// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";
import { readFileSync } from "fs";

import { fileURLToPath } from "url";

const packageJSON = fileURLToPath(new URL("package.json", import.meta.url));
const pkg = JSON.parse(readFileSync(packageJSON, "utf8"));
const versionGOFile = fileURLToPath(new URL("../pkg/version/version.go", import.meta.url));
const versionGo = readFileSync(versionGOFile, "utf8");
const versionInfo = versionGo.match(/var SemVersion = "(.*?)"/);
let backendVersion = "Error";
if (versionInfo) backendVersion = versionInfo[1] || "Error";

export default defineConfig({
  server: {
    proxy: {
      "/api/": {
        target: "http://localhost:8081/",
        changeOrigin: true
      }
    }
  },
  plugins: [sveltekit()],
  test: {
    include: ["src/**/*.{test,spec}.{js,ts}"]
  },
  define: {
    __APP_VERSION__: `${JSON.stringify(pkg.version)}`,
    __BACKEND_VERSION__: `${JSON.stringify(backendVersion)}`
  }
});
