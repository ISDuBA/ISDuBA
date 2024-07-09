<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Checkbox, Fileupload, Input, Label } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "svelte-spa-router";
  import { onMount, setContext } from "svelte";
  import BadgeInput from "./BadgeInput.svelte";
  import CustomInput from "./CustomInput.svelte";
  export let params: any = null;
  let loadError = "";
  let saveError = "";
  let files: any;
  let invalidFields: string[] = [];
  $: disableSave = job.domains.length === 0 || invalidFields.length > 0;

  const validation = {
    addInvalidField: (name: string) => (invalidFields = invalidFields.concat(name)),
    removeInvalidField: (name: string) => {
      const index = invalidFields.indexOf(name);
      if (index !== -1) invalidFields = invalidFields.toSpliced(index, 1);
    },
    getInvalidFields: () => invalidFields
  };
  setContext("validation", validation);

  type Job = {
    name: string;
    domains: string[];
    insecure: boolean;
    ignore_signature_check: boolean;
    worker: number;
    client_certificate: File | undefined;
    client_key: string | undefined;
    client_passphrase: string | undefined;
    rate: number | undefined;
    start_range: Date | undefined;
    end_range: Date | undefined;
    ignore_pattern: string | undefined;
    temporary: boolean;
    http_headers: string[][];
  };
  let job: Job = {
    name: "",
    domains: [],
    insecure: false,
    ignore_signature_check: false,
    worker: 1,
    client_certificate: undefined,
    client_key: undefined,
    client_passphrase: undefined,
    rate: undefined,
    start_range: undefined,
    end_range: undefined,
    ignore_pattern: undefined,
    temporary: false,
    http_headers: [
      ["", ""],
      ["", ""]
    ]
  };

  onMount(async () => {
    let id;
    if (params) id = params.id;
    if (id) {
      const response = await request(`/api/job/`, "GET");
      if (response.ok) {
        const result = await response.content;
        const thisJob = result.find((j: any) => {
          return j.id == id;
        });
        if (params && params.id) {
          const newJob: any = structuredClone(job);
          for (const property in thisJob) {
            newJob[property] = thisJob[property];
          }
          job = newJob;
        }
      } else if (response.error) {
        loadError = `Could not load query. ${getErrorMessage(response.error)}`;
      }
    }
  });

  async function saveJob() {
    const formData = new FormData();
    if (params?.id) formData.append("job_id", params.id);
    job.domains.forEach((domain) => {
      formData.append("domains", domain);
    });
    if (job.client_key) formData.append("client_key", job.client_key);
    if (job.client_passphrase) formData.append("client_passphrase", job.client_passphrase);
    if (job.start_range) formData.append("start_range", job.start_range.toString());
    if (job.end_range) formData.append("end_range", job.end_range.toString());
    if (job.rate) formData.append("rate", `${job.rate}`);
    if (job.ignore_pattern) formData.append("ignore_pattern", job.ignore_pattern);
    formData.append("name", job.name);
    formData.append("worker", job.worker.toString());
    formData.append("insecure", job.insecure.toString());
    formData.append("ignore_signature_check", job.ignore_signature_check.toString());
    formData.append("temporary", job.temporary.toString());
    const method = params?.id ? "PUT" : "POST";
    const response = await request(`/api/job`, method, formData);
    if (response.ok) {
      push(`/sources/`);
    } else if (response.error) {
      saveError = getErrorMessage(response.error);
    }
  }

  function onChangedHeaders() {
    const lastIndex = job.http_headers.length - 1;
    if (
      (job.http_headers[lastIndex][0].length > 0 && job.http_headers[lastIndex][1].length > 0) ||
      (lastIndex - 1 >= 0 &&
        job.http_headers[lastIndex - 1][0].length > 0 &&
        job.http_headers[lastIndex - 1][1].length > 0)
    ) {
      job.http_headers.push(["", ""]);
      job.http_headers = job.http_headers;
    }
  }

  function removeHeader(index: number) {
    if (job.http_headers.length === 1)
      job.http_headers = [
        ["", ""],
        ["", ""]
      ];
    job.http_headers = job.http_headers.toSpliced(index, 1);
  }
