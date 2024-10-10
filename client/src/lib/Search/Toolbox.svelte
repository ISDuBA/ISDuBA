<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { updateMultipleStates } from "$lib/Advisories/advisory";
  import DiffSelection from "$lib/Diff/DiffSelection.svelte";
  import { appStore } from "$lib/store";
  import { Button, Dropdown, Img, Label, Select } from "flowbite-svelte";
  import { getAllowedWorkflowChanges } from "$lib/permissions";
  import { createEventDispatcher } from "svelte";

  $: selectedDocuments =
    $appStore.app.documents?.filter((d: any) => $appStore.app.selectedDocumentIDs.has(d.id)) ?? [];
  $: allowedWorkflowStateChanges = getAllowedWorkflowChanges(
    selectedDocuments?.map((d: any) => d.state) ?? []
  );
  $: workflowOptions = allowedWorkflowStateChanges.map((c) => {
    return { name: c.to, value: c.to };
  });
  $: docA = $appStore.app.diff.docA;
  $: docB = $appStore.app.diff.docB;

  let selectedState: any;
  let dropdownOpen = false;
  const selectClass =
    "max-w-96 w-fit text-gray-900 disabled:text-gray-400 bg-gray-50 border border-gray-300 rounded-lg focus:ring-primary-500 focus:border-primary-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:disabled:text-gray-500 dark:focus:ring-primary-500 dark:focus:border-primary-500";
  const dispatch = createEventDispatcher();

  const changeWorkflowState = async () => {
    if (!selectedDocuments || selectedDocuments.length < 0) return;
    const changes: any[] = [];
    selectedDocuments?.forEach((doc: any) => {
      changes.push({
        state: selectedState,
        publisher: doc.publisher,
        tracking_id: doc.tracking_id
      });
    });
    updateMultipleStates(changes);
    dispatch("statesUpdated");
    dropdownOpen = false;
  };
</script>

<div class={`sticky bottom-0 left-0 flex flex-col items-start justify-center`}>
  <div class="flex">
    <Button
      on:click={appStore.toggleToolbox}
      class="rounded-none rounded-t-md border-b-0"
      color="light"
    >
      <span class="me-2"
        >Diff {docA
          ? `${docA?.document?.title.substring(0, 25)}${docA?.document?.title.length > 25 ? "..." : ""}`
          : ""}
        {docB
          ? ` - ${docB?.document?.title.substring(0, 25)}${docB?.document?.title.length > 25 ? "..." : ""}`
          : ""}</span
      >
      <Img src="plus-minus.svg" class="h-4 min-w-4" />
    </Button>
    <div class="mx-2 flex items-center gap-2">
      {#if appStore.isAdmin()}
        <Button
          on:click={() => {
            appStore.setDocumentsToDelete(selectedDocuments);
            appStore.setIsDeleteModalOpen(true);
          }}
          class="!p-2"
          color="light"
          disabled={!selectedDocuments || selectedDocuments.length === 0}
        >
          <i class="bx bx-trash text-red-600"></i>
        </Button>
      {/if}
      <Button class="!p-2" color="light" disabled={workflowOptions.length === 0} id="state-icon">
        <i class="bx bx-git-commit text-black-700"></i>
      </Button>
      <Dropdown
        bind:open={dropdownOpen}
        placement="top-start"
        triggeredBy="#state-icon"
        class="w-full max-w-sm divide-y divide-gray-100 rounded p-4 shadow dark:divide-gray-700 dark:bg-gray-800"
        containerClass="divide-y z-50 border border-gray-300"
      >
        <div class="flex flex-col gap-3">
          <Label>
            <span>New workflow state</span>
            <Select
              bind:value={selectedState}
              items={workflowOptions}
              placeholder="Choose..."
              defaultClass={selectClass}
            ></Select>
          </Label>
          <Button
            on:click={() => {
              changeWorkflowState();
            }}
            class="h-fit">Change</Button
          >
        </div>
      </Dropdown>
    </div>
  </div>
  {#if $appStore.app.isToolboxOpen}
    <div
      class="flex min-h-48 w-full min-w-full max-w-[700pt] items-stretch gap-6 rounded-tr-md border border-solid border-gray-300 bg-white p-4 shadow-gray-800 md:min-w-96 lg:w-auto"
    >
      <DiffSelection></DiffSelection>
    </div>
  {/if}
</div>
