<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { ASSESSING, ARCHIVED, DELETE, NEW, READ, REVIEW } from "$lib/workflow";

  interface Props {
    advisoryState: string;
  }
  let { advisoryState }: Props = $props();
</script>

{#if advisoryState === NEW}
  <i class="bx bxs-certification"></i>
{:else if advisoryState === READ}
  <i class="bx bx-show"></i>
{:else if advisoryState === ASSESSING}
  <i class="bx bxs-analyse"></i>
{:else if advisoryState === REVIEW}
  <i class="bx bx-book-open"></i>
{:else if advisoryState === ARCHIVED}
  <i class="bx bx-archive"></i>
{:else if advisoryState === DELETE}
  <i class="bx bx-trash"></i>
{/if}