</script>

<svelte:head>
  <title>Sources - {params?.id ? "Edit Job" : "New Job"}</title>
</svelte:head>

<SectionHeader title={params?.id ? "Edit Job" : "New Job"}></SectionHeader>
<form id="job-editor" on:submit={saveJob} class="max-w-[800pt]">
  <div class="mb-4 flex flex-col gap-3">
    <CustomInput
      label="Job name"
      id="name"
      minlength={3}
      placeholder="Job #1"
      required
      bind:value={job.name}
    ></CustomInput>
    <div>
      <BadgeInput
        on:edited={(event) => {
          job.domains = event.detail;
        }}
        on:submit={saveJob}
        minEntries={1}
        id="domains"
        initialEntries={job.domains}
        label="Domains"
        placeholder="example.com"
        required
      ></BadgeInput>
    </div>
    <Label class="mb-1 space-y-2">
      <span>TLS client certificate file</span>
      <Fileupload bind:files />
    </Label>
    <div class="grid gap-x-2 gap-y-2 md:grid-cols-2">
      <CustomInput id="client-key" label="TLS client private key file" bind:value={job.client_key}
      ></CustomInput>
      <CustomInput
        id="client-passphrase"
        type="password"
        label="Client passphrase"
        bind:value={job.client_passphrase}
      ></CustomInput>
      <CustomInput id="start-range" type="date" label="Start range" bind:value={job.start_range}
      ></CustomInput>
      <CustomInput id="end-range" type="date" label="End range" bind:value={job.end_range}
      ></CustomInput>
    </div>
    <CustomInput id="worker" label="Worker" type="number" step={1} bind:value={job.worker} min={1}
    ></CustomInput>
    <CustomInput id="rate" type="number" label="Rate" placeholder="5" bind:value={job.rate}
    ></CustomInput>
    <CustomInput id="ignore-pattern" label="Ignore pattern" bind:value={job.ignore_pattern}
    ></CustomInput>
    <Label>HTTP headers</Label>
    <div class="mb-3 grid items-end gap-x-2 gap-y-4 md:grid-cols-3">
      {#if job.http_headers}
        {#each job.http_headers as header, index (index)}
          <Label>
            <span class="text-gray-500">Key</span>
            <Input on:change={onChangedHeaders} bind:value={header[0]} />
          </Label>
          <Label>
            <span class="text-gray-500">Value</span>
            <Input on:change={onChangedHeaders} bind:value={header[1]} />
          </Label>
          {#if !(job.http_headers.length === 1 && job.http_headers[0][0] === "" && job.http_headers[0][1] === "")}
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
      {/if}
    </div>
    <Checkbox
      on:change={() => {
        job.insecure = !job.insecure;
      }}
      bind:checked={job.insecure}
    >
      <span>Insecure</span>
    </Checkbox>
    <Checkbox
      on:change={() => {
        job.ignore_signature_check = !job.ignore_signature_check;
      }}
      bind:checked={job.ignore_signature_check}
    >
      <span>Ignore signature check</span>
    </Checkbox>
    <Checkbox
      on:change={() => {
        job.temporary = !job.temporary;
      }}
      bind:checked={job.temporary}
    >
      <span>Temporary</span>
    </Checkbox>
  </div>
  <div class="mb-4 text-sm text-red-500">* Required fields</div>
  <Button disabled={disableSave} type="submit" color="light">
    <i class="bx bxs-save me-2"></i>
    <span>Save</span>
  </Button>
</form>
<ErrorMessage message={saveError}></ErrorMessage>
<ErrorMessage message={loadError}></ErrorMessage>
