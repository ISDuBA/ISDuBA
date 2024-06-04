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
  import { tdClass } from "$lib/table/defaults";
  const dispatch = createEventDispatcher();

  export let entries;
  let historyOnly = true;
  $: historyEntries = historyOnly
    ? entries
    : entries.filter((e: any) => {
        if (e.event_type === "add_comment") return e;
      });
</script>

<ButtonGroup class="mb-9 ml-auto mt-2 h-7">
  <Button
    size="xs"
    color="light"
    class={`h-7 py-1 text-xs ${historyOnly ? "bg-gray-200 hover:bg-gray-100" : ""}`}
    on:click={() => {
      historyOnly = true;
    }}>Full history</Button
  >
  <Button
    size="xs"
    color="light"
    class={`h-7 py-1 text-xs ${!historyOnly ? "bg-gray-200 hover:bg-gray-100" : ""}`}
    on:click={() => {
      historyOnly = false;
    }}>Comments only</Button
  >
</ButtonGroup>

<div class="max-h-96 overflow-auto">
  <Table>
    <TableBody>
      {#each historyEntries as event}
        <TableBodyRow>
          {#if event.event_type !== "add_comment"}
            <TableBodyCell {tdClass}>
              <div class="ml-1 flex flex-col">
                <div class="flex flex-row items-baseline">
                  <small class="mb-1 text-xs text-slate-400"
                    >{`${new Date(event.time).toISOString()}`}</small
                  >
                  <small class="ml-1 flex-grow">
                    {#if event.event_type === "state_change"}
                      Statechange ( {event.actor} )
                    {/if}
                    {#if event.event_type === "add_ssvc" || event.event_type === "add_sscv"}
                      SSVC ( {event.actor} )
                    {/if}
                    {#if event.event_type === "import_document"}
                      Import ( {event.actor} )
                    {/if}
                    {#if event.event_type === "change_comment"}
                      Edit comment ( {event.actor} )
                    {/if}
                  </small>
                  <div class="border-1 border p-1 text-xs text-gray-800">
                    {event.state}
                  </div>
                </div>
              </div>
            </TableBodyCell>
          {/if}
          {#if event.event_type === "add_comment"}
            <Comment
              on:commentUpdate={() => {
                dispatch("commentUpdate");
              }}
              comment={event}
            ></Comment>
          {/if}
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
</div>
