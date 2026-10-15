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
  import type { HelperToIdentifyTheProduct as HelperToIdentifyTheProduct2_0 } from "$lib/Advisories/types/csaf-2.0";
  import type { HelperToIdentifyTheProduct as HelperToIdentifyTheProduct2_1 } from "$lib/Advisories/types/csaf-2.1";
  import { appStore } from "$lib/store.svelte";

  interface Props {
    helper: HelperToIdentifyTheProduct2_0 | HelperToIdentifyTheProduct2_1;
    path: string;
  }
  let { helper, path }: Props = $props();

  const uid = $props.id();

  let csafVersion = $derived(appStore.state.webview.doc?.document.csaf_version);
  let extendedPath = $derived(`${path}/product_identification_helper`);
</script>

<div>
  <Table border={false}>
    <TableBodyRow>
      <TableBodyCell><h5>Product identification helper</h5></TableBodyCell>
    </TableBodyRow>
  </Table>
  {#if helper.cpe}
    <KeyValue keys={["cpe"]} values={[helper.cpe]} paths={[`${extendedPath}/cpe`]} />
  {/if}
  {#if helper.hashes}
    {#each helper.hashes as hash, i (`pidh-${uid}-${i}`)}
      <KeyValue
        keys={["Filename"]}
        values={[hash.filename]}
        paths={[`${extendedPath}/hashes[${i}]/filename`]}
      />
      <FileHash hash={hash.file_hashes} path={`${extendedPath}/hashes[${i}]`} />
    {/each}
  {/if}
  {#if helper.model_numbers}
    <ValueList
      label="Model numbers"
      values={helper.model_numbers}
      path={`${extendedPath}/model_numbers`}
    />
  {/if}
  {#if csafVersion === "2.0" && (helper as HelperToIdentifyTheProduct2_0).purl}
    <KeyValue
      keys={["purl"]}
      values={[(helper as HelperToIdentifyTheProduct2_0).purl]}
      paths={[`${extendedPath}/purl`]}
    />
  {:else if csafVersion === "2.1" && (helper as HelperToIdentifyTheProduct2_1).purls}
    <KeyValue
      keys={["purls"]}
      values={(helper as HelperToIdentifyTheProduct2_1).purls}
      paths={[`${extendedPath}/purls`]}
    />
  {/if}
  {#if helper.sbom_urls}
    <ValueList label="SBOM URLs" values={helper.sbom_urls} path={`${extendedPath}/sbom_urls`} />
  {/if}
  {#if helper.serial_numbers}
    <ValueList
      label="Serial numbers"
      values={helper.serial_numbers}
      path={`${extendedPath}/serial_numbers`}
    />
  {/if}
  {#if helper.skus}
    <ValueList label="SKUs" values={helper.skus} path={`${extendedPath}/skus`} />
  {/if}
  {#if helper.x_generic_uris}
    <XGenericUri x_generic_uris={helper.x_generic_uris} path={`${extendedPath}/x_generic_uris`} />
  {/if}
</div>
