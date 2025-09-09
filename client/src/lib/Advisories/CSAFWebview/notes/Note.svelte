<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import type { Note } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";

  interface Props {
    note: Note;
  }
  let { note }: Props = $props();

  let keys: string[] = [];
  let values: string[] = [];
  if (note.audience) {
    keys.push("Audience");
    values.push(note.audience);
  }
</script>

<KeyValue {keys} {values} />
<div class="ml-7">
  <h5>Text</h5>
</div>

<div class="markdown-text">
  <div class="display-markdown max-w-2/3">
    {note.text}
  </div>
</div>

<style>
  .markdown-text {
    margin-left: 1.75rem;
    padding: 0.5rem;
    border: 1px solid lightgray;
    min-width: 200px;
    overflow-x: auto;
    position: relative;
  }
</style>
