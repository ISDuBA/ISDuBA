<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import FileHash from "./FileHash.svelte";
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import ValueList from "$lib/Advisories/CSAFWebview/ValueList.svelte";
  import XGenericUri from "./XGenericURI.svelte";
  import { Table, TableBodyCell, TableBodyRow } from "flowbite-svelte";

  interface Props {
    helper: any;
    path: string;
  }
  let { helper, path }: Props = $props();

  const uid = $props.id();

  let extendedPath = $derived(`${path}/product_identification_helper`);
</script>

<div>
  <Table border={false}>
    <TableBodyRow>
      <TableBodyCell><h5>Product identification helper</h5></TableBodyCell>
    </TableBodyRow>
  </Table>
  {#if helper.cpe}
    <KeyValue keys={["cpe"]} values={helper.cpe} paths={[`${extendedPath}/cpe`]} />
  {/if}
  {#if helper.hashes}
    {#each helper.hashes as hash, i (`pidh-${uid}-${i}`)}
      <FileHash {hash} path={`${path}/hashes`} />
    {/each}
  {/if}
  {#if helper.model_numbers}
    <ValueList label="Model numbers" values={helper.model_numbers} path={`${path}/model_numbers`} />
  {/if}
  {#if helper.purl}
    <KeyValue keys={["purl"]} values={helper.purl} paths={[`${extendedPath}/purl`]} />
  {/if}
  {#if helper.sbom_urls}
    <ValueList label="SBOM URLs" values={helper.sbom_urls} path={`${path}/sbom_urls`} />
  {/if}
  {#if helper.serial_numbers}
    <ValueList
      label="Serial numbers"
      values={helper.serial_numbers}
      path={`${path}/serial_numbers`}
    />
  {/if}
  {#if helper.skus}
    <ValueList label="SKUs" values={helper.skus} path={`${path}/skus`} />
  {/if}
  {#if helper.x_generic_uris}
    <XGenericUri x_generic_uris={helper.x_generic_uris} path={`${extendedPath}/x_generic_uris`} />
  {/if}
</div>
