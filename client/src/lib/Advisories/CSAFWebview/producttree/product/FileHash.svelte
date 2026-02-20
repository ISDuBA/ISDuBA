<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import type { hashEntry } from "./producttypes";

  interface Props {
    hash: hashEntry[];
  }
  let { hash }: Props = $props();
</script>

{#each hash as hashEntry (hashEntry.value)}
  <KeyValue keys={["algorithm", "value"]} values={[hashEntry.algorithm, hashEntry.value]} />
{/each}
