<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Label, Input, Checkbox } from "flowbite-svelte";
  import { request } from "$lib/utils";

  let domains = "";
  let name = "";
  let jobId: number | null;
  let insecure: boolean = false;

  let jobLoadError = "";

  async function addJob() {
    const formData = new FormData();
    formData.append("domains", domains);
    formData.append("name", name);
    formData.append("insecure", insecure.toString());
    const response = await request(`/api/job`, "POST", formData);
    if (response.ok) {
      console.log("Success");
    } else if (response.error) {
      console.log(response.error);
    }
  }

  async function runJob() {
    const response = await request(`/api/job/` + jobId, "POST");
    if (response.ok) {
      console.log("Success");
    } else if (response.error) {
      console.log(response.error);
    }
  }

  async function getJobs() {
    const response = await request(`/api/job`, "GET");
    if (response.ok) {
      return response.content;
    } else if (response.error) {
      jobLoadError = response.error;
    }
    return [];
  }
</script>

<SectionHeader title="Sources"></SectionHeader>

<Label for="domain" class="mb-2 block">Domain</Label>
<Input id="domain" bind:value={domains} placeholder="example.com" />
<Label for="name" class="mb-2 block">Job Name</Label>
<Input id="name" bind:value={name} placeholder="Job #1" />
<Checkbox
  on:change={() => {
    insecure = !insecure;
  }}>Insecure</Checkbox
>
<br />
<Button
  on:click={() => {
    addJob();
  }}
  class="ml-auto mt-auto"
  color="primary">Add job</Button
>
<br />
Jobs:
{#if !jobLoadError}
  {#await getJobs() then jobs}
    <ul>
      {#each jobs.entries() as job}
        <li>{JSON.stringify(job, null, 4)}</li>
      {/each}
    </ul>
  {/await}
{/if}

<Label for="id" class="mb-2 block">Job ID</Label>
<Input id="id" bind:value={jobId} placeholder="42" />
<br />
<Button
  on:click={() => {
    runJob();
  }}
  class="ml-auto mt-auto"
  color="primary">Start job</Button
>
