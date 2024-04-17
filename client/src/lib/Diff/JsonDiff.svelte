<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Accordion, AccordionItem, Label } from "flowbite-svelte";
  import DiffEntry from "./DiffEntry.svelte";
  import type { JsonDiffResult, JsonDiffResultList } from "./JsonDiff";
  import LazyDiffEntry from "./LazyDiffEntry.svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";

  export let diffDocuments: any;
  export let title: string;
  let error: string;
  let diff: any;
  let urlPath: string;
  $: groupedResults = diff
    ? [
        {
          op: "add",
          changes: diff.filter((result: JsonDiffResult) => result.op === "add")
        },
        {
          op: "replace",
          changes: diff.filter((result: JsonDiffResult) => result.op === "replace")
        },
        {
          op: "remove",
          changes: diff.filter((result: JsonDiffResult) => result.op === "remove")
        }
      ]
    : [];
  $: diffDocuments, getDiff();

  const getDiff = async () => {
    urlPath = `/api/diff/${diffDocuments.docB.id}/${diffDocuments.docA.id}?word-diff=true`;
    error = "";
    const response = await request(urlPath, "GET");
    if (response.ok) {
      const result: JsonDiffResultList = response.content;
      diff = result;
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
      {#each groupedResults as result}
        {#if result.changes.length > 0}
          <AccordionItem open>
            <div slot="header" class="pl-2">
              {#if result.op === "add"}
                <div class="flex items-center gap-2">
                  <span>Added ({result.changes.length})</span>
                </div>
              {:else if result.op === "remove"}
                <div class="flex items-center gap-2">
                  <span>Removed ({result.changes.length})</span>
                </div>
              {:else}
                <div class="flex items-center gap-2">
                  <span>Edited ({result.changes.length})</span>
                </div>
              {/if}
            </div>
            {#each result.changes as change}
              <div class={getBodyClass(change.op)}>
                {#if change.value}
                  <div class="mb-1">
                    <b>
                      <code>
                        {change.path}
                      </code>
                    </b>
                  </div>
                  <DiffEntry content={change.value} operation={change.op}></DiffEntry>
                {:else}
                  <LazyDiffEntry operation={change.op} {urlPath} path={change.path}></LazyDiffEntry>
                {/if}
              </div>
            {/each}
          </AccordionItem>
        {/if}
      {/each}
    </Accordion>
  {/if}
</div>
