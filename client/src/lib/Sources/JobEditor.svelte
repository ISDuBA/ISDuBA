<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Checkbox, Input, Label } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import BadgeInput from "./BadgeInput.svelte";
  export let params: any = null;
  let disableSave = false;
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
<div class="mb-8 flex flex-col gap-4">
  <div>
    <BadgeInput
      on:edited={(event) => {
        job.domains = event.detail;
      }}
      initialEntries={job.domains}
      label="Domains"
      placeholder="example.com"
    ></BadgeInput>
  </div>
  <div>
    <Label for="name" class="mb-2 block">Job Name</Label>
    <Input id="name" bind:value={job.name} placeholder="Job #1" />
  </div>
  <div class="grid gap-x-2 gap-y-4 md:grid-cols-2">
    <div>
      <Label for="client-key" class="mb-2 block">Client Key</Label>
      <Input id="client-key" bind:value={job.clientKey} />
    </div>
    <div>
      <Label for="client-passphrase" class="mb-2 block">Client Passphrase</Label>
      <Input id="client-passphrase" type="password" bind:value={job.clientPassphrase} />
    </div>
  </div>
  <div class="grid gap-x-2 gap-y-4 md:grid-cols-2">
    <div>
      <Label for="start-range" class="mb-2 block">Start Range</Label>
      <Input id="start-range" type="date" bind:value={job.startRange} />
    </div>
    <div>
      <Label for="end-range" class="mb-2 block">End Range</Label>
      <Input id="end-range" type="date" bind:value={job.endRange} />
    </div>
  </div>
  <div>
    <Label for="rate" class="mb-2 block">Rate</Label>
    <Input id="rate" bind:value={job.rate} placeholder="1.4" />
  </div>
  <div>
    <Label for="ignore-pattern" class="mb-2 block">Ignore Pattern</Label>
    <Input id="ignore-pattern" bind:value={job.ignorePattern} />
  </div>
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
<Button disabled={disableSave} on:click={saveJob} color="light">
  <i class="bx bxs-save me-2"></i>
  <span>Save</span>
</Button>
<ErrorMessage message={saveError}></ErrorMessage>
<ErrorMessage message={loadError}></ErrorMessage>
