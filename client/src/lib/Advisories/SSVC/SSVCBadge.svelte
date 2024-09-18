<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Badge } from "flowbite-svelte";
  import { convertVectorToSSVCObject } from "./SSVCCalculator";

  export let vector: string;

  $: ssvc = convertVectorToSSVCObject(vector);
  $: style = ssvc.color ? `color: white; background-color: ${ssvc.color};` : "";
</script>

<Badge class="h-6 w-fit" title={ssvc.vector} {style}>{ssvc.label}</Badge>
