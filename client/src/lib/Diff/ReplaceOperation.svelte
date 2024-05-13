<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  export let content: string | number;
  export let isSideBySideViewActivated: boolean = true;
  $: parsedContent = isSideBySideViewActivated ? [] : parse(content);
  $: parsedSideBySideContent = isSideBySideViewActivated ? parseMixed() : [];

  const parse = (textPart: any, level = 0): any[] => {
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

  const removeAnnotation = (text: string, opening: string, ending: string) => {
    let areAnnotationsLeft = true;
    while (areAnnotationsLeft) {
      const firstIndexAdd = text.indexOf(opening);
      const secondIndexAdd = text.indexOf(ending);
      if (firstIndexAdd > -1 && secondIndexAdd > -1 && firstIndexAdd < secondIndexAdd) {
        text = text.replace(opening, "");
        text = text.replace(ending, "");
      } else {
        areAnnotationsLeft = false;
      }
    }
    return text;
  };

  const parseMixed = () => {
    const text = typeof content === "number" ? content.toString() : content;
    let added = text.replaceAll(/\[-.*?-]/g, "");
    added = removeAnnotation(added, "{+", "+}");
    let removed = text.replaceAll(/{+.*?\+}/g, "");
    removed = removeAnnotation(removed, "[-", "-]");
    return [removed, added];
  };

  const getSpanClass = (type: string) => {
    if (type === "add") return "bg-green-200";
    if (type === "remove") return "bg-red-200";
  };
</script>

<div>
  {#if parsedContent}
    {#each parsedContent as parsedPart}
      <span class={getSpanClass(parsedPart.type)}>{parsedPart.content}</span>
    {/each}
  {/if}
  {#if parsedSideBySideContent.length > 0}
    <div class="flex justify-between gap-2">
      <div class="flex w-6/12 items-center gap-1">
        <i class="bx bx-minus"></i>
        <div class="h-fit w-fit bg-red-200">{parsedSideBySideContent[0]}</div>
      </div>
      <div class="flex w-6/12 items-center gap-1">
        <i class="bx bx-plus"></i>
        <div class="h-fit w-fit bg-green-200">{parsedSideBySideContent[1]}</div>
      </div>
    </div>
  {/if}
</div>
