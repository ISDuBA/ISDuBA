<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    Button,
    Card,
    Fileupload,
    Label,
    Listgroup,
    ListgroupItem,
    Tooltip
  } from "flowbite-svelte";
  import { type UploadInfo } from "$lib/Sources/source";
  export let label;
  export let upload = async (files: FileList): Promise<UploadInfo[]> => {
    // eslint-disable-next-line no-console
    console.log(files, uploadInfo);
    return [];
  };

  export let uploadInfo: UploadInfo[] = [];

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
  <div class={`flex flex-col ${files?.length > 1 ? "mb-4" : "mb-40"}`}>
    <Label class="pb-2">{label}</Label>
    <Fileupload value="" bind:files multiple />
    <Listgroup class="mt-6">
      {#if !files}
        <ListgroupItem>No files selected</ListgroupItem>
      {:else}
        {#each files as file, i}
          {@const info = uploadInfo[i]}
          {@const color = getColor(info)}
          <ListgroupItem class={color}>{file.name}</ListgroupItem>

          {#if info?.message}
            <Tooltip>{info.message}</Tooltip>
          {/if}
        {/each}
      {/if}
    </Listgroup>
  </div>
  <Button
    on:click={async () => {
      uploadInfo = await upload(files);
    }}
    class="mt-auto ml-auto"
    color="primary">Upload</Button
  >
</Card>
