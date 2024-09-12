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
  import { request } from "$lib/request";
  import { getPublisher } from "$lib/publisher";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { onDestroy, onMount } from "svelte";
  import CIconButton from "$lib/Components/CIconButton.svelte";

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

  $: docA_ID = $appStore.app.diff.docA_ID;
  $: docB_ID = $appStore.app.diff.docB_ID;
  $: if (docA_ID || docB_ID) getDocuments();
  $: isDiffBoxOpen = $appStore.app.diff.isDiffBoxOpen;
  $: if (isDiffBoxOpen) {
    getTempDocuments();
    getDocuments();
  }
  $: disableDiffButtons =
    $appStore.app.diff.docA_ID !== undefined && $appStore.app.diff.docB_ID !== undefined;

  let freeTempDocuments = 0;
  let tempDocuments: any[] = [];
  let docA: any;
  let docB: any;
  let loadDocumentsErrorMessage: ErrorDetails | null;
  let tempDocErrorMessage: ErrorDetails | null;
  let loadTempDocsErrorMessage: ErrorDetails | null;
  let intervalID: ReturnType<typeof setTimeout> | undefined = undefined;

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
    loadDocumentsErrorMessage = null;
    if ($appStore.app.diff.docA_ID) {
      const responseDocA = await getDocument("A");
      if (responseDocA.ok) {
        if ($appStore.app.diff.docA_ID) docA = await responseDocA.content;
      } else if (responseDocA.error) {
        if (responseDocA.error === "404") {
          appStore.setDiffDocA_ID(undefined);
        } else {
          loadDocumentsErrorMessage = getErrorDetails(
            `Could not load first document.`,
            responseDocA
          );
        }
      }
    } else {
      docA = undefined;
    }
    if ($appStore.app.diff.docB_ID) {
      const responseDocB = await getDocument("B");
      if (responseDocB.ok) {
        if ($appStore.app.diff.docB_ID) docB = await responseDocB.content;
      } else if (responseDocB.error) {
        if (responseDocB.error === "404") {
          appStore.setDiffDocB_ID(undefined);
        } else {
          loadDocumentsErrorMessage = getErrorDetails(
            `Could not load second document.`,
            responseDocB
          );
        }
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
      updateExpired();
    } else if (response.error) {
      loadTempDocsErrorMessage = getErrorDetails(`Could not load temporary document.`, response);
    }
  };

  const deleteTempDocument = async (id: number) => {
    tempDocErrorMessage = null;
    const response = await request(`/api/tempdocuments/${id}`, "DELETE");
    if (response.ok) {
      if ($appStore.app.diff.docA_ID === `tempdocument${id}`) {
        appStore.setDiffDocA_ID(undefined);
      } else if ($appStore.app.diff.docB_ID === `tempdocument${id}`) {
        appStore.setDiffDocB_ID(undefined);
      }
      getTempDocuments();
    } else if (response.error) {
      tempDocErrorMessage = getErrorDetails(`Could not delete temporary document.`, response);
    }
  };

  const deleteExpiredDocFromStore = (id: number) => {
    if ($appStore.app.diff.docA_ID === `tempdocument${id}`) {
      appStore.setDiffDocA_ID(undefined);
    } else if ($appStore.app.diff.docB_ID === `tempdocument${id}`) {
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

<div
  class="fixed bottom-0 left-1 right-1 flex flex-col items-end justify-center sm:left-auto md:right-20"
>
  <Button
    on:click={appStore.toggleDiffBox}
    class="rounded-none rounded-t-md border-b-0"
    color="light"
  >
    <span class="me-2"
      >Diff {docA
        ? `${docA?.document?.title.substring(0, 25)}${docA?.document?.title.length > 25 ? "..." : ""}`
        : ""}
      {docB
        ? ` - ${docB?.document?.title.substring(0, 25)}${docB?.document?.title.length > 25 ? "..." : ""}`
        : ""}</span
    >
    <Img src="plus-minus.svg" class="h-4 min-h-2 min-w-2" />
  </Button>
  {#if $appStore.app.diff.isDiffBoxOpen}
    <div
      class="flex items-stretch gap-6 rounded-tl-md border border-solid border-gray-300 bg-white p-4 shadow-gray-800"
    >
      <div class="flex flex-col">
        <div class="mb-4 flex justify-between">
          <div class="flex items-start gap-1">
            <div class="flex min-h-28 justify-between gap-1 rounded-md pb-2 pe-3">
              {#if docA}
                <div>
                  <Button
                    on:click={() => {
                      appStore.setDiffDocA_ID(undefined);
                      docA = undefined;
                    }}
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
          <div class="flex gap-1">
            <div class="flex min-h-28 justify-between gap-1 rounded-md px-3 pb-2">
              {#if docB}
                <div>
                  <Button
                    on:click={() => {
                      appStore.setDiffDocB_ID(undefined);
                      docB = undefined;
                    }}
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
          <div class="flex h-full items-start">
            <Button
              on:click={() => push("/diff")}
              disabled={!docA ||
                !docB ||
                !$appStore.app.diff.docA_ID ||
                !$appStore.app.diff.docB_ID}
              size="sm"
              class="flex gap-x-2"
            >
              <Img src="plus-minus.svg" class="w-5 invert" />
              <span>Compare</span>
            </Button>
          </div>
        </div>
        <ErrorMessage error={loadDocumentsErrorMessage}></ErrorMessage>
        <div class="flex flex-col">
          {#if tempDocuments?.length > 0}
            <span class="mb-1">Temporary documents:</span>
            <Table>
              <TableHead>
                <TableHeadCell {padding}>Tracking ID</TableHeadCell>
                <TableHeadCell {padding}>Publisher</TableHeadCell>
                <TableHeadCell {padding}>Title</TableHeadCell>
                <TableHeadCell {padding}>Expires in</TableHeadCell>
                <TableHeadCell {padding}>File name</TableHeadCell>
                <TableHeadCell {padding}></TableHeadCell>
              </TableHead>
              <TableBody>
                {#each tempDocuments as document}
                  {@const doc = document.document}
                  {@const tempDocID = `tempdocument${document.file.id}`}
                  <TableBodyRow>
                    <TableBodyCell {tdClass}>{doc.tracking.id}</TableBodyCell>
                    <TableBodyCell {tdClass}>{doc.publisher.name}</TableBodyCell>
                    <TableBodyCell {tdClass}>
                      <span
                        class="block w-12 overflow-hidden text-ellipsis whitespace-nowrap sm:w-20 md:w-44 lg:w-60 xl:w-96 2xl:w-auto"
                        title={doc.title}>{doc.title}</span
                      >
                    </TableBodyCell>
                    <TableBodyCell {tdClass}>{doc.expired}</TableBodyCell>
                    <TableBodyCell {tdClass}>{document.file.filename}</TableBodyCell>
                    <TableBodyCell {tdClass}>
                      <div class="flex items-center">
                        <CIconButton
                          on:click={() => {
                            deleteTempDocument(document.file.id);
                          }}
                          color="red"
                          title={`delete ${doc.title}`}
                          icon="trash"
                        ></CIconButton>
                        <button
                          on:click|stopPropagation={(e) => {
                            if ($appStore.app.diff.docA_ID) {
                              appStore.setDiffDocB_ID(tempDocID);
                            } else {
                              appStore.setDiffDocA_ID(tempDocID);
                            }
                            e.preventDefault();
                          }}
                          class:invisible={!$appStore.app.diff.isDiffBoxOpen}
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
          {#if freeTempDocuments}
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
    </div>
  {/if}
</div>
