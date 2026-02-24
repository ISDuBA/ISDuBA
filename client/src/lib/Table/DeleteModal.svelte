<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { SEARCHTYPES } from "$lib/Queries/query";
  import { request } from "$lib/request";
  import { appStore } from "$lib/store.svelte";
  import { Button, Modal, Spinner } from "flowbite-svelte";

  interface Props {
    documents?: any[];
    type: SEARCHTYPES;
    onDeleted: () => void;
  }

  let { documents = [], type, onDeleted }: Props = $props();

  let isDeleteModalOpen = $derived(appStore.state.app.isDeleteModalOpen);
  $effect(() => {
    if (!isDeleteModalOpen) {
      errorMessage = null;
    }
  });
  let errorMessage: ErrorDetails | null = $state(null);
  let isLoading = $state(false);

  const deleteDocuments = async () => {
    errorMessage = null;
    let url: string;
    let failed = false;
    isLoading = true;
    for (let i = 0; i < documents.length; i++) {
      const documentToDelete = documents[i];
      if (type === SEARCHTYPES.ADVISORY) {
        url = encodeURI(
          `/api/advisory/${documentToDelete.publisher}/${documentToDelete.tracking_id}`
        );
      } else {
        url = encodeURI(`/api/documents/${documentToDelete.id}`);
      }
      const response = await request(url, "DELETE");
      if (response.error) {
        errorMessage = getErrorDetails(
          `Could not delete ${type === SEARCHTYPES.ADVISORY ? "advisory" : "document"}`,
          response
        );
        failed = true;
      }
    }
    isLoading = false;
    if (!failed) {
      onDeleted();
      appStore.setIsDeleteModalOpen(false);
    }
  };
</script>

<Modal
  size="xs"
  title={documents.length === 1 ? documents[0].title : `Delete ${type}`}
  bind:open={appStore.state.app.isDeleteModalOpen}
  autoclose
  outsideclose
>
  <div class="text-center">
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      {#if documents.length === 1}
        Are you sure you want to delete this {type === SEARCHTYPES.ADVISORY
          ? "advisory"
          : "document"}?
      {:else}
        Are you sure you want to delete the selected {type} ?
      {/if}
    </h3>
    <Button
      onclick={() => {
        deleteDocuments();
      }}
      color="red"
      class="me-2"
    >
      <span>Yes, I'm sure</span>
      <div class:invisible={!isLoading} class:ms-2={true} class={isLoading ? "loadingFadeIn" : ""}>
        <Spinner color="gray" size="4"></Spinner>
      </div>
    </Button>
    <Button color="alternative">No, cancel</Button>
  </div>
  <ErrorMessage error={errorMessage}></ErrorMessage>
</Modal>
