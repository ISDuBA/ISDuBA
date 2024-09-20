<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Button } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import DiffVersionIndicator from "$lib/Diff/DiffVersionIndicator.svelte";
  import { appStore } from "$lib/store";

  export let advisoryVersions: any;
  $: reversedAdvisoryVersions = advisoryVersions.toReversed();
  export let publisherNamespace: string;
  export let trackingID: string;
  export let selectedDocumentVersion: string;
  let diffModeActivated = false;
  let firstDocumentIndex: number | undefined;
  let secondDocumentIndex: number | undefined;
  let nextColor = "red";
  const diffButtonBaseClass = "!p-2 h-8 w-8 mb-2";
  const versionButtonClass = "text-black hover:text-black border border-solid";
  const redButtonClass = `${versionButtonClass} bg-red-100 group-hover:bg-red-300 border-red-700`;
  const greenButtonClass = `${versionButtonClass} bg-green-100 group-hover:bg-green-300 border-green-700`;
  const lightButtonClass = `${versionButtonClass} bg-white group-hover:bg-gray-200 border-gray-700`;

  const dispatch = createEventDispatcher();
  const navigateToVersion = (version: any) => {
    push(`/advisories/${publisherNamespace}/${trackingID}/documents/${version.id}`);
  };
  const toggleDiffBoxActivated = () => {
    diffModeActivated = !diffModeActivated;
    if (diffModeActivated) {
      if (reversedAdvisoryVersions[0].version === selectedDocumentVersion) {
        firstDocumentIndex = 0;
        nextColor = "green";
      } else {
        secondDocumentIndex = reversedAdvisoryVersions.findIndex(
          (advVer: any) => advVer.version === selectedDocumentVersion
        );
        if (secondDocumentIndex) {
          firstDocumentIndex = secondDocumentIndex - 1;
        }
      }
      showDiff();
    } else {
      disableDiff();
    }
  };
  const disableDiff = () => {
    if (
      secondDocumentIndex &&
      reversedAdvisoryVersions[secondDocumentIndex].version !== selectedDocumentVersion
    ) {
      navigateToVersion(reversedAdvisoryVersions[secondDocumentIndex]);
    }
    dispatch("disableDiff");
  };
  const showDiff = () => {
    if (
      firstDocumentIndex !== undefined &&
      secondDocumentIndex !== undefined &&
      nextColor === "red"
    ) {
      appStore.setDiffDocA_ID(reversedAdvisoryVersions[secondDocumentIndex].id);
      appStore.setDiffDocB_ID(reversedAdvisoryVersions[firstDocumentIndex].id);
      dispatch("selectedDiffDocuments");
    }
  };
  const selectDiffDocument = (index: number) => {
    let oldFirstDocumentIndex = firstDocumentIndex;
    let oldSecondDocumentIndex = secondDocumentIndex;
    if (nextColor === "red") {
      if (index === secondDocumentIndex) {
        secondDocumentIndex = undefined;
      }
      firstDocumentIndex = index;
      nextColor = "green";
    } else {
      if (index === firstDocumentIndex) {
        firstDocumentIndex = undefined;
      }
      secondDocumentIndex = index;
      nextColor = "red";
    }
    if (
      firstDocumentIndex !== undefined &&
      secondDocumentIndex !== undefined &&
      (oldFirstDocumentIndex !== firstDocumentIndex ||
        oldSecondDocumentIndex !== secondDocumentIndex)
    ) {
      showDiff();
    }
  };
</script>

<div class="mb-2 flex flex-row items-center gap-4">
  <div class="flex items-center">
    <div class="flex">
      <div class="me-2 flex flex-row flex-wrap gap-2">
        {#if diffModeActivated}
          {#each reversedAdvisoryVersions as version, index}
            {@const isDisabled =
              (nextColor === "red" && index === reversedAdvisoryVersions.length - 1) ||
              (nextColor === "green" &&
                ((firstDocumentIndex && index <= firstDocumentIndex) || index === 0))}
            {@const hoverIcon = nextColor === "red" ? "minus" : "plus"}
            <div class="group flex flex-col items-center">
              <Button
                disabled={isDisabled}
                class={`${diffButtonBaseClass} ${index === firstDocumentIndex ? redButtonClass : index === secondDocumentIndex ? greenButtonClass : lightButtonClass}`}
                on:click={() => {
                  selectDiffDocument(index);
                }}
                outline
                title={`Version ${version.version}`}
              >
                {version.version}
              </Button>
              {#if index === firstDocumentIndex}
                <DiffVersionIndicator
                  color="red"
                  {isDisabled}
                  icon="minus"
                  {hoverIcon}
                  permanent={true}
                ></DiffVersionIndicator>
              {:else if index === secondDocumentIndex}
                <DiffVersionIndicator
                  color="green"
                  {isDisabled}
                  icon="plus"
                  {hoverIcon}
                  permanent={true}
                ></DiffVersionIndicator>
              {:else if !isDisabled}
                {#if nextColor === "green"}
                  <DiffVersionIndicator color="gray" icon="plus" hoverIcon={undefined} {isDisabled}
                  ></DiffVersionIndicator>
                {:else}
                  <DiffVersionIndicator color="gray" icon="minus" hoverIcon={undefined} {isDisabled}
                  ></DiffVersionIndicator>
                {/if}
              {/if}
            </div>
          {/each}
        {:else}
          {#each reversedAdvisoryVersions as version}
            <Button
              class={`${diffButtonBaseClass}`}
              disabled={selectedDocumentVersion === version.version}
              on:click={() => {
                navigateToVersion(version);
              }}
              color="light"
              title={`Version ${version.version}`}
            >
              {version.version}
            </Button>
          {/each}
        {/if}
        {#if (appStore.isEditor() || appStore.isReviewer()) && advisoryVersions.length > 1}
          <Button color="light" class="flex h-8 gap-x-2 px-3" on:click={toggleDiffBoxActivated}>
            <i class="bx bx-transfer"></i>
            <span class="text-nowrap">{diffModeActivated ? "Hide" : "Show"} changes</span>
          </Button>
        {/if}
      </div>
    </div>
  </div>
</div>
