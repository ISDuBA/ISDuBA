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
  import { getAdvisoryAnchorLink } from "$lib/Advisories/advisory";

  /* eslint-disable svelte/no-at-html-tags */

  interface Props {
    colspan: number;
    doc: any;
    hits?: any[];
    index: number;
  }

  let { colspan, doc, hits = [], index }: Props = $props();

  let expanded = $state(false);
  let visibleHits = $derived(expanded ? hits : hits.slice(0, 3));
</script>

{#each visibleHits as hit, i (`hitlist-${i}`)}
  {#if hit[searchColumnName]}
    <tr
      class={index % 2 == 1
        ? "border-t border-t-gray-200 bg-white dark:border-t-gray-700 dark:bg-gray-800"
        : "border-t border-t-gray-300 bg-gray-100 dark:border-t-gray-600 dark:bg-gray-700"}
    >
      <TableBodyCell {colspan} class="px-0 py-0 whitespace-nowrap">
        <a
          aria-label="View advisory details"
          class={index % 2 == 1
            ? "block hover:bg-gray-200 dark:hover:bg-gray-600"
            : "block hover:bg-gray-200 dark:hover:bg-gray-600"}
          href={getAdvisoryAnchorLink(doc)}
        >
          <span class="block px-2 py-1">
            {@html DOMPurify.sanitize(hit[searchColumnName], {
              USE_PROFILES: { html: true }
            })}
          </span>
        </a>
        {#if !expanded && visibleHits.length < hits.length && i === visibleHits.length - 1}
          <Button
            onclick={() => {
              expanded = true;
            }}
            class="m-1 h-6 w-fit px-2 py-1"
            color="light">More ({hits.slice(3).length})</Button
          >
        {/if}
      </TableBodyCell>
    </tr>
  {/if}
{/each}
