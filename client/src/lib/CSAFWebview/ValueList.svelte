<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { P } from "flowbite-svelte";
  import { List, Li } from "flowbite-svelte";
  export let values: any = [];
  export let label = "";
</script>

<P>{label}</P>
<List>
  {#each values as value}
    <Li>{value}</Li>
  {/each}
</List>
