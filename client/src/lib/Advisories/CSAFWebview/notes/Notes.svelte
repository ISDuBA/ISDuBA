<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SingleNote from "$lib/Advisories/CSAFWebview/notes/Note.svelte";
  import type { Note } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
  import { onMount } from "svelte";
  import Collapsible from "../Collapsible.svelte";

  interface Props {
    notes: Note[];
    initOpen?: boolean;
  }
  let { notes, initOpen = false }: Props = $props();

  const uid = $props.id();

  let hasDescription = $derived(notes.some((note) => note.category === "description"));

  let openNote: boolean[] = $state([]);

  onMount(() => {
    openNote = notes.map((note) => {
      return initOpen || note.category === (hasDescription ? "description" : "summary");
    });
  });
</script>

{#if notes}
  {#each notes as note, index (`notes-${uid}-${index}`)}
    <Collapsible
      header={note.title ? `${note.category}: ${note.title}` : note.category}
      level={4}
      open={openNote[index]}
      path={`/document/notes/[${index}]`}
    >
      <SingleNote {note} path={`/document/notes/[${index}]`} />
    </Collapsible>
  {/each}
{/if}
