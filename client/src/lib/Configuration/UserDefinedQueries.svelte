<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Card, Radio, StepIndicator } from "flowbite-svelte";
  let currentSearch = {
    currentStep: 1,
    searchType: ""
  };
  const COLUMNS = [""];
  const SEARCHTYPES = {
    ADVISORY: "Advisory",
    DOCUMENT: "Document"
  };
  const STEPS = {
    CHOOSE_TYPE: 1,
    CHOOSE_COLUMNS: 2,
    ENTER_FILE_CRITERIA: 3,
    ENTER_DIRECTIONS: 4
  };

  const steps = [
    "Choose advisories or documents",
    "Choose columns",
    "Enter filter criteria",
    "Enter directions to order"
  ];

  const initCurrentSearch = () => {
    currentSearch = {
      currentStep: 1,
      searchType: ""
    };
  };

  const proceed = () => {
    if (currentSearch.currentStep < 4) currentSearch.currentStep += 1;
  };
  const back = () => {
    if (currentSearch.currentStep > 1) currentSearch.currentStep -= 1;
  };
</script>

<h2 class="mb-3 text-lg">User defined queries</h2>

<Card>
  <h5 class="mb-4 text-xl font-medium text-gray-500 dark:text-gray-400">New Query</h5>
  <StepIndicator currentStep={currentSearch.currentStep} {steps} color="gray" />
  <div class="mt-3 flex flex-col gap-3">
    {#if currentSearch.currentStep == STEPS.CHOOSE_TYPE}
      Would you like to start an...
      <Radio name="example" value={SEARCHTYPES.ADVISORY} bind:group={currentSearch.searchType}
        >Advisorysearch</Radio
      >
      <Radio name="example" value={SEARCHTYPES.DOCUMENT} bind:group={currentSearch.searchType}
        >Documentsearch</Radio
      >
    {/if}
    {#if currentSearch.currentStep == STEPS.CHOOSE_COLUMNS}{/if}
  </div>
  <div class="ml-auto">
    {#if currentSearch.currentStep > 1}
      <Button class="mt-6" color="light" on:click={back}
        ><i class="bx bx-arrow-back me-2 text-xl"></i>Back</Button
      >
    {/if}
    <Button class="mt-3" color="light" outline={true} on:click={proceed}
      ><i class="bx bx-right-arrow-alt me-2 text-xl"></i>Next</Button
    >
  </div>
</Card>
