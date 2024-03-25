<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  /* eslint-disable svelte/no-at-html-tags */
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import type { Note } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
  import { marked } from "marked";
  import DOMPurify from "dompurify";
  marked.use({ gfm: true });
  export let note: Note;
  let keys: string[] = [];
  let values: string[] = [];
  if (note.audience) {
    keys.push("Audience");
    values.push(note.audience);
  }
  if (note.title) {
    keys.push("Title");
    values.push(note.title);
  }
</script>

<KeyValue {keys} {values} />
<div class="ml-6">
  <h5>Text</h5>
</div>

<div class="markdown-text ml-6">
  <div class="display-markdown">
    {@html DOMPurify.sanitize(
      marked.parse(note.text.replace(/^[\u200B\u200C\u200D\u200E\u200F\uFEFF]/, ""))
    )}
  </div>
</div>

<style>
  .markdown-text {
    padding: 0.5rem;
    border: 1px solid lightgray;
    width: 100%;
    overflow-x: auto;
    position: relative;
  }
  .display-markdown {
  }
</style>
