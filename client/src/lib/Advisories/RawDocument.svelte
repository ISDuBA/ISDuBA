<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import CopyButton from "$lib/Components/CopyButton.svelte";
  import { appStore } from "$lib/store.svelte";
  import { ArrowInDownSquareHalf, FileDetail, X } from "@boxicons/svelte";
  import { Button, Modal } from "flowbite-svelte";

  let isDialogOpen = $state(false);
  let downloadAbortController: AbortController | undefined = $state(undefined);
  let innerWidth = $state(0);

  let json = $derived(
    appStore.state.webview.rawDoc ? JSON.stringify(appStore.state.webview.rawDoc, null, 2) : ""
  );

  const toggleDialog = () => {
    isDialogOpen = !isDialogOpen;
  };

  const downloadRawDocument = () => {
    if (downloadAbortController) {
      downloadAbortController.abort();
    }
    downloadAbortController = new AbortController();
    const file = new Blob([json], { type: "application/json" });
    let a = document.createElement("a"),
      url = URL.createObjectURL(file);
    a.href = url;
    a.download = `${appStore.state.webview.rawDoc.document.tracking.id}.json`;
    document.body.appendChild(a);
    a.click();
    setTimeout(() => {
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    }, 0);
  };
</script>

<svelte:window bind:innerWidth />

<Modal bind:open={isDialogOpen} size={innerWidth > 1300 ? "xl" : "none"} dismissable={false}>
  {#snippet header()}
    <div class="flex w-full justify-between">
      <div class="flex items-center gap-3">
        <h2 class="text-lg">Raw document</h2>
        <CopyButton errorMessage="Could not copy the document" title="Copy document" value={json} />
      </div>
      <Button onclick={toggleDialog} color="light" size="sm" title="Close dialog">
        <X />
      </Button>
    </div>
  {/snippet}
  <div>
    <pre>{json}</pre>
  </div>
</Modal>

<div class="flex">
  <Button
    onclick={toggleDialog}
    class="rounded-e-none"
    color="light"
    size="xs"
    title="View raw document"
  >
    <FileDetail />
  </Button>
  <Button
    onclick={downloadRawDocument}
    class="rounded-s-none"
    color="light"
    size="xs"
    title="Download document"
  >
    <ArrowInDownSquareHalf />
  </Button>
</div>
