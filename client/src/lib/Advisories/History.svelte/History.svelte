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
  const getEventDescription = (e: any) => {
    let msg: string = "";
    switch (e.event_type) {
      case "import_document":
        msg = "Document imported";
        if (e.actor) {
          msg += " by " + e.actor;
        }
        break;
      case "state_change":
        if (e.actor) {
          msg = e.actor + " changed status to: " + e.state;
        } else {
          msg = "Status changed to: " + e.state;
        }
        break;
      case "add_comment":
        if (e.actor) {
          msg = e.actor + " added a comment";
        } else {
          msg = "A comment has been added";
        }
        break;
      case "add_sscv":
      case "add_ssvc":
        if (e.actor) {
          msg = e.actor + " added a SŚVC";
        } else {
          msg = "A SSVC has been added";
        }
        break;
      case "change_comment":
        if (e.actor) {
          msg = e.actor + " edited a comment";
        } else {
          msg = "A comment has been edited";
        }
        break;
      default:
        msg = e.event_type;
    }
    return msg;
  };
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
                  <small class="mb-1 flex-grow text-xs text-slate-400"
                    >{`${new Date(event.time).toISOString()}`}</small
                  >
                  <div class="border-1 border p-1 text-xs text-gray-800">
                    {event.state}
                  </div>
                </div>
                <span class="mb-2">{getEventDescription(event)}</span>
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
