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
    Label
  } from "flowbite-svelte";
  export let formClass: string = "";
  export let source: Source;
  export let enableActive: boolean = false;
  export const updateSource = async () => {
    formatHeaders();
    await loadCerts();
  };

  let headers: [string, string][] = [["", ""]];
  let privateCert: FileList | undefined;
  let publicCert: FileList | undefined;

  const onChangedHeaders = () => {
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
  };

  const onChangedIgnorePatterns = () => {
    if (source.ignore_patterns.at(-1) !== "") {
      source.ignore_patterns.push("");
    }
  };

  const removeHeader = (index: number) => {
    if (headers.length === 1)
      headers = [
        ["", ""],
        ["", ""]
      ];
    headers = headers.toSpliced(index, 1);
  };

  const removePattern = (index: number) => {
    if (source.ignore_patterns.length === 1) source.ignore_patterns = [""];
    source.ignore_patterns = source.ignore_patterns.toSpliced(index, 1);
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
    onChangedHeaders();
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
  {#if enableActive}
    <Checkbox bind:checked={source.active}>Active</Checkbox>
  {/if}
  <Accordion>
    <AccordionItem
      ><span slot="header">Credentials</span>
      <Label>Private cert</Label>
      <Fileupload bind:files={privateCert}></Fileupload>
      <Button
        on:click={() => {
          source.client_cert_private = null;
        }}
        title="Remove private cert"
        class="mb-3 w-fit p-1"
        color="light"
      >
        <i class="bx bx-x"></i>
      </Button>
      <Label>Public cert</Label>
      <Fileupload bind:files={publicCert}></Fileupload>
      <Button
        on:click={() => {
          source.client_cert_public = null;
        }}
        title="Remove public cert"
        class="mb-3 w-fit p-1"
        color="light"
      >
        <i class="bx bx-x"></i>
      </Button>
      <Label>Client cert passphrase</Label>
      <Input bind:value={source.client_cert_passphrase} />
      <Button
        on:click={() => {
          source.client_cert_passphrase = null;
        }}
        title="Remove passphrase"
        class="mb-3 w-fit p-1"
        color="light"
      >
        <i class="bx bx-x"></i>
      </Button>
    </AccordionItem>
    <AccordionItem
      ><span slot="header">Advanced options</span>
      <Label>Rate</Label>
      <Input bind:value={source.rate}></Input>
      <Label>Slots</Label>
      <Input bind:value={source.slots}></Input>

      <Label>HTTP headers</Label>
      <div class="mb-3 grid items-end gap-x-2 gap-y-4 md:grid-cols-3">
        {#each headers as header, index (index)}
          <Label>
            <span class="text-gray-500">Key</span>
            <Input on:input={onChangedHeaders} bind:value={header[0]} />
          </Label>
          <Label>
            <span class="text-gray-500">Value</span>
            <Input on:input={onChangedHeaders} bind:value={header[1]} />
          </Label>
          {#if headers.length > 1}
            <Button
              on:click={() => removeHeader(index)}
              title="Remove key-value-pair"
              class="mb-3 w-fit p-1"
              color="light"
            >
              <i class="bx bx-x"></i>
            </Button>
          {:else}
            <div></div>
          {/if}
        {/each}
      </div>
      <Checkbox bind:checked={source.strict_mode}>Strict mode</Checkbox>
      <Checkbox bind:checked={source.insecure}>Insecure</Checkbox>
      <Checkbox bind:checked={source.signature_check}>Signature check</Checkbox>
      <Label>Age</Label>
      <Input placeholder="17520h" bind:value={source.age}></Input>
      <Label>Ignore patterns</Label>
      <div class="mb-3 grid items-end gap-x-2 gap-y-4 md:grid-cols-2">
        {#each source.ignore_patterns as pattern, index (index)}
          <Label>
            <Input on:input={onChangedIgnorePatterns} bind:value={pattern} />
          </Label>
          {#if source.ignore_patterns.length > 1}
            <Button
              on:click={() => removePattern(index)}
              title="Remove pattern"
              class="mb-3 w-fit p-1"
              color="light"
            >
              <i class="bx bx-x"></i>
            </Button>
          {:else}
            <div></div>
          {/if}
        {/each}
      </div>
    </AccordionItem>
  </Accordion>
  <Label>Name</Label>
  <Input bind:value={source.name}></Input>
  <br />
</form>
