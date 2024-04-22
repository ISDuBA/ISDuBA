<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Card, Radio, StepIndicator, MultiSelect } from "flowbite-svelte";

  let currentSearch = {
    currentStep: 1,
    searchType: "",
    chosenColumns: [],
    activeColumns: []
  };

  const COLUMNS = {
    ADVISORY: [
      "id",
      "tracking_id",
      "version",
      "publisher",
      "current_release_date",
      "initial_release_date",
      "title",
      "tlp",
      "cvss_v2_score",
      "cvss_v3_score",
      "ssvc",
      "four_cves",
      "state"
    ],
    DOCUMENT: [
      "id",
      "tracking_id",
      "version",
      "publisher",
      "current_release_date",
      "initial_release_date",
      "title",
      "tlp",
      "cvss_v2_score",
      "cvss_v3_score",
      "four_cves"
    ]
  };
  const SEARCHTYPES = {
    ADVISORY: "Advisory",
    DOCUMENT: "Document"
  };
  const STEPS = {
    CHOOSE_TYPE: 1,
    CHOOSE_COLUMNS: 2,
    SPECIFY_FILE_CRITERIA: 3
  };

  const steps = [
    "Choose advisories or documents",
    "Choose columns to include in the search",
    "Specify filter criteria"
  ];

  const reset = () => {
    initCurrentSearch();
  };

  const initCurrentSearch = () => {
    currentSearch = {
      currentStep: 1,
      searchType: "",
      chosenColumns: [],
      activeColumns: []
    };
  };

  const proceed = () => {
    if (currentSearch.currentStep < 3) currentSearch.currentStep += 1;
  };

  const back = () => {
    if (currentSearch.currentStep > 1) currentSearch.currentStep -= 1;
  };

  let proceedConditionMet = true;

  $: {
    if (currentSearch.currentStep === STEPS.CHOOSE_TYPE) {
      proceedConditionMet = currentSearch.searchType !== "";
    }
    if (currentSearch.currentStep === STEPS.CHOOSE_COLUMNS) {
      proceedConditionMet = currentSearch.chosenColumns.length > 0;
    }
  }

  $: {
    const generateSelectable = (el: string) => {
      return { name: el, value: el };
    };
    if (currentSearch.searchType === SEARCHTYPES.ADVISORY) {
      currentSearch.activeColumns = COLUMNS.ADVISORY.map(generateSelectable);
    }
    if (currentSearch.searchType === SEARCHTYPES.DOCUMENT) {
      currentSearch.activeColumns = COLUMNS.DOCUMENT.map(generateSelectable);
    }
  }
</script>

<h2 class="mb-3 text-lg">User defined queries</h2>

<Card size="lg">
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
    {#if currentSearch.currentStep == STEPS.CHOOSE_COLUMNS}
      <MultiSelect items={currentSearch.activeColumns} bind:value={currentSearch.chosenColumns} />
    {/if}
  </div>
  <div class="ml-auto flex gap-3">
    {#if currentSearch.currentStep > 1}
      <Button class="mt-6" color="light" on:click={back}
        ><i class="bx bx-arrow-back text-xl"></i>Back</Button
      >
      <Button class="mt-6" color="light" on:click={reset}
        ><i class="bx bx-undo text-xl"></i>Reset</Button
      >
    {/if}
    {#if currentSearch.currentStep < 3}
      <Button
        disabled={!proceedConditionMet}
        class="mt-3"
        color="light"
        outline={true}
        on:click={proceed}><i class="bx bx-right-arrow-alt text-xl"></i>Next</Button
      >
    {/if}
    {#if currentSearch.currentStep === 3}
      <Button class="mt-3" color="light" outline={true}
        ><i class="bx bx-save text-xl"></i>Finish</Button
      >
    {/if}
  </div>
</Card>
