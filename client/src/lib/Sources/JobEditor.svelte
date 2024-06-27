<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Checkbox } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import BadgeInput from "./BadgeInput.svelte";
  import CustomInput from "./CustomInput.svelte";
  export let params: any = null;
  $: disableSave = !job.name || job.name === "" || job.domains.length === 0;
  let loadError = "";
  let saveError = "";

  type Job = {
    name: string;
    domains: string[];
    insecure: boolean;
    ignoreSignatureCheck: boolean;
    worker: number;
    clientKey: string | undefined;
    clientPassphrase: string | undefined;
    rate: number | undefined;
    startRange: Date | undefined;
    endRange: Date | undefined;
    ignorePattern: string | undefined;
    temporary: boolean;
  };
  let job: Job = {
    name: "",
    domains: [],
    insecure: false,
    ignoreSignatureCheck: false,
    worker: 1,
    clientKey: undefined,
    clientPassphrase: undefined,
    rate: undefined,
    startRange: undefined,
    endRange: undefined,
    ignorePattern: undefined,
    temporary: false
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
          job = thisJob;
        }
        job.domains = JSON.parse(job.domains[0]);
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
    if (job.clientKey) formData.append("client_key", job.clientKey);
    if (job.clientPassphrase) formData.append("client_passphrase", job.clientPassphrase);
    if (job.startRange) formData.append("start_range", job.startRange.toString());
    if (job.endRange) formData.append("end_range", job.endRange.toString());
    if (job.rate) formData.append("rate", `${job.rate}`);
    if (job.ignorePattern) formData.append("ignore_pattern", job.ignorePattern);
    formData.append("name", job.name);
    formData.append("insecure", job.insecure.toString());
    formData.append("ignore_signature_check", job.ignoreSignatureCheck.toString());
    formData.append("temporary", job.temporary.toString());
    const method = params?.id ? "PUT" : "POST";
    const response = await request(`/api/job`, method, formData);
    if (response.ok) {
      push(`/sources/`);
    } else if (response.error) {
      saveError = getErrorMessage(response.error);
    }
  }
</script>

<SectionHeader title={params?.id ? "Edit Job" : "New Job"}></SectionHeader>
<form on:submit={saveJob}>
  <div class="mb-4 flex flex-col gap-4">
    <CustomInput label="Job name" id="name" placeholder="Job #1" required bind:value={job.name}
    ></CustomInput>
    <div>
      <BadgeInput
        on:edited={(event) => {
          job.domains = event.detail;
        }}
        on:submit={saveJob}
        initialEntries={job.domains}
        label="Domains"
        placeholder="example.com"
        required
      ></BadgeInput>
    </div>
    <div class="grid gap-x-2 gap-y-4 md:grid-cols-2">
      <CustomInput id="client-key" label="Client key" bind:value={job.clientKey}></CustomInput>
      <CustomInput
        id="client-passphrase"
        type="password"
        label="Client passphrase"
        bind:value={job.clientPassphrase}
      ></CustomInput>
    </div>
    <div class="grid gap-x-2 gap-y-4 md:grid-cols-2">
      <CustomInput id="start-range" type="date" label="Start range" bind:value={job.startRange}
      ></CustomInput>
      <CustomInput id="end-range" type="date" label="End range" bind:value={job.endRange}
      ></CustomInput>
    </div>
    <CustomInput id="rate" label="Rate" placeholder="1.4" bind:value={job.rate}></CustomInput>
    <CustomInput id="ignore-pattern" label="Ignore pattern" bind:value={job.ignorePattern}
    ></CustomInput>
    <Checkbox
      on:change={() => {
        job.insecure = !job.insecure;
      }}
    >
      <span>Insecure</span>
    </Checkbox>
    <Checkbox
      on:change={() => {
        job.ignoreSignatureCheck = !job.ignoreSignatureCheck;
      }}
    >
      <span>Ignore signature check</span>
    </Checkbox>
    <Checkbox
      on:change={() => {
        job.temporary = !job.temporary;
      }}
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
