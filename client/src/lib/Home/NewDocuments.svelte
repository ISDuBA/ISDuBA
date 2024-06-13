<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { request } from "$lib/utils";
  import { getErrorMessage } from "$lib/Errors/error";
  import { onMount } from "svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import CustomCard from "./CustomCard.svelte";
  import { getPublisher } from "$lib/utils";

  let documents: any[] = [];
  let newDocumentsError = "";

  const loadDocuments = async () => {
    const columns = "cvss_v3_score cvss_v2_score title publisher state";
    const response = await request(
      `/api/documents?columns=${columns}&advisories=true&query=$state new workflow =&limit=10`,
      "GET"
    );
    if (response.ok) {
      documents = await response.content.documents;
    } else if (response.error) {
      newDocumentsError = `Could not load new documents. ${getErrorMessage(response.error)}`;
    }
  };

  onMount(() => {
    loadDocuments();
  });
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="flex w-1/2 max-w-[50%] flex-col gap-4">
    <SectionHeader title="New documents"></SectionHeader>
    <div class="text-red-600">
      Attention: These are
      <span class="font-bold">advisories</span>
      for now as we are not able to fetch recently imported documents yet.
    </div>
    <div class="grid grid-cols-[repeat(auto-fit,_minmax(200pt,_1fr))] gap-6">
      {#if documents?.length && documents.length > 0}
        {#each documents as doc}
          <CustomCard>
            <div slot="top-left">
              {#if doc.cvss_v2_score}
                <div>
                  <span>CVSS v2:</span>
                  <span class:text-red-500={Number(doc.cvss_v2_score) > 5.0}>
                    {doc.cvss_v2_score}
                  </span>
                </div>
              {/if}
              {#if doc.cvss_v3_score}
                <div>
                  <span>CVSS v3:</span>
                  <span class:text-red-500={Number(doc.cvss_v3_score) > 5.0}>
                    {doc.cvss_v3_score}
                  </span>
                </div>
              {/if}
            </div>
            <span slot="top-right" class="ml-auto" title={doc.publisher}
              >{getPublisher(doc.publisher)}</span
            >
            <div class="text-black">{doc.title}</div>
          </CustomCard>
        {/each}
      {/if}
    </div>
    <ErrorMessage message={newDocumentsError}></ErrorMessage>
  </div>
{/if}
