<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { Button, Label } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import DiffVersionIndicator from "$lib/Diff/DiffVersionIndicator.svelte";
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
  $: diffButtonClass = diffModeActivated
    ? `${diffButtonBaseClass} bg-gray-800 text-white hover:bg-gray-600 focus-within:ring-transparent`
    : `${diffButtonBaseClass} bg-white text-black border border-solid border-gray-300 hover:bg-gray-200 focus-within:ring-transparent`;
  const versionButtonClass = "text-black hover:text-black border border-solid";
  const redButtonClass = `${versionButtonClass} bg-red-100 hover:bg-red-300 border-red-700`;
  const greenButtonClass = `${versionButtonClass} bg-green-100 hover:bg-green-300 border-green-700`;
  const lightButtonClass = `${versionButtonClass} bg-white hover:bg-gray-200 border-gray-700`;

  const dispatch = createEventDispatcher();
  const navigateToVersion = (version: any) => {
    push(`/advisories/${publisherNamespace}/${trackingID}/documents/${version.id}`);
  };
  const toggleDiffModeActivated = () => {
    diffModeActivated = !diffModeActivated;
    if (diffModeActivated) {
      firstDocumentIndex = reversedAdvisoryVersions.length - 2;
      secondDocumentIndex = reversedAdvisoryVersions.length - 1;
      showDiff();
    } else {
      disableDiff();
    }
  };
  const disableDiff = () => {
    dispatch("disableDiff");
  };
  const showDiff = () => {
    if (firstDocumentIndex !== undefined && secondDocumentIndex !== undefined) {
      dispatch("selectedDiffDocuments", {
        docA: reversedAdvisoryVersions[secondDocumentIndex],
        docB: reversedAdvisoryVersions[firstDocumentIndex]
      });
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

<div class="my-2">
  <Label class="mb-1">Versions</Label>
  <div class="flex items-center">
    <div class="flex">
      <div class="me-2 flex flex-row gap-2">
        {#if diffModeActivated}
          {#each reversedAdvisoryVersions as version, index}
            {@const isDisabled =
              (nextColor === "red" &&
                ((secondDocumentIndex && index >= secondDocumentIndex) ||
                  index === reversedAdvisoryVersions.length - 1)) ||
              (nextColor === "green" &&
                ((firstDocumentIndex && index <= firstDocumentIndex) || index === 0))}
            <div class="group flex flex-col items-center">
              <Button
                disabled={isDisabled}
                class={`${diffButtonBaseClass} ${index === firstDocumentIndex ? redButtonClass : index === secondDocumentIndex ? greenButtonClass : lightButtonClass}`}
                on:click={() => {
                  selectDiffDocument(index);
                }}
                outline
              >
                {version.version}
              </Button>
              {#if index === firstDocumentIndex}
                <DiffVersionIndicator color="red" {isDisabled} icon="minus" permanent={true}
                ></DiffVersionIndicator>
              {:else if index === secondDocumentIndex}
                <DiffVersionIndicator color="green" {isDisabled} icon="plus" permanent={true}
                ></DiffVersionIndicator>
              {:else if !isDisabled}
                {#if nextColor === "green"}
                  <DiffVersionIndicator color="gray" icon="plus" {isDisabled}
                  ></DiffVersionIndicator>
                {:else}
                  <DiffVersionIndicator color="gray" icon="minus" {isDisabled}
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
            >
              {version.version}
            </Button>
          {/each}
        {/if}
      </div>
      {#if advisoryVersions.length > 1}
        <Button class={diffButtonClass} on:click={toggleDiffModeActivated}>
          <i class="bx bx-transfer"></i>
        </Button>
      {/if}
    </div>
  </div>
</div>
