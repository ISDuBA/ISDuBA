<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import FileHash from "./FileHash.svelte";
  import KeyValue from "$lib/CSAFWebview/KeyValue.svelte";
  import ValueList from "$lib/CSAFWebview/ValueList.svelte";
  import XGenericUri from "./XGenericURI.svelte";
  import { Table, TableBody, TableBodyCell, TableBodyRow } from "flowbite-svelte";
  export let helper: any;
</script>

<div>
  <Table>
    <TableBodyRow>
      <TableBodyCell><h5>Product Version Helper</h5></TableBodyCell>
    </TableBodyRow>
  </Table>
  {#if helper.cpe}
    <KeyValue keys={["cpe"]} values={helper.cpe} />
  {/if}
  {#if helper.hashes}
    {#each helper.hashes as hash}
      <FileHash {hash} />
    {/each}
  {/if}
  {#if helper.model_numbers}
    <ValueList label="Model numbers" values={helper.model_numbers} />
  {/if}
  {#if helper.purl}
    <KeyValue keys={["purl"]} values={helper.purl} />
  {/if}
  {#if helper.sbom_urls}
    <ValueList label="SBOM URLs" values={helper.sbom_urls} />
  {/if}
  {#if helper.serial_numbers}
    <ValueList label="Serial numbers" values={helper.serial_numbers} />
  {/if}
  {#if helper.skus}
    <ValueList label="SKUs" values={helper.skus} />
  {/if}
  {#if helper.x_generic_uris}
    <XGenericUri x_generic_uris={helper.x_generic_uris} />
  {/if}
</div>
