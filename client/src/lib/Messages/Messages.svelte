<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store.svelte";
  import Message from "./Message.svelte";
</script>

<div style="position:absolute; bottom:1rem; right:1rem;">
  {#each appStore.state.app.errors as error (error.id)}
    <div class="mb-2">
      <Message {error}></Message>
    </div>
  {/each}
</div>
