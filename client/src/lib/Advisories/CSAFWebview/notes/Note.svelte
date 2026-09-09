<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import { getCSAFVersion } from "$lib/Advisories/docmodel";
  import { appStore } from "$lib/store.svelte";
  import SearchableText from "../SearchableText.svelte";
  import type { Note as Note2_0 } from "$lib/Advisories/types/csaf-2.0";
  import type { Note as Note2_1 } from "$lib/Advisories/types/csaf-2.1";
  import ValueList from "../ValueList.svelte";

  interface Props {
    note: Note2_0 | Note2_1;
    path: string;
  }
  let { note, path }: Props = $props();

  let csafVersion = $derived(getCSAFVersion(appStore.state.webview.doc));
  let keys: string[] = $derived(note.audience ? ["Audience"] : []);
  let values: string[] = $derived(note.audience ? [note.audience] : []);
  let paths: string[] = $derived(note.audience ? [`${path}/audience`] : []);
</script>

<KeyValue {keys} {values} {paths} />
{#if csafVersion === "2.1"}
  {@const note21 = note as Note2_1}
  <ValueList label="Group IDs" values={note21.group_ids} path={`${path}/group_id`} />
  <ValueList label="Product IDs" values={note21.product_ids} path={`${path}/product_ids`} />
{/if}
<div class="ml-7">
  <h5>Text</h5>
</div>

<div class="markdown-text">
  <div class="display-markdown max-w-2/3">
    <SearchableText text={note.text} textPath={`${path}/text`} />
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
