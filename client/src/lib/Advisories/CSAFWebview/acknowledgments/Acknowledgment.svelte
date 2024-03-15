<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import ValueList from "$lib/Advisories/CSAFWebview/ValueList.svelte";
  import type { Acknowledgement } from "../docmodel/docmodeltypes";
  export let ack: Acknowledgement;
  const keys: Array<String> = [];
  const values: string[] = [];
  if (ack.names) {
    keys.push("Names");
    values.push(ack.names.join(", "));
  }
  if (ack.organization) {
    keys.push("Organization");
    values.push(ack.organization);
  }
  if (ack.summary) {
    keys.push("Summary");
    values.push(ack.summary);
  }
</script>

<KeyValue {keys} {values} />
{#if ack.urls}
  <ValueList label="URLs" values={ack.urls} />
{/if}
