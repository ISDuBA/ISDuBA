<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { Button, Select, Label, Modal, Spinner } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Version from "$lib/Advisories/Version.svelte";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/Advisories/SSVC/SSVCCalculator.svelte";
  import Diff from "$lib/Diff/Diff.svelte";
  import { ARCHIVED, ASSESSING, DELETE, NEW, READ, REVIEW } from "$lib/workflow";
  import { canSetStateRead } from "$lib/permissions";
  import CommentTextArea from "./Comments/CommentTextArea.svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import WorkflowStates from "./WorkflowStates.svelte";
  import History from "./History.svelte/History.svelte";
  import Tlp from "./TLP.svelte";
  import SsvcBadge from "./SSVC/SSVCBadge.svelte";
  export let params: any = null;

  let document: any = {};
  let ssvcVector: string;
  let comment: string = "";
  let loadCommentsError: ErrorDetails | null = null;
  let loadEventsError: ErrorDetails | null = null;
  let loadAdvisoryVersionsError: ErrorDetails | null = null;
  let loadDocumentError: ErrorDetails | null = null;
  let loadFourCVEsError: ErrorDetails | null = null;
  let createCommentError: ErrorDetails | null = null;
  let loadDocumentSSVCError: ErrorDetails | null = null;
  let loadForwardTargetsError: ErrorDetails | null = null;
  let stateError: ErrorDetails | null = null;
  let advisoryVersions: any[] = [];
  let advisoryVersionByDocumentID: any;
  let advisoryState = "";
  let historyEntries: any = [];
  let isCommentingAllowed: boolean;
  let isSSVCediting = false;
  let position = "";
  let processRunning = false;
  let lastSuccessfulForwardTarget: number | undefined;

  $: if ([NEW, READ, ASSESSING].includes(advisoryState)) {
    isCommentingAllowed = appStore.isEditor();
  } else if ([REVIEW].includes(advisoryState)) {
    isCommentingAllowed = appStore.isEditor() || appStore.isReviewer();
  } else if ([ARCHIVED].includes(advisoryState)) {
    isCommentingAllowed = appStore.isEditor() || appStore.isAdmin();
  } else if ([DELETE].includes(advisoryState)) {
    isCommentingAllowed = appStore.isAdmin();
  } else {
    isCommentingAllowed = false;
  }

  let isCalculatingAllowed: boolean;
  $: if ([NEW, READ, ASSESSING].includes(advisoryState)) {
    isCalculatingAllowed = appStore.isEditor();
  } else {
    isCalculatingAllowed = false;
  }
  $: canSeeCommentArea =
    appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor() || appStore.isAdmin();

  const setAsReadTimeout: number[] = [];
  let isDiffOpen = false;
  let commentFocus = false;

  let availableForwardSelection: any[] = [];
  let selectedForwardTarget: number | undefined;

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
      loadAdvisoryVersionsError = getErrorDetails(`Could not load versions.`, response);
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
      loadDocumentError = getErrorDetails(`Could not load document.`, response);
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
        ssvcVector = result.documents[0].ssvc;
      }
    } else if (response.error) {
      loadDocumentSSVCError = getErrorDetails(`Could not load SSVC.`, response);
    }
  };

  const loadEvents = async () => {
    const response = await request(
      `/api/events/${params.publisherNamespace}/${params.trackingID}`,
      "GET"
    );
    if (response.ok) {
      return await response.content;
    } else if (response.error) {
      loadEventsError = getErrorDetails(`Could not load events.`, response);
      return [];
    }
  };

  const loadComments = async () => {
    const response = await request(
      `/api/comments/${params.publisherNamespace}/${params.trackingID}`,
      "GET"
    );
    if (response.ok) {
      let comments = await response.content;
      for (let i = 0; i < comments.length; i++) {
        comments[i].documentVersion = advisoryVersionByDocumentID[comments[i].document_id];
      }
      return comments;
    } else if (response.error) {
      loadEventsError = getErrorDetails(`Could not comments.`, response);
      return [];
    }
  };

  const buildHistory = async () => {
    if (!canSeeCommentArea) {
      historyEntries = [];
      return;
    }
    const comments = await loadComments();
    let events = await loadEvents();
    const commentsByTime = comments.reduce((o: any, n: any) => {
      o[`${n.time}:${n.commentator}`] = {
        message: n.message,
        id: n.id,
        documentVersion: n.documentVersion
      };
      return o;
    }, {});
    const commentsEdited = events
      .filter((e: any) => {
        return e.event_type === "change_comment";
      })
      .map((e: any) => {
        return {
          id: e.comment_id,
          time: e.time
        };
      })
      .reduce((o: any, n: any) => {
        if (!o[n.id]) o[n.id] = [];
        o[n.id].push(n.time);
        return o;
      }, {});
    events.map((e: any) => {
      if (e.event_type === "add_comment") {
        const comment = commentsByTime[`${e.time}:${e.actor}`];
        e["message"] = comment.message;
        e["comment_id"] = comment.id;
        e["documentVersion"] = comment.documentVersion;
        if (commentsEdited[comment.id]) {
          e["times"] = commentsEdited[comment.id];
        }
      }
      return e;
    });
    historyEntries = events;
  };

  async function createComment() {
    await allowEditing();
    const formData = new FormData();
    formData.append("message", comment);
    const response = await request(`/api/comments/${params.id}`, "POST", formData);
    if (response.ok) {
      comment = "";
      await loadAdvisoryState();
      await buildHistory();
    } else if (response.error) {
      createCommentError = getErrorDetails(`Could not create comment.`, response);
    }
  }

  async function sendForReview() {
    if (comment.length !== 0) {
      await createComment();
    }
    await updateState(REVIEW);
  }

  async function sendForAssessing() {
    if (comment.length !== 0) {
      await createComment();
    }
    await updateState(ASSESSING);
  }

  async function updateState(newState: string) {
    // Cancel automatic state transitions
    setAsReadTimeout.forEach((id: number) => {
      clearTimeout(id);
    });

    const response = await request(
      `/api/status/${params.publisherNamespace}/${params.trackingID}/${newState}`,
      "PUT"
    );
    if (response.ok) {
      advisoryState = newState;
      await buildHistory();
    } else if (response.error) {
      stateError = getErrorDetails(`Could not change state.`, response);
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
      stateError = getErrorDetails(`Couldn't load state.`, response);
    }
  };

  const loadFourCVEs = async () => {
    const response = await request(
      `/api/documents?advisories=false&columns=four_cves&query=$id ${params.id} integer =`,
      "GET"
    );
    if (response.ok) {
      const content = await response.content;
      let four_cves = content?.documents[0]?.four_cves;
      appStore.setFourCVEs(four_cves);
    } else if (response.error) {
      loadFourCVEsError = getErrorDetails(`Couldn't load CVEs.`, response);
    }
  };

  const loadData = async () => {
    await loadFourCVEs();
    await loadDocumentSSVC();
    await loadDocument();
    await loadAdvisoryVersions();
    await buildHistory();
    if (appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()) {
      await buildHistory();
    }
    await loadAdvisoryState();
    // Only set state to 'read' if editor opens the current version.
    if (
      advisoryState === NEW &&
      canSetStateRead(advisoryState) &&
      (advisoryVersions.length === 1 || advisoryVersions[0].version === document.tracking?.version)
    ) {
      const id: any = setTimeout(async () => {
        if (advisoryState === "new" && canSetStateRead(advisoryState)) {
          await updateState(READ);
        }
      }, 20000);
      setAsReadTimeout.push(id);
    }
  };

  async function loadMetaData() {
    await loadAdvisoryState();
    await loadDocumentSSVC();
    await buildHistory();
  }

  async function allowEditing() {
    if (advisoryState === NEW && canSetStateRead(advisoryState)) {
      await updateState(READ);
    }
  }

  const fetchForwardTargets = async () => {
    const response = await request(`/api/documents/forward`, "GET");
    if (response.ok) {
      availableForwardSelection = [];
      for (let target of response.content) {
        availableForwardSelection.push({ value: target.id, name: target.name });
      }
    } else if (response.error) {
      loadForwardTargetsError = getErrorDetails(`Couldn't load forward targets.`, response);
    }
  };

  const forwardDocument = async () => {
    processRunning = true;
    const response = await request(
      `/api/documents/forward/${params.id}/${selectedForwardTarget}`,
      "POST"
    );
    processRunning = false;
    if (response.error) {
      openForwardModal = false;
      loadForwardTargetsError = getErrorDetails(`Could not forward document`, response);
    } else {
      lastSuccessfulForwardTarget = selectedForwardTarget;
    }
  };
  onDestroy(() => {
    setAsReadTimeout.forEach((id: number) => {
      clearTimeout(id);
    });
  });

  onMount(async () => {
    if (
      appStore.isAdmin() ||
      appStore.isEditor() ||
      appStore.isImporter() ||
      appStore.isReviewer() ||
      appStore.isSourceManager()
    ) {
      await fetchForwardTargets();
    }
  });

  $: if (params) {
    loadData();
    position = params.position;
    if (!params.position) {
      const topElement = window.document.getElementById("top");
      topElement?.scrollIntoView();
      appStore.setSelectedProduct("");
      appStore.setSelectedCVE("");
    }
  }
  let openForwardModal = false;
