<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import {
    Button,
    Label,
    Timeline,
    AccordionItem,
    Accordion,
    Badge,
    Tooltip,
    Dropdown,
    DropdownItem
  } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Comment from "$lib/Advisories/Comment.svelte";
  import Version from "$lib/Advisories/Version.svelte";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/Advisories/SSVC/SSVCCalculator.svelte";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";
  import {
    ASSESSING,
    ARCHIVED,
    DELETED,
    NEW,
    READ,
    REVIEW,
    canSetStateArchived,
    canSetStateAssessing,
    canSetStateDeleted,
    canSetStateNew,
    canSetStateRead,
    canSetStateReview,
    getAllowedWorkflowChanges
  } from "$lib/permissions";
  import CommentTextArea from "./CommentTextArea.svelte";
  import { request } from "$lib/utils";
  export let params: any = null;

  let document = {};
  let ssvc: any;
  $: ssvcStyle = ssvc
    ? `color: ${ssvc.color}; border: 1pt solid ${ssvc.color}; background-color: white;`
    : "";
  let comment: string = "";
  let comments: any = [];
  let advisoryVersions: string[] = [];
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

  const loadAdvisoryVersions = async () => {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      const response = await fetch(
        `/api/documents?&columns=id version&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
        {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          }
        }
      );
      if (response.ok) {
        const result = await response.json();
        advisoryVersions = result.documents.map((doc: any) => {
          return { id: doc.id, version: doc.version };
        });
        advisoryVersionByDocumentID = advisoryVersions.reduce((acc: any, version: any) => {
          acc[version.id] = version.version;
          return acc;
        }, {});
      } else {
        appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
      }
    });
  };

  const loadDocument = async () => {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      const response = await fetch(`/api/documents/${params.id}`, {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      });
      if (response.ok) {
        const doc = await response.json();
        ({ document } = doc);
        const docModel = convertToDocModel(doc);
        appStore.setDocument(docModel);
      } else {
        appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
      }
    });
  };

  const loadDocumentSSVC = async () => {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      const response = await fetch(
        `/api/documents?columns=ssvc&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
        {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          }
        }
      );
      if (response.ok) {
        const result = await response.json();
        if (result.documents[0].ssvc) {
          ssvc = convertVectorToLabel(result.documents[0].ssvc);
        }
      } else {
        appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
      }
    });
  };

  function loadComments(): Promise<any[]> {
    return new Promise((resolve) => {
      const newComments: any = [];
      advisoryVersions.forEach((advVer: any) => {
        $appStore.app.keycloak.updateToken(5).then(async () => {
          fetch(`/api/comments/${advVer.id}`, {
            headers: {
              Authorization: `Bearer ${$appStore.app.keycloak.token}`
            }
          }).then((response) => {
            if (response.ok) {
              response.json().then((json) => {
                if (json) {
                  json.forEach((c: any) => {
                    c.documentVersion = advisoryVersionByDocumentID[c.document_id];
                  });
                  newComments.push(...json);
                }
                comments = newComments;
                resolve(newComments);
              });
            } else {
              appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
            }
          });
        });
      });
    });
  }
  async function createComment() {
    const formData = new FormData();
    formData.append("message", comment);
    const response = await request(`/api/comments/${params.id}`, "POST", formData);
    if (response) {
      comment = "";
      loadComments().then((newComments: any[]) => {
        if (newComments.length === 1) {
          loadAdvisoryState();
        }
      });
      appStore.displaySuccessMessage("Comment for advisory saved.");
    }
  }

  async function updateState(newState: string) {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      const response = await fetch(
        `/api/status/${params.publisherNamespace}/${params.trackingID}/${newState}`,
        {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          },
          method: "PUT"
        }
      );
      if (response.ok) {
        advisoryState = newState;
      } else {
        const error = await response.json();
        appStore.displayErrorMessage(`${error.error}`);
      }
    });
  }

  const loadAdvisoryState = async () => {
    $appStore.app.keycloak.updateToken(5).then(async () => {
      const response = await fetch(
        `/api/documents?advisories=true&columns=state&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
        {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          }
        }
      );
      if (response.ok) {
        const result = await response.json();
        advisoryState = result.documents[0].state;
        return result.documents[0].state;
      } else {
        appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
      }
    });
  };

  function loadMetaData() {
    loadAdvisoryState();
    loadDocumentSSVC();
  }

  onDestroy(() => {
    timeoutIDs.forEach((id: number) => {
      clearTimeout(id);
    });
  });

  onMount(async () => {
    if ($appStore.app.keycloak.authenticated) {
      loadDocumentSSVC();
      await loadDocument();
      await loadAdvisoryVersions();
      if (appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()) {
        loadComments();
      }
      const state = await loadAdvisoryState();
      // Only set state to 'read' if editor opens the current version.
      if (
        state === "new" &&
        appStore.isEditor() &&
        (advisoryVersions.length === 1 ||
          advisoryVersions[0].version === document.tracking?.version)
      ) {
        const id = setTimeout(async () => {
          await updateState("read");
          appStore.displayInfoMessage("This advisory is marked as read");
        }, 3000);
        timeoutIDs.push(id);
      }
    }
  });
</script>

<div class="flex">
  <div class="flex flex-col">
    <div class="flex flex-col">
      <div class="flex gap-2">
        <Label class="text-lg">{params.trackingID}</Label>
        <Button
          class="!p-1"
          color="light"
          disabled={getAllowedWorkflowChanges(advisoryState).length === 0}
        >
          <i class="bx bx-dots-vertical-rounded"></i>
        </Button>
        <Dropdown>
          {#if canSetStateNew(advisoryState)}
            <DropdownItem on:click={() => updateState(NEW)} class="flex items-center gap-2">
              <i class="bx bx-star text-lg"></i>
              <span>Mark as new</span>
            </DropdownItem>
          {/if}
          {#if canSetStateRead(advisoryState)}
            <DropdownItem on:click={() => updateState(READ)} class="flex items-center gap-2">
              <i class="bx bx-show text-lg"></i>
              <span>Mark as read</span>
            </DropdownItem>
          {/if}
          {#if canSetStateReview(advisoryState)}
            <DropdownItem on:click={() => updateState(REVIEW)} class="flex items-center gap-2">
              <i class="bx bx-book-open text-lg"></i>
              <span>Release for review</span>
            </DropdownItem>
          {/if}
          {#if canSetStateAssessing(advisoryState) && advisoryState === REVIEW}
            <DropdownItem on:click={() => updateState(ASSESSING)} class="flex items-center gap-2">
              <i class="bx bx-analyse text-lg"></i>
              <span>Back to assessing</span>
            </DropdownItem>
          {/if}
          {#if canSetStateArchived(advisoryState)}
            <DropdownItem on:click={() => updateState(ARCHIVED)} class="flex items-center gap-2">
              <i class="bx bx-archive text-lg"></i>
              <span>Archive</span>
            </DropdownItem>
          {/if}
          {#if canSetStateDeleted(advisoryState)}
            <DropdownItem on:click={() => updateState(DELETED)} class="flex items-center gap-2">
              <i class="bx bx-trash text-lg"></i>
              <span>Mark for deletion</span>
            </DropdownItem>
          {/if}
        </Dropdown>
      </div>
      <Label class="mb-2 text-gray-600">{params.publisherNamespace}</Label>
      <div class="flex gap-2">
        {#if advisoryState}
          <Badge class="w-fit">{advisoryState}</Badge>
          <Tooltip>Workflow state</Tooltip>
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
    ></Version>
    <Webview></Webview>
  </div>
  {#if appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()}
    <div class="ml-auto mr-3 min-w-96 max-w-96">
      <Accordion>
        <AccordionItem open>
          <span slot="header"
            ><i class="bx bx-comment-detail"></i><span class="ml-2">Comments</span></span
          >
          {#if comments?.length > 0}
            <div class="max-h-96 overflow-y-auto pl-2">
              <Timeline class="mb-4 flex flex-col-reverse">
                {#each comments as comment}
                  <Comment {comment}></Comment>
                {/each}
              </Timeline>
            </div>
          {:else}
            <div class="mb-6 text-gray-600">No comments available.</div>
          {/if}
          {#if isCommentingAllowed}
            <div class="mt-6">
              <Label class="mb-2" for="comment-textarea">New Comment:</Label>
              <CommentTextArea on:saveComment={createComment} bind:value={comment} buttonText="Send"
              ></CommentTextArea>
            </div>
          {/if}
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
