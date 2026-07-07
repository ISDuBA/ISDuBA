<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Label } from "flowbite-svelte";

  interface Props {
    accept?: string;
    browseButtonColor?: "primary" | "light";
    containerClass?: string | undefined;
    clearable?: boolean;
    disabled?: boolean;
    files?: FileList | undefined;
    id: string;
    multiple?: boolean;
    oldFile?: string | null | undefined;
    isFileReset?: boolean;
    titleClearButton?: string;
    onChanged?: (event?: any) => void;
  }

  let {
    accept = "",
    browseButtonColor = "primary",
    containerClass = undefined,
    clearable = true,
    disabled = false,
    files = $bindable(),
    id,
    multiple = true,
    oldFile = $bindable(undefined),
    isFileReset = $bindable(false),
    titleClearButton = "",
    onChanged = () => {}
  }: Props = $props();

  const onChange = (event: any) => {
    files = event.target.files;
    onChanged();
  };
</script>

<div class={`${containerClass ?? "mb-3 inline-flex w-full"}`}>
  <Button
    onclick={() => {
      document.getElementById(id)?.click();
    }}
    class="cursor-pointer rounded-none rounded-l-lg border-r-0 disabled:cursor-not-allowed"
    color={browseButtonColor}
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
  <input {multiple} onchange={onChange} {accept} {disabled} {id} type="file" />
  {#if clearable}
    <Button
      onclick={() => {
        files = undefined;
        oldFile = undefined;
        isFileReset = true;
        onChanged();
      }}
      {disabled}
      title={titleClearButton}
      class="w-fit cursor-pointer rounded-none rounded-r-lg border-l-0 p-1 disabled:cursor-not-allowed dark:border-gray-500 dark:bg-gray-600"
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
