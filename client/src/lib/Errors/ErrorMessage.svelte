<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import type { ErrorDetails } from "$lib/Errors/error";
  import { Alert } from "flowbite-svelte";
  export let error: ErrorDetails | null;

  let showDetails: boolean = false;
</script>

{#if error}
  <div class="w-fit">
    <Alert color="red" dismissable>
      <span class="text-lg"> {error.message}</span>
      {#if error.details}
        <a href={"javascript:void(0);"} on:click={() => (showDetails = !showDetails)}>
          {#if showDetails}
            <i class="bx bx-chevron-up text-2xl"></i>
          {:else}
            <i class="bx bx-chevron-down text-2xl"></i>
          {/if}
        </a>
        {#if showDetails}
          <br />
          <span class="text-lg">{error.details}</span>
        {/if}
      {/if}
    </Alert>
  </div>
{/if}
