<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import Diff from "$lib/Diff/Diff.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { request } from "$lib/utils";
  import { Button, Label, Select, TabItem, Tabs } from "flowbite-svelte";
  import JsonDiff from "./JsonDiff.svelte";

  let documents = [];
  $: selectionOfDocuments = documents.map((doc) => {
    return {
      name: `${doc.tracking_id} - Version ${doc.version}`,
      value: {
        id: doc.id,
        version: doc.version
      }
    };
  });
  let diffDocuments: any;
  let docA: any;
  let docB: any;

  const compare = () => {
    if (docA && docB) {
      diffDocuments = {
        docA: docB,
        docB: docA
      };
    }
  };

  onMount(async () => {
    const documentURL = encodeURI(`/api/documents?limit=20&columns=id version tracking_id`);
    const response = await request(documentURL, "GET");
    if (response.ok) {
      documents = response.content.documents;
    }
  });
</script>

<SectionHeader title="Comparison"></SectionHeader>
<Tabs>
  <TabItem open title="JSON diff">
    <Label class="mb-6">
      Document 1:
      <Select id="firstDoc" bind:value={docA} items={selectionOfDocuments}></Select>
    </Label>
    <Label>
      Document 2:
      <Select id="secondDoc" bind:value={docB} items={selectionOfDocuments}></Select>
    </Label>
    <Button on:click={compare} class="my-2">Compare</Button>
    {#if diffDocuments}
      <JsonDiff {diffDocuments}></JsonDiff>
    {/if}
  </TabItem>
  <TabItem title="Git diff">
    <Diff></Diff>
  </TabItem>
</Tabs>
