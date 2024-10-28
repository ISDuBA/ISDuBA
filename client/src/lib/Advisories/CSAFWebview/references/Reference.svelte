<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import type { Reference } from "$lib/pmdTypes";
  export let reference: Reference;
  let keys: string[] = [];
  let values: string[] = [];
  if (reference.category) {
    keys.push("Category");
    values.push(reference.category);
  }
  keys.push("Summary");
  values.push(reference.summary);
  keys.push("URL");
  values.push(reference.url);
</script>

<div class="">
  <KeyValue {keys} {values} />
</div>
