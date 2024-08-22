<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Spinner, TableBodyCell, Modal } from "flowbite-svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { push } from "svelte-spa-router";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { type ErrorDetails, getErrorDetails } from "$lib/Errors/error";
  import { tdClass } from "$lib/Table/defaults";
  import { request } from "$lib/request";
  import { onMount } from "svelte";
  import CustomTable from "$lib/Table/CustomTable.svelte";
  import type { Source } from "$lib/Sources/source";

  let messageError: ErrorDetails | null;
  let sourcesError: ErrorDetails | null;

  let modalOpen: boolean = false;
  let modalMessage = "";
  let modalTitle = "";
  let modalCallback: any;

  let loadingSources: boolean = false;

  let sources: Source[] = [];
  async function getMessage() {
    const response = await request("api/sources/message", "GET");
    if (response.ok) {
      return response.content;
    } else {
      messageError = getErrorDetails(`Couldn't load default message`, response);
    }
    return new Map<string, [string]>();
  }

  const getSources = async () => {
    loadingSources = true;
    const resp = await request(`/api/sources`, "GET");
    loadingSources = false;
    if (resp.ok) {
      if (resp.content.sources) {
        sources = resp.content.sources;
      } else {
        sources = [];
      }
    } else if (resp.error) {
      sourcesError = getErrorDetails(`Could not get sources`, resp);
    }
  };

  const deleteSource = async (id: number) => {
    const resp = await request(`/api/sources/${id}`, "DELETE");
    if (resp.error) {
      sourcesError = getErrorDetails(`Could not delete source`, resp);
    }
    await getSources();
  };

  onMount(() => {
    getSources();
  });
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
<SectionHeader title="Sources"></SectionHeader>
<CustomTable
  title="Sources"
  headers={[
    {
      label: "Name",
      attribute: "name"
    },
    {
      label: "URL",
      attribute: "url"
    },
    {
      label: "Active",
      attribute: "active"
    }
  ]}
>
  {#each sources as source, index (index)}
    <tr
      on:click={() => {
        push(`/sources/${source.id}`);
      }}
      on:blur={() => {}}
      on:focus={() => {}}
      class="cursor-pointer"
    >
      <TableBodyCell {tdClass}>{source.name}</TableBodyCell>
      <TableBodyCell {tdClass}>{source.url}</TableBodyCell>
      <TableBodyCell {tdClass}>{source.active}</TableBodyCell>
      <td>
        <Button
          on:click={(event) => {
            event.stopPropagation();
            modalCallback = () => {
              if (source.id) {
                deleteSource(source.id);
              }
            };
            modalMessage = "Are you sure you want to delete this source?";
            modalTitle = `Source ${source.name}`;
            modalOpen = true;
          }}
          title={`Delete source "${source.name}"`}
          class="border-0 p-2"
          color="light"
        >
          <i class="bx bx-trash text-xl text-red-500"></i>
        </Button>
      </td>
    </tr>
  {/each}
  <div slot="bottom">
    <div class:hidden={!loadingSources} class:mb-4={true}>
      Loading ...
      <Spinner color="gray" size="4"></Spinner>
    </div>
    <Button href="/#/sources/new" class="mb-2" color="primary" size="xs">
      <i class="bx bx-plus"></i>
      <span>Add source</span>
    </Button>
    <ErrorMessage error={sourcesError}></ErrorMessage>
  </div>
</CustomTable>
{#await getMessage() then resp}
  {#if resp.message}
    {resp.message}
  {/if}
{/await}
<ErrorMessage error={sourcesError}></ErrorMessage>
<ErrorMessage error={messageError}></ErrorMessage>
