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
  import type { AdvisoryVersion } from "./advisory";

  export let advisoryVersions: AdvisoryVersion[];
  $: reversedAdvisoryVersions = advisoryVersions.toReversed();
  export let publisherNamespace: string;
  export let selectedDocumentVersion: AdvisoryVersion;
  let diffModeActivated = false;
  let firstDocumentIndex: number | undefined;
  let secondDocumentIndex: number | undefined;
  let nextColor = "red";
  const diffButtonBaseClass = "!p-2 h-8 min-w-8 mb-2";
  const versionButtonClass = "dark:text-white text-black hover:text-black border border-solid";
  const redButtonClass = `${versionButtonClass} bg-red-100 group-hover:bg-red-300 border-red-700 dark:bg-red-700 dark:group-hover:bg-red-500 dark:border-red-100`;
  const greenButtonClass = `${versionButtonClass} bg-green-100 group-hover:bg-green-300 border-green-700 dark:bg-green-700 dark:group-hover:bg-green-500 dark:border-green-100`;
  const lightButtonClass = `${versionButtonClass} bg-white group-hover:bg-gray-200 border-gray-700 dark:bg-gray-800 dark:group-hover:bg-gray-600 dark:border-gray-300`;

  const dispatch = createEventDispatcher();
  const navigateToVersion = (version: any) => {
    push(
      `/advisories/${publisherNamespace}/${selectedDocumentVersion.tracking_id}/documents/${version.id}`
    );
  };
  const toggleToolboxActivated = () => {
    diffModeActivated = !diffModeActivated;
    if (diffModeActivated) {
      if (
        reversedAdvisoryVersions[0].version === selectedDocumentVersion.version &&
        reversedAdvisoryVersions[0].tracking_status === selectedDocumentVersion.tracking_status
      ) {
        firstDocumentIndex = 0;
        nextColor = "green";
      } else {
        secondDocumentIndex = reversedAdvisoryVersions.findIndex(
          (advVer: any) =>
            advVer.version === selectedDocumentVersion.version &&
            advVer.tracking_status === selectedDocumentVersion.tracking_status
        );
        // If any version but the first is selected we expect that the most interesting
        // case is to compare the current version with the previous one.
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
      reversedAdvisoryVersions[secondDocumentIndex].version !== selectedDocumentVersion.version
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
                <!-- Show tracking status only if there are at least to documents with same version number -->
                {#if (index > 0 && reversedAdvisoryVersions[index - 1].version === version.version) || (index < reversedAdvisoryVersions.length - 1 && reversedAdvisoryVersions[index + 1].version === version.version)}
                  &nbsp;({version.tracking_status})
                {/if}
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
          {#each reversedAdvisoryVersions as version, index}
            <Button
              class={`${diffButtonBaseClass}`}
              disabled={selectedDocumentVersion.version === version.version &&
                selectedDocumentVersion.tracking_status === version.tracking_status}
              on:click={() => {
                navigateToVersion(version);
              }}
              color="light"
              title={`Version ${version.version}`}
            >
              {version.version}
              <!-- Show tracking status only if there are at least to documents with same version number -->
              {#if (index > 0 && reversedAdvisoryVersions[index - 1].version === version.version) || (index < reversedAdvisoryVersions.length - 1 && reversedAdvisoryVersions[index + 1].version === version.version)}
                &nbsp;({version.tracking_status})
              {/if}
            </Button>
          {/each}
        {/if}
        {#if (appStore.isEditor() || appStore.isReviewer()) && advisoryVersions.length > 1}
          <Button color="light" class="flex h-8 gap-x-2 px-3" on:click={toggleToolboxActivated}>
            <i class="bx bx-transfer"></i>
            <span class="text-nowrap">{diffModeActivated ? "Hide" : "Show"} changes</span>
          </Button>
        {/if}
      </div>
    </div>
  </div>
</div>
