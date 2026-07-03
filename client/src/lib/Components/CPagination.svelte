<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { ArrowToLeft, ArrowToRight, ChevronsLeft, ChevronsRight } from "@boxicons/svelte";
  import { PaginationItem } from "flowbite-svelte";

  interface Props {
    currentPage: number;
    numberOfPages: number;
    onChange: (event: any) => void;
    onFirst: (event: any) => void;
    onPrevious: (event: any) => void;
    onNext: (event: any) => void;
    onLast: (event: any) => void;
  }

  let {
    currentPage = $bindable(0),
    numberOfPages,
    onChange,
    onFirst,
    onPrevious,
    onNext,
    onLast
  }: Props = $props();

  const paginationItemClass =
    "text-gray-500 bg-white hover:bg-gray-100 hover:text-gray-700 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white";
  const paginationItemDeactivatedClass =
    "text-gray-400 bg-gray-50 dark:border-gray-700 dark:text-gray-400 dark:bg-gray-700 cursor-not-allowed";
</script>

<div class="flex flex-row flex-wrap items-center">
  <div class:flex={true} class:mr-3={true}>
    <PaginationItem
      class={currentPage === 1 ? paginationItemDeactivatedClass : paginationItemClass}
      onclick={onFirst}
    >
      <ArrowToLeft />
    </PaginationItem>
    <PaginationItem
      class={currentPage === 1 ? paginationItemDeactivatedClass : paginationItemClass}
      onclick={onPrevious}
    >
      <ChevronsLeft />
    </PaginationItem>
  </div>

  <div class="flex flex-row flex-wrap items-center">
    <input
      class={`${numberOfPages < 10000 ? "w-16" : "w-20"} cursor-pointer border pr-1 text-right dark:bg-gray-800 dark:text-white`}
      onchange={onChange}
      bind:value={currentPage}
    />
    <span class="mr-3 ml-2 w-max text-nowrap">of {numberOfPages} pages</span>
  </div>

  <div class:flex={true}>
    <PaginationItem
      class={currentPage === numberOfPages ? paginationItemDeactivatedClass : paginationItemClass}
      onclick={onNext}
    >
      <ChevronsRight />
    </PaginationItem>
    <PaginationItem
      class={currentPage === numberOfPages ? paginationItemDeactivatedClass : paginationItemClass}
      onclick={onLast}
    >
      <ArrowToRight />
    </PaginationItem>
  </div>
</div>
