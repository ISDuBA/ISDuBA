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
    Textarea,
    Timeline,
    AccordionItem,
    Accordion,
    Badge,
    Tooltip,
    Modal
  } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Comment from "$lib/Advisories/Comment.svelte";
  import Version from "$lib/Advisories/Version.svelte";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/Advisories/SSVC/SSVCCalculator.svelte";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";
  import JsonDiff from "$lib/Diff/JsonDiff.svelte";
  import type { JsonDiffResultList } from "$lib/Diff/JsonDiff";
  export let params: any = null;

  let document = {};
  let ssvc: any;
  $: ssvcStyle = ssvc
    ? `color: ${ssvc.color}; border: 1pt solid ${ssvc.color}; background-color: white;`
    : "";
  let comment: string = "";
  $: count = comment.length;
  let comments: any = [];
  let advisoryVersions: string[] = [];
  let advisoryState: string;
  const timeoutIDs: number[] = [];
  let diff: any;
  let isDiffOpen = false;

  const loadAdvisoryVersions = async () => {
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
    } else {
      appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
    }
  };

  const loadDocument = async () => {
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
  };

  const loadDocumentSSVC = async () => {
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
  };

  function loadComments(): Promise<any[]> {
    return new Promise((resolve) => {
      const newComments: any = [];
      advisoryVersions.forEach((advVer: any) => {
        fetch(`/api/comments/${advVer.id}`, {
          headers: {
            Authorization: `Bearer ${$appStore.app.keycloak.token}`
          }
        }).then((response) => {
          if (response.ok) {
            response.json().then((json) => {
              if (json) {
                json.forEach((c: any) => {
                  c.documentID = advVer.id;
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
  }
  async function createComment() {
    const formData = new FormData();
    formData.append("message", comment);
    const response = await fetch(`/api/comments/${params.id}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      },
      method: "POST",
      body: formData
    });
    if (response.ok) {
      comment = "";
      loadComments().then((newComments: any[]) => {
        if (newComments.length === 1) {
          loadAdvisoryState();
        }
      });
      appStore.displaySuccessMessage("Comment for advisory saved.");
    } else {
      const error = await response.json();
      appStore.displayErrorMessage(`${error.error}`);
    }
  }

  async function updateState(newState: string) {
    await fetch(`/api/status/${params.publisherNamespace}/${params.trackingID}/${newState}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      },
      method: "PUT"
    });
  }

  const loadAdvisoryState = async () => {
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
  };

  function loadMetaData() {
    loadAdvisoryState();
    loadDocumentSSVC();
  }

  const compareLatest = async () => {
    const firstDoc = advisoryVersions[advisoryVersions.length - 2].id;
    const secondDoc = advisoryVersions[advisoryVersions.length - 1].id;
    const response = await fetch(`/api/diff/${firstDoc}/${secondDoc}?word-diff=true`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      }
    });
    if (response.ok) {
      const result: JsonDiffResultList = await response.json();
      diff = {
        docA: advisoryVersions[advisoryVersions.length - 2].version,
        docB: advisoryVersions[advisoryVersions.length - 1].version,
        result: result
      };
      isDiffOpen = true;
    } else {
      appStore.displayErrorMessage(`${response.status}. ${response.statusText}`);
    }
  };

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
      if (state === "new") {
        const id = setTimeout(async () => {
          await updateState("read");
          appStore.displayInfoMessage("This advisory is marked as read");
          advisoryState = "read";
        }, 3000);
        timeoutIDs.push(id);
      }
    }
  });
</script>

<div class="flex">
  <div class="flex flex-col">
    <div class="flex flex-col">
      <Label class="text-lg">{params.trackingID}</Label>
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
    {#if advisoryVersions.length > 1}
      {#if advisoryVersions[0].version === document.tracking?.version}
        <Button class="w-fit" color="light" on:click={compareLatest}
          ><i class="bx bx-transfer me-2 text-lg"></i>Latest Changes</Button
        >
      {/if}
    {/if}
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
          {#if appStore.isEditor() || appStore.isReviewer()}
            <div class="mt-6">
              <Label class="mb-2" for="comment-textarea">New Comment:</Label>
              <Textarea bind:value={comment} class="mb-2" id="comment-textarea">
                <div slot="footer" class="flex items-start justify-between">
                  <Button on:click={createComment} disabled={count > 10000 || count === 0}
                    >Send</Button
                  >
                  <Label class={count < 10000 ? "text-gray-600" : "font-bold text-red-600"}
                    >{`${count}/10000`}</Label
                  >
                </div>
              </Textarea>
            </div>
          {/if}
        </AccordionItem>
      </Accordion>
      <Accordion class="mt-3">
        <AccordionItem open>
          <span slot="header"><i class="bx bx-calculator"></i><span class="ml-2">SSVC</span></span>
          <SsvcCalculator
            vectorInput={ssvc?.vector}
            documentID={params.id}
            on:updateSSVC={loadMetaData}
          ></SsvcCalculator>
        </AccordionItem>
      </Accordion>
    </div>
    <Modal bind:open={isDiffOpen}>
      <JsonDiff {diff}></JsonDiff>
    </Modal>
  {/if}
</div>
