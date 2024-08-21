<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Accordion, AccordionItem, Button, ButtonGroup, Label } from "flowbite-svelte";
  import DiffEntry from "./DiffEntry.svelte";
  import type { JsonDiffResult } from "./Diff";
  import LazyEntry from "./LazyEntry.svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { appStore } from "$lib/store";
  import { onMount } from "svelte";

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
  $: addChanges = diff ? diff.filter((result: JsonDiffResult) => result.op === "add") : [];
  $: removeChanges = diff ? diff.filter((result: JsonDiffResult) => result.op === "remove") : [];
  $: replaceChanges = diff ? diff.filter((result: JsonDiffResult) => result.op === "replace") : [];
  $: diffDocuments, getDiff();
  $: docA_ID = $appStore.app.diff.docA_ID;
  $: docA_ID, compare();
  $: docB_ID = $appStore.app.diff.docB_ID;
  $: docB_ID, compare();

  const compare = async () => {
    if ($appStore.app.diff.docA_ID && $appStore.app.diff.docB_ID) {
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
    }
  };

  const getDocument = async (letter: string) => {
    const docID = letter === "A" ? $appStore.app.diff.docA_ID : $appStore.app.diff.docB_ID;
    const endpoint = docID?.startsWith("tempdocument") ? "tempdocuments" : "documents";
    const id = docID?.startsWith("tempdocument") ? docID.replace("tempdocument", "") : docID;
    return request(`/api/${endpoint}/${id}`, "GET");
  };

  onMount(async () => {
    compare();
  });

  const getDiff = async () => {
    if (!diffDocuments) return;
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
      return `${bodyClass} bg-green-100`;
    } else if (operation === "remove") {
      return `${bodyClass} bg-red-100`;
    } else {
      return `${bodyClass} bg-gray-100`;
    }
  };
</script>

<svelte:head>
  <title>Compare</title>
</svelte:head>

<div>
  <ErrorMessage {error}></ErrorMessage>
  {#if diff}
    {#if showTitle}
      <Label class="text-lg">{title}</Label>
    {/if}
    <span class={`${title ? "text-gray-700" : "text-sm text-gray-500"}`}>{diff.length} changes</span
    >
    <Accordion flush multiple class={title ? "mt-8" : "mt-1"}>
      <AccordionItem
        paddingFlush="pt-0 pb-3"
        bind:open={isAddSectionOpen}
        class="justify-start gap-x-4 text-gray-700"
      >
        <div slot="header">
          <div class="flex items-center gap-2">
            <span>Added ({addChanges.length})</span>
          </div>
        </div>
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
      <AccordionItem
        paddingFlush="py-3"
        bind:open={isRemoveSectionOpen}
        class="justify-start gap-x-4 text-gray-700"
      >
        <div slot="header">
          <div class="flex items-center gap-2">
            <span>Removed ({removeChanges.length})</span>
          </div>
        </div>
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
      <AccordionItem
        paddingFlush="py-3"
        bind:open={isEditedSectionOpen}
        class="justify-start gap-x-4 text-gray-700"
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
