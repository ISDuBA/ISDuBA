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
  import type { JsonDiffResult } from "./JsonDiff";

  export let diff: any;
  $: groupedResults = [
    {
      op: "add",
      changes: diff.result.filter((result: JsonDiffResult) => result.op === "add")
    },
    {
      op: "replace",
      changes: diff.result.filter((result: JsonDiffResult) => result.op === "replace")
    },
    {
      op: "remove",
      changes: diff.result.filter((result: JsonDiffResult) => result.op === "remove")
    }
  ];

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
  <Label class="text-lg">Compare Version {diff.docA} with Version {diff.docB}</Label>
  <span>{diff.result.length} changes</span>
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
              <div class="mb-1">
                <b>
                  <code>
                    {change.path}
                  </code>
                </b>
              </div>
              {#if change.value}
                <DiffEntry content={change.value} operation={change.op}></DiffEntry>
              {/if}
            </div>
          {/each}
        </AccordionItem>
      {/if}
    {/each}
  </Accordion>
</div>
