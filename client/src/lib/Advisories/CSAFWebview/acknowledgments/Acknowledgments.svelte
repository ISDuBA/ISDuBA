<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Acknowledgment from "./Acknowledgment.svelte";
  import type { Acknowledgment as Acknowledgment2_0 } from "$lib/Advisories/types/csaf-2.0";
  import type { Acknowledgment as Acknowledgment2_1 } from "$lib/Advisories/types/csaf-2.1";
  interface Props {
    acknowledgments: Acknowledgment2_0[] | Acknowledgment2_1[];
  }
  let { acknowledgments }: Props = $props();
</script>

{#if acknowledgments}
  {#each acknowledgments as ack, i (`ack-${i}`)}
    <Acknowledgment {ack} path={`/document/acknowledgments[${i}]`} />
  {/each}
{/if}
