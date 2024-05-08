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
  import type { JsonDiffResult } from "./JsonDiff";
  import LazyDiffEntry from "./LazyDiffEntry.svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  export let diffDocuments: any;
  export let title: string | undefined;
  let error: string;
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

  const getDiff = async () => {
    urlPath = `/api/diff/${diffDocuments.docB.id}/${diffDocuments.docA.id}?word-diff=true`;
    error = "";
    const response = await request(urlPath, "GET");
    if (response.ok) {
      diff = response.content;
    } else if (response.error) {
      error = response.error;
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

<div>
  <ErrorMessage message={error}></ErrorMessage>
  {#if diff}
    {#if title}
      <Label class="text-lg">{title}</Label>
    {/if}
    <span class={`${title ? "text-gray-700" : "text-lg text-black"}`}
      >{diff.length} changes{title ? "" : ":"}</span
    >
    <Accordion flush multiple>
      <AccordionItem bind:open={isAddSectionOpen} class="justify-start gap-x-4 text-gray-700">
        <div slot="header" class="pl-2">
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
      <AccordionItem bind:open={isRemoveSectionOpen} class="justify-start gap-x-4 text-gray-700">
        <div slot="header" class="pl-2">
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
              <LazyDiffEntry operation={change.op} {urlPath} path={change.path}></LazyDiffEntry>
            {/if}
          </div>
        {/each}
      </AccordionItem>
      <AccordionItem bind:open={isEditedSectionOpen} class="justify-start gap-x-4 text-gray-700">
        <div slot="header" class="pl-2">
          <div class="flex items-center gap-2">
            <span>Edited ({replaceChanges.length})</span>
            <ButtonGroup>
              <Button
                color="light"
                class={`${isSideBySideViewActivated === true ? pressedButtonClass : ""}`}
                on:click={(event) => {
                  event.stopPropagation();
                  isSideBySideViewActivated = true;
                }}
              >
                Side-by-side
              </Button>
              <Button
                color="light"
                class={`${isSideBySideViewActivated === false ? pressedButtonClass : ""}`}
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
              <LazyDiffEntry operation={change.op} {urlPath} path={change.path}></LazyDiffEntry>
            {/if}
          </div>
        {/each}
      </AccordionItem>
    </Accordion>
  {/if}
</div>
