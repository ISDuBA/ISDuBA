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
  import { tick } from "svelte";

  interface Props {
    cancel: () => any;
    label: any;
    upload?: any;
    uploadInfo?: UploadInfo[];
  }

  let {
    cancel,
    label,
    upload = async (
      files: FileList,
      updateCallback: ((uploadInfo: UploadInfo[]) => void) | undefined
    ): Promise<UploadInfo[]> => {
      // eslint-disable-next-line no-console
      console.log(files, uploadInfo, updateCallback);
      return [];
    },
    uploadInfo = $bindable([])
  }: Props = $props();

  const getColor = (uploadInfo: UploadInfo) => {
    let success = uploadInfo?.success;
    if (success !== undefined) {
      return success ? "text-green-600" : "text-red-600";
    }
    return "";
  };
  let files: FileList | undefined = $state(undefined);
  let filesCache: FileList | undefined = $state(undefined);
  let isUploading = $state(false);
  let showFileInput = $state(true);
  $effect(() => {
    if (files) {
      uploadInfo = [];
    }
  });
</script>

<Card size="lg" class="p-4">
  <div class={`flex flex-col gap-4 ${files?.length && files.length > 1 ? "mb-4" : "mb-40"}`}>
    <div>
      <Label class="pb-2">{label}</Label>
      {#if showFileInput}
        <Fileupload
          wrapperClass="cursor-pointer disabled:cursor-not-allowed !p-0 dark:text-gray-400"
          class="file:bg-primary-800"
          value=""
          bind:files
          multiple
          accept=".json"
          onchange={() => {
            filesCache = undefined;
          }}
        />
      {/if}
    </div>
    <div class="flex items-center justify-end gap-2">
      {#if isUploading}
        <div class="flex w-fit gap-2">
          <span>Uploading ...</span>
          <div class="w-fit min-w-8">
            <span class="min-w-16">{uploadInfo?.length}</span>/<span class="min-w-16"
              >{files?.length ?? 1}</span
            >
          </div>
        </div>
      {/if}
      {#if isUploading}
        <Button
          onclick={() => {
            cancel();
          }}
          color="red">Cancel</Button
        >
      {/if}
      <Button
        onclick={async () => {
          isUploading = true;
          setTimeout(async () => {
            if (files) {
              filesCache = files;
              uploadInfo = await upload(files, (info: UploadInfo[]) => {
                uploadInfo = info;
              });
            }
            files = undefined;
            // This is a hack. The file input has to be re-added to the DOM. Otherwise it would not change
            // the label "x files selected." to its initial value even if we set files to undefined.
            showFileInput = false;
            await tick();
            showFileInput = true;
            isUploading = false;
          });
        }}
        color="primary"
        disabled={isUploading || !files || files.length === 0}>Upload</Button
      >
    </div>
    {#if filesCache}
      <Listgroup class="mt-6">
        {#each filesCache as file, i}
          {@const info = uploadInfo[i]}
          {@const color = getColor(info)}
          <ListgroupItem>
            <div class="flex items-center gap-1">
              {#if info?.success}
                <i class={`bx bx-check-circle ${color}`}></i>
              {:else if info}
                <i class={`bx bx-x-circle ${color}`}></i>
              {/if}
              <div class={`font-bold text-black dark:text-white`}>{file.name}</div>
            </div>
            {#if info?.message}
              <div>{info.message}</div>
            {/if}
          </ListgroupItem>
        {/each}
      </Listgroup>
    {:else if files}
      <Listgroup class="mt-6">
        {#each files as file}
          <ListgroupItem>
            <div class="flex items-center gap-1">
              <div class={`font-bold text-black dark:text-white`}>{file.name}</div>
            </div>
          </ListgroupItem>
        {/each}
      </Listgroup>
    {/if}
  </div>
</Card>
