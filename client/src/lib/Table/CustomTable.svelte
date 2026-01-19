<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Table, TableBody, TableHead, TableHeadCell } from "flowbite-svelte";
  import { tablePadding, type TableHeader } from "./defaults";
  import type { Snippet } from "svelte";
  interface Props {
    title?: string;
    headers?: TableHeader[];
    hoverable?: boolean;
    stickyHeaders?: boolean;
    striped?: boolean;
    containerClass?: string;
    tableContainerClass?: string;
    tableClass?: string;
    tableHeadSlot?: Snippet;
    bottomSlot?: Snippet;
    headerRightSlot?: Snippet;
    mainSlot: Snippet;
    topSlot?: Snippet;
  }
  let {
    title = undefined,
    headers = undefined,
    hoverable = true,
    stickyHeaders = false,
    striped = true,
    containerClass = "mb-6",
    tableClass = "h-100 w-fit border-separate border-spacing-0",
    tableContainerClass = "relative",
    tableHeadSlot = undefined,
    bottomSlot = undefined,
    headerRightSlot = undefined,
    mainSlot,
    topSlot = undefined
  }: Props = $props();
  let orderBy = $state("");
</script>

<div class={containerClass}>
  {#if title}
    <SectionHeader {title}>
      {#snippet rightSlot()}
        <div>
          {#if headerRightSlot}
            {@render headerRightSlot()}
          {/if}
        </div>
      {/snippet}
    </SectionHeader>
  {/if}
  {#if topSlot}
    {@render topSlot()}
  {/if}
  <Table
    classes={{
      div: tableContainerClass
    }}
    class={tableClass}
    {hoverable}
    border={false}
    {striped}
  >
    <TableHead
      class={stickyHeaders ? "sticky top-[0] bg-white dark:bg-gray-800" : "dark:bg-gray-800"}
    >
      {#if tableHeadSlot}
        {@render tableHeadSlot()}
      {:else if headers}
        {#each headers as header}
          <TableHeadCell class={header.class ?? tablePadding} onclick={() => {}}>
            <span>{header.label}</span>
            <i
              class:bx={true}
              class:bx-caret-up={orderBy == header.attribute}
              class:bx-caret-down={orderBy == `-${header.attribute}`}
            ></i>
            {#if header.progressDuration}
              <div class="mt-1 h-1 min-h-1">
                <div class="progressmeter">
                  <span class="w-full"
                    ><span
                      style="animation-duration: {header.progressDuration}s"
                      class="infiniteprogress bg-primary-500"
                    ></span></span
                  >
                </div>
              </div>
            {/if}
          </TableHeadCell>
        {/each}
      {/if}
    </TableHead>
    <TableBody>
      {@render mainSlot()}
    </TableBody>
  </Table>
  <div class="mt-2">
    {#if bottomSlot}
      {@render bottomSlot()}
    {/if}
  </div>
</div>
