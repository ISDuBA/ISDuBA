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
    const encoder = new TextEncoder();
    const encodedText: Uint8Array = encoder.encode(text);

    // Encode text and delimiters and join the arrays so we can search the delimiters via Regex
    const joinedEncodedText: string = encodedText.join("");
    const encodedDelim0 = encoder.encode(delims[0]);
    const encodedDelim1 = encoder.encode(delims[1]);
    const matchRegex = new RegExp(`${encodedDelim0.join("")}.*?${encodedDelim1.join("")}`, "g");
    // Find out positions of delimiters
    let regexMatches;
    const positions = [];
    while ((regexMatches = matchRegex.exec(joinedEncodedText)) !== null) {
      positions.push([regexMatches.index, regexMatches["0"].length]);
    }
    // Split with the help of the positions
    const encodedSplits: any[] = [];
    let lastPos = 0;
    for (let i = 0; i < positions.length; i++) {
      const pos = positions[i];
      const term = joinedEncodedText.slice(pos[0], pos[0] + pos[1]);
      // Don't use the term to split the text although it would be easier because the method could find
      // other occurrences that were not considered by the backend.
      encodedSplits.push(joinedEncodedText.slice(lastPos, pos[0]), term);
      lastPos = pos[0] + pos[1];
      if (i === positions.length - 1) {
        encodedSplits.push(joinedEncodedText.slice(pos[0] + pos[1]));
      }
    }

    // Since the joined arrays don't have the information about the length of each number
    // we take the original encoded text array go through it and check to which split each
    // integer belongs.
    let splitIndex = 0;
    const decodedSplits: string[] = [];
    let currentSplit = encodedSplits[0];
    let decoded = "";
    const decoder = new TextDecoder();
    encodedText.forEach((t, index) => {
      if (index === encodedText.length - 1) {
        decoded = decoded + decoder.decode(new Uint8Array([t]));
        decodedSplits.push(decoded);
      } else {
        if (!currentSplit.startsWith(`${t}`)) {
          decodedSplits.push(decoded);
          decoded = "";
          splitIndex++;
          if (encodedSplits[splitIndex]) {
            currentSplit = encodedSplits[splitIndex];
          }
        }
        decoded = decoded + decoder.decode(new Uint8Array([t]));
        currentSplit = currentSplit.replace(`${t}`, "");
      }
    });
    return decodedSplits;
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
