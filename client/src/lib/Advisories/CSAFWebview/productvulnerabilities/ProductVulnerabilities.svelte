<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { ProductStatusSymbol } from "./productvulnerabilitiestypes";
  import {
    Table,
    TableHead,
    TableHeadCell,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    Button
  } from "flowbite-svelte";
  const tdClass = "whitespace-nowrap py-1 px-2 font-normal";
  const tablePadding = "px-2";
  let renderAllCVEs = false;
  let PRODUCT_AND_TOTAL = 2;
  let CVES_TO_RENDER = 1;
  let MIN_CVES = PRODUCT_AND_TOTAL + CVES_TO_RENDER;
  let headerColumns: string[] = [];
  let productLines: string[][];

  $: if ($appStore.webview.doc) {
    const vulnerabilities = [...$appStore.webview.doc.productVulnerabilities];
    // eslint-disable-next-line  @typescript-eslint/no-non-null-assertion
    headerColumns = vulnerabilities.shift()!;
    productLines = vulnerabilities;
  }
  /**
   * openProduct opens the according product given via href.
   * @param e
   */
  const openProduct = (e: Event) => {
    // eslint-disable-next-line  @typescript-eslint/no-non-null-assertion
    let product: string = (e.target as Element).getAttribute("href")!;
    appStore.setProductTreeSectionVisible();
    appStore.setSelectedProduct(product);
    appStore.unshiftHistory((e.target as Element).id);
    e.preventDefault();
  };

  /**
   * openCVE opens the CVE given via href.
   * @param e
   */
  const openCVE = (e: Event) => {
    // eslint-disable-next-line  @typescript-eslint/no-non-null-assertion
    let CVE: string = (e.target as Element).getAttribute("href")!;
    appStore.setSelectedCVE(CVE);
    appStore.unshiftHistory((e.target as Element).id);
    appStore.setVulnerabilitiesSectionVisible();
    e.preventDefault();
  };
</script>

<div class="crosstable-overview mt-3 flex flex-row">
  {#if productLines.length > 0}
    <div class="flex w-3/4 flex-col">
      <div class="crosstable">
        <Button
          color="light"
          size="sm"
          class={`mb-3 h-7 py-1 text-xs ${renderAllCVEs ? "bg-gray-200 hover:bg-gray-100" : ""}`}
          on:click={() => {
            renderAllCVEs = !renderAllCVEs;
          }}>All CVEs</Button
        >
        <Table noborder>
          <TableHead>
            {#each headerColumns as column, index}
              {#if index == 0}
                <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}
                  >{column}</TableHeadCell
                >
              {:else if index == 1}
                <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}
                  >{column}</TableHeadCell
                >
              {:else if !renderAllCVEs && index < MIN_CVES}
                <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}
                  ><a id={crypto.randomUUID()} on:click={openCVE} href={column}>{column}</a
                  ></TableHeadCell
                >
              {:else if renderAllCVEs}
                <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}
                  ><a id={crypto.randomUUID()} on:click={openCVE} href={column}>{column}</a
                  ></TableHeadCell
                >
              {/if}
            {/each}
          </TableHead>
          <TableBody>
            {#each productLines as line}
              <TableBodyRow>
                {#each line as column, index}
                  {#if index < 1}
                    <TableBodyCell {tdClass}
                      ><a
                        title={$appStore.webview.doc?.productsByID[column]}
                        id={crypto.randomUUID()}
                        on:click={openProduct}
                        href={column}
                        >{`${$appStore.webview.doc?.productsByID[column].substring(0, 20)}...`} ({column})</a
                      ></TableBodyCell
                    >
                  {:else if column === "N.A" && !renderAllCVEs && index < MIN_CVES}
                    <TableBodyCell {tdClass}>{column}</TableBodyCell>
                  {:else if column === "N.A" && renderAllCVEs && index < MIN_CVES}
                    <TableBodyCell {tdClass}>{column}</TableBodyCell>
                  {:else if !renderAllCVEs && index < MIN_CVES}
                    <TableBodyCell {tdClass}>
                      {#if column === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                        <i class="bx bx-heart" />
                        <i class="bx b-minus" />
                      {:else}
                        <i
                          class:bx={true}
                          class:bx-x={column === ProductStatusSymbol.KNOWN_AFFECTED}
                          class:bx-check={column === ProductStatusSymbol.FIXED}
                          class:bx-error={column === ProductStatusSymbol.UNDER_INVESTIGATION}
                          class:bx-minus={column === ProductStatusSymbol.NOT_AFFECTED}
                          class:bx-heart={column === ProductStatusSymbol.RECOMMENDED}
                        />
                      {/if}
                    </TableBodyCell>
                  {:else if renderAllCVEs}
                    <TableBodyCell {tdClass}>
                      {#if column === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                        <i class="bx bx-heart" />
                        <i class="bx b-minus" />
                      {:else}
                        <i
                          class:bx={true}
                          class:bx-x={column === ProductStatusSymbol.KNOWN_AFFECTED}
                          class:bx-check={column === ProductStatusSymbol.FIXED}
                          class:bx-error={column === ProductStatusSymbol.UNDER_INVESTIGATION}
                          class:bx-minus={column === ProductStatusSymbol.NOT_AFFECTED}
                          class:bx-heart={column === ProductStatusSymbol.RECOMMENDED}
                        />
                      {/if}
                    </TableBodyCell>
                  {/if}
                {/each}
              </TableBodyRow>
            {/each}
          </TableBody>
        </Table>
      </div>
    </div>
    <div class="ml-6 flex w-1/4 flex-col">
      <h5>Legend</h5>
      <div class="flex flex-col">
        <div class="flex flex-row">
          <i class="bx bx-check" />
          <span class="ml-2 text-nowrap">Fixed</span>
        </div>
        <div class="flex flex-row">
          <i class="bx bx-error" /><span class="ml-2 text-nowrap">Under investigation</span>
        </div>
        <div class="flex flex-row">
          <i class="bx bx-x" /><span class="ml-2 text-nowrap">Known affected</span>
        </div>
        <div class="flex flex-row">
          <i class="bx bx-minus" /><span class="ml-2 text-nowrap">Not affected</span>
        </div>
        <div class="flex flex-row">
          <i class="bx bx-heart" /><span class="ml-2 text-nowrap">Recommended</span>
        </div>
      </div>
    </div>
  {/if}
</div>
