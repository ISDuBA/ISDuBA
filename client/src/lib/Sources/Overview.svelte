<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { Button, Label, Input } from "flowbite-svelte";
  import { request } from "$lib/utils";

  let domains = "";

  async function importDocuments() {
    const response = await request(`/api/download/?domains=${domains}`, "GET");
    if (response.ok) {
      console.log("Success");
    } else if (response.error) {
      console.log(response.error);
    }
  }
</script>

<SectionHeader title="Sources"></SectionHeader>

<div class="mb-6">
  <Label for="domain" class="mb-2 block">Domain</Label>
  <Input id="domain" bind:value={domains} placeholder="example.com" />
  <Button
    on:click={() => {
      importDocuments();
    }}
    class="ml-auto mt-auto"
    color="primary">Import</Button
  >
</div>
