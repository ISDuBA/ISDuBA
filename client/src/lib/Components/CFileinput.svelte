<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Label } from "flowbite-svelte";
  import { createEventDispatcher } from "svelte";

  export let containerClass: string | undefined = undefined;
  export let clearable = true;
  export let disabled = false;
  export let files: FileList | undefined = undefined;
  export let id: string;
  export let multiple = true;
  export let oldFile: string | null | undefined = undefined;
  export let isFileReset: boolean = false;
  export let titleClearButton = "";

  const dispatch = createEventDispatcher();

  const onChange = (event: any) => {
    files = event.target.files;
    dispatch("change");
  };
</script>

<div class={`${containerClass ?? "mb-3 inline-flex w-full"}`}>
  <Button
    on:click={() => {
      document.getElementById(id)?.click();
    }}
    class="rounded-none rounded-l-lg border border-r-0 dark:border-gray-700 dark:bg-gray-800"
    color="primary"
    {disabled}>Browse...</Button
  >
  <Label
    class={`flex min-h-full w-full min-w-48 items-center border border-gray-300 ps-4 dark:border-gray-500 dark:bg-gray-600 ${clearable ? "" : "rounded-r-lg"}`}
    for={id}
  >
    {#if files}
      {#if files.length > 1}
        <span>{files.length} files selected</span>
      {:else if files.length > 0}
        <span>{files.item(0)?.name}</span>
      {/if}
    {:else if oldFile && !isFileReset}
      <span>{oldFile}</span>
    {:else}
      <span>No file selected</span>
    {/if}
  </Label>
  <input {multiple} on:change={onChange} {disabled} {id} type="file" />
  {#if clearable}
    <Button
      on:click={() => {
        files = undefined;
        oldFile = undefined;
        isFileReset = true;
        dispatch("change");
      }}
      title={titleClearButton}
      class="w-fit rounded-none rounded-r-lg border-l-0 p-1 dark:border-gray-500 dark:bg-gray-600"
      color="light"
    >
      <i class="bx bx-x"></i>
    </Button>
  {/if}
</div>

<style>
  input[type="file"] {
    display: none;
  }
</style>
