<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import CBadge from "$lib/Components/CBadge.svelte";
  import type { Snippet } from "svelte";

  interface Props {
    color: any;
    label: string;
    tooltip: string;
    icon: Snippet;
    onClick: (event: any) => void;
  }
  let { color, label, tooltip, icon, onClick }: Props = $props();
</script>

<button
  class="h-fit w-fit rounded-xs border-0 p-0 hover:bg-transparent"
  onclick={(event: any) => onClick(event)}
>
  <CBadge title={tooltip} class="flex w-fit gap-1" {color}>
    {@render icon()}
    <span>{label}</span>
  </CBadge>
</button>
