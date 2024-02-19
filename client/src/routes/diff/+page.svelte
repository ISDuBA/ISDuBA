<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import RouteGuard from "$lib/RouteGuard.svelte";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import Diff from "$lib/Diff.svelte";

  let diff: string;
  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      fetch("advisory.diff").then((response) => {
        response.text().then((text) => {
          diff = text;
        });
      });
    }
  });
</script>

<RouteGuard>
  <h1 class="text-lg">Comparison</h1>
  <Diff {diff}></Diff>
</RouteGuard>
