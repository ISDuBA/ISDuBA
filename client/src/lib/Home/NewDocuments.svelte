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
  import { Card } from "flowbite-svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";

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
      console.log(documents);
    } else if (response.error) {
      newDocumentsError = `Could not load new documents. ${getErrorMessage(response.error)}`;
    }
  };

  const getPublisherAbbr = (publisher: string) => {
    switch (publisher) {
      case "Red Hat Product Security":
        return "RH";
      case "Siemens ProductCERT":
        return "SI";
      case "Bundesamt für Sicherheit in der Informationstechnik":
        return "BSI";
      case "SICK PSIRT":
        return "SCK";
    }
  };

  onMount(() => {
    loadDocuments();
  });
</script>

{#if $appStore.app.isUserLoggedIn}
  <div class="mt-8 flex flex-wrap gap-4">
    <div class="flex flex-col">
      <SectionHeader title="New documents"></SectionHeader>
      <div class="flex flex-wrap gap-4">
        <div class="text-red-600">
          Attention: These are
          <span class="font-bold">advisories</span>
          for now as we are not able to fetch recently imported documents.
        </div>
        <div class="flex flex-row flex-wrap gap-6">
          {#each documents as doc}
            <Card padding="md">
              <div class="mb-3 flex flex-row text-xs">
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
                <span class="ml-auto" title={doc.publisher}>{getPublisherAbbr(doc.publisher)}</span>
              </div>
              <hr class="mb-3" />
              <div class="text-black">{doc.title}</div>
            </Card>
          {/each}
        </div>
      </div>
      <ErrorMessage message={newDocumentsError}></ErrorMessage>
    </div>
    <div class="flex flex-col md:max-w-[50%]">
      <SectionHeader title="Recent activity"></SectionHeader>
    </div>
  </div>
{/if}
