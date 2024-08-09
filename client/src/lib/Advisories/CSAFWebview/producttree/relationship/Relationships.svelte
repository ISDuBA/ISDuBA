<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { Relationship } from "$lib/types";
  import Relation from "./Relation.svelte";
  export let relationships: Relationship[];
  export let basePath = "";
</script>

{#each relationships as relation}
  <Relation {basePath} {relation} />
{/each}
