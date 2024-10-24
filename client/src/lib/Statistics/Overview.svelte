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
</script>

<SectionHeader title="Statistics"></SectionHeader>
<hr class="mb-6" />
<ImportStats
  chartType="scatter"
  height="400pt"
  showRangeSelection
  showLegend
  showModeToggle
  stepsInMinutes={180}
  title="Advisories and Documents"
  axes={[{ label: "", types: ["totals"] }]}
></ImportStats>
<ImportStats
  height="400pt"
  showRangeSelection
  showLegend
  showModeToggle
  stepsInMinutes={240}
  title="Imports and CVEs"
  axes={[
    { label: "Imports", types: ["imports", "importFailuresCombined"] },
    { label: "CVEs", types: ["cve"] }
  ]}
></ImportStats>
<ImportStats
  height="400pt"
  showRangeSelection
  isStacked
  showLegend
  showModeToggle
  stepsInMinutes={240}
  title="Critical value of imports"
  axes={[{ label: "", types: ["critical"] }]}
></ImportStats>
