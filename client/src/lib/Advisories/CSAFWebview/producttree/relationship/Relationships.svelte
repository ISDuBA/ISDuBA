<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { Relationship } from "$lib/pmdTypes";
  import Relation from "./Relation.svelte";

  interface Props {
    relationships: Relationship[];
    basePath: string;
  }
  let { basePath = "", relationships }: Props = $props();

  const uid = $props.id();
</script>

{#each relationships as relation, i (`relationships-${uid}-${i}`)}
  <Relation {basePath} {relation} />
{/each}
