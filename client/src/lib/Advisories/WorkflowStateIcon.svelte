<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { ASSESSING, ARCHIVED, DELETE, NEW, READ, REVIEW } from "$lib/workflow";
  import { Analyze, Archive, BookOpen, Certification, Eye, Trash } from "@boxicons/svelte";

  interface Props {
    advisoryState: string;
  }
  let { advisoryState }: Props = $props();
</script>

{#if advisoryState === NEW}
  <Certification pack="filled" />
{:else if advisoryState === READ}
  <Eye />
{:else if advisoryState === ASSESSING}
  <Analyze pack="filled" />
{:else if advisoryState === REVIEW}
  <BookOpen />
{:else if advisoryState === ARCHIVED}
  <Archive />
{:else if advisoryState === DELETE}
  <Trash />
{/if}
