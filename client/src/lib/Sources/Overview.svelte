<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import SectionHeader from "$lib/SectionHeader.svelte";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import { request } from "$lib/utils";

  let messageError: ErrorDetails | null;
  async function getMessage() {
    const response = await request("api/sources/message", "GET");
    if (response.ok) {
      return response.content;
    } else {
      messageError = getErrorDetails(`Couldn't load default message`, response);
    }
    return new Map<string, [string]>();
  }
</script>

<SectionHeader title="Sources"></SectionHeader>

{#await getMessage() then resp}
  {#if resp.message}
    {resp.message}
  {/if}
{/await}
<ErrorMessage error={messageError}></ErrorMessage>
