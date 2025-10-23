<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Accordion, AccordionItem, Button, ButtonGroup, Label, Spinner } from "flowbite-svelte";
  import DiffEntry from "./DiffEntry.svelte";
  import type { JsonDiffResult } from "./Diff";
  import LazyEntry from "./LazyEntry.svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { appStore } from "$lib/store.svelte";

  interface Props {
    showTitle?: boolean;
  }

  let { showTitle = true }: Props = $props();
  let title = $state("");
  let diffDocuments: any;
  let error: ErrorDetails | null = $state(null);
  let diff: any = $state();
  let urlPath: string = $state("");
  let isAddSectionOpen = $state(false);
  let isRemoveSectionOpen = $state(false);
  let isEditedSectionOpen = $state(true);
  let isSideBySideViewActivated = $state(true);
  let pressedButtonClass = "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700";
  let accordionItemDefaultClass =
    "flex justify-start items-center gap-x-4 text-gray-700 font-semibold w-full";
  let isLoading = $state(false);

  const getPartOfTitle = (document: any, showTrackingStatus: boolean) => {
    return `${document.tracking.id} (Version ${document.tracking.version}${showTrackingStatus ? ", " + document.tracking.status : ""})`;
  };

  const compare = async () => {
    isLoading = true;
    const responseDocA = await getDocument("A");
    const responseDocB = await getDocument("B");
    if (responseDocA.ok && responseDocB.ok) {
      const documentA = await responseDocA.content;
      const documentB = await responseDocB.content;
      diffDocuments = {
        docA: documentB,
        docB: documentA
      };
      const showTrackingStatus =
        diffDocuments.docB.document.tracking.status !== diffDocuments.docA.document.tracking.status;
      const from = getPartOfTitle(diffDocuments.docB.document, showTrackingStatus);
      const to = getPartOfTitle(diffDocuments.docA.document, showTrackingStatus);
      title = `Changes from ${from} to ${to}`;
    }
    await getDiff();
    isLoading = false;
  };

  const getDocument = async (letter: string) => {
    const docID =
      letter === "A" ? appStore.state.app.diff.docA_ID : appStore.state.app.diff.docB_ID;
    const endpoint = docID?.startsWith("tempdocument") ? "tempdocuments" : "documents";
    const id = docID?.startsWith("tempdocument") ? docID.replace("tempdocument", "") : docID;
    return request(`/api/${endpoint}/${id}`, "GET");
  };

  const getDiff = async () => {
    urlPath = `/api/diff/${appStore.state.app.diff.docB_ID}/${appStore.state.app.diff.docA_ID}?word-diff=true`;
    error = null;
    const response = await request(urlPath, "GET");
    if (response.ok) {
      diff = response.content;
    } else if (response.error) {
      error = getErrorDetails(`Couldn't load diffs.`, response);
    }
  };

  const getBodyClass = (operation: string) => {
    let bodyClass = "mb-4 p-2";
    if (operation === "add") {
      return `${bodyClass} bg-green-100 dark:bg-[#1a363c]`;
    } else if (operation === "remove") {
      return `${bodyClass} bg-red-100 dark:bg-[#412732]`;
    } else {
      return `${bodyClass} bg-gray-100 dark:bg-gray-700`;
    }
  };
  let addChanges = $derived(
    diff ? diff.filter((result: JsonDiffResult) => result.op === "add") : []
  );
  let removeChanges = $derived(
    diff ? diff.filter((result: JsonDiffResult) => result.op === "remove") : []
  );
  let replaceChanges = $derived(
    diff ? diff.filter((result: JsonDiffResult) => result.op === "replace") : []
  );
  let docA_ID = $derived(appStore.state.app.diff.docA_ID);
  let docB_ID = $derived(appStore.state.app.diff.docB_ID);
  $effect(() => {
    if (docA_ID && docB_ID) compare();
  });
</script>

<svelte:head>
  <title>Compare</title>
</svelte:head>

<div>
  {#if isLoading}
    <div>
      Loading diff ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
  {/if}
  <ErrorMessage {error}></ErrorMessage>
  {#if diff}
    {#if showTitle}
      <Label class="text-lg">{title}</Label>
    {/if}
    <span class={`${title ? "text-gray-700 dark:text-gray-300" : "text-sm text-gray-500"}`}
      >{diff.length} changes</span
    >
    <Accordion flush multiple class={title ? "mt-8" : "mt-1"}>
      <AccordionItem bind:open={isAddSectionOpen} class={accordionItemDefaultClass}>
        {#snippet header()}
          <div>
            <div class="flex items-center gap-2">
              <span>Added ({addChanges.length})</span>
            </div>
          </div>
        {/snippet}
        {#each addChanges as change}
          <div class={getBodyClass("add")}>
            {#if change.value}
              <div class="mb-1 text-sm font-bold">
                <code>
                  {change.path}
                </code>
              </div>
              <DiffEntry content={change.value} {isSideBySideViewActivated} operation={change.op}
              ></DiffEntry>
            {/if}
          </div>
        {/each}
      </AccordionItem>
      <AccordionItem bind:open={isRemoveSectionOpen} class={accordionItemDefaultClass}>
        {#snippet header()}
          <div>
            <div class="flex items-center gap-2">
              <span>Removed ({removeChanges.length})</span>
            </div>
          </div>
        {/snippet}
        {#each removeChanges as change}
          <div class={getBodyClass("remove")}>
            {#if change.value}
              <div class="mb-1 text-sm font-bold">
                <code>
                  {change.path}
                </code>
              </div>
              <DiffEntry content={change.value} {isSideBySideViewActivated} operation={change.op}
              ></DiffEntry>
            {:else}
              <LazyEntry operation={change.op} {urlPath} path={change.path}></LazyEntry>
            {/if}
          </div>
        {/each}
      </AccordionItem>
      <AccordionItem bind:open={isEditedSectionOpen} class={accordionItemDefaultClass}>
        {#snippet header()}
          <div>
            <div class="flex items-center gap-2">
              <span>Edited ({replaceChanges.length})</span>
              <ButtonGroup>
                <Button
                  color="light"
                  class={`py-1 text-xs ${isSideBySideViewActivated === true ? pressedButtonClass : ""}`}
                  onclick={(event: any) => {
                    event.stopPropagation();
                    isSideBySideViewActivated = true;
                  }}
                >
                  Side-by-side
                </Button>
                <Button
                  color="light"
                  class={`py-1 text-xs ${isSideBySideViewActivated === false ? pressedButtonClass : ""}`}
                  onclick={(event: any) => {
                    event.stopPropagation();
                    isSideBySideViewActivated = false;
                  }}
                >
                  Inline
                </Button>
              </ButtonGroup>
            </div>
          </div>
        {/snippet}
        {#each replaceChanges as change}
          <div class={getBodyClass("replace")}>
            {#if change.value}
              <div class="mb-1 text-sm font-bold dark:text-gray-200">
                <code>
                  {change.path}
                </code>
              </div>
              <DiffEntry content={change.value} {isSideBySideViewActivated} operation={change.op}
              ></DiffEntry>
            {:else}
              <LazyEntry operation={change.op} {urlPath} path={change.path}></LazyEntry>
            {/if}
          </div>
        {/each}
      </AccordionItem>
    </Accordion>
  {/if}
</div>
