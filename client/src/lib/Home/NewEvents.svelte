<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { tdClass, tablePadding } from "$lib/table/defaults";
</script>

<div class="border border-red-500">
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
</div>
