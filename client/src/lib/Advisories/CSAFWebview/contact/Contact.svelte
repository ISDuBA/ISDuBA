<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { cellStyleKey, cellStyleValue } from "$lib/Advisories/classes";
  import SearchableText from "$lib/Advisories/CSAFWebview/SearchableText.svelte";
  import type { Contact } from "$lib/Advisories/types/csaf-2.1";

  interface Props {
    contact: Contact;
    path: string;
  }
  let { contact, path }: Props = $props();

  const contactPath = $derived(`${path}/contact`);
</script>

{#if contact}
  {#if contact.details}
    <div class={cellStyleKey}>Publisher contact details</div>
    <div class={cellStyleValue}>
      <SearchableText text={contact.details} textPath={`${contactPath}/details`} />
    </div>
  {/if}
  {#if contact.email}
    <div class={cellStyleKey}>Publisher contact email</div>
    <div class={cellStyleValue}>
      <SearchableText text={contact.email} textPath={`${contactPath}/email`} />
    </div>
  {/if}
  {#if contact.public_openpgp_key_url}
    <div class={cellStyleKey}>Publisher contact OpenPGP key</div>
    <div class={cellStyleValue}>
      <SearchableText
        text={contact.public_openpgp_key_url}
        textPath={`${contactPath}/public_openpgp_key_url`}
      />
    </div>
  {/if}
  {#if contact.url}
    <div class={cellStyleKey}>Publisher contact URL</div>
    <div class={cellStyleValue}>
      <SearchableText text={contact.url} textPath={`${contactPath}/url`} />
    </div>
  {/if}
{/if}
