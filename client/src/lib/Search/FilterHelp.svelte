<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { request } from "$lib/request";
  import { Spinner } from "flowbite-svelte";
  import { onMount } from "svelte";
  import DOMPurify from "dompurify";

  let html: string | null = $state(null);
  let isLoading = $state(false);
  let error: ErrorDetails | null = $state(null);

  const loadFilterHelp = async () => {
    isLoading = true;
    const response = await request("/api/documents/filter_help", "GET");
    if (response.ok) {
      let tmpHtml = DOMPurify.sanitize(response.content, { USE_PROFILES: { html: true } });
      tmpHtml = tmpHtml.replaceAll(":white_check_mark:", "<i class='bx bx-check'></i>");
      tmpHtml = tmpHtml.replaceAll(":x:", "<i class='bx bx-x'></i>");
      html = tmpHtml;
    } else {
      error = getErrorDetails("Could not load help", response);
    }
    isLoading = false;
  };

  onMount(() => {
    loadFilterHelp();
  });
</script>

{#if isLoading}
  <div class="flex h-full w-full items-center justify-center overflow-scroll">
    <Spinner color="gray" size="10"></Spinner>
  </div>
{/if}
<!-- eslint-disable-next-line svelte/no-at-html-tags -->
{@html html}
<ErrorMessage {error} />

