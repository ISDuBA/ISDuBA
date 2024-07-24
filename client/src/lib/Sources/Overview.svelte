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
  import { tdClass } from "$lib/Table/defaults";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

  let sourceError = "";
  let loadingSources = false;
  let sources: any[] = [
    {
      name: "RedHat",
      source: "redhat.com",
      interval: "every 12 hours",
      lastRun: new Date(Date.now() - 1000 * 60)
    }
  ];
  let aggregators: any[] = [
    {
      name: "BSI Lister",
      source: "https://wid.cert-bund.de/.well-known/csaf-aggregator/aggregator.json",
      interval: "every 12 hours",
      lastRun: new Date(Date.now() - 1000 * 60 * 60)
    }
  ];
  let modalOpen = false;
  let modalMessage: string;
  let modalTitle: string;
  let modalCallback: any;

  const getRelativeTime = (date: Date) => {
    const now = Date.now();
    const passedTime = now - date.getTime();
    if (passedTime < 60000) {
      return "<1 min ago";
    } else if (passedTime < 3600000) {
      return `${Math.floor(passedTime / 60000)} min ago`;
    } else if (passedTime < 86400000) {
      return `${Math.floor(passedTime / 3600000)} hours ago`;
    } else {
      return `${Math.floor(passedTime / 86400000)} days ago`;
    }
  };
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
  title="Sources"
  headers={[
    {
      label: "Name",
      attribute: "name"
    },
    {
      label: "Domain/PMD",
      attribute: "source"
    },
    {
      label: "Interval",
      attribute: "interval"
    },
    {
      label: "Last run",
      attribute: "lastRun"
    },
    {
      label: "Actions",
      attribute: undefined
    }
  ]}
>
  {#each sources as source, index (index)}
    <tr on:blur={() => {}} on:focus={() => {}} class="cursor-pointer">
      <TableBodyCell {tdClass}>{source.name}</TableBodyCell>
      <TableBodyCell {tdClass}>{source.source}</TableBodyCell>
      <TableBodyCell {tdClass}>{source.interval}</TableBodyCell>
      <TableBodyCell {tdClass}>{getRelativeTime(source.lastRun)}</TableBodyCell>
      <td>
        <Button
          on:click={() => {
            push(`/sources/job/${source.id}`);
          }}
          title={`Edit Source "${source.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-edit-alt text-xl"></i>
        </Button>
        <Button
          on:click={(event) => {
            event.stopPropagation();
            modalMessage = "Are you sure you want to delete this source?";
            modalTitle = `Source ${source.name}`;
            modalOpen = true;
          }}
          title={`Remove Source "${source.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-trash text-xl text-red-500"></i>
        </Button>
        <Button
          on:click={(event) => {
            event.stopPropagation();
          }}
          title={`Download from source "${source.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-play-circle text-xl"></i>
        </Button>
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingSources} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button href="/#/sources/job/new" class="mb-2" color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>Add new source</span>
    </Button>
    <ErrorMessage message={sourceError}></ErrorMessage>
  </div>
</CustomTable>

<CustomTable
  title="Aggregators"
  headers={[
    {
      label: "Name",
      attribute: "name"
    },
    {
      label: "URL",
      attribute: "source"
    },
    {
      label: "Interval",
      attribute: "interval"
    },
    {
      label: "Last checked",
      attribute: "lastRun"
    },
    {
      label: "Actions",
      attribute: undefined
    }
  ]}
>
  {#each aggregators as aggregator, index (index)}
    <tr on:blur={() => {}} on:focus={() => {}} class="cursor-pointer">
      <TableBodyCell {tdClass}>{aggregator.name}</TableBodyCell>
      <TableBodyCell {tdClass}>{aggregator.source}</TableBodyCell>
      <TableBodyCell {tdClass}>{aggregator.interval}</TableBodyCell>
      <TableBodyCell {tdClass}>{getRelativeTime(aggregator.lastRun)}</TableBodyCell>
      <td>
        <Button
          on:click={() => {
            push(`/sources/job/${aggregator.id}`);
          }}
          title={`Edit Source "${aggregator.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-edit-alt text-xl"></i>
        </Button>
        <Button
          on:click={(event) => {
            event.stopPropagation();
            modalMessage = "Are you sure you want to delete this source?";
            modalTitle = `Source ${aggregator.name}`;
            modalOpen = true;
          }}
          title={`Remove Source "${aggregator.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-trash text-xl text-red-500"></i>
        </Button>
        <Button
          on:click={(event) => {
            event.stopPropagation();
          }}
          title={`Download from source "${aggregator.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-play-circle text-xl"></i>
        </Button>
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingSources} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button href="/#/sources/job/new" class="mb-2" color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>Add new lister</span>
    </Button>
    <ErrorMessage message={sourceError}></ErrorMessage>
  </div>
</CustomTable>
