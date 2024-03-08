<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Diff from "$lib/Diff/Diff.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";

  let diff: string;
  onMount(async () => {
    if ($appStore.app.keycloak.authenticated) {
      fetch("advisory.diff", {
        headers: {
          Authorization: `Bearer ${$appStore.app.keycloak.token}`
        }
      }).then((response) => {
        response.text().then((text) => {
          diff = text;
        });
      });
    }
  });
</script>

<SectionHeader title="Comparison"></SectionHeader>
<Diff {diff}></Diff>
