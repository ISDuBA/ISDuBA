<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import {
    Button,
    Dropzone,
    Img,
    P,
    Spinner,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { request } from "$lib/request";
  import { getPublisher } from "$lib/publisher";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { onDestroy, onMount } from "svelte";
  import CIconButton from "$lib/Components/CIconButton.svelte";
  import { getRelativeTime } from "$lib/time";

  $: docA_ID = $appStore.app.diff.docA_ID;
  $: docB_ID = $appStore.app.diff.docB_ID;
  $: if (docA_ID) updateDocumentA();
  $: if (docB_ID) updateDocumentB();
  $: docA = $appStore.app.diff.docA;
  $: docB = $appStore.app.diff.docB;
  $: isToolboxOpen = $appStore.app.isToolboxOpen;
  $: if (isToolboxOpen) {
    getTempDocuments();
    getDocuments();
  }
  $: disableDiffButtons = docA_ID !== undefined && docB_ID !== undefined;

  let freeTempDocuments = 0;
  let tempDocuments: any[] = [];
  let loadDocumentAErrorMessage: ErrorDetails | null;
  let loadDocumentBErrorMessage: ErrorDetails | null;
  let tempDocErrorMessage: ErrorDetails | null;
  let loadTempDocsErrorMessage: ErrorDetails | null;
  let intervalID: ReturnType<typeof setTimeout> | undefined = undefined;
  let isLoadingDocA = false;
  let isLoadingDocB = false;

  const tdClass = "pe-5 py-0 whitespace-nowrap font-medium";
  const padding = "pe-5 pt-2";

  onMount(() => {
    intervalID = setInterval(() => {
      if (tempDocuments.length > 0) {
        updateExpired();
      }
    }, 60000);
  });

  onDestroy(() => {
    clearInterval(intervalID);
  });

  const resetTempDocsErrorMessages = () => {
    tempDocErrorMessage = null;
    loadTempDocsErrorMessage = null;
  };

  const dropHandle = (event: any) => {
    event.preventDefault();
    if (event.dataTransfer.items) {
      [...event.dataTransfer.items].forEach(async (item) => {
        if (freeTempDocuments === 0) {
          tempDocErrorMessage = getErrorDetails(
            "You reached the maximal number of temporary documents."
          );
          return;
        }
        if (item.kind === "file") {
          const file = item.getAsFile();
          await uploadFile(file);
        }
      });
      getTempDocuments();
    }
  };

  const handleChange = async (event: any) => {
    const files = event.target.files;
    for (let i = 0; i < files.length; i++) {
      const file = files[i];
      if (freeTempDocuments === 0) {
        tempDocErrorMessage = getErrorDetails(
          "You reached the maximal number of temporary documents."
        );
        break;
      }
      await uploadFile(file);
    }
    getTempDocuments();
  };

  const uploadFile = (file: File): Promise<void> => {
    resetTempDocsErrorMessages();
    return new Promise((resolve) => {
      const formData = new FormData();
      formData.append("file", file);
      request("/api/tempdocuments", "POST", formData).then((response) => {
        if (response.ok) {
          freeTempDocuments = freeTempDocuments - 1;
          resolve();
        } else if (response.error) {
          tempDocErrorMessage = getErrorDetails(`Could not upload file`, response);
          resolve();
        }
      });
    });
  };

  const getDocuments = async () => {
    updateDocumentA();
    updateDocumentB();
  };

  const updateDocumentA = async () => {
    loadDocumentAErrorMessage = null;
    if (docA_ID) {
      isLoadingDocA = true;
      const responseDocA = await fetchDocument("A");
      if (responseDocA.ok) {
        if (docA_ID) {
          const content = await responseDocA.content;
          appStore.setDiffDocA(content);
        }
      } else if (responseDocA.error) {
        if (responseDocA.error === "404") {
          appStore.setDiffDocA_ID(undefined);
        } else {
          loadDocumentAErrorMessage = getErrorDetails(
            `Could not load first document.`,
            responseDocA
          );
        }
      }
      isLoadingDocA = false;
    } else {
      appStore.setDiffDocA(undefined);
    }
  };

  const updateDocumentB = async () => {
    loadDocumentBErrorMessage = null;
    if (docB_ID) {
      isLoadingDocB = true;
      const responseDocB = await fetchDocument("B");
      if (responseDocB.ok) {
        if (docB_ID) {
          const content = await responseDocB.content;
          appStore.setDiffDocB(content);
        }
      } else if (responseDocB.error) {
        if (responseDocB.error === "404") {
          appStore.setDiffDocB_ID(undefined);
        } else {
          loadDocumentBErrorMessage = getErrorDetails(
            `Could not load second document.`,
            responseDocB
          );
        }
      }
      isLoadingDocB = false;
    } else {
      appStore.setDiffDocB(undefined);
    }
  };

  const fetchDocument = async (letter: string) => {
    const docID = letter === "A" ? docA_ID : docB_ID;
    const endpoint = docID?.startsWith("tempdocument") ? "tempdocuments" : "documents";
    const id = docID?.startsWith("tempdocument") ? docID.replace("tempdocument", "") : docID;
    return request(`/api/${endpoint}/${id}`, "GET");
  };

  const getTempDocuments = async () => {
    loadTempDocsErrorMessage = null;
    const response = await request("/api/tempdocuments", "GET");
    if (response.ok) {
      const result = await response.content;
      freeTempDocuments = result.free;
      const newTempDocuments: any[] = [];
      for (let i = 0; i < result.files.length; i++) {
        const file = result.files[i];
        const response2 = await request(`/api/tempdocuments/${file.id}`, "GET");
        if (response2.ok) {
          const result2 = await response2.content;
          newTempDocuments.push({
            file: file,
            document: result2.document
          });
        }
      }
      tempDocuments = newTempDocuments;
      updateExpired();
    } else if (response.error) {
      loadTempDocsErrorMessage = getErrorDetails(`Could not load temporary document.`, response);
    }
  };

  const deleteTempDocument = async (id: number) => {
    resetTempDocsErrorMessages();
    const response = await request(`/api/tempdocuments/${id}`, "DELETE");
    if (response.ok) {
      if (docA_ID === `tempdocument${id}`) {
        appStore.setDiffDocA_ID(undefined);
      } else if (docB_ID === `tempdocument${id}`) {
        appStore.setDiffDocB_ID(undefined);
      }
      getTempDocuments();
    } else if (response.error) {
      tempDocErrorMessage = getErrorDetails(`Could not delete temporary document.`, response);
    }
  };

  const deleteExpiredDocFromStore = (id: number) => {
    if (docA_ID === `tempdocument${id}`) {
      appStore.setDiffDocA_ID(undefined);
    } else if (docB_ID === `tempdocument${id}`) {
      appStore.setDiffDocB_ID(undefined);
    }
  };

  const updateExpired = () => {
    let didDocExpire = false;
    tempDocuments.forEach((doc) => {
      const expiredDate = new Date(doc.file.expired);
      if (expiredDate.getTime() < Date.now()) {
        didDocExpire = true;
        deleteExpiredDocFromStore(doc.file.id);
      } else {
        doc.document.expired = getRelativeTime(expiredDate);
      }
    });
    if (didDocExpire) getTempDocuments();
    else tempDocuments = tempDocuments;
  };
</script>

<div class="flex w-full flex-col">
  <div class="mb-4 flex flex-col justify-between gap-2 lg:flex-row">
    <div class="flex items-start gap-1">
      <div class="flex min-h-28 justify-between gap-1 rounded-md pb-2 pe-3 lg:w-96 lg:max-w-96">
        {#if isLoadingDocA}
          <div>
            Loading ...
            <Spinner color="gray" size="4"></Spinner>
          </div>
        {:else if docA}
          <div>
            <Button
              on:click={() => {
                appStore.setDiffDocA_ID(undefined);
                appStore.setDiffDocA(undefined);
              }}
              color="light"
              class="border-0 p-1"
            >
              <i class="bx bx-x text-lg"></i>
            </Button>
          </div>
          <div class="lg:flex lg:flex-col">
            <div>{docA.document.tracking.id}</div>
            <div class="md:mb-1" title={docA.document.title}>
              {docA.document.title}
            </div>
            <span class="text-sm text-gray-600">{getPublisher(docA.document.publisher.name)}</span>
            <span class="text-sm text-gray-600">Version: {docA.document.tracking.version}</span>
          </div>
        {:else}
          <div class="flex flex-col gap-2">
            <P italic>Select a document or upload local ones.</P>
          </div>
        {/if}
      </div>
    </div>
    <div class="flex gap-1">
      <div class="flex min-h-28 justify-between gap-1 rounded-md pb-2 lg:w-96 lg:max-w-96">
        {#if isLoadingDocB}
          <div>
            Loading ...
            <Spinner color="gray" size="4"></Spinner>
          </div>
        {:else if docB}
          <div>
            <Button
              on:click={() => {
                appStore.setDiffDocB_ID(undefined);
                appStore.setDiffDocB(undefined);
              }}
              color="light"
              class="border-0 p-1"
            >
              <i class="bx bx-x text-lg"></i>
            </Button>
          </div>
          <div class="lg:flex lg:flex-col">
            <div>{docB.document.tracking.id}</div>
            <div class="md:mb-1" title={docB.document.title}>
              {docB.document.title}
            </div>
            <span class="text-sm text-gray-600">{getPublisher(docB.document.publisher.name)}</span>
            <span class="text-sm text-gray-600">Version: {docB.document.tracking.version}</span>
          </div>
        {:else}
          <div class:flex={true} class:flex-col={true} class:gap-2={true} class:invisible={!docA}>
            <P italic>Select a document or upload local ones.</P>
          </div>
        {/if}
      </div>
    </div>
    <div class="flex h-full items-start">
      <Button
        on:click={() => push("/diff")}
        disabled={!docA || !docB || !docA_ID || !docB_ID}
        size="sm"
        class="flex gap-x-2"
      >
        <Img src="plus-minus.svg" class="w-5 invert" />
        <span>Compare</span>
      </Button>
    </div>
  </div>
  <ErrorMessage error={loadDocumentAErrorMessage}></ErrorMessage>
  <ErrorMessage error={loadDocumentBErrorMessage}></ErrorMessage>
  <div class="flex flex-col">
    {#if tempDocuments?.length > 0}
      <span class="mb-1">Temporary documents:</span>
      <Table>
        <TableHead>
          <TableHeadCell {padding}></TableHeadCell>
          <TableHeadCell {padding}>Tracking ID</TableHeadCell>
          <TableHeadCell {padding}>Publisher</TableHeadCell>
          <TableHeadCell {padding}>Title</TableHeadCell>
          <TableHeadCell {padding}>Expires in</TableHeadCell>
          <TableHeadCell {padding}>File name</TableHeadCell>
        </TableHead>
        <TableBody>
          {#each tempDocuments as document}
            {@const doc = document.document}
            {@const tempDocID = `tempdocument${document.file.id}`}
            <TableBodyRow>
              <TableBodyCell {tdClass}>
                <div class="flex items-center">
                  <CIconButton
                    on:click={() => {
                      deleteTempDocument(document.file.id);
                    }}
                    color="red"
                    title={`delete ${doc.title} - ${doc.tracking.id}`}
                    icon="trash"
                  ></CIconButton>
                  <button
                    on:click|stopPropagation={(e) => {
                      if (docA_ID) {
                        appStore.setDiffDocB_ID(tempDocID);
                      } else {
                        appStore.setDiffDocA_ID(tempDocID);
                      }
                      e.preventDefault();
                    }}
                    class:invisible={!$appStore.app.isToolboxOpen}
                    disabled={docA_ID === tempDocID || docB_ID === tempDocID || disableDiffButtons}
                    title={`compare ${doc.title} - ${doc.tracking.id}`}
                  >
                    <Img
                      src="plus-minus.svg"
                      class={`${
                        docA_ID === tempDocID || docB_ID === tempDocID || disableDiffButtons
                          ? "invert-[70%]"
                          : ""
                      } min-h-4 min-w-6`}
                    />
                  </button>
                </div>
              </TableBodyCell>
              <TableBodyCell {tdClass}>{doc.tracking.id}</TableBodyCell>
              <TableBodyCell {tdClass}>{doc.publisher.name}</TableBodyCell>
              <TableBodyCell {tdClass}>
                <span
                  class="block overflow-hidden text-ellipsis whitespace-nowrap md:w-44 lg:w-60 xl:w-96"
                  title={doc.title}>{doc.title}</span
                >
              </TableBodyCell>
              <TableBodyCell {tdClass}>{doc.expired}</TableBodyCell>
              <TableBodyCell {tdClass}>{document.file.filename}</TableBodyCell>
            </TableBodyRow>
          {/each}
          <TableBodyRow></TableBodyRow>
        </TableBody>
      </Table>
    {/if}
    {#if freeTempDocuments}
      <Dropzone
        on:drop={dropHandle}
        on:dragover={(event) => {
          event.preventDefault();
        }}
        on:change={handleChange}
        multiple
        class="mb-2 ms-1 h-16"
      >
        <i class="bx bx-upload text-xl text-gray-500"></i>
        <p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
          Upload temporary documents ({freeTempDocuments} free {freeTempDocuments > 1
            ? "slots"
            : "slot"} left)
        </p>
      </Dropzone>
    {/if}
    <ErrorMessage error={tempDocErrorMessage}></ErrorMessage>
    <ErrorMessage error={loadTempDocsErrorMessage}></ErrorMessage>
  </div>
</div>
