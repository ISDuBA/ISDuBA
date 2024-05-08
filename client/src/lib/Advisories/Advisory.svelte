<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { Label, Timeline, AccordionItem, Accordion, Badge, Tooltip } from "flowbite-svelte";
  import { onDestroy } from "svelte";
  import { appStore } from "$lib/store";
  import Comment from "$lib/Advisories/Comment.svelte";
  import Version from "$lib/Advisories/Version.svelte";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/Advisories/SSVC/SSVCCalculator.svelte";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";
  import JsonDiff from "$lib/Diff/JsonDiff.svelte";
  import {
    ASSESSING,
    ARCHIVED,
    DELETED,
    NEW,
    READ,
    REVIEW,
    canSetStateRead,
    allowedToChangeWorkflow
  } from "$lib/permissions";
  import CommentTextArea from "./CommentTextArea.svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import Event from "$lib/Advisories/Event.svelte";
  export let params: any = null;

  let document: any = {};
  let ssvc: any;
  $: ssvcStyle = ssvc
    ? `color: ${ssvc.color}; border: 1pt solid ${ssvc.color}; background-color: white;`
    : "";
  let comment: string = "";
  let comments: any = [];
  let events: any = [];
  let loadCommentsError: string;
  let loadEventsError: string;
  let createCommentError: string;
  let advisoryVersions: any[] = [];
  let advisoryVersionByDocumentID: any;
  let advisoryState: string;
  let isCommentingAllowed: boolean;
  $: if ([READ, ASSESSING].includes(advisoryState)) {
    isCommentingAllowed = appStore.isEditor() || appStore.isReviewer();
  } else {
    isCommentingAllowed = false;
  }
  let isCalculatingAllowed: boolean;
  $: if ([READ, ASSESSING].includes(advisoryState)) {
    isCalculatingAllowed = appStore.isEditor() || appStore.isReviewer();
  } else {
    isCalculatingAllowed = false;
  }

  const timeoutIDs: number[] = [];
  let diffDocuments: any;
  let isDiffOpen = false;

  const loadAdvisoryVersions = async () => {
    const response = await request(
      `/api/documents?&columns=id version tracking_id&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
      "GET"
    );
    if (response.ok) {
      const result = await response.content;
      advisoryVersions = result.documents.map((doc: any) => {
        return { id: doc.id, version: doc.version, tracking_id: doc.tracking_id };
      });
      advisoryVersionByDocumentID = advisoryVersions.reduce((acc: any, version: any) => {
        acc[version.id] = version.version;
        return acc;
      }, {});
    } else if (response.error) {
      appStore.displayErrorMessage(response.error);
    }
  };

  const loadDocument = async () => {
    const response = await request(`/api/documents/${params.id}`, "GET");
    if (response.ok) {
      const result = await response.content;
      ({ document } = result);
      const docModel = convertToDocModel(result);
      appStore.setDocument(docModel);
    } else if (response.error) {
      appStore.displayErrorMessage(response.error);
    }
  };

  const loadDocumentSSVC = async () => {
    const response = await request(
      `/api/documents?columns=ssvc&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
      "GET"
    );
    if (response.ok) {
      const result = await response.content;
      if (result.documents[0].ssvc) {
        ssvc = convertVectorToLabel(result.documents[0].ssvc);
      }
    } else if (response.error) {
      appStore.displayErrorMessage(response.error);
    }
  };

  function loadEvents(): Promise<any[]> {
    return new Promise((resolve) => {
      const newEvents: any = [];
      loadEventsError = "";
      advisoryVersions.forEach(async (advVer: any) => {
        const response = await request(`/api/events/${advVer.id}`, "GET");
        if (response.ok) {
          const loadedEvents = response.content;
          if (loadedEvents) {
            loadedEvents.forEach((c: any) => {
              c.documentVersion = advisoryVersionByDocumentID[c.document_id];
            });
            newEvents.push(...loadedEvents);
          }
          events = newEvents;
          resolve(newEvents);
        } else if (response.error) {
          loadEventsError = response.error;
        }
      });
    });
  }

  function loadComments(): Promise<any[]> {
    return new Promise((resolve) => {
      const newComments: any = [];
      loadCommentsError = "";
      advisoryVersions.forEach(async (advVer: any) => {
        const response = await request(`/api/comments/${advVer.id}`, "GET");
        if (response.ok) {
          const loadedComments = response.content;
          if (loadedComments) {
            loadedComments.forEach((c: any) => {
              c.documentVersion = advisoryVersionByDocumentID[c.document_id];
            });
            newComments.push(...loadedComments);
          }
          comments = newComments;
          resolve(newComments);
        } else if (response.error) {
          loadCommentsError = response.error;
        }
      });
    });
  }
  async function createComment() {
    const formData = new FormData();
    formData.append("message", comment);
    const response = await request(`/api/comments/${params.id}`, "POST", formData);
    if (response.ok) {
      comment = "";
      loadComments().then((newComments: any[]) => {
        if (newComments.length === 1) {
          loadAdvisoryState();
        }
      });
      appStore.displaySuccessMessage("Comment for advisory saved.");
    } else if (response.error) {
      createCommentError = response.error;
    }
    await loadEvents();
  }

  async function sendForReview() {
    if (comment.length !== 0) {
      await createComment();
    }
    await updateState(REVIEW);
  }

  async function updateState(newState: string) {
    const response = await request(
      `/api/status/${params.publisherNamespace}/${params.trackingID}/${newState}`,
      "PUT"
    );
    if (response.ok) {
      advisoryState = newState;
    } else if (response.error) {
      appStore.displayErrorMessage(response.error);
    }
    await loadEvents();
  }

  const loadAdvisoryState = async () => {
    const response = await request(
      `/api/documents?advisories=true&columns=state&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
      "GET"
    );
    if (response.ok) {
      const result = response.content;
      advisoryState = result.documents[0].state;
      return result.documents[0].state;
    } else if (response.error) {
      appStore.displayErrorMessage(response.error);
    }
  };

  const loadData = async () => {
    await loadDocumentSSVC();
    await loadDocument();
    await loadAdvisoryVersions();
    if (appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()) {
      await loadComments();
      await loadEvents();
    }
    const state = await loadAdvisoryState();
    // Only set state to 'read' if editor opens the current version.
    if (
      state === "new" &&
      appStore.isEditor() &&
      (advisoryVersions.length === 1 || advisoryVersions[0].version === document.tracking?.version)
    ) {
      const id = setTimeout(async () => {
        if (canSetStateRead(advisoryState)) {
          await updateState(READ);
          appStore.displayInfoMessage("This advisory is marked as read");
        }
      }, 3000);
      timeoutIDs.push(id);
    }
  };

  function loadMetaData() {
    loadAdvisoryState();
    loadDocumentSSVC();
  }

  const onSelectedDiffDocuments = async (event: any) => {
    diffDocuments = {
      docA: event.detail.docA,
      docB: event.detail.docB
    };
    isDiffOpen = true;
  };

  const updateStateIfAllowed = async (state: string) => {
    if (allowedToChangeWorkflow(appStore.getRoles(), advisoryState, state)) {
      await updateState(state);
    }
  };

  const getBadgeColor = (state: string, currentState: string) => {
    if (state === currentState) {
      return "green";
    } else if (allowedToChangeWorkflow(appStore.getRoles(), currentState, state)) {
      return "primary";
    } else {
      return "dark";
    }
  };

  onDestroy(() => {
    timeoutIDs.forEach((id: number) => {
      clearTimeout(id);
    });
  });

  $: if (params) {
    loadData();
  }
</script>

<svelte:head>
  <title>{params.trackingID}</title>
</svelte:head>

<div
  class="flex h-screen max-h-full flex-wrap justify-between gap-x-4 gap-y-8 overflow-y-scroll xl:flex-nowrap"
>
  <div class="flex max-h-full w-full grow flex-col gap-y-2 overflow-y-scroll px-2">
    <div class="flex flex-col">
      <div class="flex gap-2">
        <Label class="text-lg">{params.trackingID}</Label>
      </div>
      <Label class="mb-2 text-gray-600">{params.publisherNamespace}</Label>
      <div class="flex gap-2">
        {#if advisoryState}
          <a
            href={"javascript:void(0);"}
            class="inline-flex"
            on:click={() => updateStateIfAllowed(NEW)}
          >
            <Badge class="w-fit" color={getBadgeColor(NEW, advisoryState)}>{NEW}</Badge>
          </a>
          <Tooltip>Mark as new</Tooltip>
          <a
            href={"javascript:void(0);"}
            class="inline-flex"
            on:click={() => updateStateIfAllowed(READ)}
          >
            <Badge class="w-fit" color={getBadgeColor(READ, advisoryState)}>{READ}</Badge>
          </a>
          <Tooltip>Mark as read</Tooltip>
          <a
            href={"javascript:void(0);"}
            class="inline-flex"
            on:click={() => updateStateIfAllowed(ASSESSING)}
          >
            <Badge class="w-fit" color={getBadgeColor(ASSESSING, advisoryState)}>{ASSESSING}</Badge>
          </a>
          <Tooltip>Back to assessing</Tooltip>
          <a
            href={"javascript:void(0);"}
            class="inline-flex"
            on:click={() => updateStateIfAllowed(REVIEW)}
          >
            <Badge class="w-fit" color={getBadgeColor(REVIEW, advisoryState)}>{REVIEW}</Badge>
          </a>
          <Tooltip>Release for review</Tooltip>
          <a
            href={"javascript:void(0);"}
            class="inline-flex"
            on:click={() => updateStateIfAllowed(ARCHIVED)}
          >
            <Badge class="w-fit" color={getBadgeColor(ARCHIVED, advisoryState)}>{ARCHIVED}</Badge>
          </a>
          <Tooltip>Archive</Tooltip>
          <a
            href={"javascript:void(0);"}
            class="inline-flex"
            on:click={() => updateStateIfAllowed(DELETED)}
          >
            <Badge
              on:click={() => updateState(DELETED)}
              class="w-fit"
              color={getBadgeColor(DELETED, advisoryState)}>{DELETED}</Badge
            >
          </a>
          <Tooltip>Mark for deletion</Tooltip>
        {/if}
        {#if ssvc}
          <Badge style={ssvcStyle}>{ssvc.label}</Badge>
          <Tooltip>SSVC</Tooltip>
        {/if}
      </div>
    </div>
    <Version
      publisherNamespace={params.publisherNamespace}
      trackingID={params.trackingID}
      {advisoryVersions}
      selectedDocumentVersion={document.tracking?.version}
      on:selectedDiffDocuments={onSelectedDiffDocuments}
      on:disableDiff={() => (isDiffOpen = false)}
    ></Version>
    {#if isDiffOpen}
      <JsonDiff title={undefined} {diffDocuments}></JsonDiff>
    {:else}
      <Webview></Webview>
    {/if}
  </div>
  {#if appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()}
    <div class="mr-3 w-full min-w-96 max-w-[96%] xl:w-[50%] xl:max-w-[46%] 2xl:max-w-[33%]">
      <Accordion>
        <AccordionItem open>
          <span slot="header"
            ><i class="bx bx-comment-detail"></i><span class="ml-2">Comments</span></span
          >
          {#if comments?.length > 0}
            <div class="max-h-96 overflow-y-auto pl-2">
              <Timeline class="mb-4 flex flex-col-reverse">
                {#each comments as comment (comment.id)}
                  <Comment on:commentUpdate={loadEvents} {comment}></Comment>
                {/each}
              </Timeline>
            </div>
          {:else}
            <div class="mb-6 text-gray-600">No comments available.</div>
          {/if}
          <ErrorMessage message={loadCommentsError}></ErrorMessage>
          {#if isCommentingAllowed}
            <div class="mt-6">
              <Label class="mb-2" for="comment-textarea">New Comment:</Label>
              <CommentTextArea
                on:input={() => (createCommentError = "")}
                on:saveComment={createComment}
                on:saveForReview={sendForReview}
                bind:value={comment}
                errorMessage={createCommentError}
                buttonText="Send"
                state={advisoryState}
              ></CommentTextArea>
            </div>
          {/if}
        </AccordionItem>
      </Accordion>
      <Accordion class="mt-3">
        <AccordionItem open>
          <span slot="header"
            ><i class="bx bx-calendar-event"></i><span class="ml-2">Events</span></span
          >
          {#if events?.length > 0}
            <div class="max-h-96 overflow-y-auto pl-2">
              <Timeline class="mb-4 flex flex-col-reverse">
                {#each events as event}
                  <Event {event}></Event>
                {/each}
              </Timeline>
            </div>
          {:else}
            <div class="mb-6 text-gray-600">No events available.</div>
          {/if}
          <ErrorMessage message={loadEventsError}></ErrorMessage>
        </AccordionItem>
      </Accordion>
      <Accordion class="mt-3">
        <AccordionItem open>
          <span slot="header"><i class="bx bx-calculator"></i><span class="ml-2">SSVC</span></span>
          <SsvcCalculator
            vectorInput={ssvc?.vector}
            disabled={!isCalculatingAllowed}
            documentID={params.id}
            on:updateSSVC={loadMetaData}
          ></SsvcCalculator>
        </AccordionItem>
      </Accordion>
    </div>
  {/if}
</div>
