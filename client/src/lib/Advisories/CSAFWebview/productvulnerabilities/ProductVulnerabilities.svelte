<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount, untrack } from "svelte";
  import { appStore } from "$lib/store.svelte";
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
  import { innerLinkStyle } from "./../helpers";
  import Link from "$lib/Components/Link.svelte";
  import { AlertTriangle, Check, Heart, Minus, X } from "@boxicons/svelte";

  interface Props {
    basePath: string;
  }

  let { basePath = "" }: Props = $props();

  const uid = $props.id();

  const tdClass = "whitespace-nowrap py-1 px-2 font-normal";
  const tablePadding = "px-2";
  let renderAllCVEs = $state(false);
  let headerColumns: any[] = $state([]);
  let productLines: any[] = $state([]);

  const titleStyles = [
    "w-4 overflow-hidden",
    "w-max",
    "[padding:50%_0] h-0",
    "origin-top-left -rotate-90 mt-[50%] whitespace-nowrap block"
  ];

  onMount(() => {
    appStore.resetSelectedProduct();
  });

  $effect(() => {
    untrack(() => headerColumns);
    untrack(() => productLines);
    if (appStore.state.webview.doc) {
      const vulnerabilities = [...appStore.state.webview.doc.productVulnerabilities];

      headerColumns = vulnerabilities.shift()!;
      productLines = vulnerabilities;
    }
  });

  let fourCVEs = $derived(appStore.state.webview.four_cves);
</script>

{#snippet symbol(symbol: string)}
  {#if symbol === ProductStatusSymbol.KNOWN_AFFECTED}
    <X />
  {:else if symbol === ProductStatusSymbol.FIXED}
    <Check />
  {:else if symbol === ProductStatusSymbol.UNDER_INVESTIGATION}
    <AlertTriangle />
  {:else if symbol === ProductStatusSymbol.NOT_AFFECTED}
    <Minus />
  {:else if symbol === ProductStatusSymbol.RECOMMENDED}
    <Heart />
  {/if}
{/snippet}

<div class="crosstable-overview mt-3 mb-3 flex flex-col">
  {#if productLines.length > 0}
    <div class="mt-3 mb-3 flex flex-row">
      <div class="flex flex-wrap items-baseline gap-4 text-sm">
        <div class="flex flex-row items-baseline">
          <Check />
          <span class="ml-1 text-nowrap">Fixed</span>
        </div>
        <div class="flex flex-row items-baseline">
          <AlertTriangle />
          <span class="ml-1 text-nowrap">Under investigation</span>
        </div>
        <div class="flex flex-row items-baseline">
          <X />
          <span class="ml-1 text-nowrap">Known affected</span>
        </div>
        <div class="flex flex-row items-baseline">
          <Minus />
          <span class="ml-1 text-nowrap">Not affected</span>
        </div>
        <div class="flex flex-row items-baseline">
          <Heart />
          <span class="ml-1 text-nowrap">Recommended</span>
        </div>
        {#if productLines[0].length > 6}
          <div class="ml-auto flex flex-row items-baseline">
            <Button
              color="light"
              size="sm"
              class={`mr-3 h-7 py-1 text-xs ${renderAllCVEs ? "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700" : ""}`}
              onclick={() => {
                renderAllCVEs = !renderAllCVEs;
              }}><span class="text-nowrap">All CVEs ({productLines[0].length - 2})</span></Button
            >
          </div>
        {/if}
      </div>
    </div>
    <div class="crosstable mx-auto flex flex-row lg:mx-0">
      <Table border={false} striped={true}>
        <TableHead>
          {#each headerColumns as column, index (`productvulnerabilities-${uid}-${index}`)}
            {#if index == 0}
              <TableHeadCell
                class="sticky left-0 z-30 bg-white align-bottom font-normal text-nowrap dark:bg-gray-800"
                padding={tablePadding}>{column.content}</TableHeadCell
              >
            {:else if index == 1}
              <TableHeadCell
                class="bg-white font-normal text-nowrap dark:bg-gray-800"
                padding={tablePadding}
              >
                <div class={titleStyles[0]}>
                  <div class={titleStyles[1]}>
                    <div class={titleStyles[2]}>
                      <div class={titleStyles[3]}>{column.content}</div>
                    </div>
                  </div>
                </div>
              </TableHeadCell>
            {:else if !renderAllCVEs && fourCVEs.includes(column.name)}
              <TableHeadCell
                class="bg-white font-normal text-nowrap dark:bg-gray-800"
                padding={tablePadding}
              >
                <div class={titleStyles[0]}>
                  <div class={titleStyles[1]}>
                    <div class={titleStyles[2]}>
                      <div class={titleStyles[3]}>
                        <Link
                          class={innerLinkStyle}
                          id={crypto.randomUUID()}
                          href={basePath + "cve-" + encodeURIComponent(column.content)}
                          >{column.content}</Link
                        >
                      </div>
                    </div>
                  </div>
                </div>
              </TableHeadCell>
            {:else if renderAllCVEs}
              <TableHeadCell
                class="bg-white font-normal text-nowrap dark:bg-gray-800"
                padding={tablePadding}
              >
                <div class={titleStyles[0]}>
                  <div class={titleStyles[1]}>
                    <div class={titleStyles[2]}>
                      <div class={titleStyles[3]}>
                        <Link
                          class={innerLinkStyle}
                          id={crypto.randomUUID()}
                          href={basePath + "cve-" + encodeURIComponent(column.content)}
                          >{column.content}
                        </Link>
                      </div>
                    </div>
                  </div>
                </div>
              </TableHeadCell>
            {/if}
          {/each}
        </TableHead>
        <TableBody>
          {#each productLines as line, i (`productvulnerabilities-2-${uid}-${i}`)}
            <TableBodyRow>
              {#each line as column, j (`productvulnerabilities-3-${uid}-${j}`)}
                {#if column.name === "Product"}
                  <TableBodyCell class={tdClass + " sticky left-0 bg-inherit"}>
                    <div class="max-w-1/2 min-w-56 text-wrap break-all whitespace-normal">
                      <Link
                        id={crypto.randomUUID()}
                        href={basePath + "product-" + encodeURIComponent(column.content)}
                        class={innerLinkStyle}
                        >{appStore.state.webview.doc?.productsByID[column.content]}
                        ({column.content.length > 20
                          ? column.content.substring(0, 20) + "..."
                          : column.content})</Link
                      >
                    </div>
                  </TableBodyCell>
                {:else if column.content === "N.A" && ((!renderAllCVEs && fourCVEs.includes(column.name)) || column.name === "Total")}
                  <TableBodyCell class={tdClass}>{column.content}</TableBodyCell>
                {:else if column.content === "N.A" && renderAllCVEs && (fourCVEs.includes(column.name) || column.name === "Total")}
                  <TableBodyCell class={tdClass}>{column.content}</TableBodyCell>
                {:else if !renderAllCVEs && (fourCVEs.includes(column.name) || column.name === "Total")}
                  <TableBodyCell class={tdClass}>
                    {#if column.content === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                      <Heart />
                      <Minus />
                    {:else}
                      <!-- May contain more than one status and thus more than one character -->
                      {#each column.content as char, k (`productvulnerabilities-4-${uid}-${k}`)}
                        {@render symbol(char)}
                      {/each}
                    {/if}
                  </TableBodyCell>
                {:else if renderAllCVEs}
                  <TableBodyCell class={tdClass}>
                    {#if column.content === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                      <Heart />
                      <Minus />
                    {:else}
                      {@render symbol(column.content)}
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
