<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
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
  let headerColumns: any[] = [];
  let productLines: any[] = [];

  const titleStyles = [
    "w-4 overflow-hidden",
    "w-max",
    "[padding:50%_0] h-0",
    "origin-top-left -rotate-90 mt-[50%] -translate-full whitespace-nowrap block"
  ];

  export let basePath = "";

  onMount(() => {
    appStore.setProductTreeSectionInVisible();
    appStore.resetSelectedProduct();
  });

  $: if ($appStore.webview.doc) {
    const vulnerabilities = [...$appStore.webview.doc.productVulnerabilities];

    headerColumns = vulnerabilities.shift()!;
    productLines = vulnerabilities;
  }

  $: fourCVEs = $appStore.webview.four_cves;
</script>

<div class="crosstable-overview mb-3 mt-3 flex flex-col">
  {#if productLines.length > 0}
    <div class="mb-3 mt-3 flex flex-row">
      <div class="flex flex-wrap items-baseline gap-4">
        <div class="flex flex-row items-baseline">
          <i class="bx bx-check" />
          <span class="ml-1 text-nowrap">Fixed</span>
        </div>
        <div class="flex flex-row items-baseline">
          <i class="bx bx-error" /><span class="ml-1 text-nowrap">Under investigation</span>
        </div>
        <div class="flex flex-row items-baseline">
          <i class="bx bx-x" /><span class="ml-1 text-nowrap">Known affected</span>
        </div>
        <div class="flex flex-row items-baseline">
          <i class="bx bx-minus" /><span class="ml-1 text-nowrap">Not affected</span>
        </div>
        <div class="flex flex-row items-baseline">
          <i class="bx bx-heart" /><span class="ml-1 text-nowrap">Recommended</span>
        </div>
        {#if productLines[0].length > 6}
          <div class="ml-auto flex flex-row items-baseline">
            <Button
              color="light"
              size="sm"
              class={`mr-3 h-7 py-1 text-xs ${renderAllCVEs ? "bg-gray-200 hover:bg-gray-100" : ""}`}
              on:click={() => {
                renderAllCVEs = !renderAllCVEs;
              }}><span class="text-nowrap">All CVEs ({productLines[0].length - 2})</span></Button
            >
          </div>
        {/if}
      </div>
    </div>
    <div class="crosstable flex flex-row">
      <Table noborder striped={true}>
        <TableHead>
          {#each headerColumns as column, index}
            {#if index == 0}
              <TableHeadCell
                class="sticky left-0 z-30 text-nowrap bg-white align-bottom font-normal"
                padding={tablePadding}>{column.content}</TableHeadCell
              >
            {:else if index == 1}
              <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}>
                <div class={titleStyles[0]}>
                  <div class={titleStyles[1]}>
                    <div class={titleStyles[2]}>
                      <div class={titleStyles[3]}>{column.content}</div>
                    </div>
                  </div>
                </div>
              </TableHeadCell>
            {:else if !renderAllCVEs && fourCVEs.includes(column.name)}
              <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}>
                <div class={titleStyles[0]}>
                  <div class={titleStyles[1]}>
                    <div class={titleStyles[2]}>
                      <div class={titleStyles[3]}>
                        <a
                          id={crypto.randomUUID()}
                          href={basePath + "cve-" + encodeURIComponent(column.content)}
                          >{column.content}</a
                        >
                      </div>
                    </div>
                  </div>
                </div>
              </TableHeadCell>
            {:else if renderAllCVEs}
              <TableHeadCell class="text-nowrap font-normal" padding={tablePadding}>
                <div class={titleStyles[0]}>
                  <div class={titleStyles[1]}>
                    <div class={titleStyles[2]}>
                      <div class={titleStyles[3]}>
                        <a
                          id={crypto.randomUUID()}
                          href={basePath + "cve-" + encodeURIComponent(column.content)}
                          >{column.content}
                        </a>
                      </div>
                    </div>
                  </div>
                </div>
              </TableHeadCell>
            {/if}
          {/each}
        </TableHead>
        <TableBody>
          {#each productLines as line}
            <TableBodyRow>
              {#each line as column}
                {#if column.name === "Product"}
                  <TableBodyCell tdClass={tdClass + " sticky left-0 bg-inherit"}>
                    <div class="max-w-1/2 min-w-56 whitespace-normal text-wrap">
                      <a
                        title={$appStore.webview.doc?.productsByID[column.content]}
                        id={crypto.randomUUID()}
                        href={basePath + "product-" + encodeURIComponent(column.content)}
                        >{$appStore.webview.doc?.productsByID[column.content]}
                        ({column.content.length > 20
                          ? column.content.substring(0, 20) + "..."
                          : column.content})</a
                      >
                    </div>
                  </TableBodyCell>
                {:else if column.content === "N.A" && ((!renderAllCVEs && fourCVEs.includes(column.name)) || column.name === "Total")}
                  <TableBodyCell {tdClass}>{column.content}</TableBodyCell>
                {:else if column.content === "N.A" && renderAllCVEs && (fourCVEs.includes(column.name) || column.name === "Total")}
                  <TableBodyCell {tdClass}>{column.content}</TableBodyCell>
                {:else if !renderAllCVEs && (fourCVEs.includes(column.name) || column.name === "Total")}
                  <TableBodyCell {tdClass}>
                    {#if column.content === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                      <i class="bx bx-heart" />
                      <i class="bx b-minus" />
                    {:else}
                      <i
                        class:bx={true}
                        class:bx-x={column.content === ProductStatusSymbol.KNOWN_AFFECTED}
                        class:bx-check={column.content === ProductStatusSymbol.FIXED}
                        class:bx-error={column.content === ProductStatusSymbol.UNDER_INVESTIGATION}
                        class:bx-minus={column.content === ProductStatusSymbol.NOT_AFFECTED}
                        class:bx-heart={column.content === ProductStatusSymbol.RECOMMENDED}
                      />
                    {/if}
                  </TableBodyCell>
                {:else if renderAllCVEs}
                  <TableBodyCell {tdClass}>
                    {#if column.content === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                      <i class="bx bx-heart" />
                      <i class="bx b-minus" />
                    {:else}
                      <i
                        class:bx={true}
                        class:bx-x={column.content === ProductStatusSymbol.KNOWN_AFFECTED}
                        class:bx-check={column.content === ProductStatusSymbol.FIXED}
                        class:bx-error={column.content === ProductStatusSymbol.UNDER_INVESTIGATION}
                        class:bx-minus={column.content === ProductStatusSymbol.NOT_AFFECTED}
                        class:bx-heart={column.content === ProductStatusSymbol.RECOMMENDED}
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
  {/if}
</div>
