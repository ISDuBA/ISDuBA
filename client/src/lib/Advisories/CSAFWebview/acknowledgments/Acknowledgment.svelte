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
  import type { Acknowledgment as Acknowledgment2_0 } from "$lib/Advisories/types/csaf-2.0";
  import type { Acknowledgment as Acknowledgment2_1 } from "$lib/Advisories/types/csaf-2.1";
  interface Props {
    ack: Acknowledgment2_0 | Acknowledgment2_1;
    path: string;
  }
  let { ack, path }: Props = $props();

  const { keys, values, paths } = $derived.by(() => {
    const keyArray: Array<string> = [],
      valueArray: Array<string | Array<string>> = [],
      pathArray: Array<string> = [];
    if (ack.names) {
      keyArray.push("Names");
      valueArray.push(ack.names);
      pathArray.push(`${path}/names`);
    }
    if (ack.organization) {
      keyArray.push("Organization");
      valueArray.push(ack.organization);
      pathArray.push(`${path}/organization`);
    }
    if (ack.summary) {
      keyArray.push("Summary");
      valueArray.push(ack.summary);
      pathArray.push(`${path}/summary`);
    }
    return { keys: keyArray, values: valueArray, paths: pathArray };
  });
</script>

<KeyValue {keys} {values} {paths} />
{#if ack.urls}
  <ValueList label="URLs" values={ack.urls} path={`${path}/urls`} />
{/if}
