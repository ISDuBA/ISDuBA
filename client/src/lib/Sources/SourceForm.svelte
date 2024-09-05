<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { Source } from "$lib/Sources/source";
  import {
    Accordion,
    AccordionItem,
    Button,
    Checkbox,
    Fileupload,
    Input,
    Label,
    NumberInput
  } from "flowbite-svelte";
  export let formClass: string = "";
  export let source: Source;
  export let enableActive: boolean = false;
  export const updateSource = async () => {
    formatHeaders();
    await loadCerts();
  };

  export let inputChange = () => {};

  let headers: [string, string][] = [["", ""]];
  let privateCert: FileList | undefined;
  let publicCert: FileList | undefined;

  const onChangedHeaders = (e: Event | undefined) => {
    const lastIndex = headers.length - 1;
    if (
      (headers[lastIndex][0].length > 0 && headers[lastIndex][1].length > 0) ||
      (lastIndex - 1 >= 0 &&
        headers[lastIndex - 1][0].length > 0 &&
        headers[lastIndex - 1][1].length > 0)
    ) {
      headers.push(["", ""]);
      headers = headers;
    }
    if (e) {
      inputChange();
    }
  };

  const onChangedIgnorePatterns = () => {
    if (source.ignore_patterns.at(-1) !== "") {
      source.ignore_patterns.push("");
    }
    inputChange();
  };

  const removeHeader = (index: number) => {
    if (headers.length === 1)
      headers = [
        ["", ""],
        ["", ""]
      ];
    headers = headers.toSpliced(index, 1);
    inputChange();
  };

  const removePattern = (index: number) => {
    if (source.ignore_patterns.length === 1) source.ignore_patterns = [""];
    source.ignore_patterns = source.ignore_patterns.toSpliced(index, 1);
    inputChange();
  };

  $: if (source.headers) {
    parseHeaders();
  }

  const parseHeaders = () => {
    headers = [];
    for (const header of source.headers) {
      let h = header.split(":");
      headers.push([h[0], h[1]]);
    }
    if (headers.length === 0) {
      headers.push(["", ""]);
    }
    onChangedHeaders(undefined);
  };

  const formatHeaders = () => {
    source.headers = [];
    for (const header of headers) {
      if (header[0] !== "" && header[1] !== "") source.headers.push(`${header[0]}:${header[1]}`);
    }
  };

  const loadCerts = async () => {
    if (privateCert) {
      source.client_cert_private = await privateCert.item(0)?.text();
    }
    if (publicCert) {
      source.client_cert_public = await publicCert.item(0)?.text();
    }
  };
</script>

<form class={formClass}>
  <Label>Name</Label>
  <Input class="mb-3" on:input={inputChange} bind:value={source.name}></Input>
  {#if enableActive}
    <Checkbox class="mb-3" on:change={inputChange} bind:checked={source.active}>Active</Checkbox>
  {/if}
  <Accordion>
    <AccordionItem
      ><span slot="header">Credentials</span>
      <Label>Private cert</Label>
      <div class="mb-3 inline-flex w-full">
        <Fileupload
          class="rounded-none rounded-l-lg"
          on:change={inputChange}
          bind:files={privateCert}
        ></Fileupload>
        <Button
          on:click={() => {
            source.client_cert_private = null;
            privateCert = undefined;
            inputChange();
          }}
          title="Remove private cert"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
      <Label>Public cert</Label>
      <div class="mb-3 inline-flex w-full">
        <Fileupload
          class="rounded-none rounded-l-lg"
          on:change={inputChange}
          bind:files={publicCert}
        ></Fileupload>
        <Button
          on:click={() => {
            source.client_cert_public = null;
            publicCert = undefined;
            inputChange();
          }}
          title="Remove public cert"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
      <Label>Client cert passphrase</Label>
      <div class="mb-3 inline-flex w-full">
        <Input
          class="rounded-none rounded-l-lg"
          on:input={inputChange}
          bind:value={source.client_cert_passphrase}
        />
        <Button
          on:click={() => {
            source.client_cert_passphrase = null;
          }}
          title="Remove passphrase"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
    </AccordionItem>
    <AccordionItem
      ><span slot="header">Advanced options</span>
      <div class="mb-3 grid w-full gap-x-2 gap-y-4 md:grid-cols-3">
        <div>
          <Label>Age</Label>
          <Input placeholder="17520h" on:input={inputChange} bind:value={source.age}></Input>
        </div>
        <div>
          <Label>Rate</Label>
          <NumberInput placeholder={1} on:input={inputChange} bind:value={source.rate}></NumberInput>
        </div>
        <div>
          <Label>Slots</Label>
          <NumberInput placeholder={2} on:input={inputChange} bind:value={source.slots}></NumberInput>
        </div>
      </div>

      <Label>Options</Label>
      <div class="mb-3 flex w-full gap-4">
        <Checkbox on:change={inputChange} bind:checked={source.strict_mode}>Strict mode</Checkbox>
        <Checkbox on:change={inputChange} bind:checked={source.insecure}>Insecure</Checkbox>
        <Checkbox on:change={inputChange} bind:checked={source.signature_check}
          >Signature check</Checkbox
        >
      </div>

      <Label>Ignore patterns</Label>
      {#each source.ignore_patterns as pattern, index (index)}
        <div class="mb-3 inline-flex w-full">
          <Label class="grow">
            <Input
              class="rounded-none rounded-l-lg"
              on:input={onChangedIgnorePatterns}
              bind:value={pattern}
            />
          </Label>
          <Button
            on:click={() => removePattern(index)}
            title="Remove pattern"
            class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
            color="light"
            disabled={source.ignore_patterns.length <= 1}
          >
            <i class="bx bx-x"></i>
          </Button>
        </div>
      {/each}

      <Label>HTTP headers</Label>
      {#each headers as header, index (index)}
        <div class="mb-3 grid w-full grid-cols-[1fr_auto] sm:grid-cols-[1fr_1fr_auto]">
          <Label class="col-span-2 col-start-1 row-start-1 sm:col-span-1">
            <span class="text-gray-500">Key</span>
            <span class="text-gray-500 sm:collapse">- Value</span>
          </Label>
          <Input
            class="col-span-2 row-start-2 rounded-none rounded-t-lg sm:col-span-1 sm:rounded-l-lg sm:rounded-tr-none"
            on:input={onChangedHeaders}
            bind:value={header[0]}
          />
          <Label class="collapse col-span-2 col-start-1 row-start-1 sm:visible sm:col-start-2">
            <span class="text-gray-500">Value</span>
          </Label>
          <Input
            class="row-start-3 rounded-none rounded-bl-lg border-t-0 sm:row-start-2 sm:rounded-bl-none sm:border-l-0 sm:border-t"
            on:input={onChangedHeaders}
            bind:value={header[1]}
          />
          {#if headers.length > 1}
            <Button
              on:click={() => removeHeader(index)}
              title="Remove key-value-pair"
              class="row-start-3 h-full w-fit rounded-none rounded-br-lg border-l-0 border-t-0 p-1 sm:row-start-2 sm:rounded-tr-lg sm:border-t"
              color="light"
            >
              <i class="bx bx-x"></i>
            </Button>
          {:else}
            <Button
              title="Remove key-value-pair"
              class=" row-start-3 h-full w-fit rounded-none rounded-br-lg border-l-0 border-t-0 p-1 sm:row-start-2 sm:rounded-tr-lg sm:border-t"
              color="light"
              disabled={true}
            >
              <i class="bx bx-x"></i>
            </Button>
          {/if}
        </div>
      {/each}
    </AccordionItem>
  </Accordion>
  <br />
</form>
