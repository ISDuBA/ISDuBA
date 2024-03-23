<script lang="ts">
  import { onMount } from "svelte";

  export let content: string;
  let parsedContent: object[];

  const parse = (text: string, level = 0): object[] => {
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
