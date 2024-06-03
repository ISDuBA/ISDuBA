<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { Button, ButtonGroup, Label, P, Timeline, TimelineItem } from "flowbite-svelte";
  import Comment from "$lib/Advisories/Comments/Comment.svelte";
  import { createEventDispatcher } from "svelte";
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

<ButtonGroup class="mb-6 ml-auto h-7">
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
  <Timeline class="mb-4 ml-1 flex flex-col">
    {#each historyEntries as event}
      {#if event.event_type !== "add_comment"}
        <TimelineItem classLi="mb-4 ms-4" date={`${new Date(event.time).toISOString()}`}>
          <P class="mb-2">{getEventDescription(event)}</P>
          <Label class="text-xs text-slate-400">State: {event.state}</Label>
        </TimelineItem>
      {/if}
      {#if event.event_type === "add_comment"}
        <Comment
          on:commentUpdate={() => {
            dispatch("commentUpdate");
          }}
          comment={event}
        ></Comment>
      {/if}
    {/each}
  </Timeline>
</div>
