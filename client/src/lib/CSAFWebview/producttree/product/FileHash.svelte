<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/CSAFWebview/KeyValue.svelte";
  import type { hashEntry } from "./producttypes";

  export let hash: hashEntry[];
</script>

{#each hash as hashEntry}
  <KeyValue keys={["algorithm", "value"]} values={[hashEntry.algorithm, hashEntry.value]} />
{/each}
