<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2025 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2025 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, ButtonGroup } from "flowbite-svelte";
  import { SEARCHTYPES } from "$lib/Queries/query";

  interface Props {
    advisoryButtonVisible?: boolean;
    documentButtonVisible?: boolean;
    eventButtonVisible?: boolean;
    advisoryButtonDisabled?: boolean;
    documentButtonDisabled?: boolean;
    eventButtonDisabled?: boolean;
    selected?: SEARCHTYPES;
    onSelect: (type: SEARCHTYPES) => void;
  }

  let {
    advisoryButtonVisible = true,
    documentButtonVisible = true,
    eventButtonVisible = true,
    advisoryButtonDisabled = false,
    documentButtonDisabled = false,
    eventButtonDisabled = false,
    selected = SEARCHTYPES.ADVISORY,
    onSelect
  }: Props = $props();

  const baseClasses = "py-1 text-xs";
  const selectedClass = "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700";
</script>

<ButtonGroup class="h-7">
  {#if advisoryButtonVisible}
    <Button
      disabled={advisoryButtonDisabled}
      color="light"
      class={`${baseClasses} ${selected === SEARCHTYPES.ADVISORY ? selectedClass : ""}`}
      onclick={() => {
        onSelect(SEARCHTYPES.ADVISORY);
      }}>Advisories</Button
    >
  {/if}
  {#if documentButtonVisible}
    <Button
      disabled={documentButtonDisabled}
      color="light"
      class={`${baseClasses} ${selected === SEARCHTYPES.DOCUMENT ? selectedClass : ""}`}
      onclick={() => {
        onSelect(SEARCHTYPES.DOCUMENT);
      }}>Documents</Button
    >
  {/if}
  {#if eventButtonVisible}
    <Button
      disabled={eventButtonDisabled}
      color="light"
      class={`${baseClasses} ${selected === SEARCHTYPES.EVENT ? selectedClass : ""}`}
      onclick={() => {
        onSelect(SEARCHTYPES.EVENT);
      }}>Events</Button
    >
  {/if}
</ButtonGroup>
