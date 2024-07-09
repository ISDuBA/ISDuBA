<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Modal, Spinner, TableBodyCell } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { request } from "$lib/utils";
  import { tdClass } from "$lib/Table/defaults";
  import { getErrorMessage } from "$lib/Errors/error";
  import { onDestroy, onMount } from "svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import { TASK_STATE_RUNNING, TASK_STATE_ABORTED } from "./sources";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  let cronLoadError = "";
  let jobError = "";
  let taskLoadError = "";
  let loadingCrons = false;
  let loadingJobs = false;
  let loadingTasks = false;
  let crons: any[] = [];
  let jobs: any[] = [];
  let tasks: any[] = [];
  let modalOpen = false;
  let modalMessage: string;
  let modalTitle: string;
  let modalCallback: any;
  let intervalId: number;

  onMount(() => {
    getJobs();
    getTasks();
    getCrons();
    intervalId = setInterval(() => {
      getTasks();
    }, 60000);
  });

  onDestroy(() => {
    clearInterval(intervalId);
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
      jobError = getErrorMessage(response.error);
    }
    loadingJobs = false;
  }

  async function deleteJob(id: number) {
    const response = await request(`/api/job/${id}`, "DELETE");
    if (response.ok) {
      const index = jobs.findIndex((job) => job.id === id);
      jobs = jobs.toSpliced(index, 1);
    } else if (response.error) {
      const index = jobs.findIndex((job) => job.id === id);
      jobs = jobs.toSpliced(index, 1);
      jobError = getErrorMessage(response.error);
    }
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

  async function downloadTaskLog(taskID: number) {
    const response = await request(`/api/task/` + taskID, "GET");
    if (response.ok) {
      let blob = new Blob([response.content]);
      let url = URL.createObjectURL(blob);
      let link = document.createElement("a");
      link.setAttribute("download", "task" + taskID + ".log");
      link.href = url;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(url);
    } else if (response.error) {
      taskLoadError = getErrorMessage(response.error);
    }
  }

  async function cancelTask(id: number) {
    const response = await request(`/api/task/${id}`, "DELETE");
    if (response.ok) {
      const index = tasks.findIndex((tasks) => tasks.task_id === id);
      tasks[index].status = TASK_STATE_ABORTED;
    } else if (response.error) {
      taskLoadError = getErrorMessage(response.error);
    }
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

  async function deleteCron(id: number) {
    const response = await request(`/api/cron/${id}`, "DELETE");
    if (response.ok) {
      const index = crons.findIndex((cron) => cron.id === id);
      crons = crons.toSpliced(index, 1);
    } else if (response.error) {
      jobError = getErrorMessage(response.error);
    }
  }
</script>

<svelte:head>
  <title>Sources</title>
</svelte:head>

<Modal size="xs" title={modalTitle} bind:open={modalOpen} autoclose outsideclose>
  <div class="text-center">
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {modalMessage}
    </h3>
    <Button
      on:click={() => {
        modalCallback();
      }}
      color="red"
      class="me-2">Yes, I'm sure</Button
    >
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>

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
    },
    {
      label: "Actions",
      attribute: undefined
    }
  ]}
>
  {#each jobs as job, index (index)}
    <tr
      on:click={() => {
        push(`/sources/job/${job.id}`);
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
        <Button
          on:click={(event) => {
            event.stopPropagation();
            modalCallback = () => deleteJob(job.id);
            modalMessage = "Are you sure you want to delete this job?";
            modalTitle = `Job ${job.name}`;
            modalOpen = true;
          }}
          title={`Remove Job "${job.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-trash text-xl text-red-500"></i>
        </Button>
        <Button
          on:click={(event) => {
            event.stopPropagation();
            runJob(job.id);
          }}
          title={`Run Job "${job.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-play-circle text-xl"></i>
        </Button>
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingJobs} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button href="/#/sources/job/new" class="mb-2" color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>Add job</span>
    </Button>
    <ErrorMessage message={jobError}></ErrorMessage>
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
    },
    {
      label: "Actions",
      attribute: undefined
    }
  ]}
>
  <div slot="header-right">
    <Button
      on:click={getTasks}
      disabled={loadingTasks}
      title="Refresh tasks"
      class="!p-1"
      color="light"
      size="xs"
    >
      <i class={`bx bx-sync scale-x-[-1] text-lg ${loadingTasks ? "animate-spin" : ""}`}></i>
    </Button>
  </div>
  {#each tasks as task, index (index)}
    <tr on:click={() => {}} on:blur={() => {}} on:focus={() => {}} class="cursor-default">
      <TableBodyCell {tdClass}>{task.task_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{task.job_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{task.created}</TableBodyCell>
      <TableBodyCell {tdClass}>{task.status}</TableBodyCell>
      <TableBodyCell {tdClass}>
        <Button
          title="Download log"
          class="border-0 p-2"
          color="light"
          on:click={() => {
            downloadTaskLog(task.task_id);
          }}
        >
          <i class="bx bx-download text-lg"></i>
        </Button>
      </TableBodyCell>
      <td>
        {#if task.status === TASK_STATE_RUNNING}
          <Button
            on:click={(event) => {
              event.stopPropagation();
              modalCallback = () => cancelTask(task.task_id);
              modalMessage = "Are you sure you want to cancel this task?";
              modalTitle = `Task ${task.task_id}`;
              modalOpen = true;
            }}
            title={`Cancel task "${task.task_id}"`}
            class="border-0 p-2"
            color="light"
          >
            <i class="bx bx-stop-circle text-xl"></i>
          </Button>
        {/if}
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <ErrorMessage message={taskLoadError}></ErrorMessage>
  </div>
</CustomTable>

<CustomTable
  title="Schedules"
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
      label: "Cron expression",
      attribute: "cron_timing"
    },
    {
      label: "Actions",
      attribute: undefined
    }
  ]}
>
  {#each crons as cron, index (index)}
    <tr
      on:click={() => {
        // TODO: Un-comment following line and add 'class="cursor-pointer"'
        // when editing is implemented in backend
        // push(`/sources/schedule/${cron.cron_id}`);
      }}
      on:blur={() => {}}
      on:focus={() => {}}
    >
      <TableBodyCell {tdClass}>{cron.cron_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{cron.job_id}</TableBodyCell>
      <TableBodyCell {tdClass}>{cron.name}</TableBodyCell>
      <TableBodyCell {tdClass}>{cron.cron_timing}</TableBodyCell>
      <td>
        <Button
          on:click={(event) => {
            event.stopPropagation();
            modalCallback = () => deleteCron(cron.cron_id);
            modalMessage = "Are you sure you want to delete this schedule?";
            modalTitle = `Schedule ${cron.name}`;
            modalOpen = true;
          }}
          title={`Remove Schedule "${cron.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-trash text-xl text-red-500"></i>
        </Button>
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingCrons} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button href="/#/sources/schedule/new" class="mb-2" color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>Add schedule</span>
    </Button>
    <ErrorMessage message={cronLoadError}></ErrorMessage>
  </div>
</CustomTable>
