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
  import Comment from "$lib/Advisories/Comments/Comment.svelte";
  import { createEventDispatcher } from "svelte";

  const dispatch = createEventDispatcher();
  export let entries;
  export let state = "";
  let fullHistory = false;
  const tdClass = "py-2 px-2";

  $: historyEntries = fullHistory
    ? entries
    : entries.filter((e: any) => {
        if (e.event_type === "add_comment") return e;
      });
</script>

<div class="flex max-h-[34rem] flex-col overflow-auto p-1">
  <div class="mb-9 mt-3 flex flex-row">
    <ButtonGroup class="ml-auto h-7">
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${!fullHistory ? "bg-gray-200 hover:bg-gray-100" : ""}`}
        on:click={() => {
          fullHistory = false;
        }}>Comments</Button
      >
      <Button
        size="xs"
        color="light"
        class={`h-7 py-1 text-xs ${fullHistory ? "bg-gray-200 hover:bg-gray-100" : ""}`}
        on:click={() => {
          fullHistory = true;
        }}>History</Button
      >
    </ButtonGroup>
  </div>

  <Table>
    <TableBody>
      {#each historyEntries as event}
        <TableBodyRow>
          {#if event.event_type !== "add_comment"}
            <TableBodyCell {tdClass}>
              <div class="flex flex-col">
                <div class="flex flex-row items-baseline">
                  <small class="mb-1 w-32 text-xs text-slate-400" title={event.time}
                    >{`${event.time.replace("T", " ").split(".")[0]}`}</small
                  >
                  <small class="ml-1 flex-grow">
                    {#if event.event_type === "state_change"}
                      Statechange ({event.actor})
                    {/if}
                    {#if event.event_type === "add_ssvc" || event.event_type === "add_sscv"}
                      SSVC added ({event.actor})
                    {/if}
                    {#if event.event_type === "import_document"}
                      Import ({event.actor})
                    {/if}
                    {#if event.event_type === "change_comment"}
                      Edit comment ({event.actor})
                    {/if}
                    {#if event.event_type === "change_ssvc" || event.event_type === "change_sscv"}
                      SSVC changed ({event.actor})
                    {/if}
                  </small>
                  {#if /state_change|import_document/.test(event.event_type)}
                    <div class="border-1 border p-1 text-xs text-gray-800">
                      {event.state}
                    </div>
                  {/if}
                </div>
              </div>
            </TableBodyCell>
          {/if}
          {#if event.event_type === "add_comment"}
            <Comment
              {state}
              on:commentUpdate={() => {
                dispatch("commentUpdate");
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
    <span class="ml-auto mr-auto text-gray-400">{fullHistory ? "No entries" : "No comments"}</span>
  {/if}
</div>
