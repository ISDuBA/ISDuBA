<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import {
    Button,
    Table,
    TableBody,
    TableBodyCell,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { request } from "$lib/utils";
  import { tablePadding, tdClass } from "$lib/table/defaults";

  let orderBy = "";
  let jobLoadError = "";

  async function runJob(id: number) {
    const response = await request(`/api/job/${id}`, "POST");
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
<Button href="/#/sources/new" class="mb-10" color="primary">
  <i class="bx bx-plus"></i>
  <span>Add job</span>
</Button>
<SectionHeader title="Jobs"></SectionHeader>
{#if !jobLoadError}
  {#await getJobs() then jobs}
    <Table hoverable={true} noborder={true}>
      <TableHead>
        <TableHeadCell padding={tablePadding} on:click={() => {}}>
          <span>ID</span>
          <i
            class:bx={true}
            class:bx-caret-up={orderBy == "name"}
            class:bx-caret-down={orderBy == "-name"}
          ></i>
        </TableHeadCell>
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Name<i
            class:bx={true}
            class:bx-caret-up={orderBy == "name"}
            class:bx-caret-down={orderBy == "-name"}
          ></i></TableHeadCell
        >
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Domains<i
            class:bx={true}
            class:bx-caret-up={orderBy == "description"}
            class:bx-caret-down={orderBy == "-description"}
          ></i>
        </TableHeadCell>
        <TableHeadCell padding={tablePadding} on:click={() => {}}
          >Insecure<i
            class:bx={true}
            class:bx-caret-up={orderBy == "description"}
            class:bx-caret-down={orderBy == "-description"}
          ></i>
        </TableHeadCell>
        <TableHeadCell padding={tablePadding} on:click={() => {}}>
          <span>Ignore signature Checkbox</span>
          <i
            class:bx={true}
            class:bx-caret-up={orderBy == "description"}
            class:bx-caret-down={orderBy == "-description"}
          ></i>
        </TableHeadCell>
        <TableHeadCell padding={tablePadding} on:click={() => {}}>
          <span>Worker</span>
          <i
            class:bx={true}
            class:bx-caret-up={orderBy == "description"}
            class:bx-caret-down={orderBy == "-description"}
          ></i>
        </TableHeadCell>
        <TableHeadCell></TableHeadCell>
      </TableHead>
      <TableBody>
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
            <TableBodyCell {tdClass}>{JSON.parse(job.domains).join(", ")}</TableBodyCell>
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
                <i class="bx bx-play-circle"></i>
              </button>
            </td>
          </tr>
        {/each}
      </TableBody>
    </Table>
  {/await}
{/if}
