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
  import { fetchDiffEntry, type JsonDiffResult, type JsonDiffResultWrapper } from "./Diff";
  import LazyEntry from "./LazyEntry.svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { appStore } from "$lib/store";
  import DiffGroupedEntry from "./DiffGroupedEntry.svelte";

  export let showTitle = true;
  let title = "";
  let diffDocuments: any;
  let error: ErrorDetails | null;
  let diff: any;
  let urlPath: string;
  let isAddSectionOpen = false;
  let isRemoveSectionOpen = false;
  let isEditedSectionOpen = true;
  let isSideBySideViewActivated = true;
  let pressedButtonClass = "bg-gray-200 hover:bg-gray-100";
  let accordionItemDefaultClass =
    "flex justify-start items-center gap-x-4 text-gray-700 font-semibold w-full";
  let textFlushOpen = "text-gray-500 dark:text-white";
  let isLoading = false;
  let isLoadingRemovedChanges = false;
  $: addChanges = diff ? diff.filter((result: JsonDiffResult) => result.op === "add") : [];
  let groupedAddChanges: JsonDiffResultWrapper[] = [];
  $: if (addChanges) {
    groupChanges(addChanges).then((result) => {
      groupedAddChanges = result;
    });
  }
  let groupedRemoveChanges: JsonDiffResultWrapper[] = [];
  $: removeChanges = diff ? diff.filter((result: JsonDiffResult) => result.op === "remove") : [];
  $: if (removeChanges) {
    isLoadingRemovedChanges = true;
    groupChanges(removeChanges).then((result) => {
      isLoadingRemovedChanges = false;
      groupedRemoveChanges = result;
    });
  }
  $: replaceChanges = diff ? diff.filter((result: JsonDiffResult) => result.op === "replace") : [];
  $: docA_ID = $appStore.app.diff.docA_ID;
  $: docB_ID = $appStore.app.diff.docB_ID;
  $: if (docA_ID && docB_ID) compare();

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
      const from = `${diffDocuments.docB.document.tracking.id} (Version ${diffDocuments.docB.document.tracking.version})`;
      const to = `${diffDocuments.docA.document.tracking.id} (Version ${diffDocuments.docA.document.tracking.version})`;
      title = `Changes from ${from} to ${to}`;
    }
    await getDiff();
    isLoading = false;
  };

  const groupChanges = async (changes: JsonDiffResult[]) => {
    const list: JsonDiffResultWrapper[] = [];
    const groupedChanges: any = {};
    for (let i = 0; i < changes.length; i++) {
      const change: any = structuredClone(changes[i]);
      const splitted = change.path.split("/");
      const lastPart = Number(splitted[splitted.length - 1]);
      let containsNestedObject = false;
      if (!change.value && change.op === "remove") {
        const response = await fetchDiffEntry(urlPath, change.op, change.path);
        if (response.ok) {
          change.value = response.content;
        } else if (response.error) {
          error = getErrorDetails(`Could not load entry.`, response);
        }
      }
      if (change.value && typeof change.value === "object") {
        containsNestedObject = Object.keys(change.value).some((key) => {
          return (
            (typeof change.value[key] === "object" && !Array.isArray(change.value[key])) ||
            (Array.isArray(change.value[key]) && typeof change.value[key][0] === "object")
          );
        });
      }
      if (isNaN(lastPart) || containsNestedObject || typeof change.value === "string") {
        list.push({
          result: change
        });
      } else {
        const key = splitted.slice(0, -1).join("/");
        if (!groupedChanges[key]) groupedChanges[key] = [];
        groupedChanges[key].push(change);
      }
    }
    Object.keys(groupedChanges).forEach((key) => {
      list.push({ result: groupedChanges[key] });
    });
    list.sort((a, b) => {
      const pathA = Array.isArray(a.result) ? a.result[0].path : a.result.path;
      const pathB = Array.isArray(b.result) ? b.result[0].path : b.result.path;
      return pathA.localeCompare(pathB, "en");
    });
    return list;
  };

  const getDocument = async (letter: string) => {
    const docID = letter === "A" ? $appStore.app.diff.docA_ID : $appStore.app.diff.docB_ID;
    const endpoint = docID?.startsWith("tempdocument") ? "tempdocuments" : "documents";
    const id = docID?.startsWith("tempdocument") ? docID.replace("tempdocument", "") : docID;
    return request(`/api/${endpoint}/${id}`, "GET");
  };

  const getDiff = async () => {
    urlPath = `/api/diff/${$appStore.app.diff.docB_ID}/${$appStore.app.diff.docA_ID}?word-diff=true`;
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
      <AccordionItem
        paddingFlush="pt-0 pb-3"
        bind:open={isAddSectionOpen}
        defaultClass={accordionItemDefaultClass}
        {textFlushOpen}
      >
        <div slot="header">
          <div class="flex items-center gap-2">
            <span>Added ({addChanges.length})</span>
          </div>
        </div>
        {#each groupedAddChanges as change}
          {#if !Array.isArray(change.result)}
            <div class={getBodyClass("add")}>
              {#if change.result.value}
                <div class="mb-1 text-sm font-bold">
                  <code>
                    {change.result.path}
                  </code>
                </div>
                <DiffEntry
                  content={change.result.value}
                  {isSideBySideViewActivated}
                  operation={change.result.op}
                ></DiffEntry>
              {/if}
            </div>
          {:else}
            <div class={getBodyClass("add")}>
              <DiffGroupedEntry {change}></DiffGroupedEntry>
            </div>
          {/if}
        {/each}
      </AccordionItem>
      <AccordionItem
        paddingFlush="py-3"
        bind:open={isRemoveSectionOpen}
        defaultClass={accordionItemDefaultClass}
        {textFlushOpen}
      >
        <div slot="header">
          <div class="flex items-center gap-2">
            <span>Removed ({removeChanges.length})</span>
            {#if isLoadingRemovedChanges}
              <span
                class:invisible={!isLoadingRemovedChanges}
                class={isLoadingRemovedChanges ? "loadingFadeIn" : ""}
              >
                <Spinner color="gray" size="4"></Spinner>
              </span>
            {/if}
          </div>
        </div>
        {#each groupedRemoveChanges as change}
          {#if !Array.isArray(change.result)}
            <div class={getBodyClass("remove")}>
              {#if change.result.value}
                <div class="mb-1 text-sm font-bold">
                  <code>
                    {change.result.path}
                  </code>
                </div>
                <DiffEntry
                  content={change.result.value}
                  {isSideBySideViewActivated}
                  operation={change.result.op}
                ></DiffEntry>
              {:else}
                <LazyEntry operation={change.result.op} {urlPath} path={change.result.path}
                ></LazyEntry>
              {/if}
            </div>
          {:else}
            <div class={getBodyClass("remove")}>
              <DiffGroupedEntry {change} isClosable></DiffGroupedEntry>
            </div>
          {/if}
        {/each}
      </AccordionItem>
      <AccordionItem
        paddingFlush="py-3"
        bind:open={isEditedSectionOpen}
        defaultClass={accordionItemDefaultClass}
        {textFlushOpen}
      >
        <div slot="header">
          <div class="flex items-center gap-2">
            <span>Edited ({replaceChanges.length})</span>
            <ButtonGroup>
              <Button
                color="light"
                class={`py-1 text-xs ${isSideBySideViewActivated === true ? pressedButtonClass : ""}`}
                on:click={(event) => {
                  event.stopPropagation();
                  isSideBySideViewActivated = true;
                }}
              >
                Side-by-side
              </Button>
              <Button
                color="light"
                class={`py-1 text-xs ${isSideBySideViewActivated === false ? pressedButtonClass : ""}`}
                on:click={(event) => {
                  event.stopPropagation();
                  isSideBySideViewActivated = false;
                }}
              >
                Inline
              </Button>
            </ButtonGroup>
          </div>
        </div>
        {#each replaceChanges as change}
          <div class={getBodyClass("replace")}>
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
    </Accordion>
  {/if}
</div>
