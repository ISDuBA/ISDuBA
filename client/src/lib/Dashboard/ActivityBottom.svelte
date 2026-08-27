<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SSVCBadge from "$lib/Advisories/SSVC/SSVCBadge.svelte";
  import { GalleryVerticalEnd, Message as MessageIcon } from "@boxicons/svelte";

  interface Props {
    comments?: number;
    ssvc?: string;
    versions?: number;
  }

  let { comments = undefined, ssvc = undefined, versions = undefined }: Props = $props();
</script>

<div class="flex items-center gap-4 text-slate-400">
  {#if comments !== undefined}
    <div class="flex items-center gap-1">
      <MessageIcon />
      <span>{comments}</span>
    </div>
  {/if}
  {#if versions !== undefined}
    <div class="flex items-center gap-1">
      <GalleryVerticalEnd />
      <span>{versions}</span>
    </div>
  {/if}
  {#if ssvc}
    <SSVCBadge vector={ssvc}></SSVCBadge>
  {/if}
</div>
