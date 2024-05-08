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
</script>

<TimelineItem classLi="mb-4 ms-4" date={`${new Date(event.time).toISOString()}`}>
  <P class="mb-2">
    {event.event_type}
    <small>({event.actor})</small>
  </P>
  <Label class="text-xs text-slate-400">State: {event.state}</Label>
</TimelineItem>
