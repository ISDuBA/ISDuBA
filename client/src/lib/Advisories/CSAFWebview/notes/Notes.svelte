<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Collapsible from "$lib/Advisories/CSAFWebview/Collapsible.svelte";
  import SingleNote from "$lib/Advisories/CSAFWebview/notes/Note.svelte";
  import type { Note } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
  export let notes: Note[];
  export let open: boolean = false;
  $: hasDescription = notes.some((note) => note.category === "description");
</script>

{#if notes}
  {#each notes as note}
    <Collapsible
      header={note.title ? `${note.category}: ${note.title}` : note.category}
      level={4}
      open={open || note.category === (hasDescription ? "description" : "summary")}
    >
      <SingleNote {note} />
    </Collapsible>
  {/each}
{/if}
