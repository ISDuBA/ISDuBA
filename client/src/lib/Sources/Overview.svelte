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
  import { type Source, fetchSources } from "$lib/Sources/source";
  import { appStore } from "$lib/store";
  import CIconButton from "$lib/Components/CIconButton.svelte";

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
    const result = await fetchSources(true);
    loadingSources = false;
    if (result.ok) {
      sources = result.value;
    } else {
      sources = [];
      sourcesError = result.error;
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
{#if appStore.isEditor() || appStore.isSourceManager()}
  <CustomTable
    title="CSAF Provider"
    headers={[
      {
        label: "Name",
        attribute: "name"
      },
      {
        label: "Domain/PMD",
        attribute: "url"
      },
      {
        label: "Active",
        attribute: "active"
      },
      {
        label: "Downloading",
        attribute: "downloading"
      },
      {
        label: "Waiting",
        attribute: "waiting"
      }
    ]}
  >
    {#each sources as source, index (index)}
      <tr
        on:click={() => {
          if (appStore.isSourceManager()) {
            push(`/sources/${source.id}`);
          }
        }}
        on:blur={() => {}}
        on:focus={() => {}}
        class={appStore.isSourceManager() ? "cursor-pointer" : ""}
      >
        <TableBodyCell {tdClass}>{source.name}</TableBodyCell>
        <TableBodyCell {tdClass}>{source.url}</TableBodyCell>
        <TableBodyCell {tdClass}
          ><i class={"bx " + (source.active ? "bxs-circle" : "bx-circle")}></i></TableBodyCell
        >
        <TableBodyCell {tdClass}>{source.stats?.downloading}</TableBodyCell>
        <TableBodyCell {tdClass}>{source.stats?.waiting}</TableBodyCell>
        <td>
          <CIconButton
            on:click={() => {
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
            color="red"
            icon="trash"
          ></CIconButton>
        </td>
      </tr>
    {/each}
    <div slot="bottom">
      <div
        class:invisible={!loadingSources}
        class={loadingSources ? "loadingFadeIn" : ""}
        class:mb-4={true}
      >
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
{/if}
{#await getMessage() then resp}
  {#if resp.message}
    {resp.message}
  {/if}
{/await}
<ErrorMessage error={sourcesError}></ErrorMessage>
<ErrorMessage error={messageError}></ErrorMessage>

<br />
{#if appStore.isImporter()}
  <Button href="/#/sources/upload" class="my-2" color="primary" size="xs">
    <i class="bx bx-upload"></i>
    <span>Upload documents</span>
  </Button>
{/if}
