<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Input, Label, Select } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  export let params: any = null;
  let disableSave = false;
  let loadError = "";
  let saveError = "";
  let jobs: any[] = [];

  type Schedule = {
    name: string;
    jobID: number | undefined;
    cronTiming: string;
  };
  let schedule: Schedule = {
    name: "",
    jobID: undefined,
    cronTiming: ""
  };

  onMount(async () => {
    let id;
    if (params) id = params.id;
    if (id) {
      const response = await request(`/api/cron`, "GET");
      if (response.ok) {
        const result = await response.content;
        const thisSchedule = result.find((s: any) => {
          return s.id == id;
        });
        if (params && params.id) {
          schedule = thisSchedule;
        }
      } else if (response.error) {
        loadError = `Could not load schedule. ${getErrorMessage(response.error)}`;
      }
    }
    const jobsResponse = await request(`/api/job`, "GET");
    if (jobsResponse.ok) {
      jobs = await jobsResponse.content;
    } else if (jobsResponse.error) {
      loadError = `Could not load all necessary information. ${getErrorMessage(jobsResponse.error)}`;
    }
  });

  async function saveSchedule() {
    const formData = new FormData();
    if (params?.id) formData.append("cron_id", params.id);
    formData.append("name", schedule.name);
    if (schedule.cronTiming) formData.append("cron_timing", schedule.cronTiming);
    if (schedule.jobID) formData.append("job_id", schedule.jobID.toString());
    const method = params?.id ? "PUT" : "POST";
    const response = await request(`/api/cron`, method, formData);
    if (response.ok) {
      push(`/sources/`);
    } else if (response.error) {
      saveError = getErrorMessage(response.error);
    }
  }
</script>

<SectionHeader title={params?.id ? "Edit Schedule" : "New Schedule"}></SectionHeader>
<div class="mb-8 flex flex-col gap-4">
  <div>
    <Label for="name" class="mb-2 block">Schedule name</Label>
    <Input id="name" bind:value={schedule.name} placeholder="Schedule #1" />
  </div>
  <div>
    <Label for="job-id" class="mb-2 block">Job</Label>
    <Select id="job-id" bind:value={schedule.jobID} placeholder="Choose job ...">
      {#each jobs as { id, name }}
        <option value={id}>{name}</option>
      {/each}
    </Select>
  </div>
  <div>
    <Label for="cron-timing" class="mb-2 block">Cron expression</Label>
    <Input id="cron-timing" bind:value={schedule.cronTiming} placeholder="0 0 1 * *" />
  </div>
</div>
<Button disabled={disableSave} on:click={saveSchedule} color="light">
  <i class="bx bxs-save me-2"></i>
  <span>Save</span>
</Button>
<ErrorMessage message={saveError}></ErrorMessage>
<ErrorMessage message={loadError}></ErrorMessage>