</script>

<svelte:head>
  <title>{params.trackingID}</title>
</svelte:head>

<Modal bind:open={openForwardModal}>
  <Label class="text-lg">Forward document</Label>
  <Select items={availableForwardSelection} bind:value={selectedForwardTarget}></Select>
  {#if typeof selectedForwardTarget === "number"}
    <Button disabled={processRunning} on:click={forwardDocument}>
      <span class="mr-2">Send document</span>
      {#if processRunning}
        <Spinner></Spinner>
      {:else if lastSuccessfulForwardTarget === selectedForwardTarget}
        <div class="inline-flex w-8 items-center"><i class="bx bx-check text-2xl" /></div>
      {:else}
        <div class="inline-flex w-8 items-center"><i class="bx bx-right-arrow-alt text-2xl" /></div>
      {/if}
    </Button>
  {/if}
</Modal>

<div class="grid h-full w-full grow grid-rows-[auto_minmax(100px,_1fr)] gap-y-2 px-2" id="top">
  <div class="flex w-full flex-none flex-col">
    <div class="flex gap-2">
      <Label class="text-lg">
        <span class="mr-2">{params.trackingID}</span>
        {#if $appStore.webview.doc?.tlp.label}
          <Tlp tlp={$appStore.webview.doc?.tlp.label}></Tlp>
        {/if}
      </Label>
    </div>
    <div
      class="grid grid-cols-1 justify-start gap-2 md:justify-between lg:grid-cols-[minmax(100px,_1fr)_400px]"
    >
      <Label class="mt-4 max-w-full hyphens-auto text-gray-600 [word-wrap:break-word]"
        >{params.publisherNamespace}</Label
      >
      <div class="mt-4 flex h-fit flex-row gap-2 self-center">
        <WorkflowStates {advisoryState} updateStateFn={updateState}></WorkflowStates>
      </div>
    </div>
    <div class="mb-4 mt-2" />
  </div>
  <ErrorMessage error={loadForwardTargetsError}></ErrorMessage>
  <ErrorMessage error={loadAdvisoryVersionsError}></ErrorMessage>
  <ErrorMessage error={stateError}></ErrorMessage>
  <ErrorMessage error={loadDocumentError}></ErrorMessage>
  <ErrorMessage error={loadFourCVEsError}></ErrorMessage>
  <div class={canSeeCommentArea ? "w-full lg:grid lg:grid-cols-[1fr_29rem]" : "w-full"}>
    {#if canSeeCommentArea}
      <div
        class="right-3 mr-3 flex w-full flex-col lg:order-2 lg:max-h-full lg:w-[29rem] lg:flex-none lg:overflow-auto"
      >
        <div class={isSSVCediting || commentFocus ? "w-full p-3 shadow-md" : "w-full p-3"}>
          <div class="flex flex-row items-center">
            {#if ssvcVector}
              {#if !isSSVCediting}
                <SsvcBadge vector={ssvcVector}></SsvcBadge>
              {/if}
            {/if}
            {#if advisoryState !== ARCHIVED && advisoryState !== DELETE}
              <SsvcCalculator
                bind:isEditing={isSSVCediting}
                vectorInput={ssvcVector}
                disabled={!isCalculatingAllowed}
                documentID={params.id}
                on:updateSSVC={loadMetaData}
                {allowEditing}
              ></SsvcCalculator>
            {/if}
          </div>
          {#if isCommentingAllowed && !isSSVCediting}
            <div class="mt-6">
              <Label class="mb-2" for="comment-textarea"
                >{advisoryState === ARCHIVED && appStore.isEditor()
                  ? "Reactivate with comment"
                  : "New Comment"}</Label
              >
              <CommentTextArea
                on:focus={() => {
                  commentFocus = true;
                }}
                on:blur={() => {
                  commentFocus = false;
                }}
                on:input={() => (createCommentError = null)}
                on:saveComment={createComment}
                on:saveForReview={sendForReview}
                on:saveForAssessing={sendForAssessing}
                bind:value={comment}
                errorMessage={createCommentError}
                buttonText="Send"
                state={advisoryState}
              ></CommentTextArea>
            </div>
          {/if}
        </div>
        <ErrorMessage error={loadDocumentSSVCError}></ErrorMessage>
        <div class="h-auto">
          <div class="mt-6 h-full">
            <History
              state={advisoryState}
              on:commentUpdate={() => {
                buildHistory();
              }}
              entries={historyEntries}
            >
              <div slot="additionalButtons">
                {#if availableForwardSelection.length != 0}
                  <Button
                    size="xs"
                    color="light"
                    class="h-7 py-1 text-xs"
                    on:click={() => (openForwardModal = true)}
                  >
                    Forward document</Button
                  >
                {/if}
              </div>
            </History>
          </div>
          <ErrorMessage error={loadEventsError}></ErrorMessage>
          <ErrorMessage error={loadCommentsError}></ErrorMessage>
        </div>
      </div>
    {/if}
    <div
      class={"flex h-auto flex-col lg:order-1 lg:max-h-full lg:flex-auto lg:pr-6" +
        (canSeeCommentArea ? " lg:overflow-auto" : "")}
    >
      <div class="flex flex-row">
        {#if advisoryVersions.length > 0}
          <Version
            publisherNamespace={params.publisherNamespace}
            trackingID={params.trackingID}
            {advisoryVersions}
            selectedDocumentVersion={document.tracking?.version}
            on:selectedDiffDocuments={() => (isDiffOpen = true)}
            on:disableDiff={() => (isDiffOpen = false)}
          ></Version>
        {/if}
      </div>
      <div class="flex flex-col">
        {#if isDiffOpen}
          <Diff showTitle={false}></Diff>
        {:else}
          <Webview
            widthOffset={canSeeCommentArea ? 464 : 0}
            basePath={"#/advisories/" +
              params.publisherNamespace +
              "/" +
              params.trackingID +
              "/documents/" +
              params.id +
              "/"}
            {position}
          ></Webview>
          {#if !canSeeCommentArea && availableForwardSelection.length != 0}
            <div class="my-2 flex w-full flex-row justify-end">
              <Button
                size="xs"
                color="light"
                class="h-7 py-1 text-xs"
                on:click={() => (openForwardModal = true)}
              >
                Forward document
              </Button>
            </div>
          {/if}
        {/if}
      </div>
    </div>
  </div>
</div>
