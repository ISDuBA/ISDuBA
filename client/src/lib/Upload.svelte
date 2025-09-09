<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Card, Fileupload, Label, Listgroup, ListgroupItem } from "flowbite-svelte";
  import { type UploadInfo } from "$lib/Sources/source";
  export let label;
  export let upload = async (files: FileList): Promise<UploadInfo[]> => {
    // eslint-disable-next-line no-console
    console.log(files, uploadInfo);
    return [];
  };

  export let uploadInfo: UploadInfo[] = [];
  let isUploading = false;

  const getColor = (uploadInfo: UploadInfo) => {
    let success = uploadInfo?.success;
    if (success !== undefined) {
      return success ? "text-green-600" : "text-red-600";
    }
    return "";
  };
  let files: FileList;
  $: if (files) {
    uploadInfo = [];
  }
</script>

<Card size="lg">
  <div class={`flex flex-col gap-4 ${files?.length > 1 ? "mb-4" : "mb-40"}`}>
    <div>
      <Label class="pb-2">{label}</Label>
      <Fileupload value="" bind:files multiple accept=".json" />
    </div>
    <Button
      on:click={async () => {
        isUploading = true;
        uploadInfo = await upload(files);
        isUploading = false;
      }}
      class="mt-auto ml-auto"
      color="primary"
      disabled={isUploading}>Upload</Button
    >
    <Listgroup class="mt-6">
      {#if !files}
        <ListgroupItem>No files selected</ListgroupItem>
      {:else}
        {#each files as file, i}
          {@const info = uploadInfo[i]}
          {@const color = getColor(info)}
          <ListgroupItem>
            <div class="flex items-center gap-1">
              {#if info?.success}
                <i class={`bx bx-check-circle ${color}`}></i>
              {:else if info}
                <i class={`bx bx-x-circle ${color}`}></i>
              {/if}
              <div class={`font-bold text-black`}>{file.name}</div>
            </div>
            {#if info?.message}
              <div>{info.message}</div>
            {/if}
          </ListgroupItem>
        {/each}
      {/if}
    </Listgroup>
  </div>
</Card>
