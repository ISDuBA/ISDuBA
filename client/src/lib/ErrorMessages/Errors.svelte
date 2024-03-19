<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Toast } from "flowbite-svelte";
  import { appStore } from "$lib/store";
  import { ERRORS } from "./messagetypes";

  const coloryByType = (type: string) => {
    if (type === ERRORS.ERROR) return "red";
    if (type === ERRORS.WARNING) return "yellow";
    return "green";
  };
</script>

{#each $appStore.app.errors as error}
  <Toast
    position="bottom-right"
    color={coloryByType(error.type)}
    on:close={() => {
      appStore.removeError(error.id);
    }}
  >
    <svelte:fragment slot="icon">
      <i class="bx bxs-error-alt"></i>
    </svelte:fragment>
    <span class="mb-1 text-sm font-semibold text-gray-900 dark:text-white">{error.type}</span>
    <div class="mb-2 text-sm font-normal">{error.message}</div>
  </Toast>
{/each}
