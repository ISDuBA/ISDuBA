<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { TableBodyCell } from "flowbite-svelte";
  import { searchColumnName } from "./defaults";
  import { advisorySearchState, getAdvisoryAnchorLink } from "$lib/Advisories/advisory.svelte";
  import Link from "$lib/Components/Link.svelte";
  import { splitMatches } from "$lib/utils";

  interface Props {
    doc: any;
    externalIndex: number;
    matches?: any[];
    index: number;
  }

  let { doc, externalIndex, matches = [], index }: Props = $props();
  let uid = $props.id();

  let expanded = $state(false);
  let visibleHits = $derived(expanded ? matches : matches.slice(0, 3));

  // These have to be the same as defined in the backend code.
  const delims = ["[!<", ">!]"];

  const getSplits = (text: string): string[] => {
    const matchRegex = /\[!<.*?>!\]/g;
    let regexMatches;
    const positions = [];
    while ((regexMatches = matchRegex.exec(text)) !== null) {
      positions.push([regexMatches.index, regexMatches["0"].length]);
    }
    return splitMatches(`${text}`, positions);
  };

  const getTableRowClass = (index: number) => {
    let tableRowClass = "";
    if (index === 0) {
      tableRowClass = tableRowClass + "border-t border-t-gray-300 dark:border-t-gray-600 ";
    }
    if (externalIndex % 2 == 1) {
      tableRowClass = tableRowClass + "bg-white dark:bg-gray-800";
    } else {
      tableRowClass = tableRowClass + "bg-gray-100 dark:bg-gray-700";
    }
    return tableRowClass;
  };
</script>

{#each visibleHits as match, i (`hitlist-${i}`)}
  {#if match[searchColumnName]}
    {@const splits = getSplits(match[searchColumnName])}
    <tr class={getTableRowClass(i)}>
      <TableBodyCell colspan={100} class="px-0 py-0 whitespace-nowrap">
        <Link
          ariaLabel={`Navigate directly to the ${i + 1}. match`}
          class={index % 2 == 1 ? "block" : "block"}
          href={getAdvisoryAnchorLink(doc)}
          onclick={() => {
            advisorySearchState.matchIndex = 0;
          }}
        >
          <span class="block px-2 py-1">
            {#each splits as s, index (`MatchList-${uid}-${index}`)}
              {#if (index + 1) % 2 === 0}
                {@const sanitized = s.replace(delims[0], "").replace(delims[1], "")}
                <span class="bg-yellow-200 dark:bg-yellow-800">{sanitized}</span>
              {:else}
                <span>{s}</span>
              {/if}
            {/each}
          </span>
        </Link>
      </TableBodyCell>
    </tr>
  {/if}
{/each}
