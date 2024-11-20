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
  import { tablePadding } from "./defaults";

  type TableHeader = {
    label: string;
    attribute: string | undefined;
    class?: string;
    clickCallBack?: () => void;
  };
  export let title: string | undefined = undefined;
  export let headers: TableHeader[];
  export let stickyHeaders = false;
  let orderBy = "";
</script>

<div class="mb-6">
  {#if title}
    <SectionHeader {title}>
      <div slot="right">
        <slot name="header-right"></slot>
      </div>
    </SectionHeader>
  {/if}
  <slot name="top"></slot>
  <Table divClass="relative" hoverable={true} noborder={true}>
    <TableHead theadClass={stickyHeaders ? "sticky top-[0] bg-white dark:bg-gray-800" : ""}>
      {#each headers as header}
        <TableHeadCell class={header.class ?? ""} padding={tablePadding} on:click={() => {}}>
          <span>{header.label}</span>
          <i
            class:bx={true}
            class:bx-caret-up={orderBy == header.attribute}
            class:bx-caret-down={orderBy == `-${header.attribute}`}
          ></i>
        </TableHeadCell>{/each}
    </TableHead>
    <TableBody>
      <slot></slot>
    </TableBody>
  </Table>
  <div class="mt-2">
    <slot name="bottom"></slot>
  </div>
</div>
