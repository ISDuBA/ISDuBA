<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Label, Select } from "flowbite-svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { push } from "svelte-spa-router";
  import { onMount } from "svelte";
  import CustomInput from "./CustomInput.svelte";
  export let params: any = null;
  $: disableSave =
    !schedule ||
    !schedule.name ||
    schedule.name === "" ||
    !schedule.job_id ||
    schedule.cron_timing === "";
  let loadError = "";
  let saveError = "";
  let jobs: any[] = [];

  type Schedule = {
    name: string;
    job_id: number | undefined;
    cron_timing: string;
  };
  let schedule: Schedule = {
    name: "",
    job_id: undefined,
    cron_timing: ""
  };

  onMount(async () => {
    let id;
    if (params) id = params.id;
    if (id) {
      const response = await request(`/api/cron`, "GET");
      if (response.ok) {
        const result = await response.content;
        const thisSchedule = result.find((s: any) => {
          return s.cron_id == id;
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
    if (schedule.cron_timing) formData.append("cron_timing", schedule.cron_timing);
    if (schedule.job_id) formData.append("job_id", schedule.job_id.toString());
    const method = params?.id ? "PUT" : "POST";
    const response = await request(`/api/cron`, method, formData);
    if (response.ok) {
      push(`/sources/`);
    } else if (response.error) {
      saveError = getErrorMessage(response.error);
    }
  }
</script>

<svelte:head>
  <title>Sources - {params?.id ? "Edit Schedule" : "New Schedule"}</title>
</svelte:head>

<SectionHeader title={params?.id ? "Edit Schedule" : "New Schedule"}></SectionHeader>
<form on:submit={saveSchedule} class="max-w-[800pt]">
  <div class="mb-4 flex flex-col gap-4">
    <CustomInput
      required
      label="Schedule name"
      id="name"
      placeholder="Schedule #1"
      bind:value={schedule.name}
    ></CustomInput>
    <div>
      <Label for="job-id" class="mb-2 block">
        <span>Job</span>
        <span class="text-red-500">*</span>
      </Label>
      <Select id="job-id" bind:value={schedule.job_id} placeholder="Choose job ...">
        {#each jobs as { id, name }}
          <option value={id}>{name}</option>
        {/each}
      </Select>
    </div>
    <CustomInput
      label="Cron expression"
      id="cron-timing"
      placeholder="0 0 1 * *"
      required
      bind:value={schedule.cron_timing}
    ></CustomInput>
  </div>
  <div class="mb-4 text-sm text-red-500">* Required fields</div>
  <Button disabled={disableSave} type="submit" color="light">
    <i class="bx bxs-save me-2"></i>
    <span>Save</span>
  </Button>
</form>
<ErrorMessage message={saveError}></ErrorMessage>
<ErrorMessage message={loadError}></ErrorMessage>
