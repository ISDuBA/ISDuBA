<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button } from "flowbite-svelte";
  import { push } from "svelte-spa-router";
  export let advisoryVersions: any;
  export let publisherNamespace: string;
  export let trackingID: string;
  export let selectedDocumentVersion: string;
  const navigateToVersion = (version: any) => {
    push(`/advisories/${publisherNamespace}/${trackingID}/documents/${version.id}`);
  };
</script>

<div class="my-2">
  <div class="flex items-center">
    <h2 class="mr-3">Advisory Versions</h2>
    <div class="flex flex-row-reverse gap-2">
      {#each advisoryVersions as version}
        <Button
          class="!p-2"
          disabled={selectedDocumentVersion === version.version}
          on:click={() => {
            navigateToVersion(version);
          }}
          color="light"
        >
          {version.version}
        </Button>
      {/each}
    </div>
  </div>
</div>
