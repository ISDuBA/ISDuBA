<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { push } from "svelte-spa-router";
  import { appStore } from "$lib/store";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import { onMount } from "svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import Activity from "./Activity.svelte";
  import { getPublisher } from "$lib/utils";
  import { convertVectorToLabel } from "$lib/Advisories/SSVC/SSVCCalculator";

  let documents: any[] = [];
  let newDocumentsError = "";

  const compareCrit = (a: any, b: any) => {
    if (!b.critical || a.critical > b.critical) {
      return -1;
    } else if (!a.critical || a.critical < b.critical) {
      return 1;
    }
    return 0;
  };

  const loadDocuments = async () => {
    const columns =
      "cvss_v3_score cvss_v2_score comments critical id recent title publisher ssvc state tracking_id";
    const query = "$state new workflow =";
    const sort = "-recent";
    const response = await request(
      `/api/documents?columns=${columns}&advisories=true&query=${query}&limit=6&order=${sort}`,
      "GET"
    );
    if (response.ok) {
      documents = (await response.content.documents)?.sort(compareCrit) ?? [];
    } else if (response.error) {
      newDocumentsError = `Could not load new documents. ${getErrorMessage(response.error)}`;
    }
  };

  onMount(async () => {
    await loadDocuments();
  });

  const openDocument = (doc: any) => {
    push(`/advisories/${doc.publisher}/${doc.tracking_id}/documents/${doc.id}`);
  };
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="flex w-1/2 max-w-[50%] flex-col gap-4">
    <SectionHeader title="New advisories"></SectionHeader>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if documents?.length && documents.length > 0}
        {#each documents as doc}
          <Activity on:click={() => openDocument(doc)}>
            <div slot="top-left">
              {#if doc.critical}
                <div>
                  {#if doc.cvss_v3_score && doc.cvss_v3_score === doc.critical}
                    <span>CVSS v3:</span>
                  {:else if doc.cvss_v2_score && doc.cvss_v2_score === doc.critical}
                    <span>CVSS v2:</span>
                  {:else}
                    <span>Critical:</span>
                  {/if}
                  <span class:text-red-500={Number(doc.critical) > 5.0}>
                    {doc.critical}
                  </span>
                </div>
              {/if}
            </div>
            <span slot="top-right" class="ml-auto" title={doc.publisher}
              >{getPublisher(doc.publisher)}</span
            >
            <div class="text-black">{doc.title}</div>
            <div
              slot="bottom-left"
              title={`Number of comments`}
              class="flex items-center gap-4 text-gray-500"
            >
              <div class="flex items-center gap-1">
                <i class="bx bx-comment"></i>
                <span>{doc.comments}</span>
              </div>
              {#if doc.ssvc}
                <span title="SSVC" class="rounded border border-solid border-gray-400 px-1"
                  >{convertVectorToLabel(doc.ssvc).label}</span
                >
              {/if}
              <div></div>
            </div>
            <div slot="bottom-right" class="text-gray-500">
              {#if doc.isNewAdvisory === false}
                <span>New version</span>
              {/if}
            </div>
          </Activity>
        {/each}
      {/if}
    </div>
    <ErrorMessage message={newDocumentsError}></ErrorMessage>
  </div>
{/if}
