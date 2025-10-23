<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Message from "./Message.svelte";
  import { MESSAGE } from "./messagetypes";

  interface Props {
    message?: string;
  }

  let { message = "" }: Props = $props();
</script>

<div class="my-2">
  {#if message && message.length > 0}
    <Message dismissable={false} error={{ type: MESSAGE.ERROR, message: message }}></Message>
  {/if}
</div>
