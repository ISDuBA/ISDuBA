<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell } from "flowbite-svelte";
  import { getReadableDateString } from "../../CSAFWebview/helpers";
  import SSVCBadge from "$lib/Advisories/SSVC/SSVCBadge.svelte";

  interface Props {
    ssvcData: any;
  }
  let { ssvcData }: Props = $props();

  const getLabel = () => {
    let label = "SSVC ";
    if (!ssvcData.prev_ssvc) {
      label += "added";
    } else {
      label += "changed";
    }
    label += ` (by ${ssvcData.actor})`;
    return label;
  };

  const intlFormat = new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "medium"
  });

  const tdClass = "py-2 px-2";
</script>

<TableBodyCell class={tdClass}>
  <div class="flex flex-col">
    <div class="flex flex-row items-baseline justify-between">
      <small class="w-40 text-xs text-slate-400" title={ssvcData.time}>
        {getReadableDateString(ssvcData.time, intlFormat)}
      </small>
      <div class="flex grow justify-between">
        <small class="ml-1 text-right">
          {getLabel()}
        </small>
        <span class="ml-1 text-xs text-slate-400">
          on version: {ssvcData.documentVersion}
        </span>
      </div>
    </div>
    <div class="flex flex-row items-baseline justify-between">
      <small class="text-right text-[10px] text-gray-400">
        {#if ssvcData.prev_ssvc}
          TO:
        {/if}
        {ssvcData.ssvc}
        {#if ssvcData.prev_ssvc}
          <br /> FROM: {ssvcData.prev_ssvc}
        {/if}
      </small>
      <div>
        {#if ssvcData.prev_ssvc}
          <SSVCBadge vector={ssvcData.prev_ssvc} />
          &rarr;
        {/if}
        <SSVCBadge vector={ssvcData.ssvc} />
      </div>
    </div>
  </div>
</TableBodyCell>
