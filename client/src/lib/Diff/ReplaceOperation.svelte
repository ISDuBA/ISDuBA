<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";

  export let content: string;
  let parsedContent: object[];

  const parse = (textPart: any, level = 0): object[] => {
    const text = typeof textPart === "number" ? textPart.toString() : textPart;
    if (text.length === 0) {
      return [];
    }
    const firstIndexAdd = text.indexOf("{+");
    const secondIndexAdd = text.indexOf("+}");
    const firstIndexRemove = text.indexOf("[-");
    const secondIndexRemove = text.indexOf("-]");
    const parsed: object[] = [];
    if (
      firstIndexAdd > -1 &&
      secondIndexAdd > -1 &&
      (firstIndexRemove === -1 || firstIndexAdd < firstIndexRemove)
    ) {
      const firstSplit = text.split("{+");
      const secondSplit = firstSplit[1].split("+}");
      parsed.push({ type: "plain", content: firstSplit[0] });
      parsed.push({ type: "add", content: secondSplit[0] });
      parsed.push(...parse(text.slice(secondIndexAdd + 2), level + 1));
    } else if (firstIndexRemove > -1 && secondIndexRemove > -1) {
      const firstSplit = text.split("[-");
      const secondSplit = firstSplit[1].split("-]");
      parsed.push({ type: "plain", content: firstSplit[0] });
      parsed.push({ type: "remove", content: secondSplit[0] });
      parsed.push(...parse(text.slice(secondIndexRemove + 2), level + 1));
    } else {
      parsed.push({ type: "plain", content: text });
    }
    return parsed;
  };

  onMount(() => {
    parsedContent = parse(content);
  });

  const getSpanClass = (type: string) => {
    if (type === "add") return "bg-green-200";
    if (type === "remove") return "bg-red-200";
  };
</script>

<div>
  {#if parsedContent}
    {#each parsedContent as content}
      <span class={getSpanClass(content.type)}>{content.content}</span>
    {/each}
  {/if}
</div>
