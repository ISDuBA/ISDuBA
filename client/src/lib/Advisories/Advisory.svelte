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
  import Comment from "$lib/Advisories/Comments/Comment.svelte";
  import Version from "$lib/Advisories/Version.svelte";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/Advisories/SSVC/SSVCCalculator.svelte";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";
  import JsonDiff from "$lib/Diff/JsonDiff.svelte";
  import {
    ASSESSING,
    ARCHIVED,
    DELETE,
    NEW,
    READ,
    REVIEW,
    canSetStateRead,
    allowedToChangeWorkflow
  } from "$lib/permissions";
  import CommentTextArea from "./Comments/CommentTextArea.svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import Event from "$lib/Advisories/Event.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  export let params: any = null;

  let document: any = {};
  let ssvc: any;
  $: ssvcStyle = ssvc
    ? `color: ${ssvc.color}; border: 1pt solid ${ssvc.color}; background-color: white;`
    : "";
  let comment: string = "";
  let comments: any = [];
  let events: any = [];
  let loadCommentsError = "";
  let loadEventsError = "";
  let loadAdvisoryVersionsError = "";
  let loadDocumentError = "";
  let createCommentError = "";
  let loadDocumentSSVCError = "";
  let stateError = "";
  let advisoryVersions: any[] = [];
  let advisoryVersionByDocumentID: any;
  let advisoryState = "";
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
      loadAdvisoryVersionsError = `Could not load versions. ${getErrorMessage(response.error)}`;
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
      loadDocumentError = `Could not load document. ${getErrorMessage(response.error)}`;
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
      loadDocumentSSVCError = `Could not load SSVC. ${getErrorMessage(response.error)}`;
    }
  };

  const loadEvents = async () => {
    let loadedEvents: any = [];
    if (advisoryVersions.length > 0) {
      const promises = await Promise.allSettled(
        advisoryVersions.map(async (v) => {
          return request(`/api/events/${v.id}`, "GET");
        })
      );
      const result = promises
        .filter((p: any) => p.status === "fulfilled")
        .map((p: any) => {
          return p.value;
        });
      if (promises.length != result.length) {
        loadEventsError = `Could not load all events. An error occured on the server. Please contact an administrator.`;
      }
      result.forEach((e) => {
        if (e.content !== "undefined") {
          loadedEvents = loadedEvents.concat(e.content);
        } else {
          loadEventsError = `Could not load all events. An error occured on the server. Please contact an administrator.`;
        }
      });
      events = loadedEvents;
    } else {
      loadEventsError = `Could not load events. An error occured on the server. Please contact an administrator.`;
    }
  };

  const loadComments = async () => {
    let loadedComments: any = [];
    if (advisoryVersions.length > 0) {
      const promises = await Promise.allSettled(
        advisoryVersions.map(async (v) => {
          return request(`/api/comments/${v.id}`, "GET");
        })
      );
      const result = promises
        .filter((p: any) => p.status === "fulfilled")
        .map((p: any) => {
          return p.value;
        });
      if (promises.length != result.length) {
        loadCommentsError = `Could not load all comments. An error occured on the server. Please contact an administrator.`;
      }
      result.forEach((c) => {
        if (c.content !== "undefined") {
          let comments = c.content;
          for (let i = 0; i < comments.length; i++) {
            comments[i].documentVersion = advisoryVersionByDocumentID[comments[i].document_id];
          }
          loadedComments = loadedComments.concat(comments);
        } else {
          loadCommentsError = `Could not load all comments. An error occured on the server. Please contact an administrator.`;
        }
      });
      comments = loadedComments;
    } else {
      loadCommentsError = `Could not load comments. An error occured on the server. Please contact an administrator.`;
    }
  };

  async function createComment() {
    const formData = new FormData();
    formData.append("message", comment);
    const response = await request(`/api/comments/${params.id}`, "POST", formData);
    if (response.ok) {
      comment = "";
      await loadComments();
      await loadAdvisoryState();
      await loadEvents();
    } else if (response.error) {
      createCommentError = `Could not create comment. ${getErrorMessage(response.error)}`;
    }
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
      await loadEvents();
    } else if (response.error) {
      stateError = `Could not change state. ${getErrorMessage(response.error)}`;
    }
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
      stateError = `Couldn't load state. ${getErrorMessage(response.error)}`;
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
      const id: any = setTimeout(async () => {
        if (canSetStateRead(advisoryState)) {
          await updateState(READ);
        }
      }, 20000);
      timeoutIDs.push(id);
    }
  };

  async function loadMetaData() {
    await loadAdvisoryState();
    await loadDocumentSSVC();
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
      return "dark";
    } else {
      return "none";
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
      <div class="flex flex-row flex-wrap items-end justify-start gap-y-2 md:justify-between">
        <Label class="text-gray-600">{params.publisherNamespace}</Label>
        <div class="flex h-fit flex-row gap-2">
          {#if advisoryState}
            <a
              href={"javascript:void(0);"}
              class="inline-flex"
              on:click={() => updateStateIfAllowed(NEW)}
            >
              <Badge title="Mark as new" class="w-fit" color={getBadgeColor(NEW, advisoryState)}
                >{NEW}</Badge
              >
            </a>
            <a
              href={"javascript:void(0);"}
              class="inline-flex"
              on:click={() => updateStateIfAllowed(READ)}
            >
              <Badge title="Mark as read" class="w-fit" color={getBadgeColor(READ, advisoryState)}
                >{READ}</Badge
              >
            </a>
            <a
              href={"javascript:void(0);"}
              class="inline-flex"
              on:click={() => updateStateIfAllowed(ASSESSING)}
            >
              <Badge
                title="Mark as assesing"
                class="w-fit"
                color={getBadgeColor(ASSESSING, advisoryState)}>{ASSESSING}</Badge
              >
            </a>
            <a
              href={"javascript:void(0);"}
              class="inline-flex"
              on:click={() => updateStateIfAllowed(REVIEW)}
            >
              <Badge
                title="Release for review"
                class="w-fit"
                color={getBadgeColor(REVIEW, advisoryState)}>{REVIEW}</Badge
              >
            </a>
            <a
              href={"javascript:void(0);"}
              class="inline-flex"
              on:click={() => updateStateIfAllowed(ARCHIVED)}
            >
              <Badge title="Archive" class="w-fit" color={getBadgeColor(ARCHIVED, advisoryState)}
                >{ARCHIVED}</Badge
              >
            </a>
            <a
              href={"javascript:void(0);"}
              class="inline-flex"
              on:click={() => updateStateIfAllowed(DELETE)}
            >
              <Badge
                title="Mark for deletion"
                on:click={() => updateState(DELETE)}
                class="w-fit"
                color={getBadgeColor(DELETE, advisoryState)}>{DELETE}</Badge
              >
            </a>
          {/if}
          {#if ssvc}
            <Badge style={ssvcStyle}>{ssvc.label}</Badge>
            <Tooltip>SSVC</Tooltip>
          {/if}
        </div>
      </div>
      <hr class="mb-4 mt-2" />
    </div>
    <ErrorMessage message={loadAdvisoryVersionsError}></ErrorMessage>
    <ErrorMessage message={loadDocumentSSVCError}></ErrorMessage>
    <ErrorMessage message={stateError}></ErrorMessage>
    {#if advisoryVersions.length > 0}
      <Version
        publisherNamespace={params.publisherNamespace}
        trackingID={params.trackingID}
        {advisoryVersions}
        selectedDocumentVersion={document.tracking?.version}
        on:selectedDiffDocuments={onSelectedDiffDocuments}
        on:disableDiff={() => (isDiffOpen = false)}
      ></Version>
    {/if}
    {#if isDiffOpen}
      <JsonDiff title={undefined} {diffDocuments}></JsonDiff>
    {:else}
      <Webview></Webview>
    {/if}
    <ErrorMessage message={loadDocumentError}></ErrorMessage>
  </div>
  {#if appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()}
    <div class="mr-3 w-full min-w-96 max-w-[96%] xl:w-[50%] xl:max-w-[46%] 2xl:max-w-[33%]">
      <Accordion>
        <AccordionItem open>
          <span slot="header"
            ><i class="bx bx-comment-detail"></i><span class="ml-2">Comments</span></span
          >
          {#if loadCommentsError === ""}
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
          {#if loadCommentsError === ""}
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
          {/if}
          <ErrorMessage message={loadEventsError}></ErrorMessage>
        </AccordionItem>
      </Accordion>
      <Accordion class="mt-3">
        <AccordionItem open>
          <span slot="header"><i class="bx bx-calculator"></i><span class="ml-2">SSVC</span></span>
          <ErrorMessage message={loadDocumentSSVCError}></ErrorMessage>
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
