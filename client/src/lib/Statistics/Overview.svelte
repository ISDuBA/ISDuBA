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
  height="400pt"
  enableRangeSelection
  showLegend
  stepsInMinutes={240}
  title="Imports"
  types={["cve", "imports", "importFailures"]}
></ImportStats>
<ImportStats
  height="400pt"
  enableRangeSelection
  isStacked
  showLegend
  stepsInMinutes={240}
  title="Critical value of imports"
  types={["critical"]}
></ImportStats>
