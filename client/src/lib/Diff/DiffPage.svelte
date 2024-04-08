<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Diff from "$lib/Diff/Diff.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";

  let diff: string;
  let error: string;
  onMount(async () => {
    if ($appStore.app.keycloak.authenticated) {
      error = "";
      const response = await request("advisory.diff", "GET");
      if (response.ok) {
        diff = response.content;
      } else if (response.error) {
        error = response.error;
      }
    }
  });
</script>

<SectionHeader title="Comparison"></SectionHeader>
<ErrorMessage message={error} plain={true}></ErrorMessage>
<Diff {diff}></Diff>
