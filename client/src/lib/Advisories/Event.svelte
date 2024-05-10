<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Label, P, TimelineItem } from "flowbite-svelte";

  export let event: any;

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
</script>

<TimelineItem classLi="mb-4 ms-4" date={`${new Date(event.time).toISOString()}`}>
  <P class="mb-2">{getEventDescription(event)}</P>
  <Label class="text-xs text-slate-400">State: {event.state}</Label>
</TimelineItem>
