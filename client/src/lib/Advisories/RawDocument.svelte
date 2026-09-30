<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import hljs from "highlight.js";
  import json from "highlight.js/lib/languages/json";
  import "highlight.js/styles/stackoverflow-light.css";
  import DOMPurify from "dompurify";

  hljs.registerLanguage("json", json);

  let rawDoc = $derived(
    appStore.state.webview.rawDoc ? JSON.stringify(appStore.state.webview.rawDoc, null, 2) : ""
  );

  let docWithHighlighting = $derived.by(() => {
    const html = hljs.highlight(rawDoc, { language: "json" }).value;
    return DOMPurify.sanitize(html, { USE_PROFILES: { html: true } });
  });
</script>

<svelte:head>
  {#if appStore.state.app.isDarkMode}
    <link
      rel="stylesheet"
      href="../../../node_modules/highlight.js/styles/stackoverflow-dark.min.css"
    />
  {:else}
    <link
      rel="stylesheet"
      href="../../../node_modules/highlight.js/styles/stackoverflow-light.min.css"
    />
  {/if}
</svelte:head>

<div>
  <pre>
    {@html docWithHighlighting}
  </pre>
</div>
