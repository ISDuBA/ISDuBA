<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    TableSearch,
    Pagination
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import NewAdvisories from "./NewAdvisories.svelte";
</script>

{#if $appStore.app.keycloak.authenticated}
  <SectionHeader title="New Events"></SectionHeader>
  <Table>
    <TableHead>
      <TableHeadCell padding={tablePadding}>Description</TableHeadCell>
      <TableHeadCell padding={tablePadding}>Advisory</TableHeadCell>
    </TableHead>
    <TableBody>
      <TableBodyRow>
        <TableBodyCell {tdClass}>Comment added</TableBodyCell>
        <TableBodyCell {tdClass}>Sick PSIRT SCA-2022-0032</TableBodyCell>
      </TableBodyRow>
    </TableBody>
  </Table>
  <SectionHeader title="New Advisories"></SectionHeader>
  <NewAdvisories></NewAdvisories>
{/if}
