<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import {
    Button,
    ButtonGroup,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow
  } from "flowbite-svelte";
  import Comment from "$lib/Advisories/Events/Comments/Comment.svelte";
  import { getReadableDateString } from "../CSAFWebview/helpers";
  import type { Snippet } from "svelte";
  import SSVCEntry from "$lib/Advisories/Events/SSVCEntry/SSVCEntry.svelte";

  const intlFormat = new Intl.DateTimeFormat(undefined, {
    dateStyle: "medium",
    timeStyle: "medium"
  });

  interface Props {
    entries: any;
    workflowState: string;
    onCommentUpdated: () => void;
    additionalButtons?: Snippet;
  }

  let { entries, workflowState = "", onCommentUpdated, additionalButtons }: Props = $props();

  const uid = $props.id();

  let fullHistory = $state(false);
  const tdClass = "py-2 px-2";

  let historyEntries = $derived(
    fullHistory
      ? entries
      : entries.filter((e: any) => {
          if (e.event_type === "add_comment") return e;
        })
  );
</script>

<div class="flex max-h-[470px] flex-col overflow-auto p-1">
  <Table>
    <TableBody>
      {#each historyEntries as event, i (`history-${uid}-${i}`)}
        <TableBodyRow color="default">
          {#if event.event_type !== "add_comment" && event.event_type !== "add_ssvc" && event.event_type !== "add_sscv" && event.event_type !== "change_ssvc" && event.event_type !== "change_sscv"}
            <TableBodyCell class={tdClass}>
              <div class="flex flex-col">
                <div class="flex flex-row items-baseline">
                  <small
                    class="mb-1 w-40 text-xs text-slate-400 dark:text-slate-600"
                    title={event.time}>{`${getReadableDateString(event.time, intlFormat)}`}</small
                  >
                  <small class="ml-1 flex-grow">
                    {#if event.event_type === "state_change"}
                      Statechange ({event.actor})
                    {/if}
                    {#if event.event_type === "import_document"}
                      Import ({event.actor})
                    {/if}
                    {#if event.event_type === "change_comment"}
                      Edit comment ({event.actor})
                    {/if}
                  </small>
                  {#if /state_change|import_document/.test(event.event_type)}
                    <div class="border border-1 p-1 text-xs text-gray-800 dark:text-gray-200">
                      {event.state}
                    </div>
                  {/if}
                </div>
              </div>
            </TableBodyCell>
          {/if}
          {#if event.event_type === "add_ssvc" || event.event_type === "add_sscv" || event.event_type === "change_ssvc" || event.event_type === "change_sscv"}
            <SSVCEntry ssvcData={event} />
          {/if}
          {#if event.event_type === "add_comment"}
            <Comment
              {workflowState}
              onCommentUpdated={() => {
                onCommentUpdated();
              }}
              comment={event}
              {fullHistory}
            ></Comment>
          {/if}
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
  {#if historyEntries.length === 0}
    <span class="mr-auto ml-auto text-gray-400 dark:text-gray-400"
      >{fullHistory ? "No entries" : "No comments"}</span
    >
  {/if}
  <div class="mt-6 flex flex-row justify-end gap-2">
    {#if additionalButtons}
      {@render additionalButtons()}
    {/if}
    <ButtonGroup class="h-7">
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${!fullHistory ? "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700" : ""}`}
        onclick={() => {
          fullHistory = false;
        }}>Comments</Button
      >
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${fullHistory ? "bg-gray-200 hover:bg-gray-100 dark:bg-gray-600 dark:hover:bg-gray-700" : ""}`}
        onclick={() => {
          fullHistory = true;
        }}>History</Button
      >
    </ButtonGroup>
  </div>
</div>
