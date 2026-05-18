<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Link from "$lib/Components/Link.svelte";
  import { request } from "$lib/request";
  import { onMount } from "svelte";

  let version: string | undefined = $state(undefined);

  let docLink = $derived.by(() => {
    const commitHash = version?.split("-g")?.[1];
    return `https://github.com/ISDuBA/ISDuBA/blob/${commitHash ?? "main"}/docs/search.md#filter-expressions`;
  });

  async function getVersion() {
    const response = await request("api/about", "GET");
    if (response.ok) {
      const backendInfo = response.content;
      version = backendInfo.version;
    }
  }

  onMount(() => {
    getVersion();
  });
</script>

<Link href={docLink} class="text-sm underline">
  <i class="bx bx-link"></i>
  <span>Documentation: Filter expression</span>
</Link>
