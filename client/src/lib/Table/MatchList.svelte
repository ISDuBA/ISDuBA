<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import DOMPurify from "dompurify";
  import { Button, TableBodyCell } from "flowbite-svelte";
  import { searchColumnName } from "./defaults";
  import { advisorySearchState, getAdvisoryAnchorLink } from "$lib/Advisories/advisory.svelte";
  import Link from "$lib/Components/Link.svelte";

  /* eslint-disable svelte/no-at-html-tags */

  interface Props {
    colspan: number;
    doc: any;
    matches?: any[];
    index: number;
  }

  let { colspan, doc, matches = [], index }: Props = $props();

  let expanded = $state(false);
  let visibleHits = $derived(expanded ? matches : matches.slice(0, 3));
</script>

{#each visibleHits as match, i (`hitlist-${i}`)}
  {#if match[searchColumnName]}
    <tr
      class={index % 2 == 1
        ? "border-t border-t-gray-200 bg-white dark:border-t-gray-700 dark:bg-gray-800"
        : "border-t border-t-gray-300 bg-gray-100 dark:border-t-gray-600 dark:bg-gray-700"}
    >
      <TableBodyCell {colspan} class="px-0 py-0 whitespace-nowrap">
        <Link
          ariaLabel={`Navigate directly to the ${i + 1}. match`}
          class={index % 2 == 1
            ? "block hover:bg-gray-200 dark:hover:bg-gray-600"
            : "block hover:bg-gray-200 dark:hover:bg-gray-600"}
          href={getAdvisoryAnchorLink(doc)}
          onclick={() => {
            advisorySearchState.matchIndex = i;
          }}
        >
          <span class="block px-2 py-1">
            {@html DOMPurify.sanitize(match[searchColumnName], {
              USE_PROFILES: { html: true }
            })}
          </span>
        </Link>
        {#if !expanded && visibleHits.length < matches.length && i === visibleHits.length - 1}
          <Button
            onclick={() => {
              expanded = true;
            }}
            class="m-1 h-6 w-fit px-2 py-1"
            color="light"
            title={`Show all matches for document ${visibleHits[0].tracking_id}`}
            >More ({matches.slice(3).length})</Button
          >
        {/if}
      </TableBodyCell>
    </tr>
  {/if}
{/each}
