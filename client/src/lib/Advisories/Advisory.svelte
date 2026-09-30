<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { Button, Select, Label, Modal, Spinner, ButtonGroup, RadioButton } from "flowbite-svelte";
  import { onDestroy, onMount, setContext, untrack } from "svelte";
  import { appStore } from "$lib/store.svelte";
  import Version from "$lib/Advisories/Version.svelte";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/Advisories/SSVC/SSVCCalculator.svelte";
  import Diff from "$lib/Diff/Diff.svelte";
  import { ARCHIVED, ASSESSING, DELETE, NEW, READ, REVIEW } from "$lib/workflow";
  import { canSetStateRead } from "$lib/permissions";
  import CommentTextArea from "./Events/Comments/CommentTextArea.svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import WorkflowStates from "./WorkflowStates.svelte";
  import History from "./Events/Events.svelte";
  import Tlp from "./TLP.svelte";
  import { addSlashes, selectedClass } from "$lib/utils";
  import {
    type AdvisoryVersion,
    fetchDocumentSSVC,
    fetchSearchHits,
    loadAdvisoryVersions,
    advisorySearchState,
    isResultConsistent
  } from "$lib/Advisories/advisory.svelte";
  import InconsistencyMessage from "$lib/Advisories/InconsistencyMessage.svelte";
  import SearchMatchBar from "./SearchMatchBar.svelte";
  import SearchableText from "./CSAFWebview/SearchableText.svelte";
  import {
    Check,
    AlertCircle,
    ArrowRightStroke,
    ArrowInDownSquareHalf,
    Code,
    GridRowBottom
  } from "@boxicons/svelte";
  import RawDocument from "./RawDocument.svelte";
  import type { CommentEvent, GeneralEvent, OtherEvent, SSVCEvent } from "./Events/events";
  import CopyButton from "$lib/Components/CopyButton.svelte";

  let { params } = $props();

  let oldParams: any = $state(null);
  let csafDocument: any = $state({});
  let ssvcVector: string = $state("");
  let comment: string = $state("");
  let loadCommentsError: ErrorDetails | null = $state(null);
  let loadCommentsAbortController: AbortController | null = $state(null);
  let loadEventsError: ErrorDetails | null = $state(null);
  let loadEventsAbortController: AbortController | null = $state(null);
  let loadAdvisoryVersionsError: ErrorDetails | null = $state(null);
  let loadDocumentError: ErrorDetails | null = $state(null);
  let loadFourCVEsError: ErrorDetails | null = $state(null);
  let createCommentError: ErrorDetails | null = $state(null);
  let loadDocumentSSVCError: ErrorDetails | null = $state(null);
  let loadForwardTargetsError: ErrorDetails | null = $state(null);
  let stateError: ErrorDetails | null = $state(null);
  let loadRelatedError: ErrorDetails | null = $state(null);
  let loadSearchHitsError: ErrorDetails | null = $state(null);
  let advisoryVersions: AdvisoryVersion[] = $state([]);
  let advisoryVersionByDocumentID: any = $state(undefined);
  let advisoryState: string = $state("");
  let historyEntries: Array<CommentEvent | SSVCEvent | OtherEvent> = $state([]);
  let isCommentingAllowed: boolean = $state(false);
  let isSSVCediting = $state(false);
  let position = $state("");
  let processRunning = $state(false);
  let lastSuccessfulForwardTarget: number | undefined = $state(undefined);
  let isInconsistent = $state(false);
  let documentNotFound = $state(false);
  let couldNotLoadDocument = $state(false);
  let relatedDocuments: any = $state(undefined);
  let isLoadingSearchMatches = $state(false);
  let abortControllers: AbortController[] = $state([]);
  let showRawDocument = $state(false);
  let downloadAbortController: AbortController | undefined = $state(undefined);

  $effect(() => {
    if ([NEW, READ, ASSESSING].includes(advisoryState)) {
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
  });

  let isCalculatingAllowed: boolean = $state(false);
  $effect(() => {
    if ([NEW, READ, ASSESSING].includes(advisoryState)) {
      isCalculatingAllowed = appStore.isEditor();
    } else {
      isCalculatingAllowed = false;
    }
  });
  let canSeeCommentArea = $derived(
    appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor() || appStore.isAdmin()
  );
  let encodedTrackingID = $derived(
    csafDocument?.tracking?.id
      ? encodeURIComponent(addSlashes(csafDocument.tracking.id))
      : undefined
  );
  let encodedPublisherNamespace = $derived(
    csafDocument?.publisher?.name
      ? encodeURIComponent(addSlashes(csafDocument.publisher.name))
      : undefined
  );
  let json = $derived(
    appStore.state.webview.rawDoc ? JSON.stringify(appStore.state.webview.rawDoc, null, 2) : ""
  );

  const setAsReadTimeout: number[] = [];
  let isDiffOpen = $state(false);
  let commentFocus = $state(false);

  let availableForwardSelection: any[] = $state([]);
  let selectedForwardTarget: number | undefined = $state();

  const abortAllRequests = () => {
    abortControllers.forEach((c) => {
      c.abort();
    });
  };

  const createAbortController = () => {
    const abortController = new AbortController();
    abortControllers.push(abortController);
    return abortController;
  };

  const getAdvisoryVersions = async () => {
    if (!encodedTrackingID || !encodedPublisherNamespace) return;
    const result = await loadAdvisoryVersions(encodedTrackingID, encodedPublisherNamespace);
    if (result) {
      if (result.error) {
        loadAdvisoryVersionsError = result.error;
      } else if (result.advisoryVersions) {
        advisoryVersions = result.advisoryVersions;
      }
    }
    advisoryVersionByDocumentID = advisoryVersions.reduce((acc: any, version: AdvisoryVersion) => {
      acc[version.id] = version.version;
      return acc;
    }, {});
  };

  const loadDocument = async () => {
    csafDocument = {};
    appStore.setDocument(null);
    isInconsistent = false;
    documentNotFound = false;
    couldNotLoadDocument = false;
    const abortController = createAbortController();
    const response = await request(
      `/api/documents/${params.id}`,
      "GET",
      undefined,
      abortController
    );
    if (response.ok) {
      const result = await response.content;
      if (!isResultConsistent(params, result.document)) {
        isInconsistent = true;
      }
      csafDocument = result.document;
      appStore.setRawDocument(result);
      const docModel = convertToDocModel(result);
      appStore.setDocument(docModel);
    } else if (response.error) {
      if (response.error === "AbortError") {
        return;
      }
      couldNotLoadDocument = true;
      if (response.error === "404") {
        documentNotFound = true;
      } else {
        loadDocumentError = getErrorDetails(`Could not load document.`, response);
      }
    }
  };

  const loadDocumentSSVC = async () => {
    if (params?.id) {
      const abortController = createAbortController();
      const result = await fetchDocumentSSVC(params.id, abortController);
      if (typeof result === "string") {
        ssvcVector = result;
      } else if (result?.message) {
        loadDocumentSSVCError = result;
      }
    }
  };

  const loadEvents = async () => {
    if (!csafDocument || !encodedPublisherNamespace || !encodedTrackingID) return;
    if (loadEventsAbortController) {
      loadEventsAbortController.abort();
    }
    loadEventsAbortController = createAbortController();
    const response = await request(
      `/api/events/${encodedPublisherNamespace}/${encodedTrackingID}`,
      "GET",
      undefined,
      loadEventsAbortController
    );
    if (response.ok) {
      return await response.content;
    } else if (response.error) {
      if (response.error !== "AbortError") {
        loadEventsError = getErrorDetails(`Could not load events.`, response);
      }
      return [];
    }
  };

  const loadComments = async (): Promise<CommentEvent[] | undefined> => {
    if (!csafDocument || !encodedPublisherNamespace || !encodedTrackingID) return;
    if (loadCommentsAbortController) {
      loadCommentsAbortController.abort();
    }
    loadCommentsAbortController = createAbortController();
    const response = await request(
      `/api/comments/${encodedPublisherNamespace}/${encodedTrackingID}`,
      "GET",
      undefined,
      loadCommentsAbortController
    );
    if (response.ok) {
      let comments: CommentEvent[] = await response.content;
      for (let i = 0; i < comments.length; i++) {
        comments[i].documentVersion = advisoryVersionByDocumentID[comments[i].document_id];
      }
      return comments;
    } else if (response.error) {
      if (response.error !== "AbortError") {
        loadCommentsError = getErrorDetails(`Could not load comments.`, response);
      }
      return [];
    }
  };

  // Similar structure to loadComments()
  const loadSSVCHistory = async () => {
    if (!encodedPublisherNamespace || !encodedTrackingID) return;
    const abortController = createAbortController();
    const response = await request(
      `/api/ssvc/history/${encodedPublisherNamespace}/${encodedTrackingID}`,
      "GET",
      undefined,
      abortController
    );
    if (response.ok) {
      return await response.content;
    }
    // Not found -> Empty History
    if (response.error === "404") {
      return { ssvcChanges: [] };
    }
    if (response.error) {
      if (response.error !== "AbortError") {
        loadDocumentSSVCError = getErrorDetails(`Could not load SSVC history`, response);
      }
      return { ssvcChanges: [] };
    }
  };

  const buildHistory = async () => {
    if (!canSeeCommentArea || !csafDocument || !encodedPublisherNamespace || !encodedTrackingID) {
      historyEntries = [];
      return;
    }
    const comments: CommentEvent[] | undefined = await loadComments();
    let events = await loadEvents();
    if (!events || !comments) {
      historyEntries = [];
      return;
    }
    const ssvcData = await loadSSVCHistory();

    const ssvcChanges = ssvcData?.ssvcChanges || [];

    const commentsByTime = comments.reduce((o: any, event: CommentEvent) => {
      o[`${event.time}:${event.commentator}`] = {
        commentator: event.commentator,
        message: event.message,
        id: event.id,
        documentVersion: event.documentVersion
      };
      return o;
    }, {});

    // Same logic as commentsByTime
    const ssvcByTime = ssvcChanges.reduce((o: any, event: SSVCEvent) => {
      o[`${event.changedate}:${event.actor}:${event.documents_id}`] = event;
      return o;
    }, {});

    const commentsEdited = events
      .filter((e: GeneralEvent) => {
        return e.event_type === "change_comment";
      })
      .map((e: CommentEvent) => {
        return {
          id: e.comment_id,
          time: e.time
        };
      })
      .reduce((o: any, event: CommentEvent) => {
        if (!o[event.id]) o[event.id] = [];
        o[event.id].push(event.time);
        return o;
      }, {});
    events.map((e: any) => {
      if (e.event_type === "add_comment") {
        const comment = commentsByTime[`${e.time}:${e.actor}`];
        e["commentator"] = comment.commentator;
        e["message"] = comment.message;
        e["comment_id"] = comment.id;
        e["documentVersion"] = comment.documentVersion;
        if (commentsEdited[comment.id]) {
          e["times"] = commentsEdited[comment.id];
        }
      }
      if (e.event_type === "add_sscv" || e.event_type === "change_sscv") {
        const IDsOfAllVersions = $state.snapshot(advisoryVersions).map((v) => v.id);
        const matchingVersion = IDsOfAllVersions.find((ver) => {
          return ssvcByTime[`${e.time}:${e.actor}:${ver}`] !== undefined;
        });
        if (matchingVersion) {
          const ssvcMatch = ssvcByTime[`${e.time}:${e.actor}:${matchingVersion}`];
          e["ssvc"] = ssvcMatch.ssvc;
          e["prev_ssvc"] = ssvcMatch.ssvc_prev;
          e["documentVersion"] = ssvcMatch.documents_version;
          e["documents_version"] = ssvcMatch.documents_version;
          e["documents_id"] = ssvcMatch.documents_id;
        }
      }

      return e;
    });

    historyEntries = events;
  };

  async function createComment() {
    await allowEditing();
    const formData = new FormData();
    // Clear comment before request to avoid sending duplicate comments
    let commentTmp = comment;
    comment = "";
    formData.append("message", commentTmp);
    const response = await request(`/api/comments/${params.id}`, "POST", formData);
    if (response.ok) {
      await loadAdvisoryState();
      await buildHistory();
    } else if (response.error) {
      // Restore comment on error
      comment = commentTmp;
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
      `/api/status/${encodedPublisherNamespace}/${encodedTrackingID}/${newState}`,
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
    const abortController = createAbortController();
    const response = await request(
      `/api/documents?advisories=true&columns=state&query=$tracking_id ${encodedTrackingID} = $publisher "${encodedPublisherNamespace}" = and`,
      "GET",
      undefined,
      abortController
    );
    if (response.ok) {
      const result = response.content;
      advisoryState = result.documents?.[0].state;
      return result.documents?.[0].state;
    } else if (response.error && response.error !== "AbortError") {
      stateError = getErrorDetails(`Couldn't load state.`, response);
    }
  };

  const loadFourCVEs = async () => {
    const abortController = createAbortController();
    const response = await request(
      `/api/documents?advisories=false&columns=four_cves&query=$id ${params.id} integer =`,
      "GET",
      undefined,
      abortController
    );
    if (response.ok) {
      const content = await response.content;
      let four_cves = content?.documents[0]?.four_cves;
      appStore.setFourCVEs(four_cves);
    } else if (response.error && response.error !== "AbortError") {
      loadFourCVEsError = getErrorDetails(`Couldn't load CVEs.`, response);
    }
  };

  const loadRelatedDocuments = async () => {
    if (
      !(
        appStore.isEditor() ||
        appStore.isAdmin() ||
        appStore.isAuditor() ||
        appStore.isReviewer()
      ) ||
      !params
    ) {
      return;
    }
    const abortController = createAbortController();
    const response = await request(
      `/api/documents/${params.id}/cve_related`,
      "GET",
      undefined,
      abortController
    );
    if (response.ok) {
      relatedDocuments = {};
      response.content.forEach((doc: any) => {
        if (!relatedDocuments[doc.document_id]) {
          relatedDocuments[doc.document_id] = doc;
          relatedDocuments[doc.document_id].cve = [doc.cve];
        } else if (relatedDocuments) {
          relatedDocuments[doc.document_id].cve.push(doc.cve);
        }
      });
    } else if (response.error && response.error !== "AbortError") {
      loadRelatedError = getErrorDetails(`Could not load documents.`, response);
    }
  };

  const loadData = async () => {
    advisoryState = "";
    historyEntries = [];
    advisoryVersions = [];
    ssvcVector = "";
    appStore.setDocument(null);
    await loadDocument();
    await getAdvisoryVersions();
    if (appStore.state.app.search.term) {
      isLoadingSearchMatches = true;
      const hitsResult = await fetchSearchHits(params.id);
      isLoadingSearchMatches = false;
      if (Array.isArray(hitsResult)) {
        advisorySearchState.searchMatches = hitsResult;
        if (advisorySearchState.searchMatches.length > 0 && advisorySearchState.matchIndex === -1) {
          advisorySearchState.matchIndex = 0;
        }
      } else {
        loadSearchHitsError = hitsResult;
      }
    }
    if (couldNotLoadDocument || isInconsistent) return;
    if (csafDocument) {
      await loadFourCVEs();
      await loadDocumentSSVC();
      await buildHistory();
      await loadAdvisoryState();
      loadRelatedDocuments();
      // Only set state to 'read' if editor opens the current version.
      if (
        advisoryState === NEW &&
        canSetStateRead(advisoryState) &&
        (advisoryVersions.length === 1 ||
          advisoryVersions[0].version === csafDocument.tracking?.version)
      ) {
        const id: any = setTimeout(async () => {
          if (advisoryState === "new" && canSetStateRead(advisoryState)) {
            await updateState(READ);
          }
        }, 20000);
        setAsReadTimeout.push(id);
      }
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

  const downloadRawDocument = () => {
    if (downloadAbortController) {
      downloadAbortController.abort();
    }
    downloadAbortController = new AbortController();
    const file = new Blob([json], { type: "application/json" });
    let a = document.createElement("a"),
      url = URL.createObjectURL(file);
    a.href = url;
    a.download = `${appStore.state.webview.rawDoc.document.tracking.id}.json`;
    document.body.appendChild(a);
    a.click();
    setTimeout(() => {
      document.body.removeChild(a);
      window.URL.revokeObjectURL(url);
    }, 0);
  };

  onDestroy(() => {
    appStore.setDocument(null);
    advisorySearchState.matchIndex = -1;
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

  $effect(() => {
    const old = untrack(() => oldParams);
    if (params) {
      setTimeout(() => {
        if (!old || JSON.stringify(old) !== JSON.stringify(params)) {
          abortAllRequests();
          loadData();
        }
      }, 0);
      oldParams = params;
      position = params.position;
      if (!params.position) {
        const topElement = window.document.getElementById("top");
        topElement?.scrollIntoView();
        appStore.setSelectedProduct("");
        appStore.setSelectedCVE("");
      }
    }
  });
  let openForwardModal = $state(false);

  setContext("advisory", () => relatedDocuments);
  setContext("advisoryVersions", () => advisoryVersions);
  setContext("params", () => params);
</script>

<svelte:head>
  <title>{documentNotFound ? "Document not found" : csafDocument.tracking?.id}</title>
</svelte:head>

<Modal bind:open={openForwardModal}>
  <Label class="text-lg">Forward document</Label>
  <Select items={availableForwardSelection} bind:value={selectedForwardTarget}></Select>
  {#if typeof selectedForwardTarget === "number"}
    <Button disabled={processRunning} onclick={forwardDocument}>
      <span class="mr-2">Send document</span>
      {#if processRunning}
        <Spinner></Spinner>
      {:else if lastSuccessfulForwardTarget === selectedForwardTarget}
        <div class="inline-flex w-8 items-center"><Check class="text-2xl" /></div>
      {:else}
        <div class="inline-flex w-8 items-center">
          <ArrowRightStroke class="text-2xl" />
        </div>
      {/if}
    </Button>
  {/if}
</Modal>

<div
  class="relative grid h-fit w-full grow grid-rows-[auto_minmax(100px,_1fr)] gap-y-2 px-2 lg:h-full"
  id="top"
>
  {#if documentNotFound}
    <div class="mb-2 font-bold">
      <AlertCircle aria-hidden="true" />
      <span>The URL doesn't reference any document</span>
    </div>
  {:else if isInconsistent}
    <InconsistencyMessage document={csafDocument} {params}></InconsistencyMessage>
  {:else if !couldNotLoadDocument && !isInconsistent}
    <div
      class="sticky -top-6 z-100 flex w-full flex-none flex-col bg-white pt-6 lg:static lg:pt-0 dark:bg-gray-800"
    >
      <div class="flex flex-wrap items-center gap-x-6 gap-y-1">
        <Label class="text-lg">
          <span class="mr-2">
            <SearchableText
              text={csafDocument.tracking ? csafDocument.tracking.id : ""}
              textPath="/document/tracking/id"
            />
          </span>
          {#if appStore.state.webview.doc?.tlp.label}
            <Tlp tlp={appStore.state.webview.doc?.tlp.label}></Tlp>
          {/if}
        </Label>
        <div class="flex gap-1">
          <ButtonGroup color="light" size="sm" class="h-7">
            <RadioButton
              checkedClass={selectedClass}
              title="View raw document"
              value={false}
              bind:group={showRawDocument}
            >
              <GridRowBottom />
            </RadioButton>
            <RadioButton
              checkedClass={selectedClass}
              title="View raw document"
              value={true}
              bind:group={showRawDocument}
            >
              <Code />
            </RadioButton>
          </ButtonGroup>
          <Button
            onclick={downloadRawDocument}
            class="h-7 py-1"
            color="light"
            size="xs"
            title="Download document"
          >
            <ArrowInDownSquareHalf />
          </Button>
          <CopyButton
            errorMessage="Could not copy the document"
            showBorder={true}
            title="Copy document"
            tooltipPlacement="bottomright"
            value={json}
          />
        </div>
        {#if isLoadingSearchMatches}
          <Spinner color="gray" size="4"></Spinner>
        {:else if appStore.state.app.search.term && appStore.state.webview.doc && !appStore.state.app.search.advanced}
          <SearchMatchBar />
        {/if}
      </div>
      <div
        class="grid grid-cols-1 justify-start gap-2 md:justify-between lg:grid-cols-[minmax(100px,_1fr)_500px]"
      >
        <div class="flex flex-col gap-2">
          <Label class="mt-4 max-w-full hyphens-auto text-gray-600 [word-wrap:break-word]"
            >{csafDocument.publisher ? csafDocument.publisher.name : ""}</Label
          >
        </div>
        <div class="mt-4 flex h-fit flex-row gap-2 self-center">
          <WorkflowStates {advisoryState} updateStateFn={updateState}></WorkflowStates>
        </div>
      </div>
      <div class="mt-2 mb-4"></div>
    </div>
  {/if}
  <ErrorMessage bind:error={loadForwardTargetsError}></ErrorMessage>
  <ErrorMessage bind:error={loadAdvisoryVersionsError}></ErrorMessage>
  <ErrorMessage bind:error={stateError}></ErrorMessage>
  <ErrorMessage bind:error={loadDocumentError}></ErrorMessage>
  <ErrorMessage bind:error={loadFourCVEsError}></ErrorMessage>
  <ErrorMessage bind:error={loadRelatedError}></ErrorMessage>
  <ErrorMessage bind:error={loadSearchHitsError}></ErrorMessage>
  {#if !couldNotLoadDocument && !isInconsistent}
    <div class={canSeeCommentArea ? "w-full lg:grid lg:grid-cols-[1fr_29rem]" : "w-full"}>
      {#if canSeeCommentArea}
        <div
          class="right-3 mr-3 flex w-full flex-col lg:order-2 lg:max-h-full lg:w-[29rem] lg:flex-none lg:overflow-auto"
        >
          <div class={isSSVCediting || commentFocus ? "w-full p-3 shadow-md" : "w-full p-3"}>
            <div class="flex flex-col">
              {#if advisoryState !== ARCHIVED && advisoryState !== DELETE}
                <SsvcCalculator
                  bind:isEditing={isSSVCediting}
                  vectorInput={ssvcVector}
                  disabled={!isCalculatingAllowed}
                  documentID={params.id}
                  {ssvcVector}
                  updateSSVC={loadMetaData}
                  {allowEditing}
                ></SsvcCalculator>
              {/if}
            </div>
            {#if isCommentingAllowed}
              <div class="mt-6">
                <Label class="mb-2" for="comment-textarea"
                  >{advisoryState === ARCHIVED && appStore.isEditor()
                    ? "Reactivate with comment"
                    : "New Comment"}</Label
                >
                <CommentTextArea
                  onFocus={() => {
                    commentFocus = true;
                  }}
                  onBlur={() => {
                    commentFocus = false;
                  }}
                  onInput={() => (createCommentError = null)}
                  saveComment={createComment}
                  saveForReview={sendForReview}
                  saveForAssessing={sendForAssessing}
                  bind:value={comment}
                  errorMessage={createCommentError}
                  buttonText="Send"
                  workflowState={advisoryState}
                ></CommentTextArea>
              </div>
            {/if}
          </div>
          <ErrorMessage bind:error={loadDocumentSSVCError}></ErrorMessage>
          <div class="h-auto">
            <div class="mt-6 h-full">
              <History
                workflowState={advisoryState}
                onCommentUpdated={(newComment: string, index: number) => {
                  // First update the comment locally so the user can see that editing the comment did work
                  const event: CommentEvent = historyEntries[index] as unknown as CommentEvent;
                  if (event.event_type === "add_comment") {
                    event.message = newComment;
                  } else {
                    const originalEvent: CommentEvent = historyEntries.find((e: any) => {
                      return e.event_type === "add_comment" && event.comment_id === e.comment_id;
                    }) as unknown as CommentEvent;
                    if (originalEvent) {
                      originalEvent.message = newComment;
                    }
                  }
                  // Then refresh the whole history
                  buildHistory();
                }}
                entries={historyEntries}
              >
                {#snippet additionalButtons()}
                  <div>
                    {#if availableForwardSelection.length != 0}
                      <Button
                        size="xs"
                        color="light"
                        class="h-7 py-1 text-xs"
                        onclick={() => (openForwardModal = true)}
                      >
                        Forward document</Button
                      >
                    {/if}
                  </div>
                {/snippet}
              </History>
            </div>
            <ErrorMessage bind:error={loadEventsError}></ErrorMessage>
            <ErrorMessage bind:error={loadCommentsError}></ErrorMessage>
          </div>
        </div>
      {/if}
      <div
        class={"flex h-auto flex-col lg:order-1 lg:max-h-full lg:flex-auto lg:pr-6" +
          (canSeeCommentArea ? " lg:overflow-auto" : "")}
      >
        <div class="mb-2 flex flex-col gap-2">
          {#if advisoryVersions?.length > 0}
            <Version
              publisherNamespace={csafDocument.publisher?.name}
              {advisoryVersions}
              selectedDocumentVersion={{
                id: csafDocument.id,
                tracking_id: csafDocument.tracking?.id,
                tracking_status: csafDocument.tracking?.status,
                version: csafDocument.tracking?.version
              }}
              selectedDiffDocuments={() => (isDiffOpen = true)}
              onDisabledDiff={() => (isDiffOpen = false)}
            ></Version>
          {/if}
        </div>
        <div class="flex flex-col">
          {#if isDiffOpen}
            <Diff showTitle={false}></Diff>
          {:else}
            {#if appStore.state.webview.doc}
              {#if showRawDocument}
                <RawDocument />
              {:else}
                <Webview
                  basePath={"#/advisories/" +
                    csafDocument.publisher?.name +
                    "/" +
                    csafDocument.tracking?.id +
                    "/documents/" +
                    params.id +
                    "/"}
                  {position}
                ></Webview>
              {/if}
            {:else}
              <div class="mt-32 ml-32">
                <Spinner color="gray" size="8"></Spinner>
              </div>
            {/if}
            {#if !canSeeCommentArea && availableForwardSelection.length != 0}
              <div class="my-2 flex w-full flex-row justify-end">
                <Button
                  size="xs"
                  color="light"
                  class="h-7 py-1 text-xs"
                  onclick={() => (openForwardModal = true)}
                >
                  Forward document
                </Button>
              </div>
            {/if}
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
