<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { StepIndicator } from "flowbite-svelte";
  let currentStep = 2;
  let steps = ["Step 1", "Step 2", "Step 3", "Step 4", "Step 5"];
</script>

<SectionHeader title="Configuration"></SectionHeader>
<hr class="mb-6" />
<div class="mt-3">
  <h2 class="text-lg">User defined queries</h2>
</div>
