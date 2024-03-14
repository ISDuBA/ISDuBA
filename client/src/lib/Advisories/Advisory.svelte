<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import { page } from "$app/stores";
  import { Button, Drawer, Label, Select, Textarea, Timeline } from "flowbite-svelte";
  import { sineIn } from "svelte/easing";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Comment from "$lib/Advisories/Comment.svelte";
  import Version from "$lib/Advisories/Version.svelte";
  import { getAllowedWorkflowChanges } from "$lib/permissions";
  import Webview from "$lib/Advisories/CSAFWebview/Webview.svelte";
  import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
  import SsvcCalculator from "$lib/SSVC/SSVCCalculator.svelte";
  import { convertVectorToLabel } from "$lib/SSVC/SSVCCalculator";
  export let params: any = null;

  let document = {};
  let ssvc: any;
  $: ssvcStyle = ssvc ? `color: ${ssvc.color}` : "";
  let hideComments = false;
  let comment: string = "";
  $: count = comment.length;
  let comments: any = [];
  let advisoryVersions: string[] = [];
  let showCalculator = false;

  let transitionParams = {
    x: 320,
    duration: 120,
    easing: sineIn
  };

  function toggleCalculator() {
    showCalculator = !showCalculator;
  }

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
      // Do errorhandling
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
      // Do errorhandling
    }
  };

  const loadDocumentSSVC = async () => {
    const response = await fetch(
      `/api/documents?&columns=ssvc&query=$tracking_id ${params.trackingID} = $publisher "${params.publisherNamespace}" = and`,
      {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      }
    );
    if (response.ok) {
      const result = await response.json();
      if (result.documents[0].ssvc) {
        ssvc = await convertVectorToLabel(result.documents[0].ssvc);
      }
    } else {
      // Do errorhandling
    }
  };

  function toggleComments() {
    hideComments = !hideComments;
  }
  function loadComments() {
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
          });
        } else {
          // Do errorhandling
        }
      });
    });
  }
  function createComment() {
    const formData = new FormData();
    formData.append("message", comment);
    fetch(`/api/comments/${$page.params.documentID}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      },
      method: "POST",
      body: formData
    }).then((response) => {
      if (response.ok) {
        comment = "";
        loadComments();
      } else {
        // Do errorhandling
      }
    });
  }

  function updateState(event: any) {
    const newState = event.target.value;
    fetch(`/api/status/${params.publisherNamespace}/${params.trackingID}/${newState}`, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      },
      method: "PUT"
    });
  }

  onMount(async () => {
    if ($appStore.app.keycloak.authenticated) {
      loadDocument();
      await loadAdvisoryVersions();
      if (appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()) {
        loadComments();
      }
      loadDocumentSSVC();
    }
  });
</script>

<div class="flex">
  <div>
    <div class="flex flex-col">
      <div class="flex">
        <div class="me-2 flex-col">
          <Label class="mb-4 max-w-52"
            >Workflow-State:
            <!-- TODO: Replace hard-coded state "new" with current state of document -->
            <Select
              on:change={updateState}
              items={[
                { value: "new", name: "new" },
                ...getAllowedWorkflowChanges(appStore.getRoles(), "new").map((v) => {
                  return { value: v.to, name: v.to };
                })
              ]}
              value={"new"}
              placeholder=""
              underline
            ></Select>
          </Label>
          <Label class="text-lg">
            {#if ssvc}
              <span style={ssvcStyle}>{ssvc.label}</span>
            {:else}
              <span class="text-gray-400">No SSVC</span>
            {/if}
          </Label>
          <Button on:click={toggleCalculator} class="w-40">Toggle Calculator</Button>
        </div>
        <Version
          publisherNamespace={params.publisherNamespace}
          trackingID={params.trackingID}
          {advisoryVersions}
        ></Version>
      </div>
      <Webview></Webview>
    </div>
  </div>
  {#if appStore.isEditor() || appStore.isReviewer() || appStore.isAuditor()}
    <Button
      on:click={toggleComments}
      outline={true}
      class="absolute right-2 top-2 z-[51] !p-2"
      size="lg"
    >
      <i class={hideComments ? "bx bx-chevron-left" : "bx bx-chevron-right"}></i>
    </Button>
    <Drawer
      activateClickOutside={false}
      backdrop={false}
      class="relative flex flex-col"
      placement="right"
      width="w-1/2"
      hidden={hideComments}
      transitionType="in:slide"
      {transitionParams}
    >
      {#if showCalculator}
        <SsvcCalculator documentID={params.id}></SsvcCalculator>
      {:else}
        {#if comments?.length > 0}
          <div class="overflow-y-scroll pl-2">
            <Timeline class="flex flex-col-reverse">
              {#each comments as comment}
                <Comment {comment}></Comment>
              {/each}
            </Timeline>
          </div>
        {:else}
          <span class="mb-4 text-gray-600">No comments available.</span>
        {/if}
        {#if appStore.isEditor() || appStore.isReviewer()}
          <div>
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
      {/if}
    </Drawer>
  {/if}
</div>
