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
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  import { getPublisher, request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";

  $: $appStore.app.diff.docA_ID, getDocuments();
  $: $appStore.app.diff.docB_ID, getDocuments();
  $: if ($appStore.app.diff.isDiffModeEnabled) {
    getTempDocuments();
    getDocuments();
  }
  $: disableDiffButtons =
    $appStore.app.diff.docA_ID !== undefined && $appStore.app.diff.docB_ID !== undefined;

  let freeTempDocuments = 0;
  let tempDocuments: any[] = [];
  let docA: any;
  let docB: any;
  let tempDocErrorMessage = "";

  const tdClass = "px-6 py-0 whitespace-nowrap font-medium";

  const dropHandle = (event: any) => {
    event.preventDefault();
    if (event.dataTransfer.items) {
      [...event.dataTransfer.items].forEach(async (item) => {
        if (freeTempDocuments === 0) {
          tempDocErrorMessage = "You reached the maximal number of temporary documents.";
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
        tempDocErrorMessage = "You reached the maximal number of temporary documents.";
        break;
      }
      await uploadFile(file);
    }
    getTempDocuments();
  };

  const uploadFile = (file: File): Promise<void> => {
    return new Promise((resolve) => {
      const formData = new FormData();
      formData.append("file", file);
      request("/api/tempdocuments", "POST", formData).then((response) => {
        if (response.ok) {
          freeTempDocuments = freeTempDocuments - 1;
          resolve();
        } else if (response.error) {
          tempDocErrorMessage = getErrorMessage(response.error);
          resolve();
        }
      });
    });
  };

  const getDocuments = async () => {
    if ($appStore.app.diff.docA_ID) {
      const responseDocA = await getDocument("A");
      if (responseDocA.ok) {
        docA = await responseDocA.content;
      }
    } else {
      docA = undefined;
    }
    if ($appStore.app.diff.docB_ID) {
      const responseDocB = await getDocument("B");
      if (responseDocB.ok) {
        docB = await responseDocB.content;
      }
    } else {
      docB = undefined;
    }
  };

  const getDocument = async (letter: string) => {
    const docID = letter === "A" ? $appStore.app.diff.docA_ID : $appStore.app.diff.docB_ID;
    const endpoint = docID?.startsWith("tempdocument") ? "tempdocuments" : "documents";
    const id = docID?.startsWith("tempdocument") ? docID.replace("tempdocument", "") : docID;
    return request(`/api/${endpoint}/${id}`, "GET");
  };

  const getTempDocuments = async () => {
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
    }
  };

  const deleteTempDocument = async (id: number) => {
    tempDocErrorMessage = "";
    const response = await request(`/api/tempdocuments/${id}`, "DELETE");
    if (response.ok) {
      if ($appStore.app.diff.docA_ID === `tempdocument${id}`) {
        appStore.setDiffDocA_ID(undefined);
      } else if ($appStore.app.diff.docB_ID === `tempdocument${id}`) {
        appStore.setDiffDocB_ID(undefined);
      }
      getTempDocuments();
    } else if (response.error) {
      tempDocErrorMessage = getErrorMessage(response.error);
    }
  };
</script>

<div class="fixed bottom-0 left-20 flex w-full flex-col items-center justify-center">
  <Button on:click={appStore.toggleDiffMode} class="max-w-32 rounded-none border-b-0" color="light">
    <Img src="plus-minus.svg" class="min-h-4 min-w-4" />
  </Button>
  {#if $appStore.app.diff.isDiffModeEnabled}
    <div
      class="flex items-stretch gap-6 border border-solid border-gray-300 bg-white px-2 py-4 shadow-gray-800"
    >
      <div class="flex flex-col">
        <div class="mb-4 flex">
          <div class="flex items-start gap-1">
            <div class="flex min-h-28 justify-between gap-1 rounded-md py-2 pe-3">
              {#if docA}
                <div>
                  <Button
                    on:click={() => appStore.setDiffDocA_ID(undefined)}
                    color="light"
                    class="border-0 p-1"
                  >
                    <i class="bx bx-x text-lg"></i>
                  </Button>
                </div>
                <div class="flex flex-col">
                  <div class="mb-1 max-w-96" title={docA.document.title}>{docA.document.title}</div>
                  <div class="text-gray-600">{getPublisher(docA.document.publisher.name)}</div>
                  <div class="text-gray-600">Version: {docA.document.tracking.version}</div>
                </div>
              {:else}
                <div class="flex flex-col gap-2">
                  <P italic>Select a document or upload local ones.</P>
                </div>
              {/if}
            </div>
          </div>
          <div class="flex items-center gap-1">
            <div class="flex min-h-28 justify-between gap-1 rounded-md px-3 py-2">
              {#if docB}
                <div>
                  <Button
                    on:click={() => appStore.setDiffDocB_ID(undefined)}
                    color="light"
                    class="border-0 p-1"
                  >
                    <i class="bx bx-x text-lg"></i>
                  </Button>
                </div>
                <div class="flex flex-col">
                  <div class="mb-1 max-w-96" title={docB.document.title}>{docB.document.title}</div>
                  <div class="text-gray-600">{getPublisher(docB.document.publisher.name)}</div>
                  <div class="text-gray-600">Version: {docB.document.tracking.version}</div>
                </div>
              {:else}
                <div
                  class:flex={true}
                  class:flex-col={true}
                  class:gap-2={true}
                  class:invisible={!docA}
                >
                  <P italic>Select a document or upload local ones.</P>
                </div>
              {/if}
            </div>
          </div>
          <div class="flex h-full items-center">
            <Button
              on:click={() => push("/diff")}
              disabled={!docA || !docB}
              size="sm"
              class="flex gap-x-2"
            >
              <Img src="plus-minus.svg" class="w-5 invert" />
              <span>Compare</span>
            </Button>
          </div>
        </div>
        <div class="flex flex-col">
          {#if tempDocuments?.length > 0}
            <span class="mb-1">Temporary documents:</span>
            <Table>
              <TableHead>
                <TableHeadCell padding="px-6 pt-2">Tracking ID</TableHeadCell>
                <TableHeadCell padding="px-6 pt-2">Publisher</TableHeadCell>
                <TableHeadCell padding="px-6 pt-2">Title</TableHeadCell>
                <TableHeadCell padding="px-6 pt-2">Expires</TableHeadCell>
                <TableHeadCell padding="px-6 pt-2">File name</TableHeadCell>
                <TableHeadCell padding="px-6 pt-2"></TableHeadCell>
              </TableHead>
              <TableBody>
                {#each tempDocuments as document}
                  {@const doc = document.document}
                  {@const tempDocID = `tempdocument${document.file.id}`}
                  <TableBodyRow>
                    <TableBodyCell {tdClass}>{doc.tracking.id}</TableBodyCell>
                    <TableBodyCell {tdClass}>{doc.publisher.name}</TableBodyCell>
                    <TableBodyCell {tdClass}>{doc.title.substring(0, 40)}</TableBodyCell>
                    <TableBodyCell {tdClass}
                      >{new Date(document.file.expired).toLocaleString()}</TableBodyCell
                    >
                    <TableBodyCell {tdClass}>{document.file.filename}</TableBodyCell>
                    <TableBodyCell {tdClass}>
                      <div class="flex items-center">
                        <Button
                          on:click={(e) => {
                            e.preventDefault();
                            deleteTempDocument(document.file.id);
                          }}
                          class="border-0 p-2"
                          color="light"
                          title={`delete ${doc.title}`}
                        >
                          <i class="bx bx-trash text-lg text-red-500"></i>
                        </Button>
                        <button
                          on:click|stopPropagation={(e) => {
                            $appStore.app.diff.docA_ID
                              ? appStore.setDiffDocB_ID(tempDocID)
                              : appStore.setDiffDocA_ID(tempDocID);
                            e.preventDefault();
                          }}
                          class:invisible={!$appStore.app.diff.isDiffModeEnabled}
                          disabled={$appStore.app.diff.docA_ID === tempDocID ||
                            $appStore.app.diff.docB_ID === tempDocID ||
                            disableDiffButtons}
                          title={`compare ${doc.tracking_id}`}
                        >
                          <Img
                            src="plus-minus.svg"
                            class={`${
                              $appStore.app.diff.docA_ID === tempDocID ||
                              $appStore.app.diff.docB_ID === tempDocID ||
                              disableDiffButtons
                                ? "invert-[70%]"
                                : ""
                            } min-h-4 min-w-4`}
                          />
                        </button>
                      </div>
                    </TableBodyCell>
                  </TableBodyRow>
                {/each}
                <TableBodyRow></TableBodyRow>
              </TableBody>
            </Table>
          {/if}
          {#if freeTempDocuments && tempDocuments.length <= freeTempDocuments}
            <Dropzone
              on:drop={dropHandle}
              on:dragover={(event) => {
                event.preventDefault();
              }}
              on:change={handleChange}
              multiple
              class="ms-1 h-16"
            >
              <i class="bx bx-upload text-xl text-gray-500"></i>
              <p class="mb-2 text-sm text-gray-500 dark:text-gray-400">
                Upload a temporary document
              </p>
            </Dropzone>
          {/if}
          <ErrorMessage message={tempDocErrorMessage}></ErrorMessage>
        </div>
      </div>
    </div>
  {/if}
</div>
