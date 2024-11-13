<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import ImportStats from "./ImportStats.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { DAY_MS } from "$lib/time";
</script>

<SectionHeader title="Statistics"></SectionHeader>
<hr class="mb-6" />
<ImportStats
  chartType="line"
  colors={["#3D6090", "#E6A776"]}
  height="300pt"
  initialFrom={new Date(Date.now() - DAY_MS)}
  showRangeSelection
  showLegend
  showModeToggle
  stepsInMinutes={60}
  title="Advisories and Documents"
  axes={[{ label: "", types: ["totals"] }]}
></ImportStats>
<ImportStats
  colors={["#5fcc57", "#c64f08", "#2e4899"]}
  height="300pt"
  initialFrom={new Date(Date.now() - DAY_MS)}
  showRangeSelection
  showLegend
  showModeToggle
  stepsInMinutes={60}
  title="Imports and CVEs"
  axes={[
    { label: "Imports", types: ["imports", "importFailuresCombined"] },
    { label: "CVEs", types: ["cve"] }
  ]}
></ImportStats>
<ImportStats
  height="300pt"
  initialFrom={new Date(Date.now() - DAY_MS)}
  showRangeSelection
  isStacked
  showLegend
  showModeToggle
  stepsInMinutes={60}
  title="Critical value of imports"
  axes={[{ label: "", types: ["critical"] }]}
></ImportStats>
