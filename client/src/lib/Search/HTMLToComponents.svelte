<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2026 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import Self from "./HTMLToComponents.svelte";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import {
    Heading,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell
  } from "flowbite-svelte";

  interface Props {
    element: Node;
  }

  let { element }: Props = $props();

  let tag = $derived(element.nodeName.toLowerCase());
</script>

{#snippet children(element: Node)}
  {#if element.childNodes?.length > 0}
    {#each element.childNodes as node, index (index)}
      {#if [node.ELEMENT_NODE.toString(), node.TEXT_NODE.toString()].includes(node.nodeType.toString())}
        <Self element={node} />
      {/if}
    {/each}
  {:else}
    {element.textContent}
  {/if}
{/snippet}

{#if element}
  {#if element.nodeType === Node.TEXT_NODE}
    {element.textContent}
  {:else if tag === "h1"}
    <SectionHeader title={element.textContent ?? ""}></SectionHeader>
  {:else if tag === "h2"}
    <Heading tag="h3" class="text-md">
      {element.textContent}
    </Heading>
  {:else if tag === "table"}
    <Table>
      {@render children(element)}
    </Table>
  {:else if tag === "thead"}
    <TableHead defaultRow={false}>
      {@render children(element)}
    </TableHead>
  {:else if tag === "th"}
    <TableHeadCell>
      {@render children(element)}
    </TableHeadCell>
  {:else if tag === "tbody"}
    <TableBody>
      {@render children(element)}
    </TableBody>
  {:else if tag === "tr"}
    <TableBodyRow>
      {@render children(element)}
    </TableBodyRow>
  {:else if tag === "td"}
    <TableBodyCell class="text-gray-700 dark:text-white">
      {@render children(element)}
    </TableBodyCell>
  {:else if tag === "p"}
    <p class="mb-4">
      {@render children(element)}
    </p>
  {:else if tag === "pre"}
    <pre
      class="my-2 bg-gray-100 px-2 py-3 dark:bg-gray-900 [&>code]:bg-gray-100 [&>code]:dark:bg-gray-900">{@render children(
        element
      )}</pre>
  {:else if tag === "code"}
    <code class="bg-gray-200 px-1 dark:bg-gray-700">{element.textContent}</code>
  {:else if element.childNodes?.length > 0}
    <svelte:element this={tag}>
      {#each element.childNodes as node, index (index)}
        <Self element={node} />
      {/each}
    </svelte:element>
  {/if}
{/if}
