<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Accordion, AccordionItem, Label } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import DiffEntry from "./DiffEntry.svelte";
  import type { JsonDiffResult, JsonDiffResultList } from "./JsonDiff";
  import LazyDiffEntry from "./LazyDiffEntry.svelte";
  import { onMount } from "svelte";

  export let diffDocuments: any;
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

  onMount(() => {
    getDiff();
  });

  const getDiff = async () => {
    urlPath = `/api/diff/${diffDocuments.docB.id}/${diffDocuments.docA.id}?word-diff=true`;
    const response = await fetch(urlPath, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      }
    });
    if (response.ok) {
      const result: JsonDiffResultList = await response.json();
      diff = result;
    } else {
      appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
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
  {#if diff}
    <Label class="text-lg"
      >Compare Version {diffDocuments.docB.version} with Version {diffDocuments.docA.version}
    </Label>
    <span>{diff.length} changes</span>
    <Accordion flush multiple>
      {#each groupedResults as result}
        {#if result.changes.length > 0}
          <AccordionItem open>
            <div slot="header" class="pl-2">
              {#if result.op === "add"}
                <div class="flex items-center gap-2">
                  <i class="bx bx-plus text-lg"></i>
                  <span>Added ({result.changes.length})</span>
                </div>
              {:else if result.op === "remove"}
                <div class="flex items-center gap-2">
                  <i class="bx bx-x text-lg"></i>
                  <span>Removed ({result.changes.length})</span>
                </div>
              {:else}
                <div class="flex items-center gap-2">
                  <i class="bx bx-pencil text-lg"></i>
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
