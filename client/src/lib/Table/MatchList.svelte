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

  const decodeText = (encoded: Uint8Array[]) => {
    let decodedText = "";
    const decoder = new TextDecoder();
    encoded.forEach((e) => {
      decodedText = decodedText + decoder.decode(e);
    });
    return decodedText;
  };

  const splitByDelimiterIndex = (index: number, encoded: Uint8Array[]) => {
    return [encoded.slice(0, index), encoded.slice(index)];
  };

  const findIndexOfDelim = (delim: Uint8Array, encoded: Uint8Array[]) => {
    const index = encoded.findIndex((t: Uint8Array) => {
      return t[0] === delim[0];
    });
    if (
      index !== -1 &&
      index < encoded.length - 3 &&
      encoded[index + 1][0] === delim[1] &&
      encoded[index + 2][0] === delim[2]
    ) {
      return index;
    }
    return -1;
  };

  const getSplits = (text: string): string[] => {
    const encoder = new TextEncoder();
    const encodedText: Uint8Array[] = [];
    for (let i = 0; i < text.length; i++) {
      encodedText.push(encoder.encode(text.charAt(i)));
    }

    const encodedDelim0 = encoder.encode(delims[0]);
    const encodedDelim1 = encoder.encode(delims[1]);

    const decodedSplits: string[] = [];
    let index0;
    let rest = encodedText;
    let count = 0;
    while (index0 !== -1 && count < 10) {
      index0 = findIndexOfDelim(encodedDelim0, rest);
      if (index0 !== -1) {
        const splitted0: Uint8Array[][] = splitByDelimiterIndex(index0, rest);
        const index1 = findIndexOfDelim(encodedDelim1, splitted0[1]);
        if (index1 !== -1) {
          const text0 = decodeText(splitted0[0]).replace(delims[0], "").replace(delims[1], "");
          decodedSplits.push(text0);
          const splitted1: Uint8Array[][] = splitByDelimiterIndex(index1, splitted0[1]);
          const text1 = decodeText(splitted1[0]).replace(delims[0], "").replace(delims[1], "");
          decodedSplits.push(text1);
          rest = splitted1[1];
        }
      }
      count++;
    }
    if (rest.length > 0) {
      const text2 = decodeText(rest).replace(delims[0], "").replace(delims[1], "");
      decodedSplits.push(text2);
    }

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
