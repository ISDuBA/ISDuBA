<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Spinner, TableBodyCell } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { request } from "$lib/utils";
  import { tdClass } from "$lib/table/defaults";
  import { getErrorMessage } from "$lib/Errors/error";
  import { onMount } from "svelte";
  import CustomTable from "$lib/table/CustomTable.svelte";
  import { TASK_STATE_RUNNING } from "./sources";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  let cronLoadError = "";
  let jobLoadError = "";
  let taskLoadError = "";
  let loadingCrons = false;
  let loadingJobs = false;
  let loadingTasks = false;
  let crons: any[] = [];
  let jobs: any[] = [];
  let tasks: any[] = [];

  onMount(() => {
    getJobs();
    getTasks();
    getCrons();
  });

  async function runJob(id: number) {
    const response = await request(`/api/job/${id}`, "POST");
    if (response.ok) {
      console.log("Success");
      getTasks();
    } else if (response.error) {
      console.log(response.error);
    }
  }

  async function getJobs() {
    loadingJobs = true;
    const response = await request(`/api/job`, "GET");
    if (response.ok) {
      jobs = response.content;
    } else if (response.error) {
      jobLoadError = getErrorMessage(response.error);
    }
    loadingJobs = false;
  }

  async function getTasks() {
    loadingTasks = true;
    const response = await request(`/api/task`, "GET");
    if (response.ok) {
      tasks = response.content;
    } else if (response.error) {
      taskLoadError = getErrorMessage(response.error);
    }
    loadingTasks = false;
  }

  async function cancelTask(id: number) {
    console.log("cancelTask", id);
    // TODO
  }

  async function getCrons() {
    loadingCrons = true;
    const response = await request(`/api/cron`, "GET");
    if (response.ok) {
      crons = response.content;
    } else if (response.error) {
      cronLoadError = getErrorMessage(response.error);
    }
    loadingCrons = false;
  }
</script>

<CustomTable
  title="Jobs"
  headers={[
    {
      label: "ID",
      attribute: "id"
    },
    {
      label: "Name",
      attribute: "name"
    },
    {
      label: "Domains",
      attribute: "domains"
    },
    {
      label: "Insecure",
      attribute: "insecure"
    },
    {
      label: "Ignore signature Checkbox",
      attribute: "ignore_signature_checkbox"
    },
    {
      label: "Worker",
      attribute: "worker"
    }
  ]}
>
  {#each jobs as job, index (index)}
    <tr
      on:click={() => {
        push(`/sources/${job.id}`);
      }}
      on:blur={() => {}}
      on:focus={() => {}}
      class="cursor-pointer"
    >
      <TableBodyCell {tdClass}>{job.id}</TableBodyCell>
      <TableBodyCell {tdClass}>{job.name}</TableBodyCell>
      <TableBodyCell {tdClass}>{job.domains.join(", ")}</TableBodyCell>
      <TableBodyCell {tdClass}>{job.insecure}</TableBodyCell>
      <TableBodyCell {tdClass}>{job.ignore_signature_check}</TableBodyCell>
      <TableBodyCell {tdClass}>{job.worker}</TableBodyCell>
      <td>
        <button
          title={`Run Job "${job.name}"`}
          on:click|stopPropagation={() => {
            runJob(job.id);
          }}
        >
          <i class="bx bx-play-circle text-xl"></i>
        </button>
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingJobs} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button href="/#/sources/new" class="mb-2" color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>Add job</span>
    </Button>
    <ErrorMessage message={jobLoadError}></ErrorMessage>
  </div>
</CustomTable>

<CustomTable
  title="Tasks"
  headers={[
    {
      label: "ID",
      attribute: "id"
    },
    {
      label: "Job ID",
      attribute: "job_id"
    },
    {
      label: "Created",
      attribute: "created"
    },
    {
      label: "Status",
      attribute: "status"
    }
  ]}
>
  {#each tasks as task, index (index)}
    <tr on:click={() => {}} on:blur={() => {}} on:focus={() => {}} class="cursor-pointer">
      <TableBodyCell {tdClass}>{task.task_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{task.job_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{task.created}</TableBodyCell>
      <TableBodyCell {tdClass}>{task.status}</TableBodyCell>
      <td>
        {#if task.status === TASK_STATE_RUNNING}
          <button
            title={`Cancel task "${task.task_id}"`}
            on:click|stopPropagation={() => {
              cancelTask(task.task_id);
            }}
          >
            <i class="bx bx-stop-circle text-xl"></i>
          </button>
        {/if}
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingTasks} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <ErrorMessage message={taskLoadError}></ErrorMessage>
  </div>
</CustomTable>

<CustomTable
  title="Crons"
  headers={[
    {
      label: "ID",
      attribute: "id"
    },
    {
      label: "Job ID",
      attribute: "job_id"
    },
    {
      label: "Name",
      attribute: "name"
    },
    {
      label: "Cron timing",
      attribute: "cron_timing"
    }
  ]}
>
  {#each crons as cron, index (index)}
    <tr on:click={() => {}} on:blur={() => {}} on:focus={() => {}} class="cursor-pointer">
      <TableBodyCell {tdClass}>{cron.id}</TableBodyCell>
      <TableBodyCell {tdClass}>{cron.job_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{cron.name}</TableBodyCell>
      <TableBodyCell {tdClass}>{cron.cron_timing}</TableBodyCell>
      <td> </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingCrons} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <ErrorMessage message={cronLoadError}></ErrorMessage>
  </div>
</CustomTable>
