<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Card, Fileupload, Label, Listgroup, ListgroupItem } from "flowbite-svelte";
  export let label;
  export let upload = (files: FileList) => {
    // eslint-disable-next-line no-console
    console.log(files);
  };
  let files: FileList;
</script>

<Card size="lg">
  <div class={`flex flex-col ${files?.length > 1 ? "mb-4" : "mb-40"}`}>
    <Label class="pb-2">{label}</Label>
    <Fileupload value="" bind:files multiple />
    <Listgroup class="mt-6">
      {#if !files}
        <ListgroupItem>No files selected</ListgroupItem>
      {:else}
        {#each files as file}
          <ListgroupItem>{file.name}</ListgroupItem>
        {/each}
      {/if}
    </Listgroup>
  </div>
  <Button
    on:click={() => {
      upload(files);
    }}
    class="ml-auto mt-auto"
    color="primary">Upload</Button
  >
</Card>
