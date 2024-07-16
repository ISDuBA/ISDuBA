<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Alert } from "flowbite-svelte";
  export let message = "";
</script>

{#if message}
  <div class="w-fit">
    <Alert color="red" dismissable>
      <span class="text-lg"> {message}</span>
    </Alert>
  </div>
{/if}
