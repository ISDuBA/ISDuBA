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
  import type { Acknowledgment } from "../docmodel/docmodeltypes";
  interface Props {
    ack: Acknowledgment;
  }
  let { ack }: Props = $props();

  const { keys, values } = $derived.by(() => {
    const keyArray: Array<string> = [],
      valueArray: Array<string> = [];
    if (ack.names) {
      keyArray.push("Names");
      valueArray.push(ack.names.join(", "));
    }
    if (ack.organization) {
      keyArray.push("Organization");
      valueArray.push(ack.organization);
    }
    if (ack.summary) {
      keyArray.push("Summary");
      valueArray.push(ack.summary);
    }
    return { keys: keyArray, values: valueArray };
  });
</script>

<KeyValue {keys} {values} />
{#if ack.urls}
  <ValueList label="URLs" values={ack.urls} />
{/if}
